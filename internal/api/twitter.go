package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ArronJLinton/fucci-api/internal/cache"
)

const (
	twitterAPIBaseURL = "https://api.twitter.com/2"
)

type TwitterResponse struct {
	Data []Tweet `json:"data"`
	Meta struct {
		ResultCount   int    `json:"result_count"`
		NextToken     string `json:"next_token,omitempty"`
		PreviousToken string `json:"previous_token,omitempty"`
	} `json:"meta"`
}

type Tweet struct {
	ID            string `json:"id"`
	Text          string `json:"text"`
	CreatedAt     string `json:"created_at"`
	AuthorID      string `json:"author_id"`
	PublicMetrics struct {
		RetweetCount int `json:"retweet_count"`
		ReplyCount   int `json:"reply_count"`
		LikeCount    int `json:"like_count"`
		QuoteCount   int `json:"quote_count"`
	} `json:"public_metrics"`
}

// getTeamTweets fetches recent tweets about a team
func (c *Config) getTeamTweets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queryParams := r.URL.Query()
	teamName := queryParams.Get("team_name")
	if teamName == "" {
		respondWithError(w, http.StatusBadRequest, "team_name is required")
		return
	}

	// Generate cache key
	cacheKey := fmt.Sprintf("tweets:team:%s", teamName)

	// Try to get from cache first
	var response TwitterResponse
	exists, err := c.Cache.Exists(ctx, cacheKey)
	if err != nil {
		fmt.Printf("Cache check error: %v\n", err)
	} else if exists {
		err = c.Cache.Get(ctx, cacheKey, &response)
		if err == nil {
			respondWithJSON(w, http.StatusOK, response)
			return
		}
		fmt.Printf("Cache get error: %v\n", err)
	}

	// If not in cache, fetch from Twitter API
	encodedQuery := url.QueryEscape(teamName)
	url := fmt.Sprintf("%s/tweets/search/recent?query=%s&max_results=10", twitterAPIBaseURL, encodedQuery)

	// Add debug logging
	fmt.Printf("Making Twitter API request to: %s\n", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating request: %s", err))
		return
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.TwitterAPIKey))
	fmt.Printf("Request headers: %+v\n", req.Header)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %s\n", err)
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error fetching tweets: %s", err))
		return
	}
	defer resp.Body.Close()

	// Log response status and headers
	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Headers: %+v\n", resp.Header)

	// Read and log the raw response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error reading response: %s", err))
		return
	}

	// Log the raw response
	fmt.Printf("Raw Response Body: %s\n", string(bodyBytes))

	// Create a new reader with the body bytes for json.Decoder
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Printf("Error parsing JSON: %s\n", err)
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing Twitter response: %s", err))
		return
	}

	// Log the parsed response
	fmt.Printf("Parsed Response: %+v\n", response)

	// Store in cache for 15 minutes
	err = c.Cache.Set(ctx, cacheKey, response, cache.TwitterTTL)
	if err != nil {
		fmt.Printf("Cache set error: %v\n", err)
	}

	respondWithJSON(w, http.StatusOK, response)
}

// getMatchTweets fetches recent tweets about a specific match
func (c *Config) getMatchTweets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queryParams := r.URL.Query()
	matchID := queryParams.Get("match_id")
	if matchID == "" {
		respondWithError(w, http.StatusBadRequest, "match_id is required")
		return
	}

	// Get match details to construct meaningful search query
	apiURL := fmt.Sprintf("https://api-football-v1.p.rapidapi.com/v3/fixtures?id=%s", matchID)
	headers := map[string]string{
		"Content-Type":   "application/json",
		"x-rapidapi-key": c.FootballAPIKey,
	}

	matchDetails, err := handleClientRequest[GetMatchesAPIResponse](apiURL, "GET", headers)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error fetching match details: %s", err))
		return
	}

	if len(matchDetails.Response) == 0 {
		respondWithError(w, http.StatusNotFound, "Match not found")
		return
	}

	match := matchDetails.Response[0]
	homeTeam := match.Teams.Home.Name
	awayTeam := match.Teams.Away.Name

	// Generate cache key
	cacheKey := fmt.Sprintf("tweets:match:%s", matchID)

	// Try to get from cache first
	var response TwitterResponse
	exists, err := c.Cache.Exists(ctx, cacheKey)
	if err != nil {
		fmt.Printf("Cache check error: %v\n", err)
	} else if exists {
		err = c.Cache.Get(ctx, cacheKey, &response)
		if err == nil {
			respondWithJSON(w, http.StatusOK, response)
			return
		}
		fmt.Printf("Cache get error: %v\n", err)
	}

	// Construct search query using both team names
	searchQuery := fmt.Sprintf("(%s %s) lang:en -is:retweet", homeTeam, awayTeam)
	query := url.QueryEscape(searchQuery)
	twitterURL := fmt.Sprintf("%s/tweets/search/recent?query=%s&tweet.fields=created_at,public_metrics,author_id&max_results=10",
		twitterAPIBaseURL, query)

	headers = map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", c.TwitterAPIKey),
		"Content-Type":  "application/json",
	}

	// Add debug logging
	fmt.Printf("Making Twitter API request to: %s\n", twitterURL)
	fmt.Printf("With headers: %+v\n", headers)

	resp, err := HTTPRequest("GET", twitterURL, headers, nil)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error fetching tweets: %s", err))
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing Twitter response: %s", err))
		return
	}

	// Store in cache for 5 minutes (shorter TTL for match tweets as they're more time-sensitive)
	err = c.Cache.Set(ctx, cacheKey, response, cache.MatchTweetsTTL)
	if err != nil {
		fmt.Printf("Cache set error: %v\n", err)
	}

	respondWithJSON(w, http.StatusOK, response)
}
