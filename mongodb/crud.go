package main

import (
	"os"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

//
// Mongodb Model Mapping
//
type RR	struct {
	Name  	string        `bson:"name"`
	Type  	string        `bson:"type"`
	Class  	string        `bson:"class"`
	TTL  		string        `bson:"ttl"`
	Rdata  	string        `bson:"rdata"`
}

type Zone struct {
	Domain  	string        `bson:"domain"`
	RR				[]RR					`bson:"rr"`
}

type Record 	struct {
	Id					string				`bson:"_id" json:"id"`
	Username  	string        `bson:"username"`
	Password  	string        `bson:"password"`
	Sources			[]string			`bson:"sources"`
	Forwarders	[]string			`bson:"forwarders"`
	Zones				[]Zone				`bson:"zones"`
}

//
// Usage
//
func usage() {
	fmt.Printf("Usage: %s <clientID>\n", os.Args[0])
}


//
// Main
//
func main() {

  const (
		URL   = "127.0.0.1:27017"
		DATABASE   = "saas"
		COLLECTION = "data"
	)

	// Sanity check
	if (len(os.Args) != 2) {
		usage()
		panic("Wrong argument")
	}

  // Connect to Database
  session, err := mgo.Dial(URL)
  if err != nil {
		panic(err)
	}
  fmt.Printf("Connected to %s\n",URL)
	// Read secondaries with consistence
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

// Get collection object
c := session.DB(DATABASE).C(COLLECTION)


// Get all data for that user
var rec Record
if err := c.Find(bson.M{"username" : os.Args[1]}).Select(nil).One(&rec); err != nil {
		panic(err)
	}
fmt.Printf("%v\n",rec)


}
