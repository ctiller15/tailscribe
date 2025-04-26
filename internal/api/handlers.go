package api

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/mail"
)

// Probably some sort of abstraction here. I'll figure it out eventually.
type IndexPageData struct {
	Title string
}

type SignupPageData struct {
	Title string
	SignupDetails
}

type AttributionsPageData struct {
	Title string
}

type TermsAndConditionsPageData struct {
	Title string
}

type PrivacyPolicyPageData struct {
	Title        string
	ContactEmail string
}

type ContactUsPageData struct {
	Title        string
	ContactEmail string
}

func (a *APIConfig) HandleIndex(w http.ResponseWriter, r *http.Request) {

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

func (a *APIConfig) HandleSignupPage(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles(
		"./templates/signup.html",
		"./templates/base.html",
	))

	data := IndexPageData{
		Title: "TailScribe - Sign Up",
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

type SignupDetails struct {
	Email    string
	Password string
	Valid    bool
}

func (a *APIConfig) HandlePostSignup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Posting!")
	signupDetails := SignupDetails{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	data := SignupPageData{
		Title:         "TailScribe - Sign Up",
		SignupDetails: signupDetails,
	}

	tmpl := template.Must(template.ParseFiles(
		"./templates/signup.html",
		"./templates/base.html",
	))

	fmt.Println(data)

	// Validate email.
	_, err := mail.ParseAddress(signupDetails.Email)
	if err != nil {
		signupDetails.Valid = false
		w.WriteHeader(http.StatusBadRequest)
		err = tmpl.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// Hash password.
	// Store both.
}

func (a *APIConfig) HandleAttributions(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"./templates/attributions.html",
		"./templates/base.html",
	))

	data := AttributionsPageData{
		Title: "Attributions",
	}

	err := tmpl.Execute(w, data)
	// Instead of a log fatal, probably a generic 500 page.
	if err != nil {
		log.Fatal(err)
	}
}

func (a *APIConfig) HandleTerms(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"./templates/terms_and_conditions.html",
		"./templates/base.html",
	))

	data := TermsAndConditionsPageData{
		Title: "Terms and Conditions",
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *APIConfig) HandlePrivacyPolicy(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"./templates/privacy_policy.html",
		"./templates/base.html",
	))

	data := PrivacyPolicyPageData{
		Title:        "Privacy Policy",
		ContactEmail: a.Env.ContactEmail,
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *APIConfig) HandleContactUs(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"./templates/contact_us.html",
		"./templates/base.html",
	))

	data := ContactUsPageData{
		Title:        "Contact Us",
		ContactEmail: a.Env.ContactEmail,
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}
