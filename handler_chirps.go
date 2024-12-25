package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/mohits-git/experiments/go-server/internal/auth"
	"github.com/mohits-git/experiments/go-server/internal/database"
)

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	headers := r.Header

	token, err := auth.GetBearerToken(headers)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	type requestBody struct {
		Body string `json:"body"`
	}
	var body requestBody

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if body.Body == "" || userId == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, "Body and User ID are required")
		return
	}
	if len(body.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	cleanedBody := cleanBody(body.Body)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create chirp")
		return
	}

	res := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	respondWithJson(w, http.StatusCreated, res)
}

func cleanBody(body string) string {
	bannedWords := []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Fields(body)

	for i, word := range words {
		for _, banned := range bannedWords {
			if strings.EqualFold(word, banned) {
				words[i] = "****"
			}
		}
	}

	return strings.Join(words, " ")
}

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	authorId := r.URL.Query().Get("author_id")
	var chirps []database.Chirp
	var err error

	if authorId == "" {
		chirps, err = cfg.db.GetChirps(r.Context())
	} else {
		userId, err := uuid.Parse(authorId)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		}
		chirps, err = cfg.db.GetChirpsByUserId(r.Context(), userId)
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get chirps")
		return
	}

	var res []Chirp
	for _, chirp := range chirps {
		res = append(res, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	respondWithJson(w, http.StatusOK, res)
}

func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Chirp ID is required")
		return
	}

	chirpId, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID")
		return
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), chirpId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			respondWithError(w, http.StatusNotFound, "Chirp not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJson(w, http.StatusOK, res)
}

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Chirp ID is required")
		return
	}

	chirpId, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), chirpId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			respondWithError(w, http.StatusNotFound, "Chirp not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to get chirp")
		return
	}

	if userId.String() != chirp.UserID.String() {
		respondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
