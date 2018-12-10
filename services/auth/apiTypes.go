package main

//
// MongoDB Mapping
//
type User struct {
	Id					string				`bson:"_id" json:"id"`
	Name				string				`bson:"name" json:"username"`
	Hash				string				`bson:"hash" json:"password"`
	Salt				string				`bson:"salt"`
}
