package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirpByUserID(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	chirps, err := cfg.DB.GetChirpsByUserID(r.Context(), chirpID)
	if err != nil {
		responseWithError(w, http.StatusNotFound, "No chirps found for the given chirp ID")
		return
	}

	responseWithJSON(w, http.StatusOK, chirps)
}
