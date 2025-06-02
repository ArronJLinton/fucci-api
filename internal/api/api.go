package api

import (
	"database/sql"
	"net/http"

	"github.com/ArronJLinton/fucci-api/internal/cache"
	"github.com/ArronJLinton/fucci-api/internal/database"
	"github.com/go-chi/chi"
)

type Config struct {
	DB             *database.Queries
	DBConn         *sql.DB
	FootballAPIKey string
	Cache          *cache.Cache
}

func New(c Config) http.Handler {
	router := chi.NewRouter()

	// Add health check route
	router.Get("/health", c.healthCheck)

	userRouter := chi.NewRouter()
	userRouter.Post("/create", c.handleCreateUser)

	futbolRouter := chi.NewRouter()
	futbolRouter.Get("/matches", c.getMatches)
	futbolRouter.Get("/lineup", c.getMatchLineup)
	futbolRouter.Get("/leagues", c.getLeagues)
	futbolRouter.Get("/team_standings", c.getLeagueStandingsByTeamId)
	futbolRouter.Get("/league_standings", c.getLeagueStandingsByLeagueId)

	googleRouter := chi.NewRouter()
	googleRouter.Get("/search", c.search)

	router.Mount("/users", userRouter)
	router.Mount("/futbol", futbolRouter)
	router.Mount("/google", googleRouter)

	router.Get("/healthz", handleReadiness)
	router.Get("/error", handleError)
	return router
}
