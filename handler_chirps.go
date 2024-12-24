package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/mohits-git/experiments/go-server/internal/database"
)

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	var body requestBody

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if body.Body == "" || body.UserId == uuid.Nil {
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
		UserID: body.UserId,
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
  chirps, err := cfg.db.GetChirps(r.Context())
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
