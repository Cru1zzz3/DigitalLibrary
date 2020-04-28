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

type RegistrationPageData struct {
	Username string
	Email    string
	Password string
	Hash     string
}

type SuccessfullyRegistration struct {
	Success bool
}

type IndexPageData struct {
	PageTitle string
}

func main() {

	Conn, err := db.ConnectToDB()
	if err != nil {
		log.Fatal("Error in connect to DB: " + err.Error())
	}
	defer Conn.Close()

	err = db.InsertReader("cru1zzz3", "Stanislav", "Udartsev", 22)
	if err != nil {
		log.Fatal("Error in insert new Reader: " + err.Error())
	}

	r := mux.NewRouter()

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
		//fmt.Fprintln(w, "Registration page!")

		tmpl, err := template.ParseFiles("../templates/register.html")
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
		// 	Username:  "Cru1zzz3",
		// 	Email:     "mailru",
		// 	Password:  "admin",
		// }

		registrationData := RegistrationPageData{
			Username: r.FormValue("username"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		succ := SuccessfullyRegistration{
			Success: true,
		}

		err = tmpl.Execute(w, succ)
		if err != nil {
			log.Fatal("Error in success template execute", err)
		}

		// TODO: add new user to DB
		hash, _ := HashPassword(registrationData.Password)
		registrationData.Hash = hash

		fmt.Println("Password:", registrationData.Password)
		fmt.Println("Hash:    ", registrationData.Hash)

		match := CheckPasswordHash(registrationData.Password, registrationData.Hash)
		fmt.Println("Match:   ", match)

	})

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Login page!")

	})

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

	fs := http.FileServer(http.Dir("static/")) // Directory where will be FileServer

	http.Handle("/static/", http.StripPrefix("/static", fs))

	http.ListenAndServe(":80", r)

}

//HashPassword get hash of the inputted password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//CheckPasswordHash checks password with corresponding hash and returns true or false
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
