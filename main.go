package main

import (
	"log"
	"net/http"
)



func main(){
	const filepathRoot = "."
	mux := http.NewServeMux()
	mux.Handle(
		"/app/",
		http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))),
	)
	mux.HandleFunc("/healthz", handlerHealthCheck)
	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}


func handlerHealthCheck(w http.ResponseWriter, r *http.Request){
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(http.StatusText(http.StatusOK)))
}


