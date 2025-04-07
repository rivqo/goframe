package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/example/goframe/db"
	"github.com/golang-jwt/jwt/v4"
)

type userKey struct{}

type User struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthConfig struct {
	Secret   string
	Duration time.Duration
}

type Provider struct {
	db     *db.Database
	config AuthConfig
}

// NewProvider creates a new auth provider
func NewProvider(database *db.Database, config AuthConfig) *Provider {
	return &Provider{
		db:     database,
		config: config,
	}
}

// Middleware creates a middleware that authenticates requests
func (p *Provider) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the token from the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}
			
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			
			// Parse the token
			token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(p.config.Secret), nil
			})
			
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			
			if claims, ok := token.Claims.(*Claims); ok && token.Valid {
				// Get the user from the database
				user, err := p.GetUserByID(claims.UserID)
				if err != nil {
					http.Error(w, "User not found", http.StatusUnauthorized)
					return
				}
				
				// Add the user to the context
				ctx := context.WithValue(r.Context(), userKey{}, user)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
			}
		})
	}
}

// Login authenticates a user and returns a token
func (p *Provider) Login(email, password string) (string, error) {
	// Get the user from the database
	user, err := p.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	
	// Check the password (this would use a proper password comparison in a real app)
	if !p.CheckPassword(user.ID, password) {
		return "", errors.New("invalid credentials")
	}
	
	// Create the JWT claims
	expirationTime := time.Now().Add(p.config.Duration)
	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	
	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(p.config.Secret))
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

// GetUserByID gets a user by ID from the database
func (p *Provider) GetUserByID(id uint) (*User, error) {
	// This would use the DB in a real app
	// For simplicity, we'll just return a fake user
	if id == 1 {
		return &User{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
		}, nil
	}
	return nil, errors.New("user not found")
}

// GetUserByEmail gets a user by email from the database
func (p *Provider) GetUserByEmail(email string) (*User, error) {
	// This would use the DB in a real app
	// For simplicity, we'll just return a fake user
	if email == "john@example.com" {
		return &User{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
		}, nil
	}
	return nil, errors.New("user not found")
}

// CheckPassword checks if a password is valid for a user
func (p *Provider) CheckPassword(userID uint, password string) bool {
	// This would use proper password comparison in a real app
	// For simplicity, we'll just return true for a specific password
	return password == "password"
}

// GetUser gets the user from the request context
func GetUser(ctx context.Context) *User {
	user, ok := ctx.Value(userKey{}).(*User)
	if !ok {
		return nil
	}
	return user
}

