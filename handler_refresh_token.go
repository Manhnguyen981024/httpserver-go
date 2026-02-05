package main

import (
	"net/http"
	"time"

	"github.com/Manhnguyen981024/httpserver-go/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	storedToken, err := cfg.DB.GetRefreshTokenByToken(r.Context(), refreshToken)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	if time.Now().UTC().After(storedToken.ExpiresAt.Time) {
		responseWithError(w, http.StatusUnauthorized, "Refresh token has expired")
		return
	}

	if storedToken.RevokedAt.Valid {
		responseWithError(w, http.StatusUnauthorized, "Refresh token is revoked")
		return
	}

	newToken, err := auth.MakeJWT(storedToken.UserID, cfg.secretKey, time.Duration(3600)*time.Second)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Could not generate token")
		return
	}

	type refreshTokenResponse struct {
		Token string `json:"token"`
	}

	resp := refreshTokenResponse{
		Token: newToken,
	}

	responseWithJSON(w, http.StatusOK, resp)
}
