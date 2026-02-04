package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	handlerResponse(w)
	if cfg.Flatform != "dev" {
		w.WriteHeader(http.StatusForbidden)
	} else {
		cfg.DB.DeleteAllUsers(r.Context())
		w.WriteHeader(http.StatusNoContent)
	}
	w.Write([]byte("OK"))
}
