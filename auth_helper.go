package main

import (
	"fmt"
	"net/http"

	"github.com/Rodabaugh/pragmatic-cooking/internal/auth"
	"github.com/google/uuid"
)

func (cfg apiConfig) getRequestUserID(r *http.Request) (uuid.UUID){
	var cookieUserID uuid.UUID

	accessTokenCookie, err := r.Cookie("accessToken")

	if err != nil {
		if err != http.ErrNoCookie{
			fmt.Printf("Error reading cookie: %v\n", err)
		}
	} else {
		cookieUserID, err = auth.ValidateJWT(accessTokenCookie.Value, cfg.jwtSecret)
		if err != nil{
			fmt.Printf("Error validating JWT: %v\n", err)
		}
	}

	return cookieUserID
}
