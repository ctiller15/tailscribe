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

	app := NewAPIConfig(envVars, dbQueries, logger, templateCache)

	mux := app.routes()

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
