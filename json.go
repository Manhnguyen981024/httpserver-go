package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func responseWithError(w http.ResponseWriter, httpCode int, msgString string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	resp := errorResponse{
		Error: msgString,
	}
	responseWithJSON(w, httpCode, resp)
}

func responseWithJSON(w http.ResponseWriter, code int, payload any) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		responseWithError(w, 500, "Something went wrong")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonData)
}
