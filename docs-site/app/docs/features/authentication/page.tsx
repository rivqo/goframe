import { DocPagination } from "@/components/doc-pagination"

export default function AuthenticationPage() {
  return (
    <div className="space-y-6">
      <h1>Authentication</h1>
      <p>
        GoFrame provides a robust authentication system based on JSON Web Tokens (JWT). The authentication system
        includes user registration, login, middleware for protected routes, and more.
      </p>

      <h2>Configuration</h2>
      <p>
        Authentication settings are configured in the <code>config.yaml</code> file:
      </p>

      <pre>
        <code>{`auth:
  secret: your-secret-key-here
  duration: 24h`}</code>
      </pre>

      <p>
        The <code>secret</code> is used to sign JWT tokens, and <code>duration</code> specifies how long tokens are
        valid.
      </p>

      <h2>Authentication Provider</h2>
      <p>
        The <code>auth.Provider</code> is the central component of the authentication system. It handles token
        generation, validation, and middleware.
      </p>

      <pre>
        <code>{`// auth/auth.go
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
    ID    uint   \`json:"id"\`
    Name  string \`json:"name"\`
    Email string \`json:"email"\`
}

type Claims struct {
    UserID uint \`json:"user_id"\`
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
}`}</code>
      </pre>

      <h2>Authentication Middleware</h2>
      <p>
        The authentication middleware checks for a valid JWT token in the Authorization header and adds the user to the
        request context if authenticated.
      </p>

      <pre>
        <code>{`// Middleware creates a middleware that authenticates requests
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
}`}</code>
      </pre>

      <h2>Login and Registration</h2>
      <p>
        The <code>auth.Controller</code> provides HTTP handlers for login and registration:
      </p>

      <pre>
        <code>{`// auth/auth_controller.go
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
    Email    string \`json:"email"\`
    Password string \`json:"password"\`
}

// LoginResponse represents a login response
type LoginResponse struct {
    Token string \`json:"token"\`
    User  User   \`json:"user"\`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
    Name     string \`json:"name"\`
    Email    string \`json:"email"\`
    Password string \`json:"password"\`
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
}`}</code>
      </pre>

      <h2>Getting the Current User</h2>
      <p>You can get the current authenticated user from the request context:</p>

      <pre>
        <code>{`// GetUser gets the user from the request context
func GetUser(ctx context.Context) *User {
    user, ok := ctx.Value(userKey{}).(*User)
    if !ok {
        return nil
    }
    return user
}

// Example usage in a controller
func (c *Controller) Me(w http.ResponseWriter, r *http.Request) {
    user := auth.GetUser(r.Context())
    if user == nil {
        http.Error(w, "Not authenticated", http.StatusUnauthorized)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}`}</code>
      </pre>

      <h2>Protecting Routes</h2>
      <p>You can protect routes by applying the authentication middleware:</p>

      <pre>
        <code>{`// routes/web.go
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
})`}</code>
      </pre>

      <h2>Password Hashing</h2>
      <p>GoFrame uses bcrypt for secure password hashing:</p>

      <pre>
        <code>{`// auth/user_model.go
// Create creates a new user
func (r *UserRepository) Create(name, email, password string) (*UserModel, error) {
    // Hash the password
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    user := &UserModel{
        Name:         name,
        Email:        email,
        PasswordHash: string(hash),
    }

    if err := r.repo.Create(user); err != nil {
        return nil, err
    }

    return user, nil
}

// CheckPassword checks if a password is valid for a user
func (r *UserRepository) CheckPassword(user *UserModel, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
    return err == nil
}`}</code>
      </pre>

      <DocPagination
        prev={{
          title: "Features",
          href: "/docs/features",
        }}
        next={{
          title: "Resources",
          href: "/docs/features/resources",
        }}
      />
    </div>
  )
}

