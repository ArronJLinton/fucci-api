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

func TestGetLeagueStandings(t *testing.T) {
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
			"get": "standings",
			"response": [
				{
					"league": {
						"id": 39,
						"name": "Premier League",
						"country": "England",
						"logo": "https://media.api-sports.io/football/leagues/39.png",
						"flag": "https://media.api-sports.io/flags/gb.svg",
						"season": 2024,
						"standings": [
							[
								{
									"rank": 1,
									"team": {
										"id": 40,
										"name": "Liverpool",
										"logo": "https://media.api-sports.io/football/teams/40.png"
									},
									"points": 84,
									"goalsDiff": 45,
									"group": "Premier League",
									"form": "DLDLW",
									"status": "same",
									"description": "Champions League",
									"all": {
										"played": 38,
										"win": 25,
										"draw": 9,
										"lose": 4,
										"goals": {
											"for": 86,
											"against": 41
										}
									},
									"home": {
										"played": 19,
										"win": 14,
										"draw": 4,
										"lose": 1,
										"goals": {
											"for": 42,
											"against": 16
										}
									},
									"away": {
										"played": 19,
										"win": 11,
										"draw": 5,
										"lose": 3,
										"goals": {
											"for": 44,
											"against": 25
										}
									},
									"update": "2025-05-26T00:00:00Z"
								}
							]
						]
					}
				}
			],
			"errors": [],
			"results": 1,
			"paging": {
				"current": 1,
				"total": 1
			}
		}`

		// Create a test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("x-rapidapi-key") != mockAPIKey {
				t.Errorf("Expected API key %s, got %s", mockAPIKey, r.Header.Get("x-rapidapi-key"))
			}

			// Verify query parameters
			leagueID := r.URL.Query().Get("league")
			season := r.URL.Query().Get("season")
			if leagueID != "39" {
				t.Errorf("Expected league ID 39, got %s", leagueID)
			}
			if season != "2024" {
				t.Errorf("Expected season 2024, got %s", season)
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(mockResponse))
		}))
		defer server.Close()

		// Set the base URL to the test server
		config.APIFootballBaseURL = server.URL

		// Setup request
		req := httptest.NewRequest("GET", "/fixtures/league_standings", nil)
		req.URL.RawQuery = "league_id=39&season=2024"
		rec := httptest.NewRecorder()

		// Call the function
		config.getLeagueStandingsByLeagueId(rec, req)

		// Assert the response
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		var response GetLeagueStandingsResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to parse response: %s", err)
		}

		// Verify the response structure
		if len(response.Response) != 1 {
			t.Errorf("Expected 1 league response, got %d", len(response.Response))
		}
		if len(response.Response[0].League.Standings) != 1 {
			t.Errorf("Expected 1 standings array, got %d", len(response.Response[0].League.Standings))
		}
		if len(response.Response[0].League.Standings[0]) != 1 {
			t.Errorf("Expected 1 team in standings, got %d", len(response.Response[0].League.Standings[0]))
		}
	})

	t.Run("error case - missing league_id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/fixtures/league_standings", nil)
		rec := httptest.NewRecorder()

		// Call the function
		config.getLeagueStandingsByLeagueId(rec, req)

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
		if response.Error != "league_id is required" {
			t.Errorf("Expected error message 'league_id is required', got '%s'", response.Error)
		}
	})

	t.Run("error case - missing season", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/fixtures/league_standings", nil)
		req.URL.RawQuery = "league_id=39"
		rec := httptest.NewRecorder()

		// Call the function
		config.getLeagueStandingsByLeagueId(rec, req)

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
		if response.Error != "season is required" {
			t.Errorf("Expected error message 'season is required', got '%s'", response.Error)
		}
	})
}

func TestCacheExpiration(t *testing.T) {
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

	t.Run("Cache expiration", func(t *testing.T) {
		// Mock HTTP response from the external API
		mockResponse := `{
			"get": "standings",
			"response": [
				{
					"league": {
						"id": 39,
						"name": "Premier League",
						"country": "England",
						"logo": "https://media.api-sports.io/football/leagues/39.png",
						"flag": "https://media.api-sports.io/flags/gb.svg",
						"season": 2024,
						"standings": [
							[
								{
									"rank": 1,
									"team": {
										"id": 40,
										"name": "Liverpool",
										"logo": "https://media.api-sports.io/football/teams/40.png"
									},
									"points": 84,
									"goalsDiff": 45,
									"group": "Premier League",
									"form": "DLDLW",
									"status": "same",
									"description": "Champions League",
									"all": {
										"played": 38,
										"win": 25,
										"draw": 9,
										"lose": 4,
										"goals": {
											"for": 86,
											"against": 41
										}
									},
									"home": {
										"played": 19,
										"win": 14,
										"draw": 4,
										"lose": 1,
										"goals": {
											"for": 42,
											"against": 16
										}
									},
									"away": {
										"played": 19,
										"win": 11,
										"draw": 5,
										"lose": 3,
										"goals": {
											"for": 44,
											"against": 25
										}
									},
									"update": "2025-05-26T00:00:00Z"
								}
							]
						]
					}
				}
			],
			"errors": [],
			"results": 1,
			"paging": {
				"current": 1,
				"total": 1
			}
		}`

		// Create a test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("x-rapidapi-key") != mockAPIKey {
				t.Errorf("Expected API key %s, got %s", mockAPIKey, r.Header.Get("x-rapidapi-key"))
			}

			// Verify query parameters
			leagueID := r.URL.Query().Get("league")
			season := r.URL.Query().Get("season")
			if leagueID != "39" {
				t.Errorf("Expected league ID 39, got %s", leagueID)
			}
			if season != "2024" {
				t.Errorf("Expected season 2024, got %s", season)
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(mockResponse))
		}))
		defer server.Close()

		// Set the base URL to the test server
		config.APIFootballBaseURL = server.URL

		// First request - should hit the API
		req1 := httptest.NewRequest("GET", "/fixtures/league_standings", nil)
		req1.URL.RawQuery = "league_id=39&season=2024"
		rec1 := httptest.NewRecorder()
		config.getLeagueStandingsByLeagueId(rec1, req1)
		if rec1.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec1.Code)
		}

		// Second request - should hit the cache
		req2 := httptest.NewRequest("GET", "/fixtures/league_standings", nil)
		req2.URL.RawQuery = "league_id=39&season=2024"
		rec2 := httptest.NewRecorder()
		config.getLeagueStandingsByLeagueId(rec2, req2)
		if rec2.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec2.Code)
		}

		// Verify that the responses are identical
		if rec1.Body.String() != rec2.Body.String() {
			t.Error("cached response differs from original response")
		}
	})
}
