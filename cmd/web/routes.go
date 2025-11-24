package main

import (
	"net/http"

	"github.com/ctiller15/tailscribe/ui"
)

func (app *APIConfig) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	middleware := app.logRequest

	mux.HandleFunc("GET /{$}", middleware(app.HandleIndex))
	mux.HandleFunc("GET /signup", middleware(app.HandleSignupPage))
	mux.HandleFunc("POST /signup", middleware(app.HandlePostSignup))
	mux.HandleFunc("POST /login", middleware(app.HandlePostLogin))
	mux.HandleFunc("POST /logout", middleware(app.HandlePostLogout))
	mux.HandleFunc("/attributions", middleware(app.HandleAttributions))
	mux.HandleFunc("/terms", middleware(app.HandleTerms))
	mux.HandleFunc("/privacy", middleware(app.HandlePrivacyPolicy))
	mux.HandleFunc("/contact", middleware(app.HandleContactUs))

	authMiddleware := func(next authorizedHandler) http.HandlerFunc { return middleware(app.CheckAuthMiddleware(next)) }

	mux.Handle("GET /dashboard/add_new_pet", authMiddleware(app.HandleGetAddNewPet))

	return mux
}
