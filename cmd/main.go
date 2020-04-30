package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	db "github.com/Cru1zzz3/DigitalLibrary/database"
	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"
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

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/readers", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to digital library web site!")
		err := db.SelectReader(w)
		if err != nil {
			log.Fatal("Can's select Reader from DB")
		}

	})

	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
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

		err = t.ExecuteTemplate(w, "index", nil)
		if err != nil {
			log.Fatal("Error in execute index template ", err)
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

		err = t.ExecuteTemplate(w, "index", searchResults)
		if err != nil {
			log.Fatal("Error search template execute", err)
		}

	})

	r.HandleFunc("/registration", func(w http.ResponseWriter, r *http.Request) {

		//RegistrationPageData ...
		type RegistrationPageData struct {
			Nickname string
			Email    string
			Password string
			Hash     string
		}

		//SuccessfullyRegistration ...
		type SuccessfullyRegistration struct {
			Success  bool
			Nickname string
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

		if r.Method != http.MethodPost {
			err = t.ExecuteTemplate(w, "index", SuccessfullyRegistration{})
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

		succ := SuccessfullyRegistration{
			Success:  true,
			Nickname: registrationData.Nickname,
		}

		err = t.ExecuteTemplate(w, "index", succ)
		if err != nil {
			log.Fatal("Error in success template execute", err)
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

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		type loginPage struct {
			PageTitle string
			Nickname  string
			Password  string
			Hash      string
		}

		type LoginErrors struct {
			NotRegistered bool
			AuthError     bool
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

		// Will continue if submit button will be pressed!
		if r.Method != http.MethodPost {

			err = t.ExecuteTemplate(w, "index", loginPageData)
			if err != nil {
				log.Fatal("Error in execute login template", err)
			}

			return
		}

		loginPageData = loginPage{
			Nickname: r.FormValue("Nickname"),
			Password: r.FormValue("Password"),
		}

		registered, hashFromDB, err := db.LoginUser(loginPageData.Nickname)
		if err != nil {
			log.Fatal("Error in check user login", err)
		}
		if !registered {
			err = t.ExecuteTemplate(w, "index", LoginErrors{true, false})
			if err != nil {
				log.Fatal("Error in user not registered message", err)
			}
		} else {
			if err != nil {
				log.Fatal("Error in main.go checkTmpl.Execute", err)
			}

			err := bcrypt.CompareHashAndPassword([]byte(hashFromDB), []byte(loginPageData.Password))
			if err != nil {
				log.Println("Error in compare hashing", err)

				err = t.ExecuteTemplate(w, "index", LoginErrors{false, true})
				if err != nil {
					log.Fatal("Error in user not registered message", err)
				}

				return
			}

			err = t.ExecuteTemplate(w, "index", LoginErrors{false, false})
			if err != nil {
				log.Fatal("Error in user non error registered message", err)
			}

		}

		// err = tmpl.Execute(w, loginPageData)
		// if err != nil {
		// 	log.Fatal("Error in execute login template", err)
		// }

	})

	// fs := http.FileServer(http.Dir("static/")) // Directory where will be FileServer
	// http.Handle("/static/", http.StripPrefix("/static", fs))

	http.ListenAndServe(":80", r)

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
