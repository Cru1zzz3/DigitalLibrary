package main

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/Cru1zzz3/DigitalLibrary/database"

	"github.com/gorilla/mux"
)

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
		fmt.Fprintln(w, "Welcome to digital library web site!")

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
