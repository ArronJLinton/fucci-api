package api

import (
	"net/http"

	"github.com/ArronJLinton/fucci-api/internal/database"
	"github.com/go-chi/chi"
)

type Config struct {
	DB             *database.Queries
	FootballAPIKey string
}

func New(c Config) http.Handler {
	router := chi.NewRouter()
	userRouter := chi.NewRouter()
	userRouter.Post("/create", c.handleCreateUser)

	futbolRouter := chi.NewRouter()
	futbolRouter.Get("/matches", c.getMatches)
	futbolRouter.Get("/lineup", c.getMatchLineup)

	router.Mount("/users", userRouter)
	router.Mount("/futbol", futbolRouter)

	router.Get("/healthz", handleReadiness)
	router.Get("/error", handleError)
	return router
}
