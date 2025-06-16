package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	type respUser struct {
		ID         uuid.UUID `json:"id"`
		Created_At time.Time `json:"created_at"`
		Updated_At time.Time `json:"updated_at"`
		Email      string    `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode params", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, respUser{
		ID:         user.ID,
		Created_At: user.CreatedAt,
		Updated_At: user.UpdatedAt,
		Email:      user.Email,
	})
}
