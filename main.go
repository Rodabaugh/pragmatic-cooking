package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Rodabaugh/pragmatic-cooking/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	platform string
	db       *database.Queries
	mgAPIKey string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Using enviroment variables.")
	} else {
		fmt.Println("Loaded .env file.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	platform := os.Getenv("PLATFORM")
	if platform != "dev" && platform != "prod" {
		log.Fatal("PLATFORM must be set to either dev or prod")
	}

	mgAPIKey := os.Getenv("MG_API_KEY")
	if mgAPIKey == "" {
		log.Fatal("A Mailgun API key is required")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		platform: platform,
		db:       dbQueries,
		mgAPIKey: mgAPIKey,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /admin/healthz", readinessEndpoint)

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		MainPage(&apiCfg).Render(r.Context(), w)
	})

	mux.HandleFunc("GET /user", func(w http.ResponseWriter, r *http.Request) {
		NewUserPage().Render(r.Context(), w)
	})

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting pragmatic-cooking in %s mode on port: %s\n", platform, port)
	log.Fatal(server.ListenAndServe())
}
