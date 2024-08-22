package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ArronJLinton/fucci-api/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the enviroment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to Database - ", err)
	}

	config := apiConfig{
		DB: database.New(conn),
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
	v1Router.Get("/healthz", handleReadiness)
	v1Router.Get("/error", handleError)

	apiRouter := chi.NewRouter()
	userRouter := chi.NewRouter()
	userRouter.Post("/create", config.handleCreateUser)

	apiRouter.Mount("/users", userRouter)
	v1Router.Mount("/api", apiRouter)
	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	fmt.Printf("Server starting on port %v", portString)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
