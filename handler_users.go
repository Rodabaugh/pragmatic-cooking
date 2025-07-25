package main

import (
	"encoding/json"
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
		UserName string `json:"user_name"`
		UserEmail  string `json:"user_email"`
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
		Name:      params.UserName,
		EmailAddr: params.UserEmail,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create user", err)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		respondWithJSON(w, http.StatusCreated, response{
			User: User{
				ID:        user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				Name:      user.Name,
				EmailAddr: user.EmailAddr,
			},
		})
	} else {
		Created().Render(r.Context(), w)
	}
}
