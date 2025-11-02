package api

import (
	"database/sql"
	"encoding/json"
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

	err = a.createAndAttachSessionCookies(&w, user)
	if err != nil {
		log.Fatal(err)
		signupDetails.Valid = false
		w.WriteHeader(http.StatusBadRequest)
		err = tmpl.Execute(w, signupPageData)
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
	loginPageData *LoginPageData,
	status int) error {

	loginDetails.Valid = false
	w.WriteHeader(status)

	err := tmpl.Execute(w, loginPageData)

	return err
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
		err = RejectPostLogin(w, tmpl, &loginDetails, &loginPageData, http.StatusUnauthorized)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	valid := auth.CheckPasswordHash(loginDetails.Password, user.Password.String)

	if !valid {
		err = RejectPostLogin(w, tmpl, &loginDetails, &loginPageData, http.StatusUnauthorized)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err = a.createAndAttachSessionCookies(&w, user)
	if err != nil {
		log.Fatal(err)
		err = RejectPostLogin(w, tmpl, &loginDetails, &loginPageData, http.StatusInternalServerError)
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

func (a *APIConfig) HandleAddNewPet(w http.ResponseWriter, r *http.Request, user_id int) {
	tmpl := template.Must(template.ParseFiles(
		"./templates/new_pet.html",
		"./templates/base.html",
	))

	data := BasePageData{
		Title: "Add New Pet",
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
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
		w.Write(dat)
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
	w.Write(dat)
}
