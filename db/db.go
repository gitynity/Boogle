package db

import (
	"gopkg.in/mgo.v2"
)

func ConnectUser() *mgo.Collection {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}
	c := session.DB("go-web-dev").C("users")
	return c
}

func ConnectBook() *mgo.Collection {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}
	c := session.DB("go-web-dev").C("books")
	return c
}

// disconnect from the database
func Disconnect(c *mgo.Collection) {
	c.Database.Session.Close()
}
