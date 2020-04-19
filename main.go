package main

import (
	"fmt"
	"gorilla/mux"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to digital library web site!", "Your token:", r.URL.Query().Get("token"))
	})

	fs := http.FileServer(http.Dir("static/")) // Directory where will be FileServer

	http.Handle("/static/", http.StripPrefix("/static", fs))

	http.ListenAndServe(":80", nil)
}
