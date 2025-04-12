package main

import (
	"html/template"
	"log"
	"net/http"
)

type IndexPageData struct {
	Title string
}

func handleIndex(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles(
		"./templates/index.html",
		"./templates/base.html",
	))

	data := IndexPageData{
		Title: "TailScribe",
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	addr := ":8080"

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handleIndex)

	server := http.Server{
		Handler: mux,
		Addr:    addr,
	}

	log.Println("listening on", addr)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
