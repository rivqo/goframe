package auth

import (
	"encoding/json"
	"net/http"
)

// Controller handles authentication-related HTTP requests
type Controller struct {
	provider *Provider
	repo     *UserRepository
}

// NewController creates a new auth controller
func NewController(provider *Provider, repo *UserRepository) *Controller {
	return &Controller{
		provider: provider,
		repo:     repo,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login handles login requests
func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	// Authenticate the user
	token, err := c.provider.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	
	// Get the user
	userModel, err := c.repo.FindByEmail(req.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}
	
	// Create the response
	resp := LoginResponse{
		Token: token,
		User: User{
			ID:    userModel.ID,
			Name:  userModel.Name,
			Email: userModel.Email,
		},
	}
	
	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Register handles registration requests
func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	// Check if the user already exists
	_, err := c.repo.FindByEmail(req.Email)
	if err == nil {
		http.Error(w, "Email already in use", http.StatusConflict)
		return
	}
	
	// Create the user
	userModel, err := c.repo.Create(req.Name, req.Email, req.Password)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	
	// Authenticate the user
	token, err := c.provider.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}
	
	// Create the response
	resp := LoginResponse{
		Token: token,
		User: User{
			ID:    userModel.ID,
			Name:  userModel.Name,
			Email: userModel.Email,
		},
	}
	
	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// Me handles requests to get the current user
func (c *Controller) Me(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r.Context())
	if user == nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}
	
	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

