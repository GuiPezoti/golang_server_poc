package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirpDB, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve all the chirps", err)
		return
	}
	chirps := []Chirp{}
	for _, dbChirp := range chirpDB {
    chirps = append(chirps, Chirp{
        ID:        dbChirp.ID,
        CreatedAt: dbChirp.CreatedAt,
        UpdatedAt: dbChirp.UpdatedAt,
        UserID:    dbChirp.UserID,
        Body:      dbChirp.Body,
    })
	}
	respondWithJSON(w, http.StatusOK, chirps)
}