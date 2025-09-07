package main

import (
	"net/http"
	"time"

	"github.com/austinthieu/chirpy/internal/auth"
)

func (cfg *apiConfig) handleRefresh(rw http.ResponseWriter, req *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(rw, http.StatusUnauthorized, "Refresh token not found", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(rw, http.StatusUnauthorized, "Couldn't get user for refresh token", err)
		return
	}
	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Couldn't make JWT token", err)
		return
	}

	respondWithJSON(rw, http.StatusOK, response{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handleRevoke(rw http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Error getting getting refresh_token", err)
		return
	}

	err = cfg.db.RevokeToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Error revoking session", err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
