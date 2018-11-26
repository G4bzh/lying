package main

import (
	"fmt"
  "net/http"
  "log"
  "io/ioutil"
	"encoding/json"

	"github.com/globalsign/mgo/bson"
  "github.com/gorilla/context"
)


//
//  /{id}/forwarders handler
//
func GetForwarders(w http.ResponseWriter, r *http.Request, db *string, col *string) {
    // Retrieve ID
    id := context.Get(r, "clientID")

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


func SetForwarders(w http.ResponseWriter, r *http.Request, db *string, col *string) {
    // Retrieve ID
    id := context.Get(r, "clientID")

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		b, _ := ioutil.ReadAll(r.Body)
		var f Record

		if err := json.Unmarshal(b, &f); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "{\"msg\":\"Invalid JSON\"}")
			log.Printf("setForwarders for %s : %v", id, err)
			return
		}
		log.Printf("setForwarders for %s : Got %v", id, f)

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
