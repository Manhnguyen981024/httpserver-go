package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Manhnguyen981024/httpserver-go/internal/auth"
	"github.com/Manhnguyen981024/httpserver-go/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	userId, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	errDecoder := decoder.Decode(&params)
	if errDecoder != nil {
		log.Printf("Error decoding JSON: %s", errDecoder)
		responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		responseWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	newChirp, err := cfg.DB.CreateChirps(r.Context(), database.CreateChirpsParams{
		ID:     uuid.New(),
		Body:   processProfaneWords(params.Body),
		UserID: userId,
	})

	responseWithJSON(w, http.StatusCreated, newChirp)
}

func processProfaneWords(words string) string {
	cleanedWords := []string{}
	for _, v := range strings.Fields(words) {
		if strings.ToLower(v) == "kerfuffle" || strings.ToLower(v) == "sharbert" || strings.ToLower(v) == "fornax" {
			cleanedWords = append(cleanedWords, "****")
		} else {
			cleanedWords = append(cleanedWords, v)
		}
	}
	return strings.Join(cleanedWords, " ")
}
