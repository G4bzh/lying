package main

import (
	"os"
	"fmt"
	"sort"
	"flag"
	"text/template"

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
// Templates
//
const zonesTmpl = `{{range .RRs}}
{{.Name}} {{.TTL}} {{.Class}} {{.Type}} {{.Rdata}}{{end}}
`

const namedTmpl = `options {
        directory "/var/cache/bind";
        forwarders {
              {{range .Forwarders}}
              {{.}};{{end}}
        };

        query-source address * port 53;
        listen-on { any; };

        auth-nxdomain no;    # conform to RFC1035
        listen-on-v6 { any; };

        zone-statistics yes;

        response-policy { zone "rpz"; };
};

logging {
      channel "querylog" { file "/var/log/dns-query.log"; print-time yes; print-category yes; print-severity yes; };
      category queries { querylog; };
};

{{range .Zones}}
zone "{{.Domain}}" {
    type master;
    file "{{.Domain}}.txt";
};
{{end}}
`

//
// Main
//
func main() {

	idPtr := flag.String("id","","Client ID")
	urlPtr := flag.String("url","127.0.0.1:27017","MongoDB URL")
	dbPtr := flag.String("db","saas","MongoDB Database")
	colPtr := flag.String("col","data","MongoDB Collection")

	flag.Parse()

  // Connect to Database
  session, err := mgo.Dial(*urlPtr)
  if err != nil {
		panic(err)
	}
  fmt.Printf("Connected to %s\n",*urlPtr)
	// Read secondaries with consistence
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	// Get collection object
	c := session.DB(*dbPtr).C(*colPtr)




	// Get all data for given user
	var rec Record
	if err = c.Find(bson.M{"username" : *idPtr}).Select(nil).One(&rec); err != nil {
			panic(err)
		}

	// Render named.conf
	f, err := os.OpenFile("named.conf", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	tmpl := template.Must(template.New("named").Parse(namedTmpl))
	if err = tmpl.Execute(f, rec); err != nil {
		panic(err)
	}
	f.Close()


	// Sort and render zones
	for _, z := range rec.Zones {

		sort.SliceStable(z.RRs,z.RRs.rrSort)

		f, err = os.OpenFile(z.Domain + ".txt" , os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}
		tmpl = template.Must(template.New("zones").Parse(zonesTmpl))
		if err = tmpl.Execute(f, z); err != nil {
			panic(err)
		}
		f.Close()

	}

}
