package handlers

import (
	"html/template"
	"log"
	"net/http"
)

// May be an easy interface. Every single page should have a title.
type IndexPageData struct {
	Title string
}

type AttributionsPageData struct {
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

func HandleAttributions(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"./templates/attributions.html",
		"./templates/base.html",
	))

	data := AttributionsPageData{
		Title: "Attributions",
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}
