package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gitynity/Boogle/db"
	"github.com/gitynity/Boogle/login"
)

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

var bodyBg = `<style>body {
	background-image: url("https://previews.123rf.com/images/moongdo/moongdo1809/moongdo180901759/108645686-color-pencil-with-leaves-on-yellow-background-business-concept-copyspace.jpg");
	background-repeat: no-repeat;
	background-size: cover;
}</style>`

func searchBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("search")
	query = strings.TrimSpace(query)
	query = strings.Replace(query, " ", "+", -1)
	if query == "" {
		http.Error(w, "Please enter a search query", http.StatusBadRequest)
		return
	}
	// get the books
	url := "https://www.googleapis.com/books/v1/volumes?q=" + query
	resp, err := http.Get(url)
	if err != nil {
		_, err := w.Write([]byte("Error"))
		if err != nil {
			return
		}
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	b := books{}
	err = json.NewDecoder(resp.Body).Decode(&b)
	if err != nil {
		_, err := w.Write([]byte("Error"))
		if err != nil {
			return
		}
	}
	_, err = w.Write([]byte(bodyBg + "<h1>Your searched books are:</h1>"))
	if err != nil {
		return
	}
	// create section for each book
	for _, book := range b.Items {
		_, err := w.Write([]byte("<section>" + "<h2>" + book.VolumeInfo.Title + "</h2>" + "</h3>" + "<p>" + book.VolumeInfo.Description + "</p>" + "<img src='" + book.VolumeInfo.ImageLinks.Thumbnail + "'/>" + "<br>" + "<a href='" + book.VolumeInfo.InfoLink + "'>More Info</a>" + "</section>"))
		if err != nil {
			return
		}
	}
}

func home(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte(bodyBg + "<h1>Login</h1>" + "<form action='/checklogin'>" + "<input type='text' name='username' placeholder='Username'/>" + "<br>" + "<input type='password' name='password' placeholder='Password'/>" + "<br>" + "<input type='submit' value='Login'/>" + "</form>" + "<center>" + "<h1 style='font-family: 'Google Sans', sans-serif;'>Boogle</h1>" + "<form action='/books'>" + "<input type='text' name='search' placeholder='Search for books'/>" + "<input type='submit' value='Search'/>" + "</form>" + "</center>"))
	if err != nil {
		return
	}
}

func checklogin(w http.ResponseWriter, r *http.Request) {
	// get the username and password
	// using cookies and not take password in query.
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	c := db.ConnectUser()
	defer db.Disconnect(c)
	// check if the username and password are correct
	if login.Login(username, password, c) {
		// if correct, redirect to the home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		// if not correct, redirect to the login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/books", searchBooks)
	http.HandleFunc("/login", checklogin)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
