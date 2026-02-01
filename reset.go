package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	handlerResponse(w)
	cfg.resetFileserverHits()
	w.Write([]byte("OK"))
}
