package main

import (
	"fmt"
	"sort"
	"text/template"
  "net/http"
  "log"

	"github.com/globalsign/mgo/bson"
  "github.com/gorilla/mux"
)

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
      channel "querylog" { stderr; print-time yes; print-category yes; print-severity yes; };
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
//  /{id}/config handler
//
func GetConfig(w http.ResponseWriter, r *http.Request, db *string, col *string) {
    // Retrieve ID
    id := mux.Vars(r)["id"]

    w.Header().Set("Content-Type", "text/plain; charset=utf-8")

    // Get collection object
    c := session.DB(*db).C(*col)

    // Query to get config for id
    var cf Record
  	if err := c.Find(bson.M{"_id" : id}).Select(bson.M{"zones.domain": 1, "forwarders": 1}).One(&cf); err != nil {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "Id or Zones not found")
  			log.Printf("getConfig for %s : %v", id, err)
        return
  		}
    log.Printf("getConfig for %s : Got %v", id, cf)

    // Render template to ResponseWriter
    tmpl := template.Must(template.New("config").Parse(configTmpl))
  	if err := tmpl.Execute(w, cf); err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintf(w, "Failed to execute template")
      log.Printf("getConfig for %s : %v", id, err)
      return
  	}
}


//
//  /{id}/config/zones handler
//
func GetConfigZones(w http.ResponseWriter, r *http.Request, db *string, col *string) {
    // Retrieve ID
    id := mux.Vars(r)["id"]

    w.Header().Set("Content-Type", "text/plain; charset=utf-8")

    // Get collection object
    c := session.DB(*db).C(*col)

    // Query to get zones for id
    var z Record
  	if err := c.Find(bson.M{"_id" : id}).Select(bson.M{"zones.domain": 1}).One(&z); err != nil {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "Id or Zones not found")
  			log.Printf("getConfigZones for %s : %v", id, err)
        return
  		}
    log.Printf("getConfigZones for %s : Got %v", id, z)

    // Send zones
		for _, d := range z.Zones {
    	fmt.Fprintf(w, "%s\n", d.Domain)
		}
}


//
//  /{id}/config/zone/{zone} handler
//
func GetConfigZone(w http.ResponseWriter, r *http.Request, db *string, col *string) {
    // Retrieve ID & Zone
    id := mux.Vars(r)["id"]
    zone := mux.Vars(r)["zone"]

    w.Header().Set("Content-Type", "text/plain; charset=utf-8")

    // Get collection object
    c := session.DB(*db).C(*col)

    // Query to get zone detail for id
    var z Record
    if err := c.Find(bson.M{"_id" : id, "zones.domain" : zone}).Select(bson.M{"zones.$":1}).One(&z); err != nil {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "Id or Zones not found")
        log.Printf("getConfigZone %s for %s : %v", zone, id, err)
        return
      }
    log.Printf("getConfigZone %s for %s : Got data", zone, id)

    for _, d := range z.Zones {
      // Sort data
  		sort.SliceStable(d.RRs,d.RRs.rrSort)

      // Render template to ResponseWriter
      tmpl := template.Must(template.New("zone").Parse(zoneTmpl))
  		if err := tmpl.Execute(w, d); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Failed to execute template")
        log.Printf("getConfigZone %s for %s : %v", zone, id, err)
        return
  		}

  	}
}
