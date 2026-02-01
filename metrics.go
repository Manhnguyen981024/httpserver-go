package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerAdminMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	adminPage := fmt.Sprintf(`<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
	</html>`, cfg.getFileserverHits())
	w.Write([]byte(adminPage))
}
