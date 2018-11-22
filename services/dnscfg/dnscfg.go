package main

import (
	"os"
	"fmt"
	"sort"
	"flag"
	"strings"
	"text/template"
  "net/http"
  "log"
	"io/ioutil"
	"encoding/json"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
  "github.com/gorilla/mux"
	"github.com/gorilla/context"
	jwt "github.com/dgrijalva/jwt-go"

)



//
//  Mongodb Mapping
//

type RR	struct {
	Name  	string        `bson:"name" json:"name"`
	Type  	string        `bson:"type" json:"type"`
	Class  	string        `bson:"class" json:"class"`
	TTL  		int			      `bson:"ttl" json:"ttl"`
	Rdata  	string        `bson:"rdata" json:"rdata"`
}

type RRs 	[]RR

type Zone struct {
	Domain  	string        `bson:"domain" json:"domain"`
	RRs				RRs 					`bson:"rr" json:"rrs,omitempty"`
}

type Zones 	[]Zone

type Record 	struct {
	Id					string				`bson:"_id" json:"id"`
	Forwarders	[]string			`bson:"forwarders" json:"forwarders`
	Zones				Zones					`bson:"zones" json:"zones,omitempty"`
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
// Types
//

type handlerF func(w http.ResponseWriter, r *http.Request)


//
// Globals
//

var session *mgo.Session
var signature []byte
var issuer string

//
// Auth Middleware
//
func middlewareAuth(next handlerF) handlerF {
  return handlerF(func(w http.ResponseWriter, r *http.Request) {

		// Get header
		hdr := r.Header.Get("Authorization")

		// Check header is present
		if ( hdr == "" )	{
			w.WriteHeader(http.StatusForbidden)
	    fmt.Fprintf(w, "Not Allowed")
	    log.Printf("middlewareAuth : No header")
	    return
		}

		// Look for Bearer
		hdrvalue := strings.Fields(hdr)
		if ( len(hdrvalue) != 2 ) {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "Not Allowed")
			log.Printf("middlewareAuth : Wrong header value %q", hdrvalue)
			return
		}

		if ( hdrvalue[0] != "Bearer" ) {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "Not Allowed")
			log.Printf("middlewareAuth : No Bearer keyword, Got %s", hdrvalue[0])
			return
		}

		// Finally, get token
		tokenString := hdrvalue[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				// Check signing method
				log.Printf("middlewareAuth : Unexpected signing method %v", token.Header["alg"])
				return nil, fmt.Errorf("")
		 	}
	    return signature, nil
		})

		if (err != nil)	{
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "Not Allowed")
			log.Printf("middlewareAuth : Unable to parse token  %v", err)
			return
		}

		// Check validity ans issuer
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if ( claims["iss"] == issuer ){
				// OK, go to next handler
				context.Set(r, "clientID", claims["jti"])
				log.Printf("middlewareAuth : Auth successfull")
				next(w, r)
				return
			}
		}

		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Not Allowed")
		log.Printf("middlewareAuth : Invalid  Token or Claims")

  })
}

//
//  /{id}/config/zones handler
//
func getConfigZones(w http.ResponseWriter, r *http.Request, db *string, col *string) {
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
  			log.Printf("getZones for %s : %v", id, err)
        return
  		}
    log.Printf("getZones for %s : Got %v", id, z)

    // Send zones
		for _, d := range z.Zones {
    	fmt.Fprintf(w, "%s\n", d.Domain)
		}
}


//
//  /{id}/config handler
//
func getConfig(w http.ResponseWriter, r *http.Request, db *string, col *string) {
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
//  /{id}/config/zone/{zone} handler
//
func getConfigZone(w http.ResponseWriter, r *http.Request, db *string, col *string) {
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
        log.Printf("getZone %s for %s : %v", zone, id, err)
        return
      }
    log.Printf("getZone %s for %s : Got data", zone, id)

    for _, d := range z.Zones {
      // Sort data
  		sort.SliceStable(d.RRs,d.RRs.rrSort)

      // Render template to ResponseWriter
      tmpl := template.Must(template.New("zone").Parse(zoneTmpl))
  		if err := tmpl.Execute(w, d); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Failed to execute template")
        log.Printf("getZone %s for %s : %v", zone, id, err)
        return
  		}

  	}
}

//
//  /{id} handler
//
func getID(w http.ResponseWriter, r *http.Request, db *string, col *string) {
	// Retrieve ID
	id := mux.Vars(r)["id"]

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Get collection object
	c := session.DB(*db).C(*col)

	// Query to get zone detail for id
	var i Record
	if err := c.Find(bson.M{"_id" : id}).Select(bson.M{"_id":1}).One(&i); err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "{\"msg\":\"Id not found\"}")
		log.Printf("getID %s : %v", id, err)
		return
	}

	fmt.Fprintf(w, "{\"msg\":\"OK\"}")
	log.Printf("getID %s : %v", id, i)
}


//
//  /{id} handler
//
func setID(w http.ResponseWriter, r *http.Request, db *string, col *string) {
    // Retrieve ID & Zone
    id := mux.Vars(r)["id"]

    w.Header().Set("Content-Type", "application/json; charset=utf-8")

    // Get collection object
    c := session.DB(*db).C(*col)

    // Create empty document
    var i Record
		i.Id = id
    if err := c.Insert(&i); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "{\"msg\":\"Error inserting Id\"}")
        log.Printf("setID %s : %v", id, err)
				return
      }

		fmt.Fprintf(w, "{\"msg\":\"OK\"}")
    log.Printf("setID %s : %v", id, i)
}

//
//  /{id} handler
//
func removeID(w http.ResponseWriter, r *http.Request, db *string, col *string) {
    // Retrieve ID & Zone
    id := mux.Vars(r)["id"]

    w.Header().Set("Content-Type", "application/json; charset=utf-8")

    // Get collection object
    c := session.DB(*db).C(*col)

    // Remove document
    if err := c.Remove(bson.M{"_id" : id}); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "{\"msg\":\"Error removing Id\"}")
        log.Printf("removeID %s : %v", id, err)
				return
      }

		fmt.Fprintf(w, "{\"msg\":\"OK\"}")
    log.Printf("removeID %s : OK", id)
}


//
//  /{id}/forwarders handler
//
func getForwarders(w http.ResponseWriter, r *http.Request, db *string, col *string) {
    // Retrieve ID
    id := mux.Vars(r)["id"]

    w.Header().Set("Content-Type", "application/json; charset=utf-8")

    // Get collection object
    c := session.DB(*db).C(*col)

    // Query to get forwarders for id
    var f Record
  	if err := c.Find(bson.M{"_id" : id}).Select(bson.M{"forwarders": 1}).One(&f); err != nil {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "{\"msg\":\"Id not found\"}")
  			log.Printf("getForwarders for %s : %v", id, err)
        return
  		}
    log.Printf("getForwarders for %s : Got %v", id, f)

    // Send forwarders
		b, err := json.Marshal(f.Forwarders)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "{\"msg\":\"Invalid JSON\"}")
			log.Printf("getForwarders for %s : %v", id, err)
			return
		}
    fmt.Fprintf(w, "%s", string(b))

}


func setForwarders(w http.ResponseWriter, r *http.Request, db *string, col *string) {
    // Retrieve ID
    id := mux.Vars(r)["id"]
		b, _ := ioutil.ReadAll(r.Body)
		var f Record

		if err := json.Unmarshal(b, &f); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "{\"msg\":\"Invalid JSON\"}")
			log.Printf("setForwarders for %s : %v", id, err)
			return
		}
		log.Printf("setForwarders for %s : Got %v", id, f)


    w.Header().Set("Content-Type", "application/json; charset=utf-8")

    // Get collection object
    c := session.DB(*db).C(*col)

    // Update forwarders for id
  	if err := c.Update(bson.M{"_id" : id}, bson.M{"$set": bson.M{"forwarders": f.Forwarders}}); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "{\"msg\":\"Error setting forwarders\"}")
  			log.Printf("setForwarders for %s : %v", id, err)
        return
  		}
    log.Printf("setForwarders for %s : set %v", id, f.Forwarders)
    fmt.Fprintf(w, "{\"msg\":\"OK\"}")

}

//
//  /{id}/zones handler
//
func getZones(w http.ResponseWriter, r *http.Request, db *string, col *string) {
    // Retrieve ID
    id := mux.Vars(r)["id"]

    w.Header().Set("Content-Type", "application/json; charset=utf-8")

    // Get collection object
    c := session.DB(*db).C(*col)

    // Query to get zones for id
    var z Record
  	if err := c.Find(bson.M{"_id" : id}).Select(bson.M{"zones.domain": 1}).One(&z); err != nil {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "{\"msg\":\"Id not found\"}")
  			log.Printf("getZones for %s : %v", id, err)
        return
  		}
    log.Printf("getZones for %s : Got %v", id, z)

    // Send Zones
		b, err := json.Marshal(z.Zones)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "{\"msg\":\"Invalid JSON\"}")
			log.Printf("getZones for %s : %v", id, err)
			return
		}
    fmt.Fprintf(w, "%s", string(b))

}

//
//  /{id}/zone/{zone} handler
//
func getZone(w http.ResponseWriter, r *http.Request, db *string, col *string) {
    // Retrieve ID & zone
    id := mux.Vars(r)["id"]
		zone := mux.Vars(r)["zone"]

    w.Header().Set("Content-Type", "application/json; charset=utf-8")

		// Check ClientID from authMiddleware
		if ( context.Get(r, "clientID") != id ) {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "{\"msg\":\"Not Allowed\"}")
			log.Printf("getZone : ClientID Mismatch (%s != %s)", context.Get(r, "clientID"), id)
			return
		}

    // Get collection object
    c := session.DB(*db).C(*col)

    // Query to get forwarders for id
    var z Record
  	if err := c.Find(bson.M{"_id" : id, "zones.domain" : zone}).Select(bson.M{"zones.$":1}).One(&z); err != nil {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "{\"msg\":\"Id or Zone not found\"}")
  			log.Printf("getZone %s for %s : %v", zone, id, err)
        return
  		}
    log.Printf("getZone %s for %s : Got %v", zone, id, z)

    // Send zone[0]
		b, err := json.Marshal(z.Zones[0])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "{\"msg\":\"Invalid JSON\"}")
			log.Printf("getZone %s for %s : %v", zone, id, err)
			return
		}
    fmt.Fprintf(w, "%s", string(b))

}


func setZone(w http.ResponseWriter, r *http.Request, db *string, col *string) {
    // Retrieve ID & zone
    id := mux.Vars(r)["id"]
		zone := mux.Vars(r)["zone"]

		b, _ := ioutil.ReadAll(r.Body)
		var u Zone

		if err := json.Unmarshal(b, &u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "{\"msg\":\"Invalid JSON\"}")
			log.Printf("setZone %s for %s : %v", zone, id, err)
			return
		}
		log.Printf("setZone %s for %s : Got %v", zone, id, u)


    w.Header().Set("Content-Type", "application/json; charset=utf-8")

    // Get collection object
    c := session.DB(*db).C(*col)

    // Query to remove zone for id if it exists
  	if err := c.Update(bson.M{"_id" : id}, bson.M{"$pull": bson.M{"zones": bson.M{"domain": u.Domain}}}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintf(w, "{\"msg\":\"Error while pulling data\"}")
			log.Printf("setZone %s for %s : %v", zone, id, err)
      return
		}

		// Query to add zone for id
		if err := c.Update(bson.M{"_id" : id}, bson.M{"$push": bson.M{"zones": u}}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintf(w, "{\"msg\":\"Error while pushing data\"}")
			log.Printf("setZone %s for %s : %v", zone, id, err)
      return
		}

    log.Printf("setZone %s for %s : set %v", zone, id, u)
    fmt.Fprintf(w, "%s", string(b))

}



func removeZone(w http.ResponseWriter, r *http.Request, db *string, col *string) {
    // Retrieve ID & zone
    id := mux.Vars(r)["id"]
		zone := mux.Vars(r)["zone"]

    w.Header().Set("Content-Type", "application/json; charset=utf-8")

    // Get collection object
    c := session.DB(*db).C(*col)

    // Query to remove zone for id if it exists
  	if err := c.Update(bson.M{"_id" : id}, bson.M{"$pull": bson.M{"zones": bson.M{"domain": zone}}}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintf(w, "{\"msg\":\"Error while pulling data\"}")
			log.Printf("removeZone %s for %s : %v", zone, id, err)
      return
		}

    log.Printf("removeZone %s for %s :OK", zone, id)
    fmt.Fprintf(w, "{\"msg\":\"OK\"}")

}

//
// Main
//
func main() {

	var err error
	var env string

	urlPtr := flag.String("url","127.0.0.1:27017","MongoDB URL")
	dbPtr := flag.String("db","dnscfg","MongoDB Database")
	colPtr := flag.String("col","data","MongoDB Collection")

	flag.Parse()

	if env = os.Getenv("JWT_ISSUER"); env == "" {
		panic("No JWT_ISSUER env variable.")
		return
	}
	issuer = env

	if env = os.Getenv("JWT_SIGNATURE"); env == "" {
		panic("No JWT_SIGNATURE env variable.")
		return
	}
	signature = []byte(env)

  r := mux.NewRouter()

	// Connect to Database
	session, err = mgo.Dial(*urlPtr)
	if err != nil {
		panic(err)
	}
	log.Printf("dnscfg : Connected to %s\n", *urlPtr)
	// Read secondaries with consistence
	session.SetMode(mgo.Monotonic, true)



  r.HandleFunc("/{id}/config/zones", func(w http.ResponseWriter, r *http.Request) {
      getConfigZones(w, r, dbPtr, colPtr)
    }).Methods("GET")

  r.HandleFunc("/{id}/config", func(w http.ResponseWriter, r *http.Request) {
      getConfig(w, r, dbPtr, colPtr)
    }).Methods("GET")

  r.HandleFunc("/{id}/config/zone/{zone}", func(w http.ResponseWriter, r *http.Request) {
      getConfigZone(w, r,  dbPtr, colPtr)
    }).Methods("GET")

	r.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
      getID(w, r, dbPtr, colPtr)
    }).Methods("GET")

	r.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
			setID(w, r, dbPtr, colPtr)
		}).Methods("POST")

	r.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
			removeID(w, r, dbPtr, colPtr)
		}).Methods("DELETE")

	r.HandleFunc("/{id}/forwarders", func(w http.ResponseWriter, r *http.Request) {
      getForwarders(w, r, dbPtr, colPtr)
    }).Methods("GET")

	r.HandleFunc("/{id}/forwarders", func(w http.ResponseWriter, r *http.Request) {
      setForwarders(w, r, dbPtr, colPtr)
    }).Methods("POST")

	r.HandleFunc("/{id}/zones", func(w http.ResponseWriter, r *http.Request) {
      getZones(w, r,  dbPtr, colPtr)
    }).Methods("GET")

	r.HandleFunc("/{id}/zone/{zone}", middlewareAuth(func(w http.ResponseWriter, r *http.Request) {
      getZone(w, r, dbPtr, colPtr)
    })).Methods("GET")

	r.HandleFunc("/{id}/zone/{zone}", func(w http.ResponseWriter, r *http.Request) {
      setZone(w, r, dbPtr, colPtr)
    }).Methods("POST")

	r.HandleFunc("/{id}/zone/{zone}", func(w http.ResponseWriter, r *http.Request) {
      removeZone(w, r, dbPtr, colPtr)
    }).Methods("DELETE")

	defer session.Close()
  log.Fatal(http.ListenAndServe(":8053", r))

}
