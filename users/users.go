package users

import (
	"github.com/gitynity/Boogle/db"
	"gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"
)

type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Username string        `bson:"username"`
	Password string        `bson:"password"`
	Books    []string      `bson:"books,omitempty"`
}

func CheckUser(username string, c *mgo.Collection) bool {
	if c == nil {
		c = db.ConnectUser()
	}
	count, err := c.Find(bson.M{"username": username}).Count()
	if err != nil {
		panic(err)
	}
	if count > 0 {
		return true
	}
	return false
}

func CheckPassword(username string, password string, c *mgo.Collection) bool {
	if c == nil {
		c = db.ConnectUser()
	}
	count, err := c.Find(bson.M{"username": username, "password": password}).Count()
	if err != nil {
		panic(err)
	}
	if count > 0 {
		return true
	}
	return false
}

func InsertUser(username string, password string, c *mgo.Collection) {
	if c == nil {
		c = db.ConnectUser()
	}
	err := c.Insert(&User{Username: username, Password: password})
	if err != nil {
		panic(err)
	}
}

func GetUser(username string, c *mgo.Collection) User {
	if c == nil {
		c = db.ConnectUser()
	}
	var user User
	err := c.Find(bson.M{"username": username}).One(&user)
	if err != nil {
		panic(err)
	}
	return user
}
