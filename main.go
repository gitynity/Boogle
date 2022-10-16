package main

import (
	"encoding/json"
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
	w.Write([]byte(bodyBg))
	w.Write([]byte("<h1>Your searched books are:</h1>"))
	//create section for each book
	for _, book := range b.Items {
		w.Write([]byte("<section>"))
		w.Write([]byte("<h2>" + book.VolumeInfo.Title + "</h2>"))
		w.Write([]byte("<h3>By: " + book.VolumeInfo.Authors[0] + "</h3>"))
		w.Write([]byte("<p>" + book.VolumeInfo.Description + "</p>"))
		w.Write([]byte("<img src='" + book.VolumeInfo.ImageLinks.Thumbnail + "'/>"))
		//new line
		w.Write([]byte("<br>"))
		w.Write([]byte("<a href='" + book.VolumeInfo.InfoLink + "'>More Info</a>"))
		w.Write([]byte("</section>"))
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(bodyBg))
	//login form
	w.Write([]byte("<h1>Login</h1>"))
	//input form username and password
	w.Write([]byte("<form action='/checklogin'>"))
	w.Write([]byte("<input type='text' name='username' placeholder='Username'/>"))
	w.Write([]byte("<br>"))
	w.Write([]byte("<input type='password' name='password' placeholder='Password'/>"))
	w.Write([]byte("<br>"))
	w.Write([]byte("<input type='submit' value='Login'/>"))
	w.Write([]byte("</form>"))
	//center the form
	w.Write([]byte("<center>"))
	w.Write([]byte("<h1 style='font-family: 'Google Sans', sans-serif;'>Boogle</h1>"))
	w.Write([]byte("<form action='/books'>"))
	w.Write([]byte("<input type='text' name='search' placeholder='Search for books'/>"))
	w.Write([]byte("<input type='submit' value='Search'/>"))
	w.Write([]byte("</form>"))
	w.Write([]byte("</center>"))
}

func checklogin(w http.ResponseWriter, r *http.Request) {
	//get the username and password
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
