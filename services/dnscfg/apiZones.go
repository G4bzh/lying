package main

import (
	"fmt"
  "net/http"
  "log"
  "io/ioutil"
	"encoding/json"

	"github.com/globalsign/mgo/bson"
  "github.com/gorilla/context"
  "github.com/gorilla/mux"
)


//
//  /{id}/zones handler
//
func GetZones(w http.ResponseWriter, r *http.Request, db *string, col *string) {
    // Retrieve ID
    id := mux.Vars(r)["id"]

    w.Header().Set("Content-Type", "application/json; charset=utf-8")

		// Check ClientID from authMiddleware
		if ( context.Get(r, "clientID") != id ) {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "{\"msg\":\"Not Allowed\"}")
			log.Printf("getZones : ClientID Mismatch (%s != %s)", context.Get(r, "clientID"), id)
			return
		}

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
func GetZone(w http.ResponseWriter, r *http.Request, db *string, col *string) {
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


func SetZone(w http.ResponseWriter, r *http.Request, db *string, col *string) {
    // Retrieve ID & zone
    id := mux.Vars(r)["id"]
		zone := mux.Vars(r)["zone"]

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		// Check ClientID from authMiddleware
		if ( context.Get(r, "clientID") != id ) {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "{\"msg\":\"Not Allowed\"}")
			log.Printf("setZone : ClientID Mismatch (%s != %s)", context.Get(r, "clientID"), id)
			return
		}

		b, _ := ioutil.ReadAll(r.Body)
		var u Zone

		if err := json.Unmarshal(b, &u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "{\"msg\":\"Invalid JSON\"}")
			log.Printf("setZone %s for %s : %v", zone, id, err)
			return
		}
		log.Printf("setZone %s for %s : Got %v", zone, id, u)


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



func RemoveZone(w http.ResponseWriter, r *http.Request, db *string, col *string) {
    // Retrieve ID & zone
    id := mux.Vars(r)["id"]
		zone := mux.Vars(r)["zone"]

    w.Header().Set("Content-Type", "application/json; charset=utf-8")

		// Check ClientID from authMiddleware
		if ( context.Get(r, "clientID") != id ) {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "{\"msg\":\"Not Allowed\"}")
			log.Printf("removeZone : ClientID Mismatch (%s != %s)", context.Get(r, "clientID"), id)
			return
		}

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
