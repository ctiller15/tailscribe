package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ctiller15/tailscribe/internal/api"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading .env file: %v.\n Proceeding without env file...", err)
	}

	envVars := api.NewEnvVars()
	apiCfg := api.NewAPIConfig(envVars)

	fs := http.FileServer(http.Dir("assets/"))

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", apiCfg.HandleIndex)
	mux.HandleFunc("/attributions", apiCfg.HandleAttributions)
	mux.HandleFunc("/terms", apiCfg.HandleTerms)
	mux.HandleFunc("/privacy", apiCfg.HandlePrivacyPolicy)
	mux.HandleFunc("/contact", apiCfg.HandleContactUs)

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
