package main

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)


type User struct {
	Username   string        `bson:"username"`
	Password   string        `bson:"password"`
}


func main() {

  const (
		URL   = "127.0.0.1:27017"
		DATABASE   = "saas"
		COLLECTION = "data"
	)


  // Connect to Database
  session, err := mgo.Dial(URL)
  if err != nil {
		panic(err)
	}
  fmt.Printf("Connected to %s\n",URL)
	defer session.Close()

// Get collection object
c := session.DB(DATABASE).C(COLLECTION)

// Get a user
var u User
if err := c.Find(nil).Select(bson.M{"username": 1, "password":1}).One(&u); err != nil {
		panic(err)
	}
fmt.Println(u.Username)

// Get all data for that user
var all []interface{}
if err := c.Find(bson.M{"username" : u.Username}).Select(nil).All(&all); err != nil {
		panic(err)
	}
fmt.Printf("%v\n",all.forwarders)


}
