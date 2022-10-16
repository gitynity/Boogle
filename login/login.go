package login

import (
	"github.com/gitynity/Boogle/db"
	"github.com/gitynity/Boogle/users"
	"gopkg.in/mgo.v2"
)

func Login(username string, password string, c *mgo.Collection) bool {
	if c == nil {
		c = db.ConnectUser()
	}
	userexist := users.CheckUser(username, c)
	if userexist {
		authenticated_user := users.CheckPassword(username, password, c)
		if authenticated_user {
			return true
		}
	}
	return false
}
