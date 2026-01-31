package main

import (
	"net/http"
)

type apiConfig struct{
	fileserverHits atomic.Int32
}

func (a *apiConfig) incrementFileserverHits(){
	a.fileserverHits.Add(1)
}
func (a *apiConfig) getFileserverHits() int32{
	return a.fileserverHits.Load()
}
func (a *apiConfig) resetFileserverHits(){
	a.fileserverHits.Store(0)
}

func  main()  {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("./app"))))
	mux.HandleFunc("/healthz", handlerReadiness)
	server.ListenAndServe()
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}