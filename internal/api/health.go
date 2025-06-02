package api

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func (c *Config) healthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "healthy",
		"time":   time.Now().UTC().Format(time.RFC3339),
		"db":     "not connected",
		"cache":  "not connected",
	}

	// Check database connection
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := c.DBConn.PingContext(ctx); err != nil {
		response["db"] = fmt.Sprintf("not connected: %v", err)
	} else {
		response["db"] = "connected"
	}

	// Check Redis connection using EXISTS
	exists, err := c.Cache.Exists(ctx, "health-check")
	if err != nil {
		response["cache"] = fmt.Sprintf("not connected: %v", err)
	} else {
		response["cache"] = "connected"
		response["cache_hit"] = fmt.Sprintf("%v", exists)
	}

	// Always return 200 OK for health checks
	respondWithJSON(w, http.StatusOK, response)
}
