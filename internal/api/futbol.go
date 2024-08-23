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

func (config *Config) handleGetMatches(w http.ResponseWriter, r *http.Request) {
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
	thirdPartyURL := "https://api-football-v1.p.rapidapi.com/v3/fixtures?date=2024-08-23" // Replace with actual URL
	req, err := http.NewRequest(http.MethodGet, thirdPartyURL, nil)
	if err != nil {
		http.Error(w, "Failed to create request to third-party service", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")      // Set headers if needed
	req.Header.Set("x-rapidapi-key", config.FootballAPIKey) // Set headers if needed

	// Optional: Add headers from the incoming request to the third-party request
	// for key, values := range r.Header {
	// 	for _, value := range values {
	// 		req.Header.Add(key, value)
	// 	}
	// }

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request to third-party service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Step 3: Handle the response from the third-party service
	// responseBody, err := ioutil.ReadAll(resp.Body)
	responseBody := json.NewDecoder(resp.Body)
	// TODO: // replace with actualy struct
	data := APIResponse{}
	err = responseBody.Decode(&data)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to read response from football api service: %s", err))
		return
	}

	// Step 4: Send the response back to the client
	respondWithJSON(w, http.StatusOK, data)
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
