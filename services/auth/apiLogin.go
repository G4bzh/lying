package main

import (
	"fmt"
  "log"
  "time"
  "net/http"
  "crypto/sha256"
  "encoding/json"
  "encoding/hex"
  "io/ioutil"

  jwt "github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
)


//
// Login
//
func GetLogin(w http.ResponseWriter, r *http.Request, db *string, col *string) {

  var u User

  w.Header().Set("Content-Type", "application/json; charset=utf-8")

  // Get Post Data
  b, _ := ioutil.ReadAll(r.Body)
  if err := json.Unmarshal(b, &u); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Fprintf(w, "{\"msg\":\"Invalid JSON\"}")
    log.Printf("setLogin : %v", err)
    return
  }
  log.Printf("setLogin : Got %v", u)

  // Compute SHA256
  h := sha256.New()
  h.Write([]byte(u.Hash))
  log.Printf("setLogin: hash %s", hex.EncodeToString(h.Sum(nil)))

  // Get collection object
  c := session.DB(*db).C(*col)

  // Fetch user record
  if err := c.Find(bson.M{"_id" : u.Id}).Select(nil).One(&u); err != nil {
      w.WriteHeader(http.StatusNotFound)
      fmt.Fprintf(w, "{\"msg\":\"Id not found or wrong password\"}")
      log.Printf("setLogin : %v", err)
      return
    }
  log.Printf("setLogin : Got %v", u)

  // Compare hashes
  if ( u.Hash != hex.EncodeToString(h.Sum(nil)) ) {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, "{\"msg\":\"Id not found or wrong password\"}")
    log.Printf("setLogin : hashes mismatch")
    return
  }

  log.Printf("setLogin : %v login successfull", u.Id)

  // Create Claim for token
  claims := &jwt.StandardClaims {
    ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
    Issuer: issuer,
    Id: u.Id,
  }

  // Generate then sign token
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  tokenString, err := token.SignedString(signature)
  if (err != nil) {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Fprintf(w, "{\"msg\":\"Token error\"}")
    log.Printf("setLogin : %v", err)
    return
  }
  log.Printf("setLogin :Token %s", tokenString)
	fmt.Fprintf(w, "{\"token\":\"%s\"}", tokenString)

}

//
// Logup
//
func SetLogin(w http.ResponseWriter, r *http.Request, db *string, col *string) {

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  log.Printf("GetLogin : Not Implemented")
	fmt.Fprintf(w, "{\"msg\":\"Not Implemented\"}")

}
