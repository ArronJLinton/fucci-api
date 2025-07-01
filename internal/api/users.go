package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ArronJLinton/fucci-api/internal/database"
)

func (config *Config) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := database.CreateUserParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	user, err := config.DB.CreateUser(r.Context(), database.CreateUserParams{
		Firstname: params.Firstname,
		Lastname:  params.Lastname,
		Email:     params.Email,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error creating user: %s", err))
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

func (config *Config) handleListAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := config.DB.ListUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error listing users: %s", err))
		return
	}
	respondWithJSON(w, http.StatusOK, users)
}
