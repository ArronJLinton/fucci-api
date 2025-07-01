package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ArronJLinton/fucci-api/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// PlayerProfileService provides business logic for player profiles
type PlayerProfileService struct {
	DB *database.Queries
}

// CreatePlayerProfileRequest represents the JSON request for creating a player profile
type CreatePlayerProfileRequest struct {
	UserID    int32      `json:"user_id"`
	TeamID    *uuid.UUID `json:"team_id"`
	Position  string     `json:"position"`
	Age       int32      `json:"age"`
	Country   string     `json:"country"`
	HeightCm  int32      `json:"height_cm"`
	Pace      int32      `json:"pace"`
	Shooting  int32      `json:"shooting"`
	Passing   int32      `json:"passing"`
	Stamina   int32      `json:"stamina"`
	Dribbling int32      `json:"dribbling"`
	Defending int32      `json:"defending"`
	Physical  int32      `json:"physical"`
}

// Handler: Create player profile
func (svc *PlayerProfileService) CreatePlayerProfile(w http.ResponseWriter, r *http.Request) {
	var req CreatePlayerProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Debug: Log the request
	log.Printf("Creating player profile for user_id: %d", req.UserID)

	// Debug: Check if user exists first
	user, err := svc.DB.GetUser(r.Context(), req.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User with ID %d does not exist", req.UserID)
			http.Error(w, fmt.Sprintf("User with ID %d does not exist", req.UserID), http.StatusBadRequest)
			return
		}
		log.Printf("Error checking if user exists: %v", err)
		http.Error(w, fmt.Sprintf("Error checking user: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("User found: %s %s (ID: %d)", user.Firstname, user.Lastname, user.ID)

	// Convert to database params
	var teamID uuid.NullUUID
	if req.TeamID != nil {
		teamID = uuid.NullUUID{UUID: *req.TeamID, Valid: true}
	}

	dbParams := database.CreatePlayerProfileParams{
		UserID:    req.UserID,
		TeamID:    teamID,
		Position:  req.Position,
		Age:       req.Age,
		Country:   req.Country,
		HeightCm:  req.HeightCm,
		Pace:      req.Pace,
		Shooting:  req.Shooting,
		Passing:   req.Passing,
		Stamina:   req.Stamina,
		Dribbling: req.Dribbling,
		Defending: req.Defending,
		Physical:  req.Physical,
	}

	profile, err := svc.DB.CreatePlayerProfile(r.Context(), dbParams)
	if err != nil {
		log.Printf("Error creating player profile: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Player profile created successfully with ID: %s", profile.ID)
	json.NewEncoder(w).Encode(profile)
}

// Handler: Get player profile
func (svc *PlayerProfileService) GetPlayerProfile(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	profile, err := svc.DB.GetPlayerProfile(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(profile)
}

// Handler: Update player profile
func (svc *PlayerProfileService) UpdatePlayerProfile(w http.ResponseWriter, r *http.Request) {
	var req database.UpdatePlayerProfileParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	profile, err := svc.DB.UpdatePlayerProfile(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(profile)
}

// Handler: Delete player profile
func (svc *PlayerProfileService) DeletePlayerProfile(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := svc.DB.DeletePlayerProfile(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// RecalculateIsVerified sets is_verified=true if 3+ unique verifications exist
func (svc *PlayerProfileService) RecalculateIsVerified(ctx context.Context, profileID string) error {
	id, err := uuid.Parse(profileID)
	if err != nil {
		return err
	}
	count, err := svc.DB.CountVerificationsByPlayer(ctx, id)
	if err != nil {
		return err
	}
	isVerified := count >= 3
	return svc.DB.UpdatePlayerVerificationStatus(ctx, database.UpdatePlayerVerificationStatusParams{
		ID:         id,
		IsVerified: isVerified,
	})
}
