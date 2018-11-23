package main

import (
	"os"
	"fmt"
	"flag"
	"strings"
  "net/http"
  "log"

	"github.com/globalsign/mgo"
  "github.com/gorilla/mux"
	"github.com/gorilla/context"
	jwt "github.com/dgrijalva/jwt-go"

)


//
// Globals
//

var session *mgo.Session
var signature []byte
var issuer string

//
// Auth Middleware
//
func middlewareAuth(next HandlerF) HandlerF {
  return HandlerF(func(w http.ResponseWriter, r *http.Request) {

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



  r.HandleFunc("/v1/private/{id}/config/zones", func(w http.ResponseWriter, r *http.Request) {
      GetConfigZones(w, r, dbPtr, colPtr)
    }).Methods("GET")

  r.HandleFunc("/v1/private/{id}/config", func(w http.ResponseWriter, r *http.Request) {
      GetConfig(w, r, dbPtr, colPtr)
    }).Methods("GET")

  r.HandleFunc("/v1/private/{id}/config/zone/{zone}", func(w http.ResponseWriter, r *http.Request) {
      GetConfigZone(w, r,  dbPtr, colPtr)
    }).Methods("GET")

	r.HandleFunc("/v1/public/{id}", middlewareAuth(func(w http.ResponseWriter, r *http.Request) {
      GetID(w, r, dbPtr, colPtr)
    })).Methods("GET")

	r.HandleFunc("/v1/public/{id}", middlewareAuth(func(w http.ResponseWriter, r *http.Request) {
			GetID(w, r, dbPtr, colPtr)
		})).Methods("POST")

	r.HandleFunc("/v1/public/{id}", middlewareAuth(func(w http.ResponseWriter, r *http.Request) {
			RemoveID(w, r, dbPtr, colPtr)
		})).Methods("DELETE")

	r.HandleFunc("/v1/public/{id}/forwarders", middlewareAuth(func(w http.ResponseWriter, r *http.Request) {
      GetForwarders(w, r, dbPtr, colPtr)
    })).Methods("GET")

	r.HandleFunc("/v1/public/{id}/forwarders", middlewareAuth(func(w http.ResponseWriter, r *http.Request) {
      GetForwarders(w, r, dbPtr, colPtr)
    })).Methods("POST")

	r.HandleFunc("/v1/public/{id}/zones", middlewareAuth(func(w http.ResponseWriter, r *http.Request) {
      GetZones(w, r,  dbPtr, colPtr)
    })).Methods("GET")

	r.HandleFunc("/v1/public/{id}/zone/{zone}", middlewareAuth(func(w http.ResponseWriter, r *http.Request) {
      GetZone(w, r, dbPtr, colPtr)
    })).Methods("GET")

	r.HandleFunc("/v1/public/{id}/zone/{zone}", middlewareAuth(func(w http.ResponseWriter, r *http.Request) {
      SetZone(w, r, dbPtr, colPtr)
    })).Methods("POST")

	r.HandleFunc("/v1/public/{id}/zone/{zone}", middlewareAuth(func(w http.ResponseWriter, r *http.Request) {
      RemoveZone(w, r, dbPtr, colPtr)
    })).Methods("DELETE")

	defer session.Close()
  log.Fatal(http.ListenAndServe(":8053", r))

}
