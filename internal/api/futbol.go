package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

	fmt.Printf("Fetching matches for date: %s\n", date)

	// Generate cache key
	cacheKey := fmt.Sprintf("matches:%s", date)

	// Try to get from cache first
	var data GetMatchesAPIResponse
	exists, err := c.Cache.Exists(ctx, cacheKey)
	if err != nil {
		fmt.Printf("Cache check error: %v\n", err)
	} else if exists {
		err = c.Cache.Get(ctx, cacheKey, &data)
		if err == nil {
			fmt.Printf("Cache hit for date: %s\n", date)
			respondWithJSON(w, http.StatusOK, data)
			return
		}
		fmt.Printf("Cache get error: %v\n", err)
	}

	// If not in cache or error occurred, fetch from API
	url := fmt.Sprintf("https://api-football-v1.p.rapidapi.com/v3/fixtures?date=%s", date)
	headers := map[string]string{
		"Content-Type":   "application/json",
		"x-rapidapi-key": c.FootballAPIKey,
	}

	fmt.Printf("Making API request to: %s\n", url)
	fmt.Printf("Using API key: %s\n", c.FootballAPIKey)

	resp, err := HTTPRequest("GET", url, headers, nil)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error creating http request: %s", err))
		return
	}
	defer resp.Body.Close()

	fmt.Printf("API Response Status: %s\n", resp.Status)

	// Read the raw response body for debugging
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to read response body: %s", err))
		return
	}
	fmt.Printf("Raw API Response: %s\n", string(rawBody))

	// Create a new reader from the raw body for JSON decoding
	err = json.Unmarshal(rawBody, &data)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to parse response from football api service: %s", err))
		return
	}

	fmt.Printf("Parsed Response - Results: %d, Response Length: %d\n", data.Results, len(data.Response))

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
		fmt.Printf("Cache set error: %v\n", err)
	} else {
		fmt.Printf("Stored in cache: %s with TTL: %v\n", cacheKey, ttl)
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

	url := fmt.Sprintf("https://api-football-v1.p.rapidapi.com/v3/fixtures/lineups?fixture=%s", matchID)
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
	getLineUpData := &GetLineUpResponse{}
	err = responseBody.Decode(&getLineUpData)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to read response from football api service - %s", err))
		return
	}
	if len(getLineUpData.Response) < 2 {
		respondWithJSON(w, http.StatusOK, "No lineup data available")
		return
	}

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
		p := Player{
			ID:     p.Player.ID,
			Name:   p.Player.Name,
			Number: p.Player.Number,
			Pos:    p.Player.Pos,
			Grid:   p.Player.Grid,
			Photo:  filterByName(squad.Response[0].Players, p.Player.Name).Photo,
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
		p := Player{
			ID:     p.Player.ID,
			Name:   p.Player.Name,
			Number: p.Player.Number,
			Pos:    p.Player.Pos,
			Grid:   "",
			Photo:  filterByName(squad.Response[0].Players, p.Player.Name).Photo,
		}
		result = append(result, p)
	}
	return result
}

func (c *Config) getTeamSquad(id int32) (*GetSquadResponse, error) {
	headers := map[string]string{
		"Content-Type":   "application/json",
		"x-rapidapi-key": c.FootballAPIKey,
	}
	url := fmt.Sprintf("https://api-football-v1.p.rapidapi.com/v3/players/squads?team=%d", id)
	response, err := handleClientRequest[GetSquadResponse](url, "GET", headers)
	if err != nil {
		return nil, fmt.Errorf("error creating http request: %s", err)
	}
	return response, nil
}

func filterByName(items []Player, name string) Player {
	player := Player{}
	for _, item := range items {
		if item.Name == name {
			player = item
			break
		}
	}

	return player
}

func (c *Config) getLeagues(w http.ResponseWriter, r *http.Request) {
	url := "https://api-football-v1.p.rapidapi.com/v3/leagues?season=2024"
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
	url := fmt.Sprintf("https://api-football-v1.p.rapidapi.com/v3/standings?season=2024&team=%s", teamId)
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
	queryParams := r.URL.Query()
	leagueId := queryParams.Get("league_id")
	// TODO: Dynamically set the season year
	url := fmt.Sprintf("https://api-football-v1.p.rapidapi.com/v3/standings?league=%s&season=2024", leagueId)
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
	data := GetLeagueStandingsByLeagueIdResponse{}
	err = responseBody.Decode(&data)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing response body: %s", err))
		return
	}
	respondWithJSON(w, http.StatusOK, data.Response[0].League.Standings[0])
}
