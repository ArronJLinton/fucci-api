package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5XX error: ", msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to marshal JSON response: %v", err))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

type HealthResponse struct {
	Message string `json:"message"`
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, HealthResponse{Message: "OK. Server Ready."})
}

func handleError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Something went wrong.")
}

func HTTPRequest(method, url string, headers map[string]string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return resp, nil
}

func handleClientRequest[T any](url string, method string, headers map[string]string) (*T, error) {
	resp, err := HTTPRequest(method, url, headers, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating http request: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("non-2xx status code: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}
	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	return &result, nil
}
