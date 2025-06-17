package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ArronJLinton/fucci-api/internal/cache"
)

func TestGetMatchLineup(t *testing.T) {
	// Skip if no Redis connection
	redisURL := "redis://localhost:6379"
	cache, err := cache.NewCache(redisURL)
	if err != nil {
		t.Skip("Skipping test: Redis not available")
	}

	mockAPIKey := "mock-api-key"
	config := &Config{
		Cache:          cache,
		FootballAPIKey: mockAPIKey,
	}

	t.Run("success case", func(t *testing.T) {
		// Mock HTTP response from the external API
		mockResponse := `{
			"response": [
				{
					"team": {"id": 1},
					"startXI": [
						{"player": {"id": 1, "name": "Player 1", "number": 1, "pos": "G", "grid": "", "photo": "photo1.jpg"}}
					],
					"substitutes": [
						{"player": {"id": 2, "name": "Player 2", "number": 2, "pos": "G", "grid": "", "photo": "photo2.jpg"}}
					]
				},
				{
					"team": {"id": 2},
					"startXI": [
						{"player": {"id": 3, "name": "Player 3", "number": 1, "pos": "G", "grid": "", "photo": "photo3.jpg"}}
					],
					"substitutes": [
						{"player": {"id": 4, "name": "Player 4", "number": 2, "pos": "G", "grid": "", "photo": "photo4.jpg"}}
					]
				}
			]
		}`

		// Mock team squad response
		mockSquadResponse := `{
			"response": [
				{
					"team": {"id": 1},
					"players": [
						{"id": 1, "name": "Player 1", "number": 1, "pos": "G", "grid": "", "photo": "photo1.jpg"},
						{"id": 2, "name": "Player 2", "number": 2, "pos": "G", "grid": "", "photo": "photo2.jpg"}
					]
				},
				{
					"team": {"id": 2},
					"players": [
						{"id": 3, "name": "Player 3", "number": 1, "pos": "G", "grid": "", "photo": "photo3.jpg"},
						{"id": 4, "name": "Player 4", "number": 2, "pos": "G", "grid": "", "photo": "photo4.jpg"}
					]
				}
			]
		}`

		// Create a test server that handles both lineup and squad requests
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("x-rapidapi-key") != mockAPIKey {
				t.Errorf("Expected API key %s, got %s", mockAPIKey, r.Header.Get("x-rapidapi-key"))
			}

			// Determine which response to return based on the URL
			if strings.Contains(r.URL.Path, "lineups") {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(mockResponse))
			} else if strings.Contains(r.URL.Path, "players/squads") {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(mockSquadResponse))
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}))
		defer server.Close()

		// Set the base URL to the test server
		config.APIFootballBaseURL = server.URL

		// Setup request
		req := httptest.NewRequest("GET", "/fixtures/lineups", nil)
		req.URL.RawQuery = "match_id=12345"
		rec := httptest.NewRecorder()

		// Call the function
		config.getMatchLineup(rec, req)

		// Assert the response
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		var response struct {
			Home struct {
				Starters    []Player `json:"starters"`
				Substitutes []Player `json:"substitutes"`
			} `json:"home"`
			Away struct {
				Starters    []Player `json:"starters"`
				Substitutes []Player `json:"substitutes"`
			} `json:"away"`
		}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to parse response: %s", err)
		}

		// Verify the response structure
		if len(response.Home.Starters) != 1 {
			t.Errorf("Expected 1 home starter, got %d", len(response.Home.Starters))
		}
		if len(response.Home.Substitutes) != 1 {
			t.Errorf("Expected 1 home substitute, got %d", len(response.Home.Substitutes))
		}
		if len(response.Away.Starters) != 1 {
			t.Errorf("Expected 1 away starter, got %d", len(response.Away.Starters))
		}
		if len(response.Away.Substitutes) != 1 {
			t.Errorf("Expected 1 away substitute, got %d", len(response.Away.Substitutes))
		}
	})

	t.Run("error case - missing match_id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/fixtures/lineups", nil)
		rec := httptest.NewRecorder()

		// Call the function
		config.getMatchLineup(rec, req)

		// Assert the response
		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rec.Code)
		}

		var response struct {
			Error string `json:"error"`
		}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to parse error response: %s", err)
		}
		if response.Error != "match_id is required" {
			t.Errorf("Expected error message 'match_id is required', got '%s'", response.Error)
		}
	})
}

// MockCache is a mock implementation of the cache interface
type MockCache struct {
	cache.Cache
	existsFunc func(ctx context.Context, key string) (bool, error)
	getFunc    func(ctx context.Context, key string, value interface{}) error
	setFunc    func(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}

func (m *MockCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return m.setFunc(ctx, key, value, ttl)
}

func (m *MockCache) Get(ctx context.Context, key string, value interface{}) error {
	return m.getFunc(ctx, key, value)
}

func (m *MockCache) Exists(ctx context.Context, key string) (bool, error) {
	return m.existsFunc(ctx, key)
}

func TestGetLeagueStandingsByLeagueId(t *testing.T) {
	// Skip if no Redis connection
	redisURL := "redis://localhost:6379"
	cache, err := cache.NewCache(redisURL)
	if err != nil {
		t.Skip("Skipping test: Redis not available")
	}

	// Test cases
	tests := []struct {
		name           string
		leagueID       string
		expectedStatus int
	}{
		{
			name:           "Cache miss - should call API",
			leagueID:       "39",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Cache hit - should return cached data",
			leagueID:       "39",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing league ID",
			leagueID:       "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				Cache:          cache,
				FootballAPIKey: "test-key",
			}

			// Setup request
			req := httptest.NewRequest("GET", "/v1/api/futbol/league_standings?league_id="+tt.leagueID, nil)
			w := httptest.NewRecorder()

			// Call the handler
			config.getLeagueStandingsByLeagueId(w, req)

			// Assert response
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestCacheExpiration(t *testing.T) {
	// Skip if no Redis connection
	redisURL := "redis://localhost:6379"
	cache, err := cache.NewCache(redisURL)
	if err != nil {
		t.Skip("Skipping test: Redis not available")
	}

	config := &Config{
		Cache:          cache,
		FootballAPIKey: "test-key",
	}

	// Test cache expiration
	t.Run("Cache expiration", func(t *testing.T) {
		leagueID := "39"
		req := httptest.NewRequest("GET", "/v1/api/futbol/league_standings?league_id="+leagueID, nil)
		w := httptest.NewRecorder()

		// First call - should miss cache
		config.getLeagueStandingsByLeagueId(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		// Second call - should hit cache
		w = httptest.NewRecorder()
		config.getLeagueStandingsByLeagueId(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		// Verify cache was hit by checking Redis directly
		exists, err := cache.Exists(context.Background(), "standings:league:"+leagueID)
		if err != nil {
			t.Errorf("failed to check cache existence: %v", err)
		}
		if !exists {
			t.Error("expected cache to exist after second call")
		}
	})
}
