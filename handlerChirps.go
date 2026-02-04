package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DB.GetChirps(r.Context())
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	responseWithJSON(w, http.StatusOK, chirps)
}
