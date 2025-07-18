package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/ArronJLinton/fucci-api/internal/ai"
	"github.com/ArronJLinton/fucci-api/internal/database"
	"github.com/go-chi/chi"
)

// Debate API types
type CreateDebateRequest struct {
	MatchID     string `json:"match_id"`
	DebateType  string `json:"debate_type"` // "pre_match" or "post_match"
	Headline    string `json:"headline"`
	Description string `json:"description"`
	AIGenerated bool   `json:"ai_generated"`
}

type GenerateDebateRequest struct {
	MatchID         string `json:"match_id"`
	DebateType      string `json:"debate_type"`                // "pre_match" or "post_match"
	ForceRegenerate bool   `json:"force_regenerate,omitempty"` // Force regeneration even if cached
}

type CreateDebateCardRequest struct {
	DebateID    int32  `json:"debate_id"`
	Stance      string `json:"stance"` // "agree", "disagree", "wildcard"
	Title       string `json:"title"`
	Description string `json:"description"`
	AIGenerated bool   `json:"ai_generated"`
}

type CreateVoteRequest struct {
	DebateCardID int32  `json:"debate_card_id"`
	VoteType     string `json:"vote_type"` // "upvote", "downvote", "emoji"
	Emoji        string `json:"emoji,omitempty"`
}

type CreateCommentRequest struct {
	DebateID        int32  `json:"debate_id"`
	ParentCommentID *int32 `json:"parent_comment_id,omitempty"`
	Content         string `json:"content"`
}

type DebateResponse struct {
	ID          int32                    `json:"id"`
	MatchID     string                   `json:"match_id"`
	DebateType  string                   `json:"debate_type"`
	Headline    string                   `json:"headline"`
	Description string                   `json:"description"`
	AIGenerated bool                     `json:"ai_generated"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
	Cards       []DebateCardResponse     `json:"cards,omitempty"`
	Analytics   *DebateAnalyticsResponse `json:"analytics,omitempty"`
}

type DebateCardResponse struct {
	ID          int32         `json:"id"`
	DebateID    int32         `json:"debate_id"`
	Stance      string        `json:"stance"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	AIGenerated bool          `json:"ai_generated"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	VoteCounts  VoteCounts    `json:"vote_counts"`
	UserVote    *VoteResponse `json:"user_vote,omitempty"`
}

type VoteCounts struct {
	Upvotes   int            `json:"upvotes"`
	Downvotes int            `json:"downvotes"`
	Emojis    map[string]int `json:"emojis"`
}

type VoteResponse struct {
	ID           int32     `json:"id"`
	DebateCardID int32     `json:"debate_card_id"`
	UserID       int32     `json:"user_id"`
	VoteType     string    `json:"vote_type"`
	Emoji        string    `json:"emoji,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type CommentResponse struct {
	ID              int32     `json:"id"`
	DebateID        int32     `json:"debate_id"`
	ParentCommentID *int32    `json:"parent_comment_id,omitempty"`
	UserID          int32     `json:"user_id"`
	UserFirstName   string    `json:"user_first_name"`
	UserLastName    string    `json:"user_last_name"`
	Content         string    `json:"content"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type DebateAnalyticsResponse struct {
	ID              int32     `json:"id"`
	DebateID        int32     `json:"debate_id"`
	TotalVotes      int       `json:"total_votes"`
	TotalComments   int       `json:"total_comments"`
	EngagementScore float64   `json:"engagement_score"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Debate API handlers
func (c *Config) createDebate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateDebateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.MatchID == "" || req.DebateType == "" || req.Headline == "" {
		respondWithError(w, http.StatusBadRequest, "match_id, debate_type, and headline are required")
		return
	}

	if req.DebateType != "pre_match" && req.DebateType != "post_match" {
		respondWithError(w, http.StatusBadRequest, "debate_type must be 'pre_match' or 'post_match'")
		return
	}

	// Create debate in database
	debate, err := c.DB.CreateDebate(ctx, database.CreateDebateParams{
		MatchID:     req.MatchID,
		DebateType:  req.DebateType,
		Headline:    req.Headline,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		AiGenerated: sql.NullBool{Bool: req.AIGenerated, Valid: true},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create debate: %v", err))
		return
	}

	// Create analytics record
	_, err = c.DB.CreateDebateAnalytics(ctx, database.CreateDebateAnalyticsParams{
		DebateID:        sql.NullInt32{Int32: debate.ID, Valid: true},
		TotalVotes:      sql.NullInt32{Int32: 0, Valid: true},
		TotalComments:   sql.NullInt32{Int32: 0, Valid: true},
		EngagementScore: sql.NullString{String: "0.0", Valid: true},
	})
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to create debate analytics: %v\n", err)
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message":   "Debate created successfully",
		"debate_id": debate.ID,
	})
}

func (c *Config) getDebate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	debateIDStr := chi.URLParam(r, "id")
	debateID, err := strconv.ParseInt(debateIDStr, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid debate ID")
		return
	}

	// Get debate
	debate, err := c.DB.GetDebate(ctx, int32(debateID))
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Debate not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get debate: %v", err))
		return
	}

	// Get debate cards
	cards, err := c.DB.GetDebateCards(ctx, sql.NullInt32{Int32: debate.ID, Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get debate cards: %v", err))
		return
	}

	// Get analytics
	analytics, err := c.DB.GetDebateAnalytics(ctx, sql.NullInt32{Int32: debate.ID, Valid: true})
	if err != nil && err != sql.ErrNoRows {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get debate analytics: %v", err))
		return
	}

	// Build response
	response := DebateResponse{
		ID:          debate.ID,
		MatchID:     debate.MatchID,
		DebateType:  debate.DebateType,
		Headline:    debate.Headline,
		Description: debate.Description.String,
		AIGenerated: debate.AiGenerated.Bool,
		CreatedAt:   debate.CreatedAt.Time,
		UpdatedAt:   debate.UpdatedAt.Time,
	}

	// Add cards with vote counts
	cardIDs := make([]int32, len(cards))
	for i, card := range cards {
		cardIDs[i] = card.ID
	}

	if len(cardIDs) > 0 {
		voteCounts, err := c.DB.GetVoteCounts(ctx, cardIDs)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get vote counts: %v", err))
			return
		}

		// Build vote counts map
		voteCountsMap := make(map[int32]VoteCounts)
		for _, vc := range voteCounts {
			if vc.DebateCardID.Valid {
				counts := voteCountsMap[vc.DebateCardID.Int32]
				switch vc.VoteType {
				case "upvote":
					counts.Upvotes = int(vc.Count)
				case "downvote":
					counts.Downvotes = int(vc.Count)
				case "emoji":
					if counts.Emojis == nil {
						counts.Emojis = make(map[string]int)
					}
					if vc.Emoji.Valid {
						counts.Emojis[vc.Emoji.String] = int(vc.Count)
					}
				}
				voteCountsMap[vc.DebateCardID.Int32] = counts
			}
		}

		// Build card responses
		for _, card := range cards {
			cardResponse := DebateCardResponse{
				ID:          card.ID,
				DebateID:    card.DebateID.Int32,
				Stance:      card.Stance,
				Title:       card.Title,
				Description: card.Description.String,
				AIGenerated: card.AiGenerated.Bool,
				CreatedAt:   card.CreatedAt.Time,
				UpdatedAt:   card.UpdatedAt.Time,
				VoteCounts:  voteCountsMap[card.ID],
			}
			response.Cards = append(response.Cards, cardResponse)
		}
	}

	// Add analytics if available
	if err == nil {
		engagementScore := 0.0
		if analytics.EngagementScore.Valid {
			// Parse engagement score from string
			if score, err := strconv.ParseFloat(analytics.EngagementScore.String, 64); err == nil {
				engagementScore = score
			}
		}

		response.Analytics = &DebateAnalyticsResponse{
			ID:              analytics.ID,
			DebateID:        analytics.DebateID.Int32,
			TotalVotes:      int(analytics.TotalVotes.Int32),
			TotalComments:   int(analytics.TotalComments.Int32),
			EngagementScore: engagementScore,
			CreatedAt:       analytics.CreatedAt.Time,
			UpdatedAt:       analytics.UpdatedAt.Time,
		}
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (c *Config) getDebatesByMatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	matchID := r.URL.Query().Get("match_id")
	if matchID == "" {
		respondWithError(w, http.StatusBadRequest, "match_id parameter is required")
		return
	}

	debates, err := c.DB.GetDebatesByMatch(ctx, matchID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get debates: %v", err))
		return
	}

	// Convert to response format
	var response []DebateResponse
	for _, debate := range debates {
		response = append(response, DebateResponse{
			ID:          debate.ID,
			MatchID:     debate.MatchID,
			DebateType:  debate.DebateType,
			Headline:    debate.Headline,
			Description: debate.Description.String,
			AIGenerated: debate.AiGenerated.Bool,
			CreatedAt:   debate.CreatedAt.Time,
			UpdatedAt:   debate.UpdatedAt.Time,
		})
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (c *Config) createDebateCard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateDebateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.DebateID == 0 || req.Stance == "" || req.Title == "" {
		respondWithError(w, http.StatusBadRequest, "debate_id, stance, and title are required")
		return
	}

	if req.Stance != "agree" && req.Stance != "disagree" && req.Stance != "wildcard" {
		respondWithError(w, http.StatusBadRequest, "stance must be 'agree', 'disagree', or 'wildcard'")
		return
	}

	// Create debate card
	card, err := c.DB.CreateDebateCard(ctx, database.CreateDebateCardParams{
		DebateID:    sql.NullInt32{Int32: req.DebateID, Valid: true},
		Stance:      req.Stance,
		Title:       req.Title,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		AiGenerated: sql.NullBool{Bool: req.AIGenerated, Valid: true},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create debate card: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Debate card created successfully",
		"card_id": card.ID,
	})
}

func (c *Config) createVote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateVoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.DebateCardID == 0 || req.VoteType == "" {
		respondWithError(w, http.StatusBadRequest, "debate_card_id and vote_type are required")
		return
	}

	if req.VoteType != "upvote" && req.VoteType != "downvote" && req.VoteType != "emoji" {
		respondWithError(w, http.StatusBadRequest, "vote_type must be 'upvote', 'downvote', or 'emoji'")
		return
	}

	// Get user ID from context (you'll need to implement authentication)
	userID := int32(1) // TODO: Get from auth context

	// Create vote
	vote, err := c.DB.CreateVote(ctx, database.CreateVoteParams{
		DebateCardID: sql.NullInt32{Int32: req.DebateCardID, Valid: true},
		UserID:       sql.NullInt32{Int32: userID, Valid: true},
		VoteType:     req.VoteType,
		Emoji:        sql.NullString{String: req.Emoji, Valid: req.Emoji != ""},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create vote: %v", err))
		return
	}

	// Update analytics
	c.updateDebateAnalytics(ctx, req.DebateCardID)

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Vote created successfully",
		"vote_id": vote.ID,
	})
}

func (c *Config) createComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.DebateID == 0 || req.Content == "" {
		respondWithError(w, http.StatusBadRequest, "debate_id and content are required")
		return
	}

	// Get user ID from context (you'll need to implement authentication)
	userID := int32(1) // TODO: Get from auth context

	// Create comment
	var parentCommentID sql.NullInt32
	if req.ParentCommentID != nil && *req.ParentCommentID > 0 {
		parentCommentID = sql.NullInt32{Int32: *req.ParentCommentID, Valid: true}
	} else {
		parentCommentID = sql.NullInt32{Valid: false}
	}

	comment, err := c.DB.CreateComment(ctx, database.CreateCommentParams{
		DebateID:        sql.NullInt32{Int32: req.DebateID, Valid: true},
		ParentCommentID: parentCommentID,
		UserID:          sql.NullInt32{Int32: userID, Valid: true},
		Content:         req.Content,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create comment: %v", err))
		return
	}

	// Update analytics
	c.updateDebateAnalytics(ctx, req.DebateID)

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message":    "Comment created successfully",
		"comment_id": comment.ID,
	})
}

func (c *Config) getComments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	debateIDStr := chi.URLParam(r, "debateId")
	debateID, err := strconv.ParseInt(debateIDStr, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid debate ID")
		return
	}

	comments, err := c.DB.GetComments(ctx, sql.NullInt32{Int32: int32(debateID), Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get comments: %v", err))
		return
	}

	// Convert to response format
	var response []CommentResponse
	for _, comment := range comments {
		commentResponse := CommentResponse{
			ID:            comment.ID,
			DebateID:      comment.DebateID.Int32,
			UserID:        comment.UserID.Int32,
			UserFirstName: comment.Firstname,
			UserLastName:  comment.Lastname,
			Content:       comment.Content,
			CreatedAt:     comment.CreatedAt.Time,
			UpdatedAt:     comment.UpdatedAt.Time,
		}

		if comment.ParentCommentID.Valid {
			commentResponse.ParentCommentID = &comment.ParentCommentID.Int32
		}

		response = append(response, commentResponse)
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (c *Config) generateAIPrompt(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if c.AIPromptGenerator == nil {
		respondWithError(w, http.StatusNotImplemented, "AI prompt generation is not configured. Please set the OpenAI API key.")
		return
	}

	matchID := r.URL.Query().Get("match_id")
	promptType := r.URL.Query().Get("type") // "pre_match" or "post_match"

	if matchID == "" || promptType == "" {
		respondWithError(w, http.StatusBadRequest, "match_id and type parameters are required")
		return
	}

	if promptType != "pre_match" && promptType != "post_match" {
		respondWithError(w, http.StatusBadRequest, "type must be 'pre_match' or 'post_match'")
		return
	}

	// Get basic match information first
	matchInfo, err := c.getMatchInfo(ctx, matchID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get match info: %v", err))
		return
	}

	// Validate match status for debate type
	if err := c.validateMatchStatusForDebateType(matchInfo.Status, promptType); err != nil {
		respondWithJSON(w, http.StatusOK, map[string]string{"info": err.Error()})
		return
	}

	// Use the data aggregator to get comprehensive match data
	aggregator := NewDebateDataAggregator(c)
	matchData, err := aggregator.AggregateMatchData(ctx, c.buildMatchDataRequest(matchID, matchInfo))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to aggregate match data: %v", err))
		return
	}

	var prompt *ai.DebatePrompt
	if promptType == "pre_match" {
		prompt, err = c.AIPromptGenerator.GeneratePreMatchPrompt(ctx, *matchData)
	} else {
		prompt, err = c.AIPromptGenerator.GeneratePostMatchPrompt(ctx, *matchData)
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate AI prompt: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, prompt)
}

// validateMatchStatusForDebateType checks if the match status is appropriate for the requested debate type
func (c *Config) validateMatchStatusForDebateType(matchStatus, debateType string) error {
	// Define match status categories
	notStartedStatuses := []string{"NS", "TBD", "POSTPONED", "CANCELLED", "SUSPENDED"}
	inProgressStatuses := []string{"1H", "2H", "HT", "ET", "P", "BT"}
	finishedStatuses := []string{"FT", "AET", "PEN", "FT_PEN", "AET_PEN"}

	// Check if status is in not started category
	for _, status := range notStartedStatuses {
		if matchStatus == status {
			if debateType == "post_match" {
				return fmt.Errorf("cannot generate post_match debate for a match that hasn't started (status: %s)", matchStatus)
			}
			return nil // pre_match is allowed for not started matches
		}
	}

	// Check if status is in progress
	for _, status := range inProgressStatuses {
		if matchStatus == status {
			if debateType == "post_match" {
				return fmt.Errorf("cannot generate post_match debate for a match that is still in progress (status: %s)", matchStatus)
			}
			return nil // pre_match is allowed for in-progress matches
		}
	}

	// Check if status is finished
	for _, status := range finishedStatuses {
		if matchStatus == status {
			if debateType == "pre_match" {
				return fmt.Errorf("cannot generate pre_match debate for a finished match (status: %s)", matchStatus)
			}
			return nil // post_match is allowed for finished matches
		}
	}

	// If status doesn't match any known category, be conservative
	if debateType == "post_match" {
		return fmt.Errorf("cannot generate post_match debate for match with unknown status: %s", matchStatus)
	}

	return nil
}

func (c *Config) generateDebate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if c.AIPromptGenerator == nil {
		respondWithError(w, http.StatusNotImplemented, "AI prompt generation is not configured. Please set the OpenAI API key.")
		return
	}

	var req GenerateDebateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.MatchID == "" || req.DebateType == "" {
		respondWithError(w, http.StatusBadRequest, "match_id and debate_type are required")
		return
	}

	if req.DebateType != "pre_match" && req.DebateType != "post_match" {
		respondWithError(w, http.StatusBadRequest, "debate_type must be 'pre_match' or 'post_match'")
		return
	}

	// Validate match_id format (should be numeric)
	if _, err := strconv.ParseInt(req.MatchID, 10, 64); err != nil {
		respondWithError(w, http.StatusBadRequest, "match_id must be a valid numeric ID")
		return
	}

	// Check if debate already exists for this match and type
	existingDebates, err := c.DB.GetDebatesByMatch(ctx, req.MatchID)
	if err == nil {
		for _, existing := range existingDebates {
			if existing.DebateType == req.DebateType {
				if !req.ForceRegenerate {
					// Return existing debate
					c.getDebateByID(w, r, existing.ID)
					return
				} else {
					// Soft delete existing debate to regenerate
					err := c.DB.SoftDeleteDebate(ctx, existing.ID)
					if err != nil {
						respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to soft delete existing debate: %v", err))
						return
					}
					fmt.Printf("Regenerating debate for match %s, type %s\n", req.MatchID, req.DebateType)
				}
			}
		}
	}

	// Get basic match information
	matchInfo, err := c.getMatchInfo(ctx, req.MatchID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get match info: %v", err))
		return
	}

	// Validate match status for debate type
	if err := c.validateMatchStatusForDebateType(matchInfo.Status, req.DebateType); err != nil {
		respondWithJSON(w, http.StatusOK, map[string]string{"info": err.Error()})
		return
	}

	// Use the data aggregator to get comprehensive match data
	aggregator := NewDebateDataAggregator(c)
	matchData, err := aggregator.AggregateMatchData(ctx, c.buildMatchDataRequest(req.MatchID, matchInfo))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to aggregate match data: %v", err))
		return
	}

	// Generate AI prompt
	var prompt *ai.DebatePrompt
	if req.DebateType == "pre_match" {
		prompt, err = c.AIPromptGenerator.GeneratePreMatchPrompt(ctx, *matchData)
	} else {
		prompt, err = c.AIPromptGenerator.GeneratePostMatchPrompt(ctx, *matchData)
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate AI prompt: %v", err))
		return
	}

	// Validate prompt structure
	if prompt.Headline == "" || len(prompt.Cards) == 0 {
		respondWithError(w, http.StatusInternalServerError, "Generated prompt is invalid (missing headline or cards)")
		return
	}

	// Create the debate in the database
	debate, err := c.DB.CreateDebate(ctx, database.CreateDebateParams{
		MatchID:     req.MatchID,
		DebateType:  req.DebateType,
		Headline:    prompt.Headline,
		Description: sql.NullString{String: prompt.Description, Valid: prompt.Description != ""},
		AiGenerated: sql.NullBool{Bool: true, Valid: true},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create debate: %v", err))
		return
	}

	// Create analytics record
	_, err = c.DB.CreateDebateAnalytics(ctx, database.CreateDebateAnalyticsParams{
		DebateID:        sql.NullInt32{Int32: debate.ID, Valid: true},
		TotalVotes:      sql.NullInt32{Int32: 0, Valid: true},
		TotalComments:   sql.NullInt32{Int32: 0, Valid: true},
		EngagementScore: sql.NullString{String: "0.0", Valid: true},
	})
	if err != nil {
		fmt.Printf("Failed to create debate analytics: %v\n", err)
	}

	// Create debate cards
	var cardResponses []DebateCardResponse
	for _, card := range prompt.Cards {
		// Validate card data
		if card.Stance == "" || card.Title == "" {
			fmt.Printf("Skipping invalid card: stance=%s, title=%s\n", card.Stance, card.Title)
			continue
		}

		// Validate stance
		if card.Stance != "agree" && card.Stance != "disagree" && card.Stance != "wildcard" {
			fmt.Printf("Skipping card with invalid stance: %s\n", card.Stance)
			continue
		}

		// Create the card in the database
		dbCard, err := c.DB.CreateDebateCard(ctx, database.CreateDebateCardParams{
			DebateID:    sql.NullInt32{Int32: debate.ID, Valid: true},
			Stance:      card.Stance,
			Title:       card.Title,
			Description: sql.NullString{String: card.Description, Valid: card.Description != ""},
			AiGenerated: sql.NullBool{Bool: true, Valid: true},
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create debate card: %v", err))
			return
		}

		// Add to response
		cardResponse := DebateCardResponse{
			ID:          dbCard.ID,
			DebateID:    dbCard.DebateID.Int32,
			Stance:      dbCard.Stance,
			Title:       dbCard.Title,
			Description: dbCard.Description.String,
			AIGenerated: dbCard.AiGenerated.Bool,
			CreatedAt:   dbCard.CreatedAt.Time,
			UpdatedAt:   dbCard.UpdatedAt.Time,
			VoteCounts: VoteCounts{
				Upvotes:   0,
				Downvotes: 0,
				Emojis:    make(map[string]int),
			},
		}
		cardResponses = append(cardResponses, cardResponse)
	}

	// Ensure we have at least one card
	if len(cardResponses) == 0 {
		respondWithError(w, http.StatusInternalServerError, "No valid debate cards were created")
		return
	}

	// Build the complete response
	response := DebateResponse{
		ID:          debate.ID,
		MatchID:     debate.MatchID,
		DebateType:  debate.DebateType,
		Headline:    debate.Headline,
		Description: debate.Description.String,
		AIGenerated: debate.AiGenerated.Bool,
		CreatedAt:   debate.CreatedAt.Time,
		UpdatedAt:   debate.UpdatedAt.Time,
		Cards:       cardResponses,
		Analytics: &DebateAnalyticsResponse{
			ID:              debate.ID,
			DebateID:        debate.ID,
			TotalVotes:      0,
			TotalComments:   0,
			EngagementScore: 0.0,
			CreatedAt:       debate.CreatedAt.Time,
			UpdatedAt:       debate.UpdatedAt.Time,
		},
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Debate generated successfully",
		"debate":  response,
	})
}

// Helper function to get debate by ID (extracted from getDebate for reuse)
func (c *Config) getDebateByID(w http.ResponseWriter, r *http.Request, debateID int32) {
	ctx := r.Context()

	// Get debate
	debate, err := c.DB.GetDebate(ctx, debateID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Debate not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get debate: %v", err))
		return
	}

	// Get debate cards
	cards, err := c.DB.GetDebateCards(ctx, sql.NullInt32{Int32: debate.ID, Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get debate cards: %v", err))
		return
	}

	// Get analytics
	analytics, err := c.DB.GetDebateAnalytics(ctx, sql.NullInt32{Int32: debate.ID, Valid: true})
	if err != nil && err != sql.ErrNoRows {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get debate analytics: %v", err))
		return
	}

	// Build response
	response := DebateResponse{
		ID:          debate.ID,
		MatchID:     debate.MatchID,
		DebateType:  debate.DebateType,
		Headline:    debate.Headline,
		Description: debate.Description.String,
		AIGenerated: debate.AiGenerated.Bool,
		CreatedAt:   debate.CreatedAt.Time,
		UpdatedAt:   debate.UpdatedAt.Time,
	}

	// Add cards with vote counts
	cardIDs := make([]int32, len(cards))
	for i, card := range cards {
		cardIDs[i] = card.ID
	}

	if len(cardIDs) > 0 {
		voteCounts, err := c.DB.GetVoteCounts(ctx, cardIDs)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get vote counts: %v", err))
			return
		}

		// Build vote counts map
		voteCountsMap := make(map[int32]VoteCounts)
		for _, vc := range voteCounts {
			if vc.DebateCardID.Valid {
				counts := voteCountsMap[vc.DebateCardID.Int32]
				switch vc.VoteType {
				case "upvote":
					counts.Upvotes = int(vc.Count)
				case "downvote":
					counts.Downvotes = int(vc.Count)
				case "emoji":
					if counts.Emojis == nil {
						counts.Emojis = make(map[string]int)
					}
					if vc.Emoji.Valid {
						counts.Emojis[vc.Emoji.String] = int(vc.Count)
					}
				}
				voteCountsMap[vc.DebateCardID.Int32] = counts
			}
		}

		// Build card responses
		for _, card := range cards {
			cardResponse := DebateCardResponse{
				ID:          card.ID,
				DebateID:    card.DebateID.Int32,
				Stance:      card.Stance,
				Title:       card.Title,
				Description: card.Description.String,
				AIGenerated: card.AiGenerated.Bool,
				CreatedAt:   card.CreatedAt.Time,
				UpdatedAt:   card.UpdatedAt.Time,
				VoteCounts:  voteCountsMap[card.ID],
			}
			response.Cards = append(response.Cards, cardResponse)
		}
	}

	// Add analytics if available
	if err == nil {
		engagementScore := 0.0
		if analytics.EngagementScore.Valid {
			// Parse engagement score from string
			if score, err := strconv.ParseFloat(analytics.EngagementScore.String, 64); err == nil {
				engagementScore = score
			}
		}

		response.Analytics = &DebateAnalyticsResponse{
			ID:              analytics.ID,
			DebateID:        analytics.DebateID.Int32,
			TotalVotes:      int(analytics.TotalVotes.Int32),
			TotalComments:   int(analytics.TotalComments.Int32),
			EngagementScore: engagementScore,
			CreatedAt:       analytics.CreatedAt.Time,
			UpdatedAt:       analytics.UpdatedAt.Time,
		}
	}

	respondWithJSON(w, http.StatusOK, response)
}

// getMatchInfo gets basic match information
func (c *Config) getMatchInfo(ctx context.Context, matchID string) (*MatchInfo, error) {
	// Use configurable base URL with fallback
	baseURL := c.APIFootballBaseURL
	if baseURL == "" {
		baseURL = "https://api-football-v1.p.rapidapi.com/v3"
	}

	url := fmt.Sprintf("%s/fixtures?id=%s", baseURL, matchID)
	headers := map[string]string{
		"Content-Type":   "application/json",
		"x-rapidapi-key": c.FootballAPIKey,
	}

	resp, err := HTTPRequest("GET", url, headers, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching match info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading match info response: %w", err)
	}

	var matchResponse struct {
		Response []struct {
			Fixture struct {
				Date   string `json:"date"`
				Status struct {
					Short string `json:"short"`
				} `json:"status"`
				Venue struct {
					Name string `json:"name"`
				} `json:"venue"`
			} `json:"fixture"`
			Teams struct {
				Home struct {
					Name string `json:"name"`
				} `json:"home"`
				Away struct {
					Name string `json:"name"`
				} `json:"away"`
			} `json:"teams"`
			Goals struct {
				Home int `json:"home"`
				Away int `json:"away"`
			} `json:"goals"`
			Score struct {
				Halftime struct {
					Home int `json:"home"`
					Away int `json:"away"`
				} `json:"halftime"`
				Fulltime struct {
					Home int `json:"home"`
					Away int `json:"away"`
				} `json:"fulltime"`
				Extratime struct {
					Home *int `json:"home"`
					Away *int `json:"away"`
				} `json:"extratime"`
				Penalty struct {
					Home *int `json:"home"`
					Away *int `json:"away"`
				} `json:"penalty"`
			} `json:"score"`
			League struct {
				Name   string `json:"name"`
				Season int    `json:"season"`
			} `json:"league"`
		} `json:"response"`
	}

	err = json.Unmarshal(body, &matchResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing match info response: %w", err)
	}

	if len(matchResponse.Response) == 0 {
		return nil, fmt.Errorf("no match found with ID %s", matchID)
	}

	match := matchResponse.Response[0]

	// Determine final score based on match status
	var homeScore, awayScore int
	switch match.Fixture.Status.Short {
	case "FT", "AET", "PEN":
		homeScore = match.Score.Fulltime.Home
		awayScore = match.Score.Fulltime.Away
	case "HT":
		homeScore = match.Score.Halftime.Home
		awayScore = match.Score.Halftime.Away
	default:
		homeScore = match.Goals.Home
		awayScore = match.Goals.Away
	}

	// Handle extra time and penalties
	if match.Score.Extratime.Home != nil && match.Score.Extratime.Away != nil {
		homeScore = *match.Score.Extratime.Home
		awayScore = *match.Score.Extratime.Away
	}
	if match.Score.Penalty.Home != nil && match.Score.Penalty.Away != nil {
		homeScore = *match.Score.Penalty.Home
		awayScore = *match.Score.Penalty.Away
	}

	return &MatchInfo{
		HomeTeam:        match.Teams.Home.Name,
		AwayTeam:        match.Teams.Away.Name,
		Date:            match.Fixture.Date,
		Status:          match.Fixture.Status.Short,
		HomeScore:       homeScore,
		AwayScore:       awayScore,
		HomeGoals:       match.Goals.Home,
		AwayGoals:       match.Goals.Away,
		HomeShots:       0, // Will be populated by fetchMatchStats if available
		AwayShots:       0,
		HomePossession:  0,
		AwayPossession:  0,
		HomeFouls:       0,
		AwayFouls:       0,
		HomeYellowCards: 0,
		AwayYellowCards: 0,
		HomeRedCards:    0,
		AwayRedCards:    0,
		Venue:           match.Fixture.Venue.Name,
		League:          match.League.Name,
		Season:          fmt.Sprintf("%d", match.League.Season),
	}, nil
}

func (c *Config) getTopDebates(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limitStr := r.URL.Query().Get("limit")
	limit := int32(10) // default limit

	if limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil {
			limit = int32(l)
		}
	}

	debates, err := c.DB.GetTopDebates(ctx, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get top debates: %v", err))
		return
	}

	// Convert to response format
	var response []DebateResponse
	for _, debate := range debates {
		debateResponse := DebateResponse{
			ID:          debate.ID,
			MatchID:     debate.MatchID,
			DebateType:  debate.DebateType,
			Headline:    debate.Headline,
			Description: debate.Description.String,
			AIGenerated: debate.AiGenerated.Bool,
			CreatedAt:   debate.CreatedAt.Time,
			UpdatedAt:   debate.UpdatedAt.Time,
		}

		if debate.TotalVotes.Valid {
			engagementScore := 0.0
			if debate.EngagementScore.Valid {
				// Parse engagement score from string
				if score, err := strconv.ParseFloat(debate.EngagementScore.String, 64); err == nil {
					engagementScore = score
				}
			}

			debateResponse.Analytics = &DebateAnalyticsResponse{
				ID:              debate.ID,
				DebateID:        debate.ID,
				TotalVotes:      int(debate.TotalVotes.Int32),
				TotalComments:   int(debate.TotalComments.Int32),
				EngagementScore: engagementScore,
				CreatedAt:       debate.CreatedAt.Time,
				UpdatedAt:       debate.UpdatedAt.Time,
			}
		}

		response = append(response, debateResponse)
	}

	respondWithJSON(w, http.StatusOK, response)
}

// Helper function to update debate analytics
func (c *Config) updateDebateAnalytics(ctx context.Context, debateCardID int32) {
	// Get the debate ID from the card
	card, err := c.DB.GetDebateCard(ctx, debateCardID)
	if err != nil {
		fmt.Printf("Failed to get debate card: %v\n", err)
		return
	}

	debateID := card.DebateID.Int32

	// Get vote counts for all cards in this debate
	cards, err := c.DB.GetDebateCards(ctx, sql.NullInt32{Int32: debateID, Valid: true})
	if err != nil {
		fmt.Printf("Failed to get debate cards: %v\n", err)
		return
	}

	cardIDs := make([]int32, len(cards))
	for i, card := range cards {
		cardIDs[i] = card.ID
	}

	voteCounts, err := c.DB.GetVoteCounts(ctx, cardIDs)
	if err != nil {
		fmt.Printf("Failed to get vote counts: %v\n", err)
		return
	}

	// Calculate total votes
	totalVotes := 0
	for _, vc := range voteCounts {
		totalVotes += int(vc.Count)
	}

	// Get comment count
	commentCount, err := c.DB.GetCommentCount(ctx, sql.NullInt32{Int32: debateID, Valid: true})
	if err != nil {
		fmt.Printf("Failed to get comment count: %v\n", err)
		return
	}

	// Calculate engagement score (votes + comments * 2 for comment weight)
	engagementScore := float64(totalVotes) + float64(commentCount)*2.0

	// Update analytics
	_, err = c.DB.UpdateDebateAnalytics(ctx, database.UpdateDebateAnalyticsParams{
		DebateID:        sql.NullInt32{Int32: debateID, Valid: true},
		TotalVotes:      sql.NullInt32{Int32: int32(totalVotes), Valid: true},
		TotalComments:   sql.NullInt32{Int32: int32(commentCount), Valid: true},
		EngagementScore: sql.NullString{String: fmt.Sprintf("%.2f", engagementScore), Valid: true},
	})
	if err != nil {
		fmt.Printf("Failed to update debate analytics: %v\n", err)
	}
}

// checkDebateGenerationHealth checks if all components needed for debate generation are working
func (c *Config) checkDebateGenerationHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	health := map[string]interface{}{
		"status": "healthy",
		"components": map[string]interface{}{
			"ai_prompt_generator": c.AIPromptGenerator != nil,
			"database":            c.DB != nil,
			"football_api":        c.FootballAPIKey != "",
			"cache":               c.Cache != nil,
		},
		"timestamp": time.Now().UTC(),
	}

	// Test database connection
	if c.DB != nil {
		_, err := c.DB.GetTopDebates(ctx, 1)
		if err != nil {
			health["status"] = "unhealthy"
			health["database_error"] = err.Error()
		}
	}

	// Test cache connection
	if c.Cache != nil {
		err := c.Cache.Set(ctx, "health_check", "test", time.Minute)
		if err != nil {
			health["status"] = "unhealthy"
			health["cache_error"] = err.Error()
		}
	}

	// Test football API
	if c.FootballAPIKey != "" {
		// Try to get a simple fixture to test API
		testMatchID := "1321727" // Use a known match ID
		_, err := c.getMatchInfo(ctx, testMatchID)
		if err != nil {
			health["status"] = "unhealthy"
			health["football_api_error"] = err.Error()
		}
	}

	statusCode := http.StatusOK
	if health["status"] == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	respondWithJSON(w, statusCode, health)
}

// hardDeleteDebate handles permanent deletion of a debate (admin only)
func (c *Config) hardDeleteDebate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract debate ID from URL
	debateIDStr := chi.URLParam(r, "id")
	debateID, err := strconv.Atoi(debateIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid debate ID")
		return
	}

	// TODO: Add admin authentication check here
	// For now, we'll allow the operation but log it

	// Hard delete the debate
	err = c.DB.DeleteDebate(ctx, int32(debateID))
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Debate not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete debate: %v", err))
		return
	}

	fmt.Printf("Hard deleted debate ID: %d\n", debateID)
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Debate permanently deleted"})
}

// restoreDebate handles restoring a soft-deleted debate
func (c *Config) restoreDebate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract debate ID from URL
	debateIDStr := chi.URLParam(r, "id")
	debateID, err := strconv.Atoi(debateIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid debate ID")
		return
	}

	// Restore the debate
	err = c.DB.RestoreDebate(ctx, int32(debateID))
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Debate not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to restore debate: %v", err))
		return
	}

	fmt.Printf("Restored debate ID: %d\n", debateID)
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Debate restored successfully"})
}

// MatchInfo represents detailed information about a match
type MatchInfo struct {
	HomeTeam        string
	AwayTeam        string
	Date            string
	Status          string
	HomeScore       int
	AwayScore       int
	HomeGoals       int
	AwayGoals       int
	HomeShots       int
	AwayShots       int
	HomePossession  int
	AwayPossession  int
	HomeFouls       int
	AwayFouls       int
	HomeYellowCards int
	AwayYellowCards int
	HomeRedCards    int
	AwayRedCards    int
	Venue           string
	League          string
	Season          string
}

// buildMatchDataRequest converts MatchInfo to MatchDataRequest
func (c *Config) buildMatchDataRequest(matchID string, matchInfo *MatchInfo) MatchDataRequest {
	return MatchDataRequest{
		MatchID:         matchID,
		HomeTeam:        matchInfo.HomeTeam,
		AwayTeam:        matchInfo.AwayTeam,
		Date:            matchInfo.Date,
		Status:          matchInfo.Status,
		HomeScore:       matchInfo.HomeScore,
		AwayScore:       matchInfo.AwayScore,
		HomeGoals:       matchInfo.HomeGoals,
		AwayGoals:       matchInfo.AwayGoals,
		HomeShots:       matchInfo.HomeShots,
		AwayShots:       matchInfo.AwayShots,
		HomePossession:  matchInfo.HomePossession,
		AwayPossession:  matchInfo.AwayPossession,
		HomeFouls:       matchInfo.HomeFouls,
		AwayFouls:       matchInfo.AwayFouls,
		HomeYellowCards: matchInfo.HomeYellowCards,
		AwayYellowCards: matchInfo.AwayYellowCards,
		HomeRedCards:    matchInfo.HomeRedCards,
		AwayRedCards:    matchInfo.AwayRedCards,
		Venue:           matchInfo.Venue,
		League:          matchInfo.League,
		Season:          matchInfo.Season,
	}
}
