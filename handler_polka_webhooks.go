package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/mohits-git/experiments/go-server/internal/auth"
)

type polkaWebhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserId string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlePolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if apiKey != cfg.polkaKey {
		w.WriteHeader(http.StatusUnauthorized)
	}

	var req polkaWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userId, err := uuid.Parse(req.Data.UserId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := cfg.db.UpgradeUser(r.Context(), userId); err != nil {
		if err.Error() == "user not found" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
