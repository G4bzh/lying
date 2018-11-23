package main

import (
  "net/http"
)


//
//  Mongodb Mapping
//

type RR	struct {
	Name  	string        `bson:"name" json:"name"`
	Type  	string        `bson:"type" json:"type"`
	Class  	string        `bson:"class" json:"class"`
	TTL  		int			      `bson:"ttl" json:"ttl"`
	Rdata  	string        `bson:"rdata" json:"rdata"`
}

type RRs 	[]RR

type Zone struct {
	Domain  	string        `bson:"domain" json:"domain"`
	RRs				RRs 					`bson:"rr" json:"rrs,omitempty"`
}

type Zones 	[]Zone

type Record 	struct {
	Id					string				`bson:"_id" json:"id"`
	Forwarders	[]string			`bson:"forwarders" json:"forwarders`
	Zones				Zones					`bson:"zones" json:"zones,omitempty"`
}



//
// Types
//

type HandlerF func(w http.ResponseWriter, r *http.Request)
