package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Manhnguyen981024/httpserver-go/internal/auth"
	"github.com/Manhnguyen981024/httpserver-go/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerAddUsers(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	jsonDecoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := jsonDecoder.Decode(&params)
	if err != nil {
		log.Printf("Can not decode a json body: %v", err)
		responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Can not hash password: %v", err)
		responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	newUser, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		log.Printf("Can not insert data to DB: %v", err)
		responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	responseWithJSON(w, http.StatusCreated, newUser)
}
