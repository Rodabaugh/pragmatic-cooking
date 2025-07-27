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
	platform  string
	db        *database.Queries
	jwtSecret string
	mgAPIKey  string
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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
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
		platform:  platform,
		db:        dbQueries,
		jwtSecret: jwtSecret,
		mgAPIKey:  mgAPIKey,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /admin/healthz", readinessEndpoint)

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		MainPage(&apiCfg).Render(r.Context(), w)
	})

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("GET /user", apiCfg.handlerUserPage)

	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		LoginPage().Render(r.Context(), w)
	})

	mux.HandleFunc("GET /ingredients", func(w http.ResponseWriter, r *http.Request) {
		IngredientsPage(&apiCfg).Render(r.Context(), w)
	})

	mux.HandleFunc("GET /ingredients/{ingredientID}", apiCfg.handlerIngredientPage)
	mux.HandleFunc("POST /api/ingredients", apiCfg.handlerCreateIngredient)
	mux.HandleFunc("DELETE /api/ingredients/{ingredientID}", apiCfg.handlerDeleteIngredient)

	mux.HandleFunc("GET /recipes", func(w http.ResponseWriter, r *http.Request) {
		RecipesPage(&apiCfg).Render(r.Context(), w)
	})
		
	mux.HandleFunc("GET /recipes/{recipeID}", apiCfg.handlerRecipePage)
	mux.HandleFunc("POST /api/recipes", apiCfg.handlerCreateRecipe)
	mux.HandleFunc("DELETE /api/recipes/{recipeID}", apiCfg.handlerDeleteRecipe)
	
	mux.HandleFunc("POST /api/recipes/{recipeID}", apiCfg.handlerCreateRecipeIngredient)
	mux.HandleFunc("DELETE /api/recipes/{recipeID}/{ingredientID}", apiCfg.handlerDeleteRecipeIngredient)

	mux.HandleFunc("POST /api/logins", apiCfg.handlerLoginRequest)
	mux.HandleFunc("GET /login/{login_token}", apiCfg.handlerLogin)

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
