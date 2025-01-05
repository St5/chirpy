package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/St5/goboot-srv/internal/auth"
	"github.com/St5/goboot-srv/internal/database"
	"github.com/google/uuid"
)


type Webhook struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handleWebhook(w http.ResponseWriter, r *http.Request) {
	//Authorize request

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	if apiKey != cfg.PolkaKey {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	//Decode request
	decoder := json.NewDecoder(r.Body)
	reqWebhook := Webhook{}
	err = decoder.Decode(&reqWebhook)

	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	//Validate event
	if reqWebhook.Event != "user.upgraded" {
		respondWithError(w, 204, "Event not supported")
		return
	}

	if reqWebhook.Data.UserID == "" {
		respondWithError(w, 400, "User ID is required")
		return
	}

	userID, err := uuid.Parse(reqWebhook.Data.UserID)
	if err != nil {
		respondWithError(w, 400, "Invalid User ID format")
		return
	}

	_, err = cfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, 404, "User not found")
		return
	}


	err = cfg.db.UpdateChirpyRedByUserID(r.Context(), database.UpdateChirpyRedByUserIDParams{
		ID: userID,
		IsChirpyRed: sql.NullBool{Bool: true, Valid: true},
	})
	if err != nil {
		respondWithError(w, 404, "Something went wrong")
		return
	}

	respondWithJSON(w, 204, nil)

}