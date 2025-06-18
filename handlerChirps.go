package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/Rishan-Jadva/chirpy/internal/auth"
	"github.com/Rishan-Jadva/chirpy/internal/database"
	"github.com/google/uuid"
)

type returnChirp struct {
	ID         uuid.UUID `json:"id"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
	Body       string    `json:"body"`
	User_ID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body    string    `json:"body"`
		User_ID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanedBody := cleanWords(params.Body)

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, returnChirp{
		ID:         chirp.ID,
		Created_At: chirp.CreatedAt,
		Updated_At: chirp.UpdatedAt,
		Body:       chirp.Body,
		User_ID:    chirp.UserID,
	})
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")

	var chirps []database.Chirp
	var err error

	if authorID != "" {
		var parsedUUID uuid.UUID
		parsedUUID, err = uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't parse author id", err)
			return
		}
		chirps, err = cfg.db.GetChirpsByAuthorID(r.Context(), parsedUUID)
	} else {
		chirps, err = cfg.db.GetChirps(r.Context())
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	switch sortOrder {
	case "desc":
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	default:
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
	}

	response := make([]returnChirp, 0, len(chirps))
	for _, chirp := range chirps {
		response = append(response, returnChirp{
			ID:         chirp.ID,
			Created_At: chirp.CreatedAt,
			Updated_At: chirp.UpdatedAt,
			Body:       chirp.Body,
			User_ID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpIDstr := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDstr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID format", err)
		return
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Chirp not found", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirp by ID", err)
	}

	respondWithJSON(w, http.StatusOK, returnChirp{
		ID:         chirp.ID,
		Created_At: chirp.CreatedAt,
		Updated_At: chirp.UpdatedAt,
		Body:       chirp.Body,
		User_ID:    chirp.UserID,
	})
}

func cleanWords(body string) string {
	profaneWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}

	words := strings.Split(body, " ")

	for i, word := range words {
		if profaneWords[strings.ToLower(word)] {
			words[i] = "****"
		}
	}

	profaneLessWords := strings.Join(words, " ")
	return profaneLessWords
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or malformed token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "invalid access token", err)
		return
	}

	chirpIDstr := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDstr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID format", err)
		return
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find chirp", err)
		} else {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirp", err)
		}
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Foribbiden", err)
		return
	}

	err = cfg.db.DeleteChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
