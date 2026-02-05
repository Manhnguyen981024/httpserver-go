package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Manhnguyen981024/httpserver-go/internal/auth"
	"github.com/Manhnguyen981024/httpserver-go/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	token, err := auth.MakeJWT(user.ID, cfg.secretKey, time.Duration(3600)*time.Second)

	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Could not generate token")
		return
	}

	refreshToken, err := cfg.DB.GetRefreshTokenByUserId(r.Context(), user.ID)
	refreshTokenStr := refreshToken.Token

	if err != nil {
		refreshTokenStr, err = createNewRefreshToken(user, cfg, r)
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, "Could not generate refresh token")
			return
		}
	}
	type loginResponse struct {
		database.GetUserByEmailRow
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	resp := loginResponse{
		GetUserByEmailRow: user,
		Token:             token,
		RefreshToken:      refreshTokenStr,
	}
	responseWithJSON(w, http.StatusOK, resp)
}

func createNewRefreshToken(user database.GetUserByEmailRow, cfg *apiConfig, r *http.Request) (string, error) {
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		return "", err
	}
	tokenExpiry := time.Now().Add(1440 * time.Hour)
	_, err = cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: sql.NullTime{Time: tokenExpiry, Valid: true},
	})
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}
