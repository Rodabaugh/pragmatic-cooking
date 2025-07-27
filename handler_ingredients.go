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

type Ingredient struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Unit      string
}

func (cfg *apiConfig) handlerCreateIngredient(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		IngredientName string `json:"ingredient_name"`
		IngredientUnit string `json:"ingredient_unit"`
	}

	type response struct {
		Ingredient
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

	ingredient, err := cfg.db.CreateIngredient(r.Context(), database.CreateIngredientParams{
		Name:    params.IngredientName,
		Unit:    params.IngredientUnit,
		OwnerID: requesterID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create Ingredient", err)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		respondWithJSON(w, http.StatusCreated, response{
			Ingredient{
				ID:        ingredient.ID,
				CreatedAt: ingredient.CreatedAt,
				UpdatedAt: ingredient.UpdatedAt,
				Name:      ingredient.Name,
				Unit:      ingredient.Unit,
			},
		})
	} else {
		IngredientsList(cfg.Ingredients()).Render(r.Context(), w)
	}
}

func (cfg *apiConfig) handlerDeleteIngredient(w http.ResponseWriter, r *http.Request) {
	ingredientID, err := uuid.Parse(r.PathValue("ingredientID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
	}

	type response struct {
		Ingredient
	}

	requesterID := cfg.getRequestUserID(r)
	if requesterID == uuid.Nil {
		respondWithError(w, http.StatusUnauthorized, "User is not logged in", fmt.Errorf("User is not logged in"))
		return
	}

	ingredient, err := cfg.db.GetIngredientByID(r.Context(), ingredientID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Unable to get ingredient with that ID", err)
		return
	}

	if requesterID != ingredient.OwnerID {
		respondWithError(w, http.StatusUnauthorized, "You don't own that ingredient", err)
		return
	}

	err = cfg.db.DeleteIngrendientByID(r.Context(), ingredient.ID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to delete Ingredient", err)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		respondWithJSON(w, http.StatusNoContent, response{
			Ingredient{},
		})
	} else {
		IngredientsList(cfg.Ingredients()).Render(r.Context(), w)
	}
}

func (cfg *apiConfig) handlerIngredientPage(w http.ResponseWriter, r *http.Request) {
	ingredientID, err := uuid.Parse(r.PathValue("ingredientID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
	}

	dbIngredient, err := cfg.db.GetIngredientByID(r.Context(), ingredientID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Unable to get ingredient with that ID", err)
		return
	}

	ingredient := Ingredient{
		ID:        dbIngredient.ID,
		CreatedAt: dbIngredient.CreatedAt,
		UpdatedAt: dbIngredient.UpdatedAt,
		Name:      dbIngredient.Name,
		Unit:      dbIngredient.Unit,
	}

	fmt.Printf("Surving recipes using %s\n", ingredient.Name)
	IngredientPage(cfg, ingredient).Render(r.Context(), w)
}

func (cfg *apiConfig) Ingredients() []Ingredient {
	databaseIngredients, _ := cfg.db.GetAllIngredients(context.Background())

	ingredients := []Ingredient{}

	for _, dbIngredient := range databaseIngredients {
		ingredients = append(ingredients, Ingredient{
			ID:        dbIngredient.ID,
			CreatedAt: dbIngredient.CreatedAt,
			UpdatedAt: dbIngredient.UpdatedAt,
			Name:      dbIngredient.Name,
			Unit:      dbIngredient.Unit,
		})
	}

	return ingredients
}
