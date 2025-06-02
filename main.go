package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/ArronJLinton/fucci-api/internal/api"
	"github.com/ArronJLinton/fucci-api/internal/cache"
	"github.com/ArronJLinton/fucci-api/internal/config"
	"github.com/ArronJLinton/fucci-api/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

type Config struct {
	DB *database.Queries
}

var (
	version = "dev"
)

func main() {
	// Initialize the logger
	// TODO: What does the following code do?
	zlog, _ := zap.NewProduction(
		zap.Fields(
			zap.String("version", version),
		),
	)
	defer func() {
		_ = zlog.Sync()
	}()
	logger := otelzap.New(zlog)

	// Initialize the configuration
	c := config.InitConfig(logger)
	conn, err := sql.Open("postgres", c.DB_URL)
	if err != nil {
		log.Fatal("Failed to connect to Database - ", err)
	}

	// Initialize Redis cache
	redisCache, err := cache.NewCache(c.REDIS_URL)
	if err != nil {
		log.Fatal("Failed to connect to Redis - ", err)
	}

	router := chi.NewRouter()
	// Tells browsers how this api can be used
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", ";http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"string"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	apiRouter := api.New(api.Config{
		DB:             database.New(conn),
		FootballAPIKey: c.FOOTBALL_API_KEY,
		Cache:          redisCache,
	})
	v1Router.Mount("/api", apiRouter)
	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + c.PORT,
	}
	fmt.Printf("Server starting on port %v", c.PORT)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
