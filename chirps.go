package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/modernLifeRocko/Chirpy/internal/database"
)

type chirp struct{
	Id uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

func (cfg apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type chirpReq struct {
		Body string `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	chirp := chirpReq{}
	err := decoder.Decode(&chirp)
	if err != nil {
		log.Fatalf("Could not decode request: %s", err)
		w.WriteHeader(500)
	} 
	if len(chirp.Body) > 140 {
		w.WriteHeader(400)
		log.Fatal("Chirp is too long")
	}

	// post to chirp database 
	Chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: chirp.Body,
		UserID: chirp.UserId,
	})

	if err != nil {
		w.WriteHeader(400)
		log.Fatalf("Error creating new chirp: %s", err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write([]byte(fmt.Sprintf(
		`{
		"id": "%s",
		"created_at": "%s",
		"updated_at":"%s",
		"body":"%s",
		"user_id": "%s"
		}`,
		Chirp.ID,
		Chirp.CreatedAt,
		Chirp.UpdatedAt,
		cleanBody(Chirp.Body),
		Chirp.UserID,
		),
		))
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirplst, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		w.WriteHeader(500)
		log.Fatalf("Failed to get chirps: %s", err)
	}
	rtnChrps := make([]chirp, len(chirplst))
	for i := range chirplst {
		rtnChrps[i].Id = chirplst[i].ID
		rtnChrps[i].CreatedAt = chirplst[i].CreatedAt
		rtnChrps[i].UpdatedAt = chirplst[i].UpdatedAt
		rtnChrps[i].Body = chirplst[i].Body
		rtnChrps[i].UserId = chirplst[i].UserID
	} 

	data, err := json.Marshal(rtnChrps)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Couldnot Marshal chirps: %s", err)
	}
	w.Write(data)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	// log.SetOutput()
	fmt.Print("helllo")
	log.Print("hello")
	chirpId := r.PathValue("chirpID")
	chirpUID, err := uuid.Parse(chirpId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatalf("Failed to parse id: %s", err)
	}
	chirpFull, err := cfg.dbQueries.GetChirpByID(r.Context(), chirpUID)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Fatalf("Failed to recover chirp: %s", err)
	}
	w.WriteHeader(http.StatusOK)
	rtnChirp := chirp{
		Id: chirpFull.ID,
		CreatedAt: chirpFull.CreatedAt,
		UpdatedAt: chirpFull.UpdatedAt,
		Body: chirpFull.Body,
		UserId: chirpFull.UserID,
	}
	log.Print(rtnChirp)
	dat, err := json.Marshal(rtnChirp)
	if err != nil {
		w.WriteHeader(500)
		log.Fatalf("couldnot create json response: %s", err)
	}
	w.Write(dat)
}
