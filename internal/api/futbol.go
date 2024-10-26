package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GetMatchesParams struct {
	Date string
}

func (c *Config) getMatches(w http.ResponseWriter, r *http.Request) {
	// Step 1: Extract data from the incoming request
	decoder := json.NewDecoder(r.Body)
	params := GetMatchesParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}
	defer r.Body.Close()

	// Step 2: Make a request to the third-party service
	url := fmt.Sprintf("https://api-football-v1.p.rapidapi.com/v3/fixtures?date=%s", params.Date)
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
	data := GetMatchesAPIResponse{}
	err = responseBody.Decode(&data)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to read response from football api service: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, data)
}

func (c *Config) getMatch(w http.ResponseWriter, r *http.Request) {
}

func (c *Config) getMatchLineup(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	matchID := queryParams.Get("match_id")
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
	// read the lineups and add the player details
	homeTeamStarters := []Player{}
	for _, p := range getLineUpData.Response[0].StartXI {
		p := Player{
			ID:     p.Player.ID,
			Name:   p.Player.Name,
			Number: p.Player.Number,
			Pos:    p.Player.Pos,
			Grid:   p.Player.Grid,
			Photo:  filterByName(homeTeamSquad.Response[0].Players, p.Player.Name).Photo,
		}
		homeTeamStarters = append(homeTeamStarters, p)
	}
	homeTeamSubstitutes := []Player{}
	for _, p := range getLineUpData.Response[0].Substitutes {
		p := Player{
			ID:     p.Player.ID,
			Name:   p.Player.Name,
			Number: p.Player.Number,
			Pos:    p.Player.Pos,
			Grid:   "",
			Photo:  filterByName(homeTeamSquad.Response[0].Players, p.Player.Name).Photo,
		}
		homeTeamSubstitutes = append(homeTeamSubstitutes, p)
	}

	awayTeamStarters := []Player{}
	for _, p := range getLineUpData.Response[1].StartXI {
		p := Player{
			ID:     p.Player.ID,
			Name:   p.Player.Name,
			Number: p.Player.Number,
			Pos:    p.Player.Pos,
			Grid:   p.Player.Grid,
			Photo:  filterByName(awayTeamSquad.Response[0].Players, p.Player.Name).Photo,
		}
		awayTeamStarters = append(awayTeamStarters, p)
	}
	awayTeamSubstitutes := []Player{}
	for _, p := range getLineUpData.Response[1].Substitutes {
		p := Player{
			ID:     p.Player.ID,
			Name:   p.Player.Name,
			Number: p.Player.Number,
			Pos:    p.Player.Pos,
			Grid:   "",
			Photo:  filterByName(awayTeamSquad.Response[0].Players, p.Player.Name).Photo,
		}
		awayTeamSubstitutes = append(awayTeamSubstitutes, p)
	}

	type Lineup struct {
		Starters    []Player `json:"starters"`
		Substitutes []Player `json:"substitutes"`
	}

	response := struct {
		Home Lineup `json:"home"`
		Away Lineup `json:"away"`
	}{
		Home: Lineup{
			Starters:    homeTeamStarters,
			Substitutes: homeTeamSubstitutes,
		},
		Away: Lineup{
			Starters:    awayTeamStarters,
			Substitutes: awayTeamSubstitutes,
		},
	}
	respondWithJSON(w, http.StatusOK, response)
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
	// Step 2: Make a request to the third-party service
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

	respondWithJSON(w, http.StatusOK, data.Response)
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
