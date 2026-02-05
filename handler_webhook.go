package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Manhnguyen981024/httpserver-go/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type requestParameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId string `json:"user_id"`
		} `json:"data"`
	}
	request := &requestParameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(request)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if request.Event != "user.upgraded" {
		responseWithJSON(w, http.StatusNoContent, nil)
		return
	}

	userId, err := uuid.Parse(request.Data.UserId)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid user id")
	}

	log.Printf("User id :%v", userId)
	userData, err := cfg.DB.GetUserByID(r.Context(), userId)
	if err != nil {
		responseWithError(w, http.StatusNotFound, "User not found")
		return
	}

	_, errUp := cfg.DB.UpdateUserChirpyRed(r.Context(), database.UpdateUserChirpyRedParams{
		ID:          userData.ID,
		IsChirpyRed: true,
	})

	if errUp != nil {
		responseWithError(w, http.StatusInternalServerError, "Could not update user chirpy red status")
		return
	}

	responseWithJSON(w, http.StatusNoContent, nil)
}
