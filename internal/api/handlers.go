package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/mail"
	"time"

	"github.com/ctiller15/tailscribe/internal/auth"
	"github.com/ctiller15/tailscribe/internal/database"
	"github.com/imagekit-developer/imagekit-go/v2" // imported as imagekit
	"github.com/imagekit-developer/imagekit-go/v2/option"
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

type LoginPageData struct {
	Title string
	LoginForm
}

type AddNewPetForm struct {
	Image string
	Name  string
	Valid bool
}

type AddNewPetPageData struct {
	Title string
	AddNewPetForm
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

type AddPetForm struct {
	Image string
	Name  string
}

type NewPetPageData struct {
	Title string
	AddPetForm
}

func (a *APIConfig) HandleIndex(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles(
		"./templates/base.tmpl",
		"./templates/nav.tmpl",
		"./templates/index.tmpl",
	))

	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *APIConfig) HandleSignupPage(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles(
		"./templates/base.tmpl",
		"./templates/nav.tmpl",
		"./templates/signup.tmpl",
	))

	err := tmpl.Execute(w, nil)
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

	tmpl := template.Must(template.ParseFiles(
		"./templates/base.tmpl",
		"./templates/nav.tmpl",
		"./templates/signup.tmpl",
	))

	// Validate email.
	_, err := mail.ParseAddress(signupDetails.Email)
	if err != nil {
		// Abstract this failure state into a function
		signupDetails.Valid = false
		w.WriteHeader(http.StatusBadRequest)
		err = tmpl.Execute(w, signupDetails)
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
		err = tmpl.Execute(w, signupDetails)
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
		err = tmpl.Execute(w, signupDetails)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err = a.createAndAttachSessionCookies(&w, user)
	if err != nil {
		log.Fatal(err)
		signupDetails.Valid = false
		w.WriteHeader(http.StatusBadRequest)
		err = tmpl.Execute(w, signupDetails)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	http.Redirect(w, r, "/add_new_pet", http.StatusFound)
}

func expireCookie(w *http.ResponseWriter, cookie_name string) {
	http.SetCookie(*w, &http.Cookie{
		Name:     cookie_name,
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
}

func (a *APIConfig) createAndAttachSessionCookies(
	w *http.ResponseWriter,
	user database.User,
) error {
	tokenString, err := auth.MakeJWT(user.ID, a.Env.Secret)
	if err != nil {
		return err
	}

	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		return err
	}

	http.SetCookie(*w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		Secure:   true,
		// Domain:   "/",
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(*w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshTokenString,
		Expires:  time.Now().Add(time.Hour * 30 * 24),
		HttpOnly: true,
		// Domain:   "/",
		SameSite: http.SameSiteStrictMode,
	})

	return nil
}

func RejectPostLogin(
	w http.ResponseWriter,
	tmpl *template.Template,
	loginDetails *LoginForm,
	status int) error {

	loginDetails.Valid = false
	w.WriteHeader(status)

	err := tmpl.Execute(w, loginDetails)

	return err
}

func (a *APIConfig) HandlePostLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	loginDetails := LoginForm{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	tmpl := template.Must(template.ParseFiles(
		"./templates/base.tmpl",
		"./templates/nav.tmpl",
		"./templates/login.tmpl",
	))

	email := sql.NullString{
		String: loginDetails.Email,
		Valid:  true,
	}

	user, err := a.Db.GetUserByEmail(ctx, email)
	if err != nil {
		err = RejectPostLogin(w, tmpl, &loginDetails, http.StatusUnauthorized)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	valid := auth.CheckPasswordHash(loginDetails.Password, user.Password.String)

	if !valid {
		err = RejectPostLogin(w, tmpl, &loginDetails, http.StatusUnauthorized)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err = a.createAndAttachSessionCookies(&w, user)
	if err != nil {
		log.Fatal(err)
		err = RejectPostLogin(w, tmpl, &loginDetails, http.StatusInternalServerError)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func (a *APIConfig) HandlePostLogout(w http.ResponseWriter, r *http.Request) {
	expireCookie(&w, "token")
	expireCookie(&w, "refresh_token")

	http.Redirect(w, r, "/", http.StatusFound)
}

func (a *APIConfig) HandleAttributions(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"./templates/base.tmpl",
		"./templates/nav.tmpl",
		"./templates/attributions.tmpl",
	))

	err := tmpl.Execute(w, nil)
	// Instead of a log fatal, probably a generic 500 page.
	if err != nil {
		log.Fatal(err)
	}
}

func (a *APIConfig) HandleTerms(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"./templates/base.tmpl",
		"./templates/nav.tmpl",
		"./templates/terms_and_conditions.tmpl",
	))

	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *APIConfig) HandlePrivacyPolicy(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"./templates/base.tmpl",
		"./templates/nav.tmpl",
		"./templates/privacy_policy.tmpl",
	))

	data := PrivacyPolicyPageData{
		ContactEmail: a.Env.ContactEmail,
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *APIConfig) HandleContactUs(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"./templates/base.tmpl",
		"./templates/nav.tmpl",
		"./templates/contact_us.tmpl",
	))

	data := ContactUsPageData{
		ContactEmail: a.Env.ContactEmail,
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *APIConfig) HandleGetAddNewPet(w http.ResponseWriter, r *http.Request, user_id int) {
	tmpl := template.Must(template.ParseFiles(
		"./templates/base.tmpl",
		"./templates/nav.tmpl",
		"./templates/new_pet.tmpl",
	))

	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *APIConfig) HandlePostAddNewPet(w http.ResponseWriter, r *http.Request, user_id int) {
	ctx := r.Context()
	addNewPetForm := AddNewPetForm{
		Image: r.FormValue("image"),
		Name:  r.FormValue("name"),
	}

	addNewPetPageData := AddNewPetPageData{
		Title:         "TailScribe - Log In",
		AddNewPetForm: addNewPetForm,
	}

	// Attempt to create pet.
	imageUrl := sql.NullString{
		String: addNewPetForm.Image,
	}
	createPetParams := database.CreatePetParams{
		Name:     addNewPetForm.Name,
		Imageurl: imageUrl,
	}

	newPet, err := a.Db.CreatePet(ctx, createPetParams)

	if err != nil {
		// return previous page, etc.
		log.Printf("error creating pet: %v", err)
		tmpl := template.Must(template.ParseFiles(
			"./templates/new_pet.html",
			"./templates/base.html",
		))

		addNewPetPageData.Valid = false
		w.WriteHeader(http.StatusBadRequest)

		err = tmpl.Execute(w, addNewPetPageData)
		if err != nil {
			log.Printf("an error occurred: %v", err)
		}
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/dashboard/pet/%d", newPet.ID), http.StatusCreated)
}

func (a *APIConfig) HandleGetImageAuthParams(w http.ResponseWriter, r *http.Request, user_id int) {
	type responseStruct struct {
		Expire    int64  `json:"expire"`
		Signature string `json:"signature"`
		Token     string `json:"token"`
	}
	client := imagekit.NewClient(
		option.WithPrivateKey(a.Env.ImageKitPrivateKey),
	)

	authParams, err := client.Helper.GetAuthenticationParameters("", 0)

	if err != nil {
		log.Printf("Error getting auth parameters, %s\n", err)

		type errStruct struct {
			Error string `json:"error"`
		}

		newErr := errStruct{
			Error: "error getting auth parameters",
		}

		dat, err := json.Marshal(newErr)

		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		_, err = w.Write(dat)
		if err != nil {
			log.Printf("Error writing data: %s", err)
		}
		return
	}

	response := responseStruct{
		Expire:    authParams["expire"].(int64),
		Signature: authParams["signature"].(string),
		Token:     authParams["token"].(string),
	}

	// Pass params back via json
	// // Result: map[expire:<timestamp> signature:<hmac-signature> token:<uuid-token>]
	// Frontend uses params to build request. See https://imagekit.io/docs/integration/javascript#upload-example-and-error-handling
	// template.JS()

	// Steps for future respondwjson method

	dat, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	_, err = w.Write(dat)
	if err != nil {
		log.Printf("Error writing data: %s", err)
	}
}
