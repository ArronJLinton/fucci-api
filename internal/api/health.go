package api

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func (c *Config) healthCheck(w http.ResponseWriter, r *http.Request) {
	// Check database connection
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := c.DBConn.PingContext(ctx); err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "database connection failed")
		return
	}

	// Check Redis connection using EXISTS
	exists, err := c.Cache.Exists(ctx, "health-check")
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "cache connection failed")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"status":    "healthy",
		"time":      time.Now().UTC().Format(time.RFC3339),
		"db":        "connected",
		"cache":     "connected",
		"cache_hit": fmt.Sprintf("%v", exists),
	})
}
