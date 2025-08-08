package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) mwMetricsReset(
	next func(http.ResponseWriter, *http.Request),
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		plat := os.Getenv("PLATFORM")
		if plat != "dev" {
			log.Fatal("You don't have the permissions to reset the system")
			return
		}
		cfg.fileserverHits.Store(0)
		cfg.dbQueries.ResetUsers(r.Context())
		next(w,r)
	}
}

func (cfg *apiConfig) handlerReqCheck (w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(
		`<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>
		`,
		cfg.fileserverHits.Load()),
	))
}
