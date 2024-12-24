package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Email string `json:"email"`
	}
	var reqBody requestBody

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if reqBody.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email is required")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), reqBody.Email)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			respondWithError(w, http.StatusConflict, "Email already exists")
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	userPayload := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJson(w, http.StatusCreated, userPayload)
}
