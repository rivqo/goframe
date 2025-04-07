package controllers

import (
	"net/http"
	"time"

	"github.com/example/goframe/view"
)

// WebController handles web requests
type WebController struct {
	// Add your dependencies here
}

// NewWebController creates a new web controller
func NewWebController() *WebController {
	return &WebController{}
}

// Home handles the home page
func (c *WebController) Home(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"title":       "Welcome",
		"currentYear": time.Now().Year(),
	}
	
	view.Render(w, "pages/home", data)
}

// About handles the about page
func (c *WebController) About(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"title":       "About Us",
		"currentYear": time.Now().Year(),
	}
	
	view.Render(w, "pages/about", data)
}

// Contact handles the contact page
func (c *WebController) Contact(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"title":       "Contact Us",
		"currentYear": time.Now().Year(),
	}
	
	view.Render(w, "pages/contact", data)
}

// NotFound handles 404 errors
func (c *WebController) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	
	data := map[string]interface{}{
		"title":       "Page Not Found",
		"currentYear": time.Now().Year(),
	}
	
	view.Render(w, "errors/404", data)
}

