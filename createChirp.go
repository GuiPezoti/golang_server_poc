package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/GuiPezoti/golang_server_poc/internal/auth"
	"github.com/GuiPezoti/golang_server_poc/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID	uuid.UUID	`json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time	`json:"updated_at"`
	Body	string `json:"body"`
	UserID	uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't decode parameters", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't decode parameters", err)
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	chirpParams := database.CreateChirpParams{
		Body: profanityValidation(params.Body),
		UserID: userID,
	}

	chirpDB, err := cfg.db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID: chirpDB.ID,
		CreatedAt: chirpDB.CreatedAt,
		UpdatedAt: chirpDB.UpdatedAt,
		Body: chirpDB.Body,
		UserID: chirpDB.UserID,
	})
}

func profanityValidation(message string) string {
	profanities := []string{"kerfuffle", "sharbert", "fornax"}
	messageWords := strings.Split(message, " ")
	for i, word := range(messageWords) {
		isProfanity := slices.Contains(profanities, strings.ToLower(word))
		if isProfanity {
			messageWords[i] = "****"
		}
	}
	return strings.Join(messageWords, " ")
}
 