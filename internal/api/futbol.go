package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/ArronJLinton/fucci-api/internal/cache"
)

type GetMatchesParams struct {
	Date string
}

func (c *Config) getMatches(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Step 1: Get date from query parameters
	queryParams := r.URL.Query()
	date := queryParams.Get("date")
	if date == "" {
		respondWithError(w, http.StatusBadRequest, "date parameter is required")
		return
	}

	// Generate cache key
	cacheKey := fmt.Sprintf("matches:%s", date)

	// Try to get from cache first
	var data GetMatchesAPIResponse
	exists, err := c.Cache.Exists(ctx, cacheKey)
	if err != nil {
		log.Printf("Cache check error: %v\n", err)
	} else if exists {
		err = c.Cache.Get(ctx, cacheKey, &data)
		if err == nil {
			respondWithJSON(w, http.StatusOK, data)
			return
		}
		log.Printf("Cache get error: %v\n", err)
	}

	// If not in cache or error occurred, fetch from API
	url := fmt.Sprintf("https://api-football-v1.p.rapidapi.com/v3/fixtures?date=%s", date)
	headers := map[string]string{
		"Content-Type":   "application/json",
		"x-rapidapi-key": c.FootballAPIKey,
	}

	resp, err := HTTPRequest("GET", url, headers, nil)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error creating http request: %s", err))
		return
	}
	defer resp.Body.Close()

	// Read the raw response body for debugging
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to read response body: %s", err))
		return
	}

	// Create a new reader from the raw body for JSON decoding
	err = json.Unmarshal(rawBody, &data)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to parse response from football api service: %s", err))
		return
	}

	// Determine cache TTL based on match statuses
	ttl := cache.DefaultTTL
	if len(data.Response) > 0 {
		// Check if any matches are live
		for _, match := range data.Response {
			status := match.Fixture.Status.Short
			matchTTL := cache.GetMatchTTL(status)
			// Use the shortest TTL found (most conservative)
			if matchTTL < ttl {
				ttl = matchTTL
			}
		}
	}

	// Store in cache with determined TTL
	err = c.Cache.Set(ctx, cacheKey, data, ttl)
	if err != nil {
		log.Printf("Cache set error: %v\n", err)
	}

	respondWithJSON(w, http.StatusOK, data)
}

func (c *Config) getMatch(w http.ResponseWriter, r *http.Request) {
}

func (c *Config) getMatchLineup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queryParams := r.URL.Query()
	matchID := queryParams.Get("match_id")
	if matchID == "" {
		respondWithError(w, http.StatusBadRequest, "match_id is required")
		return
	}

	// Generate cache key
	cacheKey := fmt.Sprintf("lineup:%s", matchID)

	// Try to get from cache first
	var response struct {
		Home Lineup `json:"home"`
		Away Lineup `json:"away"`
	}
	exists, err := c.Cache.Exists(ctx, cacheKey)
	if err != nil {
		fmt.Printf("Cache check error: %v\n", err)
	} else if exists {
		err = c.Cache.Get(ctx, cacheKey, &response)
		if err == nil {
			fmt.Printf("Cache hit for lineup: %s\n", matchID)
			respondWithJSON(w, http.StatusOK, response)
			return
		}
		fmt.Printf("Cache get error: %v\n", err)
	}

	// Use configurable base URL with fallback
	baseURL := c.APIFootballBaseURL
	if baseURL == "" {
		baseURL = "https://api-football-v1.p.rapidapi.com/v3"
	}
	url := fmt.Sprintf("%s/fixtures/lineups?fixture=%s", baseURL, matchID)
	headers := map[string]string{
		"Content-Type":   "application/json",
		"x-rapidapi-key": c.FootballAPIKey,
	}
	resp, err := HTTPRequest("GET", url, headers, nil)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error creating http request: %s", err))
		return
	}
	defer resp.Body.Close()

	// Read the raw response for debugging
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to read response body: %s", err))
		return
	}

	// Create a new reader from the raw body for JSON decoding
	getLineUpData := &GetLineUpResponse{}
	err = json.Unmarshal(rawBody, &getLineUpData)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to parse response from football api service - %s", err))
		return
	}

	fmt.Printf("Number of lineup responses: %d\n", len(getLineUpData.Response))

	if len(getLineUpData.Response) < 2 {
		respondWithJSON(w, http.StatusOK, "No lineup data available")
		return
	}

	// Use the same base URL for squad requests
	homeTeamSquad, err := c.getTeamSquad(int32(getLineUpData.Response[0].Team.ID))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to get team squad: %s", err))
		return
	}
	awayTeamSquad, err := c.getTeamSquad(int32(getLineUpData.Response[1].Team.ID))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to get team squad: %s", err))
		return
	}

	// Process lineups and create response
	response = struct {
		Home Lineup `json:"home"`
		Away Lineup `json:"away"`
	}{
		Home: Lineup{
			Starters:    processPlayers(getLineUpData.Response[0].StartXI, homeTeamSquad),
			Substitutes: processSubstitutes(getLineUpData.Response[0].Substitutes, homeTeamSquad),
		},
		Away: Lineup{
			Starters:    processPlayers(getLineUpData.Response[1].StartXI, awayTeamSquad),
			Substitutes: processSubstitutes(getLineUpData.Response[1].Substitutes, awayTeamSquad),
		},
	}

	// Store in cache
	err = c.Cache.Set(ctx, cacheKey, response, cache.LineupTTL)
	if err != nil {
		fmt.Printf("Cache set error: %v\n", err)
	} else {
		fmt.Printf("Stored lineup in cache: %s with TTL: %v\n", cacheKey, cache.LineupTTL)
	}

	respondWithJSON(w, http.StatusOK, response)
}

// Helper functions to process players
func processPlayers(players []struct {
	Player Player `json:"player"`
}, squad *GetSquadResponse) []Player {
	result := make([]Player, 0, len(players))
	for _, p := range players {
		squadPlayer := filterByName(squad.Response[0].Players, p.Player)
		p := Player{
			ID:     p.Player.ID,
			Name:   p.Player.Name,
			Number: p.Player.Number,
			Pos:    p.Player.Pos,
			Grid:   p.Player.Grid,
			Photo:  squadPlayer.Photo,
		}
		result = append(result, p)
	}
	return result
}

func processSubstitutes(substitutes []struct {
	Player struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Number int    `json:"number"`
		Pos    string `json:"pos"`
		Grid   any    `json:"grid"`
		Photo  string `json:"photo"`
	} `json:"player"`
}, squad *GetSquadResponse) []Player {
	result := make([]Player, 0, len(substitutes))
	for _, p := range substitutes {
		squadPlayer := filterByName(squad.Response[0].Players, Player{
			ID:     p.Player.ID,
			Name:   p.Player.Name,
			Number: p.Player.Number,
		})

		p := Player{
			ID:     p.Player.ID,
			Name:   p.Player.Name,
			Number: p.Player.Number,
			Pos:    p.Player.Pos,
			Grid:   "",
			Photo:  squadPlayer.Photo,
		}
		result = append(result, p)
	}
	return result
}

func filterByName(items []Player, player Player) Player {
	// First try to match by ID
	if player.ID != 0 {
		for _, item := range items {
			if item.ID == player.ID {
				return item
			}
		}
	}

	// If ID match fails, try name matching with various normalizations
	normalizedSearchName := normalizeName(player.Name)
	var bestMatch Player
	var maxSimilarity float32 = 0

	for _, item := range items {
		normalizedItemName := normalizeName(item.Name)

		// Try exact match first
		if normalizedItemName == normalizedSearchName {
			return item
		}

		// Check if names contain each other
		if strings.Contains(normalizedItemName, normalizedSearchName) ||
			strings.Contains(normalizedSearchName, normalizedItemName) {
			similarity := float32(len(normalizedItemName)) / float32(len(normalizedSearchName))
			if similarity > maxSimilarity {
				maxSimilarity = similarity
				bestMatch = item
			}
		}

		// If jersey numbers match and names are similar enough, consider it a match
		if player.Number != 0 && item.Number == player.Number {
			return item
		}
	}

	// If we found a good match, return it
	if maxSimilarity > 0.7 {
		return bestMatch
	}

	// Return empty player if no match found
	return Player{}
}

func normalizeName(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)

	// Remove dots and extra spaces
	name = strings.ReplaceAll(name, ".", "")
	name = strings.Join(strings.Fields(name), " ")

	// Remove special characters
	name = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, name)

	return name
}

func (c *Config) getTeamSquad(id int32) (*GetSquadResponse, error) {
	headers := map[string]string{
		"Content-Type":   "application/json",
		"x-rapidapi-key": c.FootballAPIKey,
	}

	// Use configurable base URL with fallback
	baseURL := c.APIFootballBaseURL
	if baseURL == "" {
		baseURL = "https://api-football-v1.p.rapidapi.com/v3"
	}
	url := fmt.Sprintf("%s/players/squads?team=%d", baseURL, id)

	response, err := handleClientRequest[GetSquadResponse](url, "GET", headers)
	if err != nil {
		return nil, fmt.Errorf("error creating http request: %s", err)
	}

	// Log only if no squad data found
	if len(response.Response) == 0 {
		fmt.Printf("No squad data received for team ID: %d\n", id)
	}

	return response, nil
}

func (c *Config) getLeagues(w http.ResponseWriter, r *http.Request) {

	url := "https://api-football-v1.p.rapidapi.com/v3/leagues?season=2025"
	headers := map[string]string{
		"Content-Type":   "application/json",
		"x-rapidapi-key": c.FootballAPIKey,
	}
	resp, err := HTTPRequest("GET", url, headers, nil)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error creating http request: %s", err))
		return
	}
	defer resp.Body.Close()

	responseBody := json.NewDecoder(resp.Body)
	data := GetLeaguesResponse{}
	err = responseBody.Decode(&data)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to read response from football api service: %s", err))
		return
	}
	type LeagueInfo struct {
		Name    string `json:"name"`
		Country string `json:"country"`
		Logo    string `json:"logo"`
	}
	leagueNames := []LeagueInfo{}
	for _, l := range data.Response {
		obj := LeagueInfo{
			Name:    l.League.Name,
			Country: l.Country.Name,
			Logo:    l.League.Logo,
		}
		leagueNames = append(leagueNames, obj)
	}
	respondWithJSON(w, http.StatusOK, leagueNames)
}

func (c *Config) getLeagueStandingsByTeamId(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	teamId := queryParams.Get("team_id")
	// TODO: Dynamically set the season year
	currentYear := time.Now().Year()
	url := fmt.Sprintf("https://api-football-v1.p.rapidapi.com/v3/standings?season=%d&team=%s", currentYear, teamId)
	headers := map[string]string{
		"Content-Type":   "application/json",
		"x-rapidapi-key": c.FootballAPIKey,
	}
	resp, err := HTTPRequest("GET", url, headers, nil)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error creating http request: %s", err))
		return
	}
	defer resp.Body.Close()

	responseBody := json.NewDecoder(resp.Body)
	data := GetLeagueStandingsByTeamIdResponse{}
	err = responseBody.Decode(&data)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing response body: %s", err))
		return
	}
	respondWithJSON(w, http.StatusOK, data)
}

func (c *Config) getLeagueStandingsByLeagueId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queryParams := r.URL.Query()
	leagueId := queryParams.Get("league_id")
	if leagueId == "" {
		respondWithError(w, http.StatusBadRequest, "league_id is required")
		return
	}

	// Generate cache key
	cacheKey := fmt.Sprintf("standings:league:%s", leagueId)

	// Try to get from cache first
	var data GetLeagueStandingsByLeagueIdResponse
	exists, err := c.Cache.Exists(ctx, cacheKey)
	if err != nil {
		log.Printf("Cache check error: %v\n", err)
	} else if exists {
		err = c.Cache.Get(ctx, cacheKey, &data)
		if err == nil {
			respondWithJSON(w, http.StatusOK, data.Response[0].League.Standings[0])
			return
		}
		log.Printf("Cache get error: %v\n", err)
	}

	// Use previous year for season
	seasonYear := time.Now().Year() - 1
	url := fmt.Sprintf("https://api-football-v1.p.rapidapi.com/v3/standings?league=%s&season=%d", leagueId, seasonYear)
	headers := map[string]string{
		"Content-Type":   "application/json",
		"x-rapidapi-key": c.FootballAPIKey,
	}

	log.Printf("Making request to URL: %s", url)

	resp, err := HTTPRequest("GET", url, headers, nil)
	if err != nil {
		log.Printf("Error making request: %v", err)
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error creating http request: %s", err))
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error reading response body")
		return
	}

	// Log the raw response for debugging
	log.Printf("Raw response: %s", string(body))

	// Parse the response
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("Error parsing response: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error parsing response")
		return
	}

	// Check if we have any response data
	if len(data.Response) == 0 {
		log.Printf("No response data found for league ID: %s", leagueId)
		respondWithError(w, http.StatusNotFound, "No league standings found for the given league ID")
		return
	}

	// Check if we have any standings data
	if len(data.Response[0].League.Standings) == 0 {
		log.Printf("No standings data found for league: %s", data.Response[0].League.Name)
		respondWithError(w, http.StatusNotFound, "No standings data available for this league")
		return
	}

	// Store in cache with 6 hour TTL
	err = c.Cache.Set(ctx, cacheKey, data, cache.StandingsTTL)
	if err != nil {
		log.Printf("Cache set error: %v\n", err)
	}

	// Return the first standings array
	respondWithJSON(w, http.StatusOK, data.Response[0].League.Standings[0])
}
