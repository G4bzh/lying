package main

import (
	"fmt"
  "log"
  "time"
  "net/http"
	"crypto/rand"
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
func DoLogin(w http.ResponseWriter, r *http.Request, db *string, col *string) {

  var uin User
  var uout User

  w.Header().Set("Content-Type", "application/json; charset=utf-8")

  // Get Post Data
  b, _ := ioutil.ReadAll(r.Body)
  if err := json.Unmarshal(b, &uin); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Fprintf(w, "{\"msg\":\"Invalid JSON\"}")
    log.Printf("doLogin : %v", err)
    return
  }
  log.Printf("doLogin : Got input %v", uin)


  // Get collection object
  c := session.DB(*db).C(*col)

  // Fetch user record
  if err := c.Find(bson.M{"_id" : uin.Id}).Select(nil).One(&uout); err != nil {
      w.WriteHeader(http.StatusNotFound)
      fmt.Fprintf(w, "{\"msg\":\"Id not found or wrong password\"}")
      log.Printf("doLogin : %v", err)
      return
    }
  log.Printf("doLogin : Got output %v", uout)


  // Compute SHA256
  h := sha256.New()
  h.Write([]byte(uout.Salt + uin.Hash))
  log.Printf("doLogin: hash %s", hex.EncodeToString(h.Sum(nil)))


  // Compare hashes
  if ( uout.Hash != hex.EncodeToString(h.Sum(nil)) ) {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, "{\"msg\":\"Id not found or wrong password\"}")
    log.Printf("doLogin : hashes mismatch")
    return
  }

  log.Printf("doLogin : %v login successfull", uin.Id)

  // Create Claim for token
  claims := &jwt.StandardClaims {
    ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
    Issuer: issuer,
    Id: uin.Id,
  }

  // Generate then sign token
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  tokenString, err := token.SignedString(signature)
  if (err != nil) {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Fprintf(w, "{\"msg\":\"Token error\"}")
    log.Printf("doLogin : %v", err)
    return
  }
  log.Printf("doLogin :Token %s", tokenString)
	fmt.Fprintf(w, "{\"token\":\"%s\", \"username\":\"%s\"}", tokenString, uout.Name)

}


//
// Signup
//
func SetLogin(w http.ResponseWriter, r *http.Request, db *string, col *string) {

  var u User
  var v User

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

	if ( u.Id=="" || u.Name=="" || u.Hash=="" ) {
		w.WriteHeader(http.StatusInternalServerError)
    fmt.Fprintf(w, "{\"msg\":\"empty fields\"}")
    log.Printf("setLogin : empty fields %v", u)
    return
	}

  // Get collection object
  c := session.DB(*db).C(*col)

  // Check if already exists
  if err := c.Find(bson.M{"_id" : u.Id}).Select(bson.M{"_id":1}).One(&v); err == nil {
    w.WriteHeader(http.StatusConflict)
    fmt.Fprintf(w, "{\"msg\":\"User already exists\"}")
    log.Printf("setLogin : %s already exists", u.Id)
    return
  }

  // Generate Salt
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
		log.Printf("setLogin : Error generating salt %v", err)
    fmt.Fprintf(w, "{\"msg\":\"Failed to add user\"}")
		return
	}
  u.Salt = hex.EncodeToString(salt)

  // Generate Salted Password
  h := sha256.New()
  h.Write([]byte(u.Salt + u.Hash))
  u.Hash = hex.EncodeToString(h.Sum(nil))

  // Store user
  if err := c.Insert(&u); err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintf(w, "{\"msg\":\"Failed to add user\"}")
      log.Printf("setLogin : Error adding %v %v", u, err)
      return
    }

  log.Printf("setLogin : Added %v", u)
	fmt.Fprintf(w, "{\"msg\":\"Ok\"}")

}
