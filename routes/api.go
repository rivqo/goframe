package routes

import (
	"encoding/json"
	"net/http"

	"github.com/example/goframe/auth"
	"github.com/example/goframe/config"
	"github.com/example/goframe/router"
)

// RegisterAPIRoutes registers API routes
func RegisterAPIRoutes(r *router.Router, cfg *config.Config, authProvider *auth.Provider) {
	// Create API group
	api := r.Group("/api")

	// Add middleware
	api.Use(authProvider.Middleware())

	// Register routes
	api.Get("/user", getUserHandler)

	// Register resource routes
	// Example: RegisterResourceRoutes(api, "/users", &UserController{})
}

// getUserHandler handles GET /api/user requests
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// RegisterResourceRoutes registers RESTful routes for a resource
func RegisterResourceRoutes(g *router.RouteGroup, path string, controller ResourceController) {
	g.Get(path, controller.Index)
	g.Get(path+"/:id", controller.Show)
	g.Post(path, controller.Store)
	g.Put(path+"/:id", controller.Update)
	g.Delete(path+"/:id", controller.Destroy)
}

// ResourceController defines the interface for resource controllers
type ResourceController interface {
	Index(w http.ResponseWriter, r *http.Request)
	Show(w http.ResponseWriter, r *http.Request)
	Store(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Destroy(w http.ResponseWriter, r *http.Request)
}