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

//IndexPageData ...
type IndexPageData struct {
	PageTitle string
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
		//fmt.Fprintln(w, "Welcome to digital library web site!")

		tmpl, err := template.ParseFiles("../templates/index.html")
		if err != nil {
			log.Fatal("Error in parsing index template", err)
		}

		indexData := IndexPageData{
			PageTitle: "Welcome to main page of DigitalLibrary!",
		}

		err = tmpl.Execute(w, indexData)
		if err != nil {
			log.Fatal("Error in execute index template", err)
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

		//fmt.Fprintln(w, "Registration page!")

		tmpl, err := template.ParseFiles("../templates/registration.html")
		if err != nil {
			log.Fatal("Error in parsing register template", err)
		}

		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return // if get it will not proceed to get values from form
		}

		// tmpl := template.Must(template.ParseFiles("../register.html"))

		// data := RegistrationPageData{
		// 	PageTitle: "This is registration page",
		// 	Nickname:  "Cru1zzz3",
		// 	Email:     "mailru",
		// 	IDPassword:  "admin",
		// }

		registrationData := RegistrationPageData{
			Nickname: r.FormValue("Nickname"),
			Email:    r.FormValue("Email"),
			Password: r.FormValue("Password"),
		}

		succ := SuccessfullyRegistration{
			Success:  true,
			Nickname: registrationData.Nickname,
		}

		err = tmpl.Execute(w, succ)
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

		loginPageData := loginPage{
			PageTitle: " ",
		}

		tmpl, err := template.ParseFiles("../templates/login.html")
		if err != nil {
			log.Fatal("Error in parse login template", err)
		}

		// Will continue if submit button will be pressed!
		if r.Method != http.MethodPost {
			err = tmpl.Execute(w, loginPageData)
			if err != nil {
				log.Fatal("Error in execute login template", err)
			}
			return
		}

		loginPageData = loginPage{
			Nickname: r.FormValue("Nickname"),
			Password: r.FormValue("Password"),
		}

		checkTmpl, err := template.ParseFiles("../templates/checklogin.html")
		if err != nil {
			log.Fatal("Error in parse login template", err)
		}

		registered, hashFromDB, err := db.LoginUser(loginPageData.Nickname)
		if err != nil {
			log.Fatal("Error in check user login", err)
		}
		if !registered {
			err = checkTmpl.Execute(w, LoginErrors{true, false})
			if err != nil {
				log.Fatal("Error in user not registered message", err)
			}
		} else {
			// TODO: Compare pass and hash
			if err != nil {
				log.Fatal("Error in main.go checkTmpl.Execute", err)
			}

			err := bcrypt.CompareHashAndPassword([]byte(hashFromDB), []byte(loginPageData.Password))
			if err != nil {
				log.Println("Error in compare hashing", err)
				err = checkTmpl.Execute(w, LoginErrors{false, true})
				return
			}

			err = checkTmpl.Execute(w, LoginErrors{false, false})

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

//HashPassword get hash of the inputted IDPassword
func HashPassword(Password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(Password), 14)
	return string(bytes), err
}

//CheckPasswordHash checks IDPassword with corresponding hash and returns true or false
func CheckPasswordHash(Password, Hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(Hash), []byte(Password))
	return err == nil
}
