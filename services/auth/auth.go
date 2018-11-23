package main

import (
	"os"
  "flag"
  "log"
  "net/http"

  "github.com/globalsign/mgo"
  "github.com/gorilla/mux"
)




//
// Globals
//

var session *mgo.Session
var signature []byte
var issuer string




//
// Main
//

func main() {

	var err error
	var env string

	urlPtr := flag.String("url","127.0.0.1:27017","MongoDB URL")
	dbPtr := flag.String("db","auth","MongoDB Database")
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


	// Routes
  r.HandleFunc("/v1/login", func(w http.ResponseWriter, r *http.Request) {
      GetLogin(w, r, dbPtr, colPtr)
    }).Methods("POST")

	r.HandleFunc("/v1/login", func(w http.ResponseWriter, r *http.Request) {
      SetLogin(w, r, dbPtr, colPtr)
    }).Methods("PUT")


	defer session.Close()
  log.Fatal(http.ListenAndServe(":8080", r))

}
