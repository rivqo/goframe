package routes

import (
	"net/http"
	"time"

	"github.com/example/goframe/auth"
	"github.com/example/goframe/config"
	"github.com/example/goframe/controllers"
	"github.com/example/goframe/router"
	"github.com/example/goframe/view"
)

// RegisterWebRoutes registers web routes
func RegisterWebRoutes(r *router.Router, cfg *config.Config, authProvider *auth.Provider, authController *auth.Controller) {
	// Create web controller
	webController := controllers.NewWebController()
	
	// Register routes
	r.Get("/", webController.Home)
	r.Get("/d", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	r.Get("/about", webController.About)
	r.Get("/contact", webController.Contact)
	
	// Auth routes
	r.Post("/login", authController.Login)
	r.Post("/register", authController.Register)
	
	// Static files
	r.Static("/assets", "./public/assets")
	r.StaticFile("/favicon.ico", "./public/favicon.ico")
	r.StaticFile("/robots.txt", "./public/robots.txt")
	
	// Protected routes
	protected := r.Group("/dashboard")
	protected.Use(authProvider.Middleware())
	
	protected.Get("", func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUser(r.Context())
		if user == nil {
			http.Error(w, "Not authenticated", http.StatusUnauthorized)
			return
		}
		
		data := map[string]interface{}{
			"title":       "Dashboard",
			"currentYear": time.Now().Year(),
			"user":        user,
		}
		
		view.Render(w, "pages/dashboard", data)
	})
	
	// 404 handler
	r.NotFound(webController.NotFound)
}

