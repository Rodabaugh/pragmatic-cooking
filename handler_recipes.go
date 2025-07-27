package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Rodabaugh/pragmatic-cooking/internal/database"
	"github.com/google/uuid"
)

type Recipe struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Desc      string
	Link      string
}

func (cfg *apiConfig) handlerCreateRecipe(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		RecipeName string `json:"recipe_name"`
		RecipeDesc string `json:"recipe_desc"`
		RecipeLink string `json:"recipe_link"`
	}

	type response struct {
		Recipe
	}

	requesterID := cfg.getRequestUserID(r)
	if requesterID == uuid.Nil {
		respondWithError(w, http.StatusUnauthorized, "User is not logged in", fmt.Errorf("User is not logged in"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Was unable to decode parameters", err)
		return
	}

	recipe, err := cfg.db.CreateRecipe(r.Context(), database.CreateRecipeParams{
		Name:        params.RecipeName,
		Description: params.RecipeDesc,
		Link:        params.RecipeLink,
		OwnerID:     requesterID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create Recipe", err)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		respondWithJSON(w, http.StatusCreated, response{
			Recipe{
				ID:        recipe.ID,
				CreatedAt: recipe.CreatedAt,
				UpdatedAt: recipe.UpdatedAt,
				Name:      recipe.Name,
				Desc:      params.RecipeDesc,
				Link:      params.RecipeLink,
			},
		})
	} else {
		RecipesList(cfg.Recipes()).Render(r.Context(), w)
	}
}

func (cfg *apiConfig) handlerDeleteRecipe(w http.ResponseWriter, r *http.Request) {
	recipeID, err := uuid.Parse(r.PathValue("recipeID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
	}

	type response struct {
		Recipe
	}

	requesterID := cfg.getRequestUserID(r)
	if requesterID == uuid.Nil {
		respondWithError(w, http.StatusUnauthorized, "User is not logged in", fmt.Errorf("User is not logged in"))
		return
	}

	recipe, err := cfg.db.GetRecipeByID(r.Context(), recipeID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Unable to get recipe with that ID", err)
		return
	}

	if requesterID != recipe.OwnerID {
		respondWithError(w, http.StatusUnauthorized, "You don't own that recipe", err)
		return
	}

	err = cfg.db.DeleteRecipeByID(r.Context(), recipe.ID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to delete recipe", err)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		respondWithJSON(w, http.StatusNoContent, response{
			Recipe{},
		})
	} else {
		RecipesList(cfg.Recipes()).Render(r.Context(), w)
	}
}

func (cfg *apiConfig) handlerRecipePage(w http.ResponseWriter, r *http.Request) {
	recipeID, err := uuid.Parse(r.PathValue("recipeID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
	}

	dbRecipe, err := cfg.db.GetRecipeByID(r.Context(), recipeID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Unable to get recipe with that ID", err)
		return
	}

	recipe := Recipe{
		ID:        dbRecipe.ID,
		CreatedAt: dbRecipe.CreatedAt,
		UpdatedAt: dbRecipe.UpdatedAt,
		Name:      dbRecipe.Name,
		Desc:      dbRecipe.Description,
		Link:      dbRecipe.Link,
	}

	fmt.Printf("Surving up recipe for %s\n", recipe.Name)
	RecipePage(cfg, recipe).Render(r.Context(), w)
}

func (cfg *apiConfig) Recipes() ([]Recipe, error) {
	databaseRecipes, err := cfg.db.GetAllRecipes(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to get recipes from database", err)
	}

	recipes := []Recipe{}

	for _, dbRecipe := range databaseRecipes {
		recipes = append(recipes, Recipe{
			ID:        dbRecipe.ID,
			CreatedAt: dbRecipe.CreatedAt,
			UpdatedAt: dbRecipe.UpdatedAt,
			Name:      dbRecipe.Name,
			Desc:      dbRecipe.Description,
			Link:      dbRecipe.Link,
		})
	}

	return recipes, nil
}
