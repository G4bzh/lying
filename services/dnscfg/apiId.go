package main

import (
	"fmt"
  "net/http"
  "log"

	"github.com/globalsign/mgo/bson"
  "github.com/gorilla/context"
)


//
//  /{id} handler
//
func GetID(w http.ResponseWriter, r *http.Request, db *string, col *string) {
	// Retrieve ID
	id := context.Get(r, "clientID")

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
func SetID(w http.ResponseWriter, r *http.Request, db *string, col *string) {
    // Retrieve ID & Zone
    id := context.Get(r, "clientID")

    w.Header().Set("Content-Type", "application/json; charset=utf-8")

		// Get collection object
    c := session.DB(*db).C(*col)

    // Create empty document
    var i Record
		i.Id = id.(string)
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
func RemoveID(w http.ResponseWriter, r *http.Request, db *string, col *string) {
    // Retrieve ID & Zone
    id := context.Get(r, "clientID")

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
