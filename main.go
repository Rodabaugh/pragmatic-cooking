package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /admin/healthz", readinessEndpoint)

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		MainPage().Render(r.Context(), w)
	})

	server := &http.Server{
		Addr:	":8080",
		Handler: mux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting pragmatic-cooking on port 8080")
	log.Fatal(server.ListenAndServe())
}
