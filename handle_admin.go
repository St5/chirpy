package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {

	cfg.db.ResetAllUsers(r.Context())
	cfg.fileserverHits.Store(0)
	w.WriteHeader(200)
	w.Write([]byte("Hits reset to 0  "))
}

func (cfg *apiConfig) hadlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>`, cfg.fileserverHits.Load())))

}


func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

