package api

import (
	"database/sql"
	"net/http"

	"github.com/ArronJLinton/fucci-api/internal/ai"
	"github.com/ArronJLinton/fucci-api/internal/cache"
	"github.com/ArronJLinton/fucci-api/internal/database"
	"github.com/go-chi/chi"
)

type Config struct {
	DB                 *database.Queries
	DBConn             *sql.DB
	FootballAPIKey     string
	RapidAPIKey        string
	Cache              cache.CacheInterface
	APIFootballBaseURL string
	OpenAIKey          string
	OpenAIBaseURL      string
	AIPromptGenerator  *ai.PromptGenerator
}

func New(c Config) http.Handler {
	router := chi.NewRouter()

	// Initialize AI prompt generator if OpenAI key is provided
	if c.OpenAIKey != "" {
		c.AIPromptGenerator = ai.NewPromptGenerator(c.OpenAIKey, c.OpenAIBaseURL, c.Cache)
	}

	// Health check routes
	router.Get("/health", HandleReadiness)
	router.Get("/health/redis", c.HandleRedisHealth)
	router.Get("/health/cache-stats", c.HandleCacheStats)

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

	debateRouter := chi.NewRouter()
	debateRouter.Post("/", c.createDebate)
	debateRouter.Get("/top", c.getTopDebates)
	debateRouter.Get("/generate", c.generateAIPrompt)
	debateRouter.Get("/match", c.getDebatesByMatch)
	debateRouter.Get("/{id}", c.getDebate)
	debateRouter.Post("/cards", c.createDebateCard)
	debateRouter.Post("/votes", c.createVote)
	debateRouter.Post("/comments", c.createComment)
	debateRouter.Get("/{debateId}/comments", c.getComments)

	router.Mount("/users", userRouter)
	router.Mount("/futbol", futbolRouter)
	router.Mount("/google", googleRouter)
	router.Mount("/debates", debateRouter)

	return router
}
