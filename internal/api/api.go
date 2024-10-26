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
	futbolRouter.Get("/leagues", c.getLeagues)
	futbolRouter.Get("/team_standings", c.getLeagueStandingsByTeamId)
	futbolRouter.Get("/league_standings", c.getLeagueStandingsByLeagueId)

	router.Mount("/users", userRouter)
	router.Mount("/futbol", futbolRouter)

	router.Get("/healthz", handleReadiness)
	router.Get("/error", handleError)
	return router
}
