package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	db "github.com/Cru1zzz3/DigitalLibrary/database"
	mail "github.com/Cru1zzz3/DigitalLibrary/mailing"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

// var TemplatePath string

// TemplatePath = "../template/"

// func PrepareParseFiles(files []string) []string {
// 	for file := range files {
// 		file = append("TemplatePath", file)
// 		file = append(file, ".html")
// 	}

// 	return files
// }
var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

var TopGenresSlice []string

var PageData Data

type Data struct {
	Authenticated      bool
	Nickname           string
	LoginError         LoginErrors
	SearchResults      db.SearchStruct
	RegistrationResult SuccessfullyRegistration
	TopGenres          []string
	Other              interface{}
}

//SuccessfullyRegistration ...
type SuccessfullyRegistration struct {
	Success  bool
	Nickname string
}

type LoginErrors struct {
	NotRegistered bool
	AuthError     bool
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/readers", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to digital library web site!")
		err := db.SelectReader(w)
		if err != nil {
			log.Fatal("Can's select Reader from DB")
		}

	})

	Conn, err := db.ConnectToDB()
	if err != nil {
		log.Fatal("Error in connect to DB: " + err.Error())
	}
	defer Conn.Close()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		//var templates = template.Must(template.ParseGlob("../templates/*"))

		files := []string{
			"../templates/index.html",
			"../templates/navbar.html",
		}

		t, err := template.ParseFiles(files...)
		if err != nil {
			log.Fatal(err)
		}

		session, _ := store.Get(r, "cookie-name")

		PageData := Data{
			Authenticated: CheckIfLogin(r),
			Nickname:      GetNickname(r),
			TopGenres:     db.GetGenres(w),
		}

		// Check if user is authenticated
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			// If not auth

			err = t.ExecuteTemplate(w, "index", PageData)
			if err != nil {
				log.Fatal("Error in execute index template ", err)
			}

			return
		}

		// If auth successfull

		err = t.ExecuteTemplate(w, "index", PageData)
		if err != nil {
			log.Fatal("Error in execute index template ", err)
		}

	})

	r.HandleFunc("/partnership", func(w http.ResponseWriter, r *http.Request) {

		type Partnership struct {
			Sended bool
		}

		files := []string{
			"../templates/index.html",
			"../templates/navbar.html",
			"../templates/partnership.html",
		}

		t, err := template.ParseFiles(files...)
		if err != nil {
			log.Fatal(err)
		}

		PageData := Data{
			Authenticated: CheckIfLogin(r),
			Nickname:      GetNickname(r),
			TopGenres:     db.GetGenres(w),
		}

		if r.Method != http.MethodPost {
			err = t.ExecuteTemplate(w, "index", PageData)
			if err != nil {
				log.Fatal("Error in execute partnership template", err)
			}
			return // if get it will not proceed to get values from form
		}

		recevier := r.FormValue("Email")

		err = mail.SendMail(recevier)
		if err != nil {
			log.Fatal("Error in sending mail during partnership page")
		}

		err = t.ExecuteTemplate(w, "index", nil)
		if err != nil {
			log.Fatal("Error in execute partnership send mail template", err)
		}

	})

	r.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {

		files := []string{
			"../templates/index.html",
			"../templates/navbar.html",
			"../templates/search.html",
		}

		t, err := template.ParseFiles(files...)
		if err != nil {
			log.Fatal(err)
		}

		search := r.FormValue("search")

		//Get result collection from search
		searchResults, err := db.SelectBook(w, search)
		if err != nil {
			log.Fatal("Error in search book query", err)
		}

		// Print result collection
		//log.Println(searchResults.NameBook, searchResults.NameAuthor)

		PageData := Data{
			Authenticated: CheckIfLogin(r),
			Nickname:      GetNickname(r),
			TopGenres:     db.GetGenres(w),
			SearchResults: searchResults,
		}

		err = t.ExecuteTemplate(w, "index", PageData)
		if err != nil {
			log.Fatal("Error search template execute", err)
		}

	})

	r.HandleFunc("/genres", func(w http.ResponseWriter, r *http.Request) {

		files := []string{
			"../templates/index.html",
			"../templates/navbar.html",
			"../templates/allgenres.html",
		}

		t, err := template.ParseFiles(files...)
		if err != nil {
			log.Fatal(err)
		}

		searchResults := db.GetAllGenres(w)
		if err != nil {
			log.Fatal("Error in search all genres query", err)
		}

		err = t.ExecuteTemplate(w, "index", Data{
			Authenticated: CheckIfLogin(r),
			Nickname:      GetNickname(r),
			TopGenres:     db.GetGenres(w),
			SearchResults: searchResults},
		)
		if err != nil {
			log.Fatal("Error search template execute", err)
		}

		// fmt.Fprintf(w, "You've requested the genre: %s\n", namegenre)
		// fmt.Fprintf(w, "Your result: %s\n", searchResults)
	})

	r.HandleFunc("/genres/{namegenre}", func(w http.ResponseWriter, r *http.Request) {

		files := []string{
			"../templates/index.html",
			"../templates/navbar.html",
			"../templates/genres.html",
		}

		t, err := template.ParseFiles(files...)
		if err != nil {
			log.Fatal(err)
		}

		vars := mux.Vars(r)
		namegenre := vars["namegenre"]

		searchResults, err := db.AboutGenre(w, namegenre)
		if err != nil {
			log.Fatal("Error in search info about book query", err)
		}

		err = t.ExecuteTemplate(w, "index", Data{
			Authenticated: CheckIfLogin(r),
			Nickname:      GetNickname(r),
			TopGenres:     db.GetGenres(w),
			SearchResults: searchResults},
		)
		if err != nil {
			log.Fatal("Error search template execute", err)
		}

		// fmt.Fprintf(w, "You've requested the genre: %s\n", namegenre)
		// fmt.Fprintf(w, "Your result: %s\n", searchResults)
	})

	r.HandleFunc("/authors/{nameauthor}", func(w http.ResponseWriter, r *http.Request) {

		files := []string{
			"../templates/index.html",
			"../templates/navbar.html",
			"../templates/authors.html",
		}

		t, err := template.ParseFiles(files...)
		if err != nil {
			log.Fatal(err)
		}

		vars := mux.Vars(r)
		nameauthor := vars["nameauthor"]

		searchResults, err := db.AboutAuthor(w, nameauthor)
		if err != nil {
			log.Fatal("Error in search info about author query", err)
		}

		err = t.ExecuteTemplate(w, "index", Data{
			Authenticated: CheckIfLogin(r),
			Nickname:      GetNickname(r),
			TopGenres:     db.GetGenres(w),
			SearchResults: searchResults})
		if err != nil {
			log.Fatal("Error search template execute", err)
		}

		// fmt.Fprintf(w, "You've requested the author: %s\n", nameauthor)
		// fmt.Fprintf(w, "Your result: %s\n", searchResults)
	})

	r.HandleFunc("/books/{namebook}", func(w http.ResponseWriter, r *http.Request) {

		files := []string{
			"../templates/index.html",
			"../templates/navbar.html",
			"../templates/books.html",
		}

		t, err := template.ParseFiles(files...)
		if err != nil {
			log.Fatal(err)
		}

		vars := mux.Vars(r)
		namebook := vars["namebook"]

		searchResults, err := db.AboutBook(w, namebook)
		if err != nil {
			log.Fatal("Error in search info about book query", err)
		}

		err = t.ExecuteTemplate(w, "index", Data{
			Authenticated: CheckIfLogin(r),
			Nickname:      GetNickname(r),
			TopGenres:     db.GetGenres(w),
			SearchResults: searchResults})
		if err != nil {
			log.Fatal("Error search template execute", err)
		}

		// fmt.Fprintf(w, "You've requested the book: %s\n", namebook)
		// fmt.Fprintf(w, "Your result: %s\n", searchResults)
	})

	r.HandleFunc("/registration", func(w http.ResponseWriter, r *http.Request) {

		//RegistrationPageData ...
		type RegistrationPageData struct {
			Nickname string
			Email    string
			Password string
			Hash     string
		}

		files := []string{
			"../templates/index.html",
			"../templates/navbar.html",
			"../templates/registration.html",
		}

		t, err := template.ParseFiles(files...)
		if err != nil {
			log.Fatal(err)
		}

		PageData := Data{
			Authenticated: CheckIfLogin(r),
			Nickname:      GetNickname(r),
			TopGenres:     db.GetGenres(w),
		}

		if r.Method != http.MethodPost {
			err = t.ExecuteTemplate(w, "index", PageData)
			if err != nil {
				log.Fatal("Error in execute registration template", err)
			}

			return // if get it will not proceed to get values from form
		}

		registrationData := RegistrationPageData{
			Nickname: r.FormValue("Nickname"),
			Email:    r.FormValue("Email"),
			Password: r.FormValue("Password"),
		}

		PageData = Data{
			Authenticated:      CheckIfLogin(r),
			Nickname:           GetNickname(r),
			RegistrationResult: SuccessfullyRegistration{Success: true},
			TopGenres:          db.GetGenres(w),
		}

		err = t.ExecuteTemplate(w, "index", PageData)
		if err != nil {
			log.Fatal("Error in execute registration template", err)
		}

		hash, _ := HashPassword(registrationData.Password)
		registrationData.Hash = hash

		fmt.Println("Password:", registrationData.Password)
		fmt.Println("Hash:    ", registrationData.Hash)

		match := CheckPasswordHash(registrationData.Password, registrationData.Hash)
		fmt.Println("Match:   ", match)

		if match {
			err = db.RegisterUser(registrationData.Nickname, registrationData.Email, registrationData.Hash)
			if err != nil {
				log.Fatal("Error in register new user: " + err.Error())
			}
		}

	})

	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {

		files := []string{
			"../templates/index.html",
			"../templates/navbar.html",
		}

		t, err := template.ParseFiles(files...)
		if err != nil {
			log.Fatal(err)
		}

		session, _ := store.Get(r, "cookie-name")

		// Revoke users authentication
		session.Values["authenticated"] = false
		session.Save(r, w)

		err = t.ExecuteTemplate(w, "index", Data{Authenticated: CheckIfLogin(r)})
		if err != nil {
			log.Fatal("Error search template execute", err)
		}

	})

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		type loginPage struct {
			PageTitle string
			Nickname  string
			Password  string
			Hash      string
		}

		files := []string{
			"../templates/index.html",
			"../templates/navbar.html",
			"../templates/login.html",
		}

		t, err := template.ParseFiles(files...)
		if err != nil {
			log.Fatal(err)
		}

		loginPageData := loginPage{}

		PageData := Data{
			Authenticated: CheckIfLogin(r),
			Nickname:      GetNickname(r),
			TopGenres:     db.GetGenres(w),
		}

		// Will continue if submit button will be pressed!
		if r.Method != http.MethodPost {

			err = t.ExecuteTemplate(w, "index", PageData)
			if err != nil {
				log.Fatal("Error search template execute", err)
			}

			return
		}

		loginPageData = loginPage{
			Nickname: r.FormValue("Nickname"),
			Password: r.FormValue("Password"),
		}

		// Get cookie for current session
		session, _ := store.Get(r, "cookie-name")

		registered, hashFromDB, err := db.LoginUser(loginPageData.Nickname)
		if err != nil {
			log.Fatal("Error in check user login", err)
		}

		files = []string{
			"../templates/index.html",
			"../templates/navbar.html",
			"../templates/checklogin.html",
		}

		t, err = template.ParseFiles(files...)
		if err != nil {
			log.Fatal(err)
		}

		if !registered {

			PageData = Data{
				Authenticated: CheckIfLogin(r),
				Nickname:      GetNickname(r),
				TopGenres:     db.GetGenres(w),
				LoginError:    LoginErrors{true, false},
			}

			err = t.ExecuteTemplate(w, "index", PageData)
			if err != nil {
				log.Fatal("Error in user not registered message", err)
			}
		} else {

			err = bcrypt.CompareHashAndPassword([]byte(hashFromDB), []byte(loginPageData.Password))
			if err != nil {
				log.Println("Error in compare hashing", err)

				PageData = Data{
					Authenticated: CheckIfLogin(r),
					Nickname:      GetNickname(r),
					TopGenres:     TopGenresSlice,
					LoginError:    LoginErrors{false, true},
				}

				err = t.ExecuteTemplate(w, "index", PageData)
				if err != nil {
					log.Fatal("Error in user not registered message", err)
				}

				return
			}

			// Save cookie if user successfully logged in
			session.Values["authenticated"] = true
			session.Values["nickname"] = loginPageData.Nickname
			session.Save(r, w)

			PageData = Data{
				Authenticated: CheckIfLogin(r),
				Nickname:      GetNickname(r),
				TopGenres:     db.GetGenres(w),
				LoginError:    LoginErrors{false, false},
			}

			err = t.ExecuteTemplate(w, "index", PageData)
			if err != nil {
				log.Fatal("Error in user non error registered message", err)
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)

		}

		// err = tmpl.Execute(w, loginPageData)
		// if err != nil {
		// 	log.Fatal("Error in execute login template", err)
		// }

	})

	// fs := http.FileServer(http.Dir("static/")) // Directory where will be FileServer
	// http.Handle("/static/", http.StripPrefix("/static", fs))

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../static/"))))

	http.ListenAndServe(":80", r)

}

func GetNickname(r *http.Request) string {
	session, _ := store.Get(r, "cookie-name")

	if session.Values["nickname"] != nil {
		return session.Values["nickname"].(string)
	}
	return ""

}

func CheckIfLogin(r *http.Request) bool {
	session, _ := store.Get(r, "cookie-name")

	// Check users authentication
	if session.Values["authenticated"] == false {
		return false
	} else {
		return true
	}

}

//HashPassword get hash of the inputted Password
func HashPassword(Password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(Password), 14)
	return string(bytes), err
}

//CheckPasswordHash checks Password with corresponding hash and returns true or false
func CheckPasswordHash(Password, Hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(Hash), []byte(Password))
	return err == nil
}
