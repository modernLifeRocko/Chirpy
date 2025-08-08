package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/modernLifeRocko/Chirpy/internal/database"
)


func main(){
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Could not open database: %s", err)
	}
	const filepathRoot = "."
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries: database.New(db),
	}
	mux := http.NewServeMux()
	mux.Handle(
		"/app/",
		apiCfg.middlewareMetricsInc(
			http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))),
		),
	)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerReqCheck)
	mux.HandleFunc("POST /admin/reset", apiCfg.mwMetricsReset(apiCfg.handlerReqCheck))
	mux.HandleFunc("GET /api/healthz", handlerHealthCheck)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidate)
	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries *database.Queries
}

func handlerHealthCheck(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}


