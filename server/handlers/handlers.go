package handlers

import (
	"html/template"
	"log"
	"net/http"
)

type IndexPageData struct {
	Title string
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {

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
