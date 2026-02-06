package main

import (
	"net/http"
	"sort"

	"github.com/Manhnguyen981024/httpserver-go/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")
	sort := r.URL.Query().Get("sort")

	// query all chirp
	if s == "" {
		chirps, err := cfg.DB.GetChirps(r.Context())
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}
		sortChirps(chirps, sort)
		responseWithJSON(w, http.StatusOK, chirps)
		return
	}

	authorId, err := uuid.Parse(s)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid author id")
		return
	}
	// query chirp by author id
	chirps, err := cfg.DB.GetChirpsByUserID(r.Context(), authorId)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if sort != "" {
		sortChirps(chirps, sort)
	}
	responseWithJSON(w, http.StatusOK, chirps)
}

func sortChirps(chirps []database.Chirp, sortStr string) {
	if sortStr == "asc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
		return
	}
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
	})
}
