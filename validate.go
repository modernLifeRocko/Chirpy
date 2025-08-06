package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

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
