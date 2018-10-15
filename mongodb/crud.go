package main

import (
	"os"
	"fmt"
	"sort"
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
	TTL  		int			      `bson:"ttl"`
	Rdata  	string        `bson:"rdata"`
}

type RRs 	[]RR

type Zone struct {
	Domain  	string        `bson:"domain"`
	RRs				RRs 					`bson:"rr"`
}

type Zones 	[]Zone

type Record 	struct {
	Id					string				`bson:"_id" json:"id"`
	Username  	string        `bson:"username"`
	Password  	string        `bson:"password"`
	Sources			[]string			`bson:"sources"`
	Forwarders	[]string			`bson:"forwarders"`
	Zones				Zones					`bson:"zones"`
}

//
// Usage
//
func usage() {
	fmt.Printf("Usage: %s <username>\n", os.Args[0])
}

//
// Sort RR Slice
//
func (r RRs)rrSort(i,j int) bool {

	typeI := r[i].Type;
	typeJ := r[j].Type;

	if (typeI == "SOA")	{
		return true
	}

	if (typeJ == "SOA")	{
		return false
	}

	return typeI < typeJ

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

	// Get all data for given user
	var rec Record
	if err := c.Find(bson.M{"username" : os.Args[1]}).Select(nil).One(&rec); err != nil {
			panic(err)
		}

	// Sort Zones
	for _, z := range rec.Zones {
		sort.SliceStable(z.RRs,z.RRs.rrSort)
		fmt.Printf("Zone: %s\n",z.Domain)

		for _, r := range z.RRs {
			fmt.Printf("\t%s\t%d\t%s\t%s\t%s\n",
				r.Name,
				r.TTL,
				r.Class,
				r.Type,
				r.Rdata)
		}

	}

}
