package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		responseWithError(w, 500, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		responseWithError(w, 400, "Chirp is too long")
		return
	}

	type responseSuccess struct {
		CleanedBody string `json:"cleaned_body"`
	}
	respdata := responseSuccess{CleanedBody: processProfaneWords(params.Body)}
	responseWithJSON(w, 200, respdata)
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
