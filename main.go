package main

import (
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (a *apiConfig) incrementFileserverHits() {
	a.fileserverHits.Add(1)
}
func (a *apiConfig) getFileserverHits() int32 {
	return a.fileserverHits.Load()
}
func (a *apiConfig) resetFileserverHits() {
	a.fileserverHits.Store(0)
}

func (cfg *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.incrementFileserverHits()
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux.Handle("/app/", cfg.middlewareMetricInc(http.StripPrefix("/app/", http.FileServer(http.Dir("./app")))))
	mux.HandleFunc("GET /admin/metrics", cfg.handlerAdminMetric)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	server.ListenAndServe()
}

func handlerResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
