package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Rishan-Jadva/chirpy/internal/auth"
	"github.com/Rishan-Jadva/chirpy/internal/database"
	"github.com/google/uuid"
)

type parameters struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type respLogin struct {
	ID            uuid.UUID `json:"id"`
	Created_At    time.Time `json:"created_at"`
	Updated_At    time.Time `json:"updated_at"`
	Email         string    `json:"email"`
	IsChirpyRed   bool      `json:"is_chirpy_red"`
	Token         string    `json:"token"`
	Refresh_Token string    `json:"refresh_token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode params", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.JWTSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't add refresh token to db", err)
		return
	}

	respondWithJSON(w, http.StatusOK, respLogin{
		ID:            user.ID,
		Created_At:    user.CreatedAt,
		Updated_At:    user.UpdatedAt,
		Email:         user.Email,
		IsChirpyRed:   user.IsChirpyRed,
		Token:         accessToken,
		Refresh_Token: refreshToken,
	})
}
