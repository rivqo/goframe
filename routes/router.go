package routes

import (
	"net/http"
	"time"
	"fmt"
	"encoding/json"

	"github.com/example/goframe/auth"
	"github.com/example/goframe/config"
	"github.com/example/goframe/db"
	"github.com/example/goframe/middleware"
	"github.com/example/goframe/router"
	"github.com/example/goframe/view"
)

// InitializeRouter initializes and configures the application router
func InitializeRouter(cfg *config.Config) (*router.Router, error) {
	dbConfig := db.DatabaseConfig{
		Driver:   cfg.Database.Driver,
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		Name:     cfg.Database.Name,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
	}
	// Initialize database connection
	database, err := db.NewDatabase(&dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create new router instance
	r := router.New()
	if r == nil { // Add this check
		return nil, fmt.Errorf("router.New() returned nil") 
	}

	// Register global middleware one by one
	r.Use(middleware.Logger())
	r.Use(middleware.RateLimit(cfg.RateLimit.Requests, cfg.RateLimit.Period))
	r.Use(middleware.Recover())

	// Setup authentication system
	authProvider := auth.NewProvider(database, auth.AuthConfig{
		Secret:    cfg.Auth.Secret,
		Duration:  cfg.Auth.Duration,
	})
	userRepo := auth.NewUserRepository(database)
	authController := auth.NewController(authProvider, userRepo)

	// Register application routes
	RegisterWebRoutes(r, cfg, authProvider, authController)
	RegisterAPIRoutes(r, cfg, authProvider)

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return r, nil
}

// registerWebRoutes registers web routes
func registerWebRoutes(r *router.Router, cfg *config.Config, authProvider *auth.Provider, authController *auth.Controller) {
	// Example web route with view rendering
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		view.Render(w, "home.html", map[string]interface{}{
			"Title":   "Welcome",
			"Year":    time.Now().Year(),
			"Version": cfg.App.Version,
		})
	})
}

// registerAPIRoutes registers API routes
func registerAPIRoutes(r *router.Router, cfg *config.Config, authProvider *auth.Provider) {
	api := r.Group("/api")
	api.Use(authProvider.Middleware())

	api.Get("/user", func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUser(r.Context())
		if user == nil {
			http.Error(w, "Not authenticated", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})
}