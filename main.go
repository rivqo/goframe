package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/example/goframe/auth"
	"github.com/example/goframe/config"
	"github.com/example/goframe/db"
	"github.com/example/goframe/middleware"
	"github.com/example/goframe/router"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	dbConfig := db.DatabaseConfig{
		Driver:   cfg.Database.Driver,
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		Name:     cfg.Database.Name,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
	}

	database, err := db.NewDatabase(&dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize router
	r := router.New()

	// Apply global middleware
	r.Use(middleware.Logger())
	r.Use(middleware.RateLimit(cfg.RateLimit.Requests, cfg.RateLimit.Period))
	r.Use(middleware.Recover())

	// Set up authentication

	authConfig := auth.AuthConfig{
		Secret:   cfg.Auth.Secret,
		Duration:     cfg.Auth.Duration,
	}
	authProvider := auth.NewProvider(database, authConfig)

	// Define routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to GoFrame!")
	})

	// Create a group with auth middleware
	api := r.Group("/api")
	api.Use(authProvider.Middleware())

	api.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUser(r.Context())
		fmt.Fprintf(w, "Hello, %s!", user.Name)
	})

	// Start the server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

