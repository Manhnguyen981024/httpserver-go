package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Manhnguyen981024/httpserver-go/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	DB             *database.Queries
	Flatform       string
	secretKey      string
}

func (a *apiConfig) incrementFileserverHits() {
	a.fileserverHits.Add(1)
}
func (a *apiConfig) getFileserverHits() int32 {
	return a.fileserverHits.Load()
}

func (cfg *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.incrementFileserverHits()
		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load()
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:     ":8080",
		Handler:  mux,
		ErrorLog: log.New(os.Stdout, "[HTTP-ERROR] ", log.LstdFlags),
	}
	dbURL := os.Getenv("DB_URL")

	log.Printf("URL DB : %v ", dbURL)
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer dbConn.Close()
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		DB:             database.New(dbConn),
		Flatform:       os.Getenv("PLATFORM"),
		secretKey:      os.Getenv("SECRET_KEY"),
	}

	mux.Handle("/app/", cfg.middlewareMetricInc(http.StripPrefix("/app/", http.FileServer(http.Dir("./app")))))
	mux.HandleFunc("GET /admin/metrics", cfg.handlerAdminMetric)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/chirps", cfg.handlerChirps)
	mux.HandleFunc("GET /api/chirps/{id}", cfg.handlerGetChirpByUserID)

	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", cfg.handlerValidateChirp)
	mux.HandleFunc("POST /api/users", cfg.handlerAddUsers)
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)

	server.ListenAndServe()
}

func handlerResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
