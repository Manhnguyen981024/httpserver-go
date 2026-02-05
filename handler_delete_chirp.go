package main

import (
	"net/http"

	"github.com/Manhnguyen981024/httpserver-go/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	chirpIdStr := r.PathValue("chirpID")
	userId, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Invalid access token")
		return
	}

	chirpId, err := uuid.Parse(chirpIdStr)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Chirp id not valid ")
	}

	chirpData, err := cfg.DB.GetChirpsByChirpID(r.Context(), chirpId)
	if err != nil {
		responseWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	if chirpData.UserID != userId {
		responseWithError(w, http.StatusForbidden, "You are not authorized to delete this chirp")
		return
	}

	err = cfg.DB.DeleteChirpByChirpId(r.Context(), chirpId)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Could not delete chirp")
		return
	}

	responseWithJSON(w, http.StatusNoContent, nil)
}
