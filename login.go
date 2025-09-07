package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/austinthieu/chirpy/internal/auth"
	"github.com/austinthieu/chirpy/internal/database"
)

func (cfg *apiConfig) handleLogin(rw http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUser(req.Context(), params.Email)
	if err != nil {
		respondWithError(rw, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(rw, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Couldn't make JWT token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Couldn't make refresh token", err)
		return
	}

	_, err = cfg.db.AddNewRefreshToken(req.Context(), database.AddNewRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		RevokedAt: sql.NullTime{Time: time.Time{}, Valid: false},
	})
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Couldn't add refresh token to database", err)
		return
	}

	respondWithJSON(rw, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token:        token,
		RefreshToken: refreshToken,
	})
}
