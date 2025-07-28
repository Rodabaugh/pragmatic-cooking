package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Rodabaugh/pragmatic-cooking/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	EmailAddr string
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserName  string `json:"user_name"`
		UserEmail string `json:"user_email"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Was unable to decode parameters", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Name:          params.UserName,
		EmailAddr:     params.UserEmail,
		EmailVerified: false,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create user", err)
		return
	}

	loginToken, err := cfg.CreateLoginToken(user.ID, r)

	var welcomeMsg string

	if err != nil {
		welcomeMsg = "An account on Pragmatic.Recepies has been created using this email address."
	} else {
		loginLink := "https://app.pragmatic.cooking.login/" + loginToken
		welcomeMsg = "An account on Pragmatic.Recepies has been created using this email address. Login with the link below.\n\n" + loginLink
	}

	cfg.sendMGEmail(user.Name, user.EmailAddr, "New Pragmatic Recepies Account", welcomeMsg)

	if r.Header.Get("Accept") == "application/json" {
		respondWithJSON(w, http.StatusCreated, response{
			User{
				ID:        user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				Name:      user.Name,
				EmailAddr: user.EmailAddr,
			},
		})
	} else {
		AccountCreated().Render(r.Context(), w)
	}
}

func (cfg *apiConfig) handlerUserPage(w http.ResponseWriter, r *http.Request) {
	userID := cfg.getRequestUserID(r)

	if userID == uuid.Nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	dbUser, err := cfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		fmt.Errorf("Was unable to get the user after they logged in: %v", err)
	}

	user := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		EmailAddr: dbUser.EmailAddr,
	}

	dbUserRecipes, err := cfg.db.GetRecipesByOwner(r.Context(), userID)
	if err != nil {
		println(err)
	}

	userRecipes := []Recipe{}

	for _, dbRecipe := range dbUserRecipes {
		userRecipes = append(userRecipes, Recipe{
			ID:        dbRecipe.ID,
			CreatedAt: dbRecipe.CreatedAt,
			UpdatedAt: dbRecipe.UpdatedAt,
			Name:      dbRecipe.Name,
			Desc:      dbRecipe.Description,
			Link:      dbRecipe.Link,
		})
	}

	UserPage(user, userRecipes).Render(r.Context(), w)
}
