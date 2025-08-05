package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
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
	mux.HandleFunc("POST /api/validate_chirp", handlerValidate)
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

func handlerValidate (w http.ResponseWriter, r *http.Request){
	type params struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	parameters := params{}
	err := decoder.Decode(&parameters)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)	
		return
	}

	type returnVal struct {
		CleanChirp string `json:"cleaned_body"`
		Errormsg string `json:"error"`
	}
	rtnBody := returnVal{}
	w.Header().Set("Content-Type", "application/json")

	if len(parameters.Body) <= 140 {
		w.WriteHeader(200)
		rtnBody.CleanChirp = cleanBody(parameters.Body)
	} else {
		w.WriteHeader(400)
		rtnBody.Errormsg = "Chirp is too long"
	}

	dat, err := json.Marshal(rtnBody)

	if err != nil {
		log.Printf("Error marshaling response: %s", err)
		w.WriteHeader(500)
	}

	w.Write(dat)
}

func cleanBody(s string) string {
	r, _ := regexp.Compile(`(?i)(\b)(kerfuffle|sharbert|fornax)(\s)`) 
	return r.ReplaceAllString(s, "$1****$3")
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


