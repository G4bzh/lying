package main

//
// MongoDB Mapping
//
type User struct {
	Id					string				`bson:"_id" json:"id"`
	Hash				string				`bson:"hash" json:"password"`
	Salt				string				`bson:"salt"`
}
