package main

import (
	"net/http"

	"github.com/Manhnguyen981024/httpserver-go/internal/auth"
)

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	token, err := cfg.DB.GetRefreshTokenByToken(r.Context(), refreshToken)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	err = cfg.DB.RevokeRefreshToken(r.Context(), token.Token)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "failed to revoke token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
