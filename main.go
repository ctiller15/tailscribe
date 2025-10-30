package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/ctiller15/tailscribe/internal/api"
	"github.com/ctiller15/tailscribe/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	// Initialize environment
	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading .env file: %v.\n Proceeding without env file...", err)
	}

	envVars := api.NewEnvVars()
	dbUrl := envVars.Database.ConnectionString()

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	apiCfg := api.NewAPIConfig(envVars, dbQueries)

	// Initialize routing - break into own func first
	fs := http.FileServer(http.Dir("assets/"))

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", apiCfg.HandleIndex)
	mux.HandleFunc("GET /signup", apiCfg.HandleSignupPage)
	mux.HandleFunc("POST /signup", apiCfg.HandlePostSignup)
	mux.HandleFunc("POST /login", apiCfg.HandlePostLogin)
	mux.HandleFunc("POST /logout", apiCfg.HandlePostLogout)
	mux.HandleFunc("/attributions", apiCfg.HandleAttributions)
	mux.HandleFunc("/terms", apiCfg.HandleTerms)
	mux.HandleFunc("/privacy", apiCfg.HandlePrivacyPolicy)
	mux.HandleFunc("/contact", apiCfg.HandleContactUs)

	// Start server
	server := http.Server{
		Handler:           mux,
		Addr:              envVars.Addr,
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Println("listening on", envVars.Addr)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
