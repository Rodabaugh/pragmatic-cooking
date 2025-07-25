package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type LoginToken struct {
	Token     uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
	ExpireAt  time.Time
}

func (cfg *apiConfig) CreateLoginToken (user uuid.UUID, r *http.Request) (string, error) {

	_, err := cfg.db.GetUserByID(r.Context(), user)
	if err != nil {
		return "", fmt.Errorf("User not found")
	}

	token, err := cfg.db.CreateLoginToken(r.Context(), user)
	if err != nil {
		log.Printf("Failed to create login token.")
	}

	return token.UserID.String(), nil
}
