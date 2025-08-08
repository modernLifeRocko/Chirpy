package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct{
	Id uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding input: %s", err)
		w.WriteHeader(500)
		return
	}


	w.Header().Add("Content-Type", "application/json")
	returnedUser, err := cfg.dbQueries.CreateUser(r.Context(), params.Email)
	newUser := User{
		Id: returnedUser.ID,
		CreatedAt: returnedUser.CreatedAt,
		UpdatedAt: returnedUser.UpdatedAt,
		Email: returnedUser.Email,
	}
	if err != nil {
		log.Printf("Error creating user: %s", err)
		w.WriteHeader(500)
	}
	w.WriteHeader(201)
	dat, err := json.Marshal(newUser)
	if err != nil {
		log.Printf("Could not write Json response: %s", err)
		return
	}
	w.Write(dat)
}
