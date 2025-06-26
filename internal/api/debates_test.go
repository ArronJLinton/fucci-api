package api

import (
	"testing"

	"github.com/ArronJLinton/fucci-api/internal/cache"
)

func TestDebateDataAggregator(t *testing.T) {
	// Skip if no Redis connection
	redisURL := "redis://localhost:6379"
	cache, err := cache.NewCache(redisURL)
	if err != nil {
		t.Skip("Skipping test: Redis not available")
	}

	config := &Config{
		Cache:          cache,
		FootballAPIKey: "mock-api-key",
		RapidAPIKey:    "mock-rapid-api-key",
	}

	aggregator := NewDebateDataAggregator(config)

	t.Run("test aggregator creation", func(t *testing.T) {
		if aggregator == nil {
			t.Error("DebateDataAggregator should not be nil")
		}
		if aggregator.Config != config {
			t.Error("Config should be set correctly")
		}
	})
}

func TestCreateDebateRequest(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		req := CreateDebateRequest{
			MatchID:     "12345",
			DebateType:  "pre_match",
			Headline:    "Test Debate",
			Description: "Test Description",
			AIGenerated: true,
		}

		if req.MatchID != "12345" {
			t.Errorf("Expected MatchID to be '12345', got %s", req.MatchID)
		}
		if req.DebateType != "pre_match" {
			t.Errorf("Expected DebateType to be 'pre_match', got %s", req.DebateType)
		}
		if req.Headline != "Test Debate" {
			t.Errorf("Expected Headline to be 'Test Debate', got %s", req.Headline)
		}
		if !req.AIGenerated {
			t.Error("Expected AIGenerated to be true")
		}
	})
}

func TestCreateVoteRequest(t *testing.T) {
	t.Run("valid vote request", func(t *testing.T) {
		req := CreateVoteRequest{
			DebateCardID: 1,
			VoteType:     "upvote",
			Emoji:        "",
		}

		if req.DebateCardID != 1 {
			t.Errorf("Expected DebateCardID to be 1, got %d", req.DebateCardID)
		}
		if req.VoteType != "upvote" {
			t.Errorf("Expected VoteType to be 'upvote', got %s", req.VoteType)
		}
	})

	t.Run("emoji vote request", func(t *testing.T) {
		req := CreateVoteRequest{
			DebateCardID: 1,
			VoteType:     "emoji",
			Emoji:        "👍",
		}

		if req.Emoji != "👍" {
			t.Errorf("Expected Emoji to be '👍', got %s", req.Emoji)
		}
	})
}

func TestCreateCommentRequest(t *testing.T) {
	t.Run("valid comment request", func(t *testing.T) {
		req := CreateCommentRequest{
			DebateID: 1,
			Content:  "This is a test comment",
		}

		if req.DebateID != 1 {
			t.Errorf("Expected DebateID to be 1, got %d", req.DebateID)
		}
		if req.Content != "This is a test comment" {
			t.Errorf("Expected Content to be 'This is a test comment', got %s", req.Content)
		}
		if req.ParentCommentID != nil {
			t.Error("Expected ParentCommentID to be nil")
		}
	})

	t.Run("nested comment request", func(t *testing.T) {
		parentID := int32(5)
		req := CreateCommentRequest{
			DebateID:        1,
			ParentCommentID: &parentID,
			Content:         "This is a reply",
		}

		if req.ParentCommentID == nil {
			t.Error("Expected ParentCommentID to not be nil")
		}
		if *req.ParentCommentID != 5 {
			t.Errorf("Expected ParentCommentID to be 5, got %d", *req.ParentCommentID)
		}
	})
}

func TestDebateResponse(t *testing.T) {
	t.Run("debate response structure", func(t *testing.T) {
		response := DebateResponse{
			ID:          1,
			MatchID:     "12345",
			DebateType:  "pre_match",
			Headline:    "Test Debate",
			Description: "Test Description",
			AIGenerated: true,
		}

		if response.ID != 1 {
			t.Errorf("Expected ID to be 1, got %d", response.ID)
		}
		if response.MatchID != "12345" {
			t.Errorf("Expected MatchID to be '12345', got %s", response.MatchID)
		}
		if response.DebateType != "pre_match" {
			t.Errorf("Expected DebateType to be 'pre_match', got %s", response.DebateType)
		}
		if !response.AIGenerated {
			t.Error("Expected AIGenerated to be true")
		}
	})
}

func TestVoteCounts(t *testing.T) {
	t.Run("vote counts structure", func(t *testing.T) {
		counts := VoteCounts{
			Upvotes:   10,
			Downvotes: 5,
			Emojis: map[string]int{
				"👍": 15,
				"👎": 3,
			},
		}

		if counts.Upvotes != 10 {
			t.Errorf("Expected Upvotes to be 10, got %d", counts.Upvotes)
		}
		if counts.Downvotes != 5 {
			t.Errorf("Expected Downvotes to be 5, got %d", counts.Downvotes)
		}
		if counts.Emojis["👍"] != 15 {
			t.Errorf("Expected 👍 emoji count to be 15, got %d", counts.Emojis["👍"])
		}
		if counts.Emojis["👎"] != 3 {
			t.Errorf("Expected 👎 emoji count to be 3, got %d", counts.Emojis["👎"])
		}
	})
}

func TestDebateCardResponse(t *testing.T) {
	t.Run("debate card response structure", func(t *testing.T) {
		card := DebateCardResponse{
			ID:          1,
			DebateID:    1,
			Stance:      "agree",
			Title:       "Agree with the decision",
			Description: "This was the right call",
			AIGenerated: true,
			VoteCounts: VoteCounts{
				Upvotes:   25,
				Downvotes: 5,
			},
		}

		if card.Stance != "agree" {
			t.Errorf("Expected Stance to be 'agree', got %s", card.Stance)
		}
		if card.Title != "Agree with the decision" {
			t.Errorf("Expected Title to be 'Agree with the decision', got %s", card.Title)
		}
		if card.VoteCounts.Upvotes != 25 {
			t.Errorf("Expected Upvotes to be 25, got %d", card.VoteCounts.Upvotes)
		}
		if !card.AIGenerated {
			t.Error("Expected AIGenerated to be true")
		}
	})
}

func TestDebateAnalyticsResponse(t *testing.T) {
	t.Run("analytics response structure", func(t *testing.T) {
		analytics := DebateAnalyticsResponse{
			ID:              1,
			DebateID:        1,
			TotalVotes:      100,
			TotalComments:   25,
			EngagementScore: 150.5,
		}

		if analytics.TotalVotes != 100 {
			t.Errorf("Expected TotalVotes to be 100, got %d", analytics.TotalVotes)
		}
		if analytics.TotalComments != 25 {
			t.Errorf("Expected TotalComments to be 25, got %d", analytics.TotalComments)
		}
		if analytics.EngagementScore != 150.5 {
			t.Errorf("Expected EngagementScore to be 150.5, got %f", analytics.EngagementScore)
		}
	})
}
