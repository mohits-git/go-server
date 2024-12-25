package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mohits-git/experiments/go-server/internal/auth"
	"github.com/mohits-git/experiments/go-server/internal/database"
)

func (cfg *apiConfig) handleLoginUser(w http.ResponseWriter, r *http.Request) {
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

	token, err := auth.MakeJwt(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshTokenString,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken.Token,
	}

	respondWithJson(w, http.StatusOK, res)
}

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

  user, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
  if err != nil {
    respondWithError(w, http.StatusUnauthorized, "Unauthorized")
    return
  }

	newToken, err := auth.MakeJwt(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := struct {
		Token string `json:"token"`
	}{
		Token: newToken,
	}

	respondWithJson(w, http.StatusOK, res)
}

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {
  token, err := auth.GetBearerToken(r.Header)
  if err != nil {
    respondWithError(w, http.StatusUnauthorized, "Unauthorized")
    return
  }

  err = cfg.db.RevokeRefreshToken(r.Context(), token)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Unable to revoke token")
    return
  }

  w.WriteHeader(http.StatusNoContent)
}
