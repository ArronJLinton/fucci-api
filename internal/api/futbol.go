package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GetMatchesParams struct {
	Date string
}

type APIResponse struct {
	Get        string `json:"get"`
	Parameters struct {
		Date string `json:"date"`
	} `json:"parameters"`
	Errors  []any `json:"errors"`
	Results int   `json:"results"`
	Paging  struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"paging"`
	Response []struct {
		Fixture struct {
			ID        int       `json:"id"`
			Referee   string    `json:"referee"`
			Timezone  string    `json:"timezone"`
			Date      time.Time `json:"date"`
			Timestamp int       `json:"timestamp"`
			Periods   struct {
				First  int `json:"first"`
				Second int `json:"second"`
			} `json:"periods"`
			Venue struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				City string `json:"city"`
			} `json:"venue"`
			Status struct {
				Long    string `json:"long"`
				Short   string `json:"short"`
				Elapsed int    `json:"elapsed"`
			} `json:"status"`
		} `json:"fixture"`
		League struct {
			ID      int    `json:"id"`
			Name    string `json:"name"`
			Country string `json:"country"`
			Logo    string `json:"logo"`
			Flag    any    `json:"flag"`
			Season  int    `json:"season"`
			Round   string `json:"round"`
		} `json:"league"`
		Teams struct {
			Home struct {
				ID     int    `json:"id"`
				Name   string `json:"name"`
				Logo   string `json:"logo"`
				Winner any    `json:"winner"`
			} `json:"home"`
			Away struct {
				ID     int    `json:"id"`
				Name   string `json:"name"`
				Logo   string `json:"logo"`
				Winner any    `json:"winner"`
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
				Home any `json:"home"`
				Away any `json:"away"`
			} `json:"extratime"`
			Penalty struct {
				Home any `json:"home"`
				Away any `json:"away"`
			} `json:"penalty"`
		} `json:"score"`
	} `json:"response"`
}

func (config *Config) getMatches(w http.ResponseWriter, r *http.Request) {
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
		"x-rapidapi-key": config.FootballAPIKey,
	}
	resp, err := HTTPRequest("GET", url, headers, nil)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error creating http request: %s", err))
		return
	}
	defer resp.Body.Close()

	responseBody := json.NewDecoder(resp.Body)
	data := APIResponse{}
	err = responseBody.Decode(&data)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to read response from football api service: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, data)
}

func (config *Config) getMatch(w http.ResponseWriter, r *http.Request) {
}

func (config *Config) getMatchLineup(w http.ResponseWriter, r *http.Request) {
}
