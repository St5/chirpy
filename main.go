package main

import (
	"database/sql"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/St5/goboot-srv/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	tokenSecret    string
	PolkaKey       string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		panic(err)
	}

	secretToken := os.Getenv("TOKEN_SECRET")
	PolkaKey := os.Getenv("POLKA_KEY")

	conf := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             database.New(db),
		tokenSecret:    secretToken,
		PolkaKey:       PolkaKey,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", conf.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./public/")))))

	mux.HandleFunc("GET /admin/metrics", conf.hadlerMetrics)

	mux.HandleFunc("POST /admin/reset", conf.handlerReset)

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		w.Write([]byte("OK"))

	})

	//Users API
	mux.HandleFunc("POST /api/users", conf.handleUser)

	mux.HandleFunc("PUT /api/users", conf.handleUpdateUser)

	mux.HandleFunc("POST /api/login", conf.handleLogin)

	mux.HandleFunc("POST /api/refresh", conf.handRefresh)

	mux.HandleFunc("POST /api/revoke", conf.handleRevoke)

	//Chirps CRUD

	mux.HandleFunc("POST /api/chirps", conf.handleCreateChirp)

	mux.HandleFunc("GET /api/chirps", conf.handleGetAllChirps)

	mux.HandleFunc("GET /api/chirps/{chirpID}", conf.handleGetChirp)

	mux.HandleFunc("DELETE /api/chirps/{chirpID}", conf.handleDeleteChirp)

	//Webhooks

	mux.HandleFunc("POST /api/polka/webhooks", conf.handleWebhook)

	server := &http.Server{
		Addr:    ":8585",
		Handler: mux,
	}

	server.ListenAndServe()
}
