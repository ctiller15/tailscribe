package main

import (
	"html/template"
	"log"
	"net/http"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	tmpl.Execute(w, nil)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handleIndex)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
