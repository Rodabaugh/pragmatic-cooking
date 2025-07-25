package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Rodabaugh/pragmatic-cooking/internal/auth"
	"github.com/google/uuid"
)

type LoginToken struct {
	Token     uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
	ExpireAt  time.Time
}

func (cfg *apiConfig) CreateLoginToken(user uuid.UUID, r *http.Request) (string, error) {

	_, err := cfg.db.GetUserByID(r.Context(), user)
	if err != nil {
		return "", fmt.Errorf("User not found")
	}

	token, err := cfg.db.CreateLoginToken(r.Context(), user)
	if err != nil {
		log.Printf("Failed to create login token.")
	}

	return token.Token.String(), nil
}

func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	loginToken := r.PathValue("login_token")

	fmt.Println(loginToken)

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	uuidLoginToken, err := uuid.Parse(loginToken)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Invalid token", err)
		return
	}

	token, err := apiCfg.db.GetLoginByToken(r.Context(), uuidLoginToken)

	if err != nil {
		respondWithError(w, http.StatusForbidden, "Invalid login", err)
		return
	}

	if token.ExpireAt.After(time.Now()) {
		respondWithError(w, http.StatusForbidden, "Expired token", err)
		return
	}

	apiCfg.db.DeleteToken(r.Context(), uuidLoginToken)

	accessToken, err := auth.MakeJWT(
		token.UserID,
		apiCfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		respondWithJSON(w, http.StatusOK, response{
			Token: accessToken})
	} else {
		accessTokenCookie := http.Cookie{
			Name:     "accessToken",
			Value:    accessToken,
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, &accessTokenCookie)

		UserPage().Render(r.Context(), w)
	}

	fmt.Println("Login successful")
}
