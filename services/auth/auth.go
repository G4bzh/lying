package main

import (
	"os"
	"fmt"
  "flag"
  "log"
  "time"
	"strings"
  "net/http"
  "crypto/sha256"
  "encoding/json"
  "encoding/hex"
  "io/ioutil"

  jwt "github.com/dgrijalva/jwt-go"
  "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
  "github.com/gorilla/mux"
	"github.com/gorilla/context"
)

//
// MongoDB Mapping
//

type User struct {
	Id					string				`bson:"_id" json:"id"`
	Hash				string				`bson:"hash" json:"password"`
}

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
// Login
//
func authLogin(w http.ResponseWriter, r *http.Request, db *string, col *string) {

  var u User

  w.Header().Set("Content-Type", "text/plain; charset=utf-8")

  // Get Post Data
  b, _ := ioutil.ReadAll(r.Body)
  if err := json.Unmarshal(b, &u); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Fprintf(w, "Invalid JSON")
    log.Printf("authLogin : %v", err)
    return
  }
  log.Printf("authLogin : Got %v", u)

  // Compute SHA256
  h := sha256.New()
  h.Write([]byte(u.Hash))
  log.Printf("authLogin: hash %s", hex.EncodeToString(h.Sum(nil)))

  // Get collection object
  c := session.DB(*db).C(*col)

  // Fetch user record
  if err := c.Find(bson.M{"_id" : u.Id}).Select(nil).One(&u); err != nil {
      w.WriteHeader(http.StatusNotFound)
      fmt.Fprintf(w, "Id not found or wrong password")
      log.Printf("authLogin : %v", err)
      return
    }
  log.Printf("authLogin : Got %v", u)

  // Compare hashes
  if ( u.Hash != hex.EncodeToString(h.Sum(nil)) ) {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, "Id not found or wrong password")
    log.Printf("authLogin : hashes mismatch")
    return
  }

  log.Printf("authLogin : %v login successfull", u.Id)

  // Create Claim for token
  claims := &jwt.StandardClaims {
    ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
    Issuer: "lying",
    Id: u.Id,
  }

  // Generate then sign token
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  tokenString, err := token.SignedString(signature)
  if (err != nil) {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Fprintf(w, "Token error")
    log.Printf("authLogin : %v", err)
    return
  }
  log.Printf("authLogin :Token %s", tokenString)
  w.Write([]byte(tokenString))

}

//
// Test
//


func middlewareAuth(next handlerF) handlerF {
  return handlerF(func(w http.ResponseWriter, r *http.Request) {

		// Get header
		hdr := r.Header.Get("Authorization")

		// Check header is present
		if ( hdr == "" )	{
			w.WriteHeader(http.StatusForbidden)
	    fmt.Fprintf(w, "Not Allowed")
	    log.Printf("authTest : No header")
	    return
		}

		// Look for Bearer
		hdrvalue := strings.Fields(hdr)
		if ( len(hdrvalue) != 2 ) {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "Not Allowed")
			log.Printf("authTest : Wrong header value %q", hdrvalue)
			return
		}

		if ( hdrvalue[0] != "Bearer" ) {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "Not Allowed")
			log.Printf("authTest : No Bearer keyword, Got %s", hdrvalue[0])
			return
		}

		// Finally, get token
		tokenString := hdrvalue[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				// Check signing method
				log.Printf("authTest : Unexpected signing method %v", token.Header["alg"])
				return nil, fmt.Errorf("")
		 	}
	    return signature, nil
		})

		if (err != nil)	{
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "Not Allowed")
			log.Printf("authTest : Unable to parse token  %v", err)
			return
		}

		// Check validity ans issuer
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if ( claims["iss"] == "lying" ){
				// OK, go to next handler
				context.Set(r, "clientID", claims["jti"])
				next(w, r)
				return
			}
		}

		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Not Allowed")
		log.Printf("authTest : Invalid  Token or Claims")

  })
}


func authTest(w http.ResponseWriter, r *http.Request, db *string, col *string) {

	log.Printf("authTest : request Allowed")
	fmt.Fprintf(w,"Hello %s",context.Get(r, "clientID"))

	return
}



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
  r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
      authLogin(w, r, dbPtr, colPtr)
    }).Methods("POST")

	r.HandleFunc("/login", middlewareAuth(func(w http.ResponseWriter, r *http.Request) {
      authTest(w, r, dbPtr, colPtr)
    })).Methods("GET")


	defer session.Close()
  log.Fatal(http.ListenAndServe(":8080", r))

}
