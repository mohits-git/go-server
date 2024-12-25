package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mohits-git/experiments/go-server/internal/auth"
	"github.com/mohits-git/experiments/go-server/internal/database"
)

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var reqBody requestBody

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if reqBody.Email == "" || reqBody.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and Password are required")
		return
	}

	hashedPassword, err := auth.HashPassword(reqBody.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:    reqBody.Email,
		Password: hashedPassword,
	})
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

func (cfg *apiConfig) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		ExpiresIn int    `json:"expiresIn"`
	}
	var reqBody requestBody

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if reqBody.Email == "" || reqBody.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and Password are required")
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), reqBody.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	err = auth.ComparePassword(user.Password, reqBody.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

  expiresIn := time.Second * time.Duration(reqBody.ExpiresIn)
  if reqBody.ExpiresIn == 0 {
    expiresIn = time.Hour
  }

  token, err := auth.MakeJwt(user.ID, cfg.jwtSecret, expiresIn)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error())
    return
  }

	res := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
    Token:     token,
	}

	respondWithJson(w, http.StatusOK, res)
}
