package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gitynity/Boogle/db"
	"github.com/gitynity/Boogle/login"
)

type booksearch struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Link        string `json:"link"`
}

type books struct {
	Items []struct {
		VolumeInfo struct {
			Title               string   `json:"title"`
			Authors             []string `json:"authors"`
			Publisher           string   `json:"publisher"`
			PublishedDate       string   `json:"publishedDate"`
			Description         string   `json:"description"`
			IndustryIdentifiers []struct {
				Type       string `json:"type"`
				Identifier string `json:"identifier"`
			} `json:"industryIdentifiers"`
			PageCount        int      `json:"pageCount"`
			PrintType        string   `json:"printType"`
			Categories       []string `json:"categories"`
			AverageRating    float64  `json:"averageRating"`
			RatingsCount     int      `json:"ratingsCount"`
			MaturityRating   string   `json:"maturityRating"`
			AllowAnonLogging bool     `json:"allowAnonLogging"`
			ContentVersion   string   `json:"contentVersion"`
			ImageLinks       struct {
				SmallThumbnail string `json:"smallThumbnail"`
				Thumbnail      string `json:"thumbnail"`
			} `json:"imageLinks"`
			Language            string `json:"language"`
			PreviewLink         string `json:"previewLink"`
			InfoLink            string `json:"infoLink"`
			CanonicalVolumeLink string `json:"canonicalVolumeLink"`
		} `json:"volumeInfo"`
		SelfLink string `json:"selfLink"`
	} `json:"items"`
	Kind       string `json:"kind"`
	TotalItems int    `json:"totalItems"`
}

func searchBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("search")
	//trim the spaces at the beginning and end of the string
	query = strings.TrimSpace(query)
	//replace spaces with +
	query = strings.Replace(query, " ", "+", -1)
	if query == "" {
		http.Error(w, "Please enter a search query", http.StatusBadRequest)
		return
	}
	//get the books
	url := "https://www.googleapis.com/books/v1/volumes?q=" + query
	resp, err := http.Get(url)
	if err != nil {
		w.Write([]byte("Error"))
	}
	defer resp.Body.Close()

	b := books{}
	err = json.NewDecoder(resp.Body).Decode(&b)
	if err != nil {
		w.Write([]byte("Error"))
	}
	//get the first 10 books and send them to the client in json format of struct booksearch
	var book []booksearch
	for i := 0; i < 10; i++ {
		book = append(book, booksearch{b.Items[i].VolumeInfo.Title, b.Items[i].VolumeInfo.Description, b.Items[i].VolumeInfo.ImageLinks.Thumbnail, b.Items[i].VolumeInfo.InfoLink})
	}
	json.NewEncoder(w).Encode(book)
}

func home(w http.ResponseWriter, r *http.Request) {
	// link to index.html
	http.ServeFile(w, r, "index.html")
}

func checklogin(w http.ResponseWriter, r *http.Request) {
	//get the username and password
	// using cookies and not take password in query.
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	c := db.ConnectUser()
	defer db.Disconnect(c)
	//check if the username and password are correct
	if login.Login(username, password, c) {
		//if correct, redirect to the home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		//if not correct, redirect to the login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/books", searchBooks)
	http.HandleFunc("/login", checklogin)
	http.ListenAndServe(":8080", nil)
}
