package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Body string `json:"body"`
	}
	var body requestBody
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if body.Body == "" {
		respondWithError(w, http.StatusBadRequest, "Chirp is empty")
    return
	}
	if len(body.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	cleanedBody := cleanBody(body.Body)

	type validResponse struct {
		Valid       bool   `json:"valid"`
		CleanedBody string `json:"cleaned_body"`
	}
	res := validResponse{Valid: true, CleanedBody: cleanedBody}
	respondWithJson(w, http.StatusOK, res)
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
