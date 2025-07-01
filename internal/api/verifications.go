package api

import (
	"encoding/json"
	"net/http"

	"github.com/ArronJLinton/fucci-api/internal/database"
	"github.com/google/uuid"
)

type VerificationService struct {
	DB               *database.Queries
	PlayerProfileSvc *PlayerProfileService
}

// Handler: Add verification
func (svc *VerificationService) AddVerification(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PlayerProfileID string `json:"player_profile_id"`
		VerifierUserID  int64  `json:"verifier_user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	profileID, err := uuid.Parse(req.PlayerProfileID)
	if err != nil {
		http.Error(w, "invalid profile id", http.StatusBadRequest)
		return
	}
	params := database.CreateVerificationParams{
		PlayerProfileID: profileID,
		VerifierUserID:  int32(req.VerifierUserID),
	}
	verification, err := svc.DB.CreateVerification(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Recalculate verification status
	svc.PlayerProfileSvc.RecalculateIsVerified(r.Context(), req.PlayerProfileID)
	json.NewEncoder(w).Encode(verification)
}

// Handler: Remove verification
func (svc *VerificationService) RemoveVerification(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	// Get verification to find player_profile_id
	verification, err := svc.DB.GetVerification(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err := svc.DB.DeleteVerification(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Recalculate verification status
	svc.PlayerProfileSvc.RecalculateIsVerified(r.Context(), verification.PlayerProfileID.String())
	w.WriteHeader(http.StatusNoContent)
}

// Handler: List verifications for a player
func (svc *VerificationService) ListVerifications(w http.ResponseWriter, r *http.Request) {
	profileIDStr := r.URL.Query().Get("player_profile_id")
	profileID, err := uuid.Parse(profileIDStr)
	if err != nil {
		http.Error(w, "invalid profile id", http.StatusBadRequest)
		return
	}
	verifications, err := svc.DB.ListVerificationsByPlayer(r.Context(), profileID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(verifications)
}
