package main

import (
	"fmt"
	"sort"
	"flag"
	"text/template"
  "net/http"
  "log"
  "encoding/json"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
  "github.com/gorilla/mux"
)


//
// Common MongoDB Mapping
//

type ZoneName struct {
  Domain      string        `bson:"domain" json:"zone"`
}

//
// getZone Mongodb Mapping
//
type RR	struct {
	Name  	string        `bson:"name"`
	Type  	string        `bson:"type"`
	Class  	string        `bson:"class"`
	TTL  		int			      `bson:"ttl"`
	Rdata  	string        `bson:"rdata"`
}

type RRs 	[]RR

type Zone struct {
	Domain  	string        `bson:"domain"`
	RRs				RRs 					`bson:"rr"`
}

type Zones 	[]Zone

type ZoneQuery 	struct {
	Id					string				`bson:"_id" json:"id"`
	Zones				Zones					`bson:"zones"`
}

//
// getConfig MongoDB Mapping
//
type ConfigQuery struct {
  Id					string				`bson:"_id"`
  Forwarders	[]string			`bson:"forwarders"`
  Zones       []ZoneName    `bson:"zones"`
}


//
// getZone MongoDB Mapping
//


type ZonesQuery struct {
  Id					string				`bson:"_id" json:"id"`
  Zones       []ZoneName    `bson:"zones"`
}


//
// Sort RR Slice
//
func (r RRs)rrSort(i,j int) bool {

	typeI := r[i].Type;
	typeJ := r[j].Type;

	if (typeI == "SOA")	{
		return true
	}

	if (typeJ == "SOA")	{
		return false
	}

  if (typeI == "NS")	{
    return true
  }

  if (typeJ == "NS")	{
    return false
  }

	return typeI < typeJ

}

//
// Templates
//
const zoneTmpl = `{{range .RRs}}
{{.Name}} {{.TTL}} {{.Class}} {{.Type}} {{.Rdata}}{{end}}
`

const configTmpl = `options {
        directory "/etc/bind";
        forwarders {
              {{range .Forwarders}}
              {{.}};{{end}}
        };

        query-source address * port 53;
        listen-on { any; };

        auth-nxdomain no;    # conform to RFC1035
        listen-on-v6 { any; };

        zone-statistics yes;

				{{range .Zones}}{{ if eq .Domain  "rpz" }}
        response-policy { zone "rpz"; };{{end}}{{end}}
};

logging {
      channel "querylog" { file "/var/log/dns-query.log"; print-time yes; print-category yes; print-severity yes; };
      category queries { querylog; };
};

{{range .Zones}}
zone "{{.Domain}}" {
    type master;
    file "{{.Domain}}.txt";
};
{{end}}
`

//
//  /{id}/zones handler
//
func getZones(w http.ResponseWriter, r *http.Request, url *string, db *string, col *string) {
    // Retrieve ID
    id := mux.Vars(r)["id"]

    w.Header().Set("Content-Type", "application/json; charset=utf-8")

    // Connect to Database
    session, err := mgo.Dial(*url)
    if err != nil {
      panic(err)
    }
    log.Printf("getZones for %s : Connected to %s\n", id, *url)
    // Read secondaries with consistence
    session.SetMode(mgo.Monotonic, true)
    defer session.Close()

    // Get collection object
    c := session.DB(*db).C(*col)

    // Query to get zones for id
    var z ZonesQuery
  	if err = c.Find(bson.M{"_id" : id}).Select(bson.M{"zones.domain": 1}).One(&z); err != nil {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "Id or Zones not found")
  			log.Printf("getZones for %s : %v", id, err)
        return
  		}
    log.Printf("getZones for %s : Got %v", id, z)

    // Translate back to JSON
    data, err := json.Marshal(z.Zones)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintf(w, "Failed to json marshal")
      log.Printf("getZones for %s : %v", id, err)
      return
    }

    // Send zones
    fmt.Fprintf(w, "%s", data)
}


//
//  /{id}/config handler
//
func getConfig(w http.ResponseWriter, r *http.Request, url *string, db *string, col *string) {
    // Retrieve ID
    id := mux.Vars(r)["id"]

    w.Header().Set("Content-Type", "text/plain; charset=utf-8")

    // Connect to Database
    session, err := mgo.Dial(*url)
    if err != nil {
      panic(err)
    }
    log.Printf("getConfig for %s : Connected to %s\n", id, *url)
    // Read secondaries with consistence
    session.SetMode(mgo.Monotonic, true)
    defer session.Close()

    // Get collection object
    c := session.DB(*db).C(*col)

    // Query to get config for id
    var cf ConfigQuery
  	if err = c.Find(bson.M{"_id" : id}).Select(bson.M{"zones.domain": 1, "forwarders": 1}).One(&cf); err != nil {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "Id or Zones not found")
  			log.Printf("getConfig for %s : %v", id, err)
        return
  		}
    log.Printf("getConfig for %s : Got %v", id, cf)

    // Render template to ResponseWriter
    tmpl := template.Must(template.New("config").Parse(configTmpl))
  	if err = tmpl.Execute(w, cf); err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintf(w, "Failed to execute template")
      log.Printf("getZones for %s : %v", id, err)
      return
  	}
}


//
//  /{id}/zone{zone} handler
//
func getZone(w http.ResponseWriter, r *http.Request, url *string, db *string, col *string) {
    // Retrieve ID & Zone
    id := mux.Vars(r)["id"]
    zone := mux.Vars(r)["zone"]

    w.Header().Set("Content-Type", "text/plain; charset=utf-8")

    // Connect to Database
    session, err := mgo.Dial(*url)
    if err != nil {
      panic(err)
    }
    log.Printf("getZone %s for %s : Connected to %s\n", zone, id, *url)
    // Read secondaries with consistence
    session.SetMode(mgo.Monotonic, true)
    defer session.Close()

    // Get collection object
    c := session.DB(*db).C(*col)

    // Query to get zone detail for id
    var z ZoneQuery
    if err = c.Find(bson.M{"_id" : id, "zones.domain" : zone}).Select(bson.M{"zones.$":1}).One(&z); err != nil {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "Id or Zones not found")
        log.Printf("getZone %s for %s : %v", zone, id, err)
        return
      }
    log.Printf("getZone %s for %s : Got data", zone, id)

    for _, d := range z.Zones {
      // Sort data
  		sort.SliceStable(d.RRs,d.RRs.rrSort)

      // Render template to ResponseWriter
      tmpl := template.Must(template.New("zone").Parse(zoneTmpl))
  		if err = tmpl.Execute(w, d); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Failed to execute template")
        log.Printf("getZone %s for %s : %v", zone, id, err)
        return
  		}

  	}
}



//
// Main
//
func main() {

	urlPtr := flag.String("url","127.0.0.1:27017","MongoDB URL")
	dbPtr := flag.String("db","saas","MongoDB Database")
	colPtr := flag.String("col","data","MongoDB Collection")

	flag.Parse()

  r := mux.NewRouter()

  r.HandleFunc("/{id}/zones", func(w http.ResponseWriter, r *http.Request) {
      getZones(w, r, urlPtr, dbPtr, colPtr)
    }).Methods("GET")

  r.HandleFunc("/{id}/config", func(w http.ResponseWriter, r *http.Request) {
      getConfig(w, r, urlPtr, dbPtr, colPtr)
    }).Methods("GET")

  r.HandleFunc("/{id}/zone/{zone}", func(w http.ResponseWriter, r *http.Request) {
      getZone(w, r, urlPtr, dbPtr, colPtr)
    }).Methods("GET")


  log.Fatal(http.ListenAndServe(":8080", r))

}
