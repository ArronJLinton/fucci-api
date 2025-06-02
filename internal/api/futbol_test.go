package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMatchLineup(t *testing.T) {
	mockAPIKey := "mock-api-key"
	config := &Config{FootballAPIKey: mockAPIKey}

	t.Run("success case", func(t *testing.T) {
		// Mock HTTP response from the external API
		mockResponse := `{"response": "mocked data"}`
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("x-rapidapi-key") != mockAPIKey {
				t.Errorf("Expected API key %s, got %s", mockAPIKey, r.Header.Get("x-rapidapi-key"))
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(mockResponse))
		}))
		defer server.Close()

		// Replace the URL with the mock server URL
		url := fmt.Sprintf("%s/fixtures/lineups?fixture=12345", server.URL)
		req := httptest.NewRequest("GET", url, nil)
		req.URL.RawQuery = "match_id=12345"
		rec := httptest.NewRecorder()

		// Call the function
		config.getMatchLineup(rec, req)

		// Assert the response
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to parse response: %s", err)
		}
		if response["response"] != "mocked data" {
			t.Errorf("Expected response 'mocked data', got %v", response["response"])
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
		expectedError := "Error creating http request: match_id is required"
		if rec.Body.String() != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, rec.Body.String())
		}
	})
}
