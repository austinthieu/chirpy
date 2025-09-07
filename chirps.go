package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/austinthieu/chirpy/internal/auth"
	"github.com/austinthieu/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleCreateChirp(rw http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(rw, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(rw, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(rw, 500, "Couldn't decode parameters", err)
		return
	}

	cleanedBody, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, err.Error(), err)
		return
	}

	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
	})
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Error creating chirp: ", err)
		return
	}

	respondWithJSON(rw, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	},
	)
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	return cleanBody(body), nil
}

func cleanBody(body string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	words := strings.Split(body, " ")

	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

func (cfg *apiConfig) handleGetChirps(rw http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Error getting chirps", err)
		return
	}

	chirpsMapped := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		chirpsMapped[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	respondWithJSON(rw, http.StatusOK, chirpsMapped)
}

func (cfg *apiConfig) handleGetChirp(rw http.ResponseWriter, req *http.Request) {
	chirpID := req.PathValue("chirpID")
	ID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(rw, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}
	dbChirp, err := cfg.db.GetChirp(req.Context(), ID)
	if err != nil {
		respondWithError(rw, http.StatusNotFound, "Chirp not found", err)
		return
	}

	respondWithJSON(rw, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	})
}

func (cfg *apiConfig) handleDeleteChirp(rw http.ResponseWriter, req *http.Request) {
	idFromPath := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(idFromPath)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	chirp, err := cfg.db.GetChirp(req.Context(), chirpID)
	if err != nil {
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	if chirp.UserID != userID {
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	err = cfg.db.DeleteChirp(req.Context(), chirp.ID)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
