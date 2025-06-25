package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/ArronJLinton/fucci-api/internal/cache"
)

type GoogleNewsResponse struct {
	Status string `json:"status"`
	Items  []struct {
		Timestamp string `json:"timestamp"`
		Title     string `json:"title"`
		Snippet   string `json:"snippet"`
		Images    struct {
			Thumbnail        string `json:"thumbnail"`
			ThumbnailProxied string `json:"thumbnailProxied"`
		} `json:"images"`
		NewsUrl   string `json:"newsUrl"`
		Publisher string `json:"publisher"`
	} `json:"items"`
}

func (c *Config) search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get query parameter from request
	query := r.URL.Query().Get("q")
	if query == "" {
		respondWithError(w, http.StatusBadRequest, "Query parameter 'q' is required")
		return
	}

	// Get language parameter (optional, default to en-US)
	language := r.URL.Query().Get("lr")
	if language == "" {
		language = "en-US"
	}

	// Generate cache key
	cacheKey := fmt.Sprintf("google_news:%s:%s", query, language)

	// Try to get from cache first
	var newsResponse GoogleNewsResponse
	exists, err := c.Cache.Exists(ctx, cacheKey)
	if err != nil {
		log.Printf("Cache check error: %v\n", err)
	} else if exists {
		err = c.Cache.Get(ctx, cacheKey, &newsResponse)
		if err == nil {
			respondWithJSON(w, http.StatusOK, newsResponse)
			return
		}
		log.Printf("Cache get error: %v\n", err)
	}

	// Construct the URL
	baseURL := "https://google-news13.p.rapidapi.com/search"
	params := url.Values{}
	params.Add("keyword", query)
	params.Add("lr", language)

	// Create HTTP request
	req, err := http.NewRequest("GET", baseURL+"?"+params.Encode(), nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create request")
		return
	}

	// Add RapidAPI headers
	req.Header.Add("x-rapidapi-key", c.RapidAPIKey)
	req.Header.Add("x-rapidapi-host", "google-news13.p.rapidapi.com")

	// Make the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch news")
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to read response")
		return
	}

	// Check if response is successful
	if resp.StatusCode != http.StatusOK {
		log.Printf("API returned status %d: %s", resp.StatusCode, string(body))
		respondWithError(w, resp.StatusCode, fmt.Sprintf("News API error: %s", string(body)))
		return
	}

	// Parse JSON response
	if err := json.Unmarshal(body, &newsResponse); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to parse response")
		return
	}

	// Cache the response for 30 minutes (news data changes frequently)
	err = c.Cache.Set(ctx, cacheKey, newsResponse, cache.NewsTTL)
	if err != nil {
		log.Printf("Cache set error: %v\n", err)
	}

	// Return the response
	respondWithJSON(w, http.StatusOK, newsResponse)
}
