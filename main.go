package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)



func main(){
	const filepathRoot = "."
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
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
	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) mwMetricsReset(
	next func(http.ResponseWriter, *http.Request),
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		cfg.fileserverHits.Store(0)
		next(w,r)
	}
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
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

func handlerHealthCheck(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}


