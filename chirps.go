package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/modernLifeRocko/Chirpy/internal/database"
)


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

