package main

import (
	"fmt"
  "flag"
  "log"
  "time"
  "net/http"
  "crypto/sha256"
  "encoding/json"
  "encoding/hex"
  "io/ioutil"

  jwt "github.com/dgrijalva/jwt-go"
  "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
  "github.com/gorilla/mux"
)

//
// MongoDB Mapping
//

type User struct {
	Id					string				`bson:"_id" json:"id"`
	Hash				string				`bson:"hash" json:"password"`
}


//
// Globals
//

var session *mgo.Session

//
// Login
//
func authLogin(w http.ResponseWriter, r *http.Request, db *string, col *string) {

  var u User
  var signature = []byte("secret")

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
// Main
//

func main() {

	var err error

	urlPtr := flag.String("url","127.0.0.1:27017","MongoDB URL")
	dbPtr := flag.String("db","auth","MongoDB Database")
	colPtr := flag.String("col","data","MongoDB Collection")

	flag.Parse()

  r := mux.NewRouter()

	// Connect to Database
	session, err = mgo.Dial(*urlPtr)
	if err != nil {
		panic(err)
	}
	log.Printf("dnscfg : Connected to %s\n", *urlPtr)
	// Read secondaries with consistence
	session.SetMode(mgo.Monotonic, true)



  r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
      authLogin(w, r, dbPtr, colPtr)
    }).Methods("POST")


	defer session.Close()
  log.Fatal(http.ListenAndServe(":8080", r))

}
