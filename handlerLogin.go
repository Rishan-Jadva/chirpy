package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Rishan-Jadva/chirpy/internal/auth"
	"github.com/google/uuid"
)

type parameters struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
}

type respLogin struct {
	ID         uuid.UUID `json:"id"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
	Email      string    `json:"email"`
	Token      string    `json:"token"`
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

	duration := time.Hour
	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds <= 3600 {
		duration = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	token, err := auth.MakeJWT(user.ID, cfg.JWTSecret, duration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, respLogin{
		ID:         user.ID,
		Created_At: user.CreatedAt,
		Updated_At: user.UpdatedAt,
		Email:      user.Email,
		Token:      token,
	})
}
