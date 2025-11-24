package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"

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

	envVars := NewEnvVars()
	dbUrl := envVars.Database.ConnectionString()

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	dbQueries := database.New(db)

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	apiCfg := NewAPIConfig(envVars, dbQueries, logger, templateCache)

	// Initialize routing - break into own func first
	fs := http.FileServer(http.Dir("./ui/static/"))

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	middleware := apiCfg.logRequest

	mux.HandleFunc("GET /{$}", middleware(apiCfg.HandleIndex))
	mux.HandleFunc("GET /signup", middleware(apiCfg.HandleSignupPage))
	mux.HandleFunc("POST /signup", middleware(apiCfg.HandlePostSignup))
	mux.HandleFunc("POST /login", middleware(apiCfg.HandlePostLogin))
	mux.HandleFunc("POST /logout", middleware(apiCfg.HandlePostLogout))
	mux.HandleFunc("/attributions", middleware(apiCfg.HandleAttributions))
	mux.HandleFunc("/terms", middleware(apiCfg.HandleTerms))
	mux.HandleFunc("/privacy", middleware(apiCfg.HandlePrivacyPolicy))
	mux.HandleFunc("/contact", middleware(apiCfg.HandleContactUs))

	authMiddleware := func(next authorizedHandler) http.HandlerFunc { return middleware(apiCfg.CheckAuthMiddleware(next)) }

	mux.Handle("GET /dashboard/add_new_pet", authMiddleware(apiCfg.HandleGetAddNewPet))

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
