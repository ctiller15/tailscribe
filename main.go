package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ctiller15/tailscribe/server/handlers"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading .env file: %v.\n Proceeding without env file...", err)
	}

	addr := os.Getenv("PORT")

	fs := http.FileServer(http.Dir("assets/"))

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", handlers.HandleIndex)
	mux.HandleFunc("/attributions", handlers.HandleAttributions)

	server := http.Server{
		Handler:           mux,
		Addr:              addr,
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Println("listening on", addr)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
