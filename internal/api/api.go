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

func New(config Config) http.Handler {
	router := chi.NewRouter()
	userRouter := chi.NewRouter()
	userRouter.Post("/create", config.handleCreateUser)

	futbolRouter := chi.NewRouter()
	futbolRouter.Get("/matches", config.getMatches)

	router.Mount("/users", userRouter)
	router.Mount("/futbol", futbolRouter)

	router.Get("/healthz", handleReadiness)
	router.Get("/error", handleError)
	return router
}
