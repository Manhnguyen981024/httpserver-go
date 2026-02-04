package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Manhnguyen981024/httpserver-go/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	var req loginRequest
	jsonDecoder := json.NewDecoder(r.Body)
	if err := jsonDecoder.Decode(&req); err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := cfg.DB.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	match, err := auth.CheckPasswordHash(req.Password, user.HashedPassword)
	if err != nil || !match {
		responseWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	realSecond := req.ExpiresInSeconds
	if realSecond <= 0 {
		realSecond = 3600 // default 1 hour
	}
	token, err := auth.MakeJWT(user.ID, cfg.secretKey, time.Duration(realSecond)*time.Second)

	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Could not generate token")
		return
	}

	user.Token = token
	responseWithJSON(w, http.StatusOK, user)
}
