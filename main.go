package main

import (
	"log"
	"net/http"

	"github.com/ctiller15/tailscribe/server/handlers"
)

func main() {
	addr := ":8080"

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handlers.HandleIndex)

	server := http.Server{
		Handler: mux,
		Addr:    addr,
	}

	log.Println("listening on", addr)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
