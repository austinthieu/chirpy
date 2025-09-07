package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/austinthieu/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleUpgrade(rw http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(rw, http.StatusUnauthorized, "Couldn't get API key", err)
		return
	}

	if apiKey != cfg.apiKey {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	if params.Event != "user.upgraded" {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.db.UpgradeUser(req.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(rw, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		respondWithError(rw, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
