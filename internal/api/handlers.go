package api

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"net/mail"
	"time"

	"github.com/ctiller15/tailscribe/internal/auth"
	"github.com/ctiller15/tailscribe/internal/database"
)

type BasePageData struct {
	Title string
}

// Probably some sort of abstraction here. I'll figure it out eventually.
type IndexPageData struct {
	BasePageData
}

type SignupForm struct {
	Email    string
	Password string
	Valid    bool
}

type SignupPageData struct {
	Title string
	SignupForm
}

type LoginForm struct {
	Email    string
	Password string
	Valid    bool
}

type HtmlPageData interface {
	SignupPageData | LoginPageData
}

type LoginPageData struct {
	Title string
	LoginForm
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

type NewPetPageData struct {
	Title string
}

func (a *APIConfig) HandleIndex(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles(
		"./templates/index.html",
		"./templates/base.html",
	))

	data := BasePageData{
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

	data := BasePageData{
		Title: "TailScribe - Sign Up",
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *APIConfig) HandlePostSignup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	signupDetails := SignupForm{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	signupPageData := SignupPageData{
		Title:      "TailScribe - Sign Up",
		SignupForm: signupDetails,
	}

	tmpl := template.Must(template.ParseFiles(
		"./templates/signup.html",
		"./templates/base.html",
	))

	// Validate email.
	_, err := mail.ParseAddress(signupDetails.Email)
	if err != nil {
		// Abstract this failure state into a function
		signupDetails.Valid = false
		w.WriteHeader(http.StatusBadRequest)
		err = tmpl.Execute(w, signupPageData)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// Hash password.
	hashedPassword, err := auth.HashPassword(signupDetails.Password)
	if err != nil {
		signupDetails.Valid = false
		w.WriteHeader(http.StatusBadRequest)
		err = tmpl.Execute(w, signupPageData)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// Store both.
	createUserParams := database.CreateUserParams{
		Email: sql.NullString{
			String: signupDetails.Email,
			Valid:  true,
		},
		Password: sql.NullString{
			String: hashedPassword,
			Valid:  true,
		},
	}

	// First is user.
	user, err := a.Db.CreateUser(ctx, createUserParams)
	if err != nil {
		signupDetails.Valid = false
		w.WriteHeader(http.StatusBadRequest)
		err = tmpl.Execute(w, signupPageData)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	tokenString, err := auth.MakeJWT(user.ID, a.Env.Secret)
	if err != nil {
		signupDetails.Valid = false
		w.WriteHeader(http.StatusInternalServerError)
		err = tmpl.Execute(w, signupPageData)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		signupDetails.Valid = false
		w.WriteHeader(http.StatusInternalServerError)
		err = tmpl.Execute(w, signupPageData)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	newPetPageData := BasePageData{
		Title: "Add a new Pet",
	}

	// Create new template that points to new pet page.
	tmpl = template.Must(template.ParseFiles(
		"./templates/new_pet.html",
		"./templates/base.html",
	))

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		Secure:   true,
		Domain:   "/",
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshTokenString,
		Expires:  time.Now().Add(time.Hour * 30 * 24),
		HttpOnly: true,
		Domain:   "/",
		SameSite: http.SameSiteStrictMode,
	})
	w.WriteHeader(http.StatusCreated)
	err = tmpl.Execute(w, newPetPageData)
	if err != nil {
		log.Fatal(err)
	}
}

func expireCookie(w *http.ResponseWriter, cookie_name string) {
	http.SetCookie(*w, &http.Cookie{
		Name:     cookie_name,
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
}

func (a *APIConfig) HandlePostLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	loginDetails := LoginForm{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	loginPageData := LoginPageData{
		Title:     "TailScribe - Log In",
		LoginForm: loginDetails,
	}

	tmpl := template.Must(template.ParseFiles(
		"./templates/login.html",
		"./templates/base.html",
	))

	email := sql.NullString{
		String: loginDetails.Email,
		Valid:  true,
	}

	user, err := a.Db.GetUserByEmail(ctx, email)
	if err != nil {
		loginDetails.Valid = false
		w.WriteHeader(http.StatusUnauthorized)
		err = tmpl.Execute(w, loginPageData)
		if err != nil {
			log.Fatal(err)
		}
	}

	valid := auth.CheckPasswordHash(loginDetails.Password, user.Password.String)

	if !valid {
		loginDetails.Valid = false
		w.WriteHeader(http.StatusUnauthorized)
		err = tmpl.Execute(w, loginPageData)
		if err != nil {
			log.Fatal(err)
		}
	}

	tokenString, err := auth.MakeJWT(user.ID, a.Env.Secret)
	if err != nil {
		loginDetails.Valid = false
		w.WriteHeader(http.StatusInternalServerError)
		err = tmpl.Execute(w, loginPageData)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		loginDetails.Valid = false
		w.WriteHeader(http.StatusInternalServerError)
		err = tmpl.Execute(w, loginPageData)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		Secure:   true,
		Domain:   "/",
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshTokenString,
		Expires:  time.Now().Add(time.Hour * 30 * 24),
		HttpOnly: true,
		Domain:   "/",
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
}

func (a *APIConfig) HandlePostLogout(w http.ResponseWriter, r *http.Request) {
	expireCookie(&w, "token")
	expireCookie(&w, "refresh_token")

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
