package books

import (
	"github.com/gitynity/Boogle/db"
	"gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"
)

type Book struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Title       string        `bson:"title"`
	Author      string        `bson:"author"`
	Description string        `bson:"description"`
	ISBN        string        `bson:"isbn,omitempty"`
	InfoLink    string
}

func GetBooks(c *mgo.Collection) []Book {
	if c == nil {
		c = db.ConnectBook()
	}
	var books []Book
	err := c.Find(nil).All(&books)
	if err != nil {
		panic(err)
	}
	return books
}

func GetBook(id string, c *mgo.Collection) Book {
	if c == nil {
		c = db.ConnectBook()
	}
	var book Book
	err := c.FindId(bson.ObjectIdHex(id)).One(&book)
	if err != nil {
		panic(err)
	}
	return book
}

func InsertBook(title string, author string, description string, isbn string, infoLink string, c *mgo.Collection) {
	if c == nil {
		c = db.ConnectBook()
	}
	err := c.Insert(&Book{Title: title, Author: author, Description: description, ISBN: isbn, InfoLink: infoLink})
	if err != nil {
		panic(err)
	}
}

func DeleteBook(id string, c *mgo.Collection) {
	if c == nil {
		c = db.ConnectBook()
	}
	err := c.RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		panic(err)
	}
}

func UpdateBook(id string, title string, author string, description string, isbn string, infoLink string, c *mgo.Collection) {
	if c == nil {
		c = db.ConnectBook()
	}
	err := c.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"title": title, "author": author, "description": description, "isbn": isbn, "infoLink": infoLink}})
	if err != nil {
		panic(err)
	}
}

func AddBookToUser(username string, bookID string, c *mgo.Collection) {
	if c == nil {
		c = db.ConnectUser()
	}
	err := c.Update(bson.M{"username": username}, bson.M{"$push": bson.M{"books": bookID}})
	if err != nil {
		panic(err)
	}
}
