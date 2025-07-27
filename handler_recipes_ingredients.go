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

type RecipeIngredient struct {
	RecipeID     uuid.UUID
	IngredientID uuid.UUID
	Quantity     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (cfg *apiConfig) handlerCreateRecipeIngredient(w http.ResponseWriter, r *http.Request) {
	recipeID, err := uuid.Parse(r.PathValue("recipeID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
	}

	type parameters struct {
		RecipeIngredientID uuid.UUID `json:"ingredient_id"`
		RecipeIngredientQuantity string `json:"quantity"`
	}

	type response struct {
		RecipeIngredient
	}

	requesterID := cfg.getRequestUserID(r)
	fmt.Println(requesterID)
	if requesterID == uuid.Nil {
		respondWithError(w, http.StatusUnauthorized, "User is not logged in", fmt.Errorf("User is not logged in"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Was unable to decode parameters", err)
		return
	}

	recipeIngredient, err := cfg.db.CreateRecipeIngredient(r.Context(), database.CreateRecipeIngredientParams{
		RecipeID: recipeID,
		IngredientID: params.RecipeIngredientID,
		Quantity: params.RecipeIngredientQuantity,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create recipe ingredient", err)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		respondWithJSON(w, http.StatusCreated, response{
			RecipeIngredient{
				RecipeID:	recipeIngredient.RecipeID,
				IngredientID: recipeIngredient.IngredientID,
				Quantity: recipeIngredient.Quantity,
			},
		})
	} else {
		RecipeIngredientsList(cfg.RecipeIngredients()).Render(r.Context(), w)
	}
}

func (cfg *apiConfig) handlerDeleteRecipeIngredient(w http.ResponseWriter, r *http.Request) {
	recipeID, err := uuid.Parse(r.PathValue("recipeID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid recipe ID", err)
	}

	ingredientID, err := uuid.Parse(r.PathValue("ingredientID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ingredient ID", err)
	}

	type response struct {
		RecipeIngredient
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

	err = cfg.db.DeleteRecipeIngredient(r.Context(), database.DeleteRecipeIngredientParams{
		RecipeID: recipeID,
		IngredientID: ingredientID,
	})

	fmt.Printf("Deleted Recipe Ingredient: %v", err)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to delete recipe", err)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		respondWithJSON(w, http.StatusNoContent, response{
			RecipeIngredient{},
		})
	} else {
		RecipeIngredientsList(cfg.RecipeIngredients()).Render(r.Context(), w)
	}
}

func (cfg *apiConfig) RecipeIngredients() ([]RecipeIngredient, error) {
	databaseRecipeIngredients, err := cfg.db.GetAllRecipeIngredients(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to get recipes from database", err)
	}

	recipes := []RecipeIngredient{}

	for _, dbRecipeIngredient := range databaseRecipeIngredients {
		recipes = append(recipes, RecipeIngredient{
			RecipeID:	dbRecipeIngredient.RecipeID,
			IngredientID: dbRecipeIngredient.IngredientID,
			Quantity: dbRecipeIngredient.Quantity,
		})
	}

	return recipes, nil
}
