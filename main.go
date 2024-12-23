package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mohits-git/experiments/go-server/internal/database"
)

func main() {
  godotenv.Load()
  dbUrl := os.Getenv("DB_URL")

  db, err := sql.Open("postgres", dbUrl)
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  dbQueries := database.New(db)
  _ = dbQueries

	port := "8080"
	filepathRoot := "."

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
