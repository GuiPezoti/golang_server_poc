package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/GuiPezoti/golang_server_poc/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
		ExpiresInSeconds int `json:"expires_in_seconds"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	userDB, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	isAuth, err := auth.CheckPasswordHash(params.Password, userDB.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to check password", err)
		return
	}
	if isAuth {
		expirationTime := expirationTimeCalc(params.ExpiresInSeconds)
		jwt, err := auth.MakeJWT(
			userDB.ID,
			cfg.jwtSecret,
			expirationTime,
		)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Error generating jwt", err)
			return
		}
		type response struct {
			User
			Token string `json:"token"`
		}
		respondWithJSON(w, http.StatusOK, response{
		User: User{
        ID:        userDB.ID,
        CreatedAt: userDB.CreatedAt,
        UpdatedAt: userDB.UpdatedAt,
        Email:     userDB.Email,
    	},
		Token: jwt,
		})
	} else {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
}

func expirationTimeCalc(clientExpirationTime int)(time.Duration) {
	if clientExpirationTime == 0 || clientExpirationTime > 3600 {
		return time.Hour
	}
	return time.Duration(clientExpirationTime) * time.Second
}