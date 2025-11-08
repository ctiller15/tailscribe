package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ctiller15/tailscribe/internal/api"
	"github.com/ctiller15/tailscribe/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	// Initialize environment
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	err := godotenv.Load()
	if err != nil {
		logger.Warn("error loading .env file. \n Proceeding without env file...", slog.String("error", err.Error()))
	}

	envVars := api.NewEnvVars()
	dbUrl := envVars.Database.ConnectionString()

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	dbQueries := database.New(db)

	apiCfg := api.NewAPIConfig(envVars, dbQueries, logger)

	// Initialize routing - break into own func first
	fs := http.FileServer(http.Dir("./ui/static/"))

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("GET /{$}", apiCfg.HandleIndex)
	mux.HandleFunc("GET /signup", apiCfg.HandleSignupPage)
	mux.HandleFunc("POST /signup", apiCfg.HandlePostSignup)
	mux.HandleFunc("POST /login", apiCfg.HandlePostLogin)
	mux.HandleFunc("POST /logout", apiCfg.HandlePostLogout)
	mux.HandleFunc("/attributions", apiCfg.HandleAttributions)
	mux.HandleFunc("/terms", apiCfg.HandleTerms)
	mux.HandleFunc("/privacy", apiCfg.HandlePrivacyPolicy)
	mux.HandleFunc("/contact", apiCfg.HandleContactUs)

	mux.Handle("GET /dashboard/add_new_pet", apiCfg.CheckAuthMiddleware(apiCfg.HandleGetAddNewPet))

	// Start server
	server := http.Server{
		Handler:           mux,
		Addr:              envVars.Addr,
		ReadHeaderTimeout: 2 * time.Second,
	}

	logger.Info("listening on", slog.String("addr", envVars.Addr))

	err = server.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
