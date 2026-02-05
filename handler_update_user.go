package main

import (
	"encoding/json"
	"net/http"

	"github.com/Manhnguyen981024/httpserver-go/internal/auth"
	"github.com/Manhnguyen981024/httpserver-go/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	type requestParameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	request := &requestParameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(request)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.secretKey)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Invalid access token")
		return
	}

	hashedPassword, err := auth.HashPassword(request.Password)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Could not hash password")
		return
	}

	updatedUser, errUpdate := cfg.DB.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
		ID:             userID,
		HashedPassword: hashedPassword,
		Email:          request.Email,
	})

	if errUpdate != nil {
		responseWithError(w, http.StatusInternalServerError, "Could not update user password")
		return
	}

	responseWithJSON(w, http.StatusOK, updatedUser)
}
