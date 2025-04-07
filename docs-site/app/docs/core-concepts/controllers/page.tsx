import { DocPagination } from "@/components/doc-pagination"

export default function ControllersPage() {
  return (
    <div className="space-y-6">
      <h1>Controllers</h1>
      <p>
        Controllers are responsible for handling requests and returning responses. They are the central part of your
        application's logic and help organize your code by grouping related request handling logic into a single class.
      </p>

      <h2>Basic Controllers</h2>
      <p>A basic controller is a struct with methods that handle HTTP requests:</p>
      <pre>
        <code>{`// controllers/user_controller.go
package controllers

import (
    "net/http"
    "encoding/json"
    
    "github.com/example/goframe/models"
)

// UserController handles user-related HTTP requests
type UserController struct {
    repo *models.UserRepository
}

// NewUserController creates a new user controller
func NewUserController(repo *models.UserRepository) *UserController {
    return &UserController{
        repo: repo,
    }
}

// Index returns all users
func (c *UserController) Index(w http.ResponseWriter, r *http.Request) {
    users, err := c.repo.FindAll()
    if err != nil {
        http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

// Show returns a single user  "application/json")
    json.NewEncoder(w).Encode(users)
}

// Show returns a single user
func (c *UserController) Show(w http.ResponseWriter, r *http.Request) {
    // Get ID from URL
    id := router.GetParam(r, "id")
    
    // Find the user
    user, err := c.repo.FindByID(id)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    
    // Return the user
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

// Store creates a new user
func (c *UserController) Store(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    if err := c.repo.Create(&user); err != nil {
        http.Error(w, "Failed to create user", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

// Update updates a user
func (c *UserController) Update(w http.ResponseWriter, r *http.Request) {
    // Get ID from URL
    id := router.GetParam(r, "id")
    
    // Find the user
    user, err := c.repo.FindByID(id)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    
    // Update the user
    if err := json.NewDecoder(r.Body).Decode(user); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    if err := c.repo.Update(user); err != nil {
        http.Error(w, "Failed to update user", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

// Destroy deletes a user
func (c *UserController) Destroy(w http.ResponseWriter, r *http.Request) {
    // Get ID from URL
    id := router.GetParam(r, "id")
    
    // Find the user
    user, err := c.repo.FindByID(id)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    
    // Delete the user
    if err := c.repo.Delete(user); err != nil {
        http.Error(w, "Failed to delete user", http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}`}</code>
      </pre>

      <h2>Resource Controllers</h2>
      <p>
        Resource controllers provide a convenient way to build RESTful controllers around resources. You can create a
        resource controller using the <code>make:controller</code> command with the <code>--resource</code> flag:
      </p>
      <pre>
        <code>{`goframe make:controller Post --resource`}</code>
      </pre>

      <p>This will generate a controller with the following methods:</p>

      <ul>
        <li>
          <code>Index</code> - Display a listing of the resource
        </li>
        <li>
          <code>Show</code> - Display a specific resource
        </li>
        <li>
          <code>Store</code> - Store a new resource
        </li>
        <li>
          <code>Update</code> - Update a specific resource
        </li>
        <li>
          <code>Destroy</code> - Delete a specific resource
        </li>
      </ul>

      <p>You can register all the resource routes with a single line:</p>

      <pre>
        <code>{`// Register resource routes
routes.RegisterResourceRoutes(api, "/posts", postController)`}</code>
      </pre>

      <h2>Dependency Injection</h2>
      <p>GoFrame encourages the use of dependency injection to manage your controller's dependencies:</p>

      <pre>
        <code>{`// Initialize repositories
userRepo := models.NewUserRepository(db)
postRepo := models.NewPostRepository(db)

// Initialize controllers with dependencies
userController := controllers.NewUserController(userRepo)
postController := controllers.NewPostController(postRepo, userRepo)

// Register routes
r.Get("/users", userController.Index)
r.Get("/users/:id", userController.Show)
// ...`}</code>
      </pre>

      <h2>Controller Middleware</h2>
      <p>You can apply middleware to controller methods:</p>

      <pre>
        <code>{`// Apply middleware to a controller method
r.Get("/admin/dashboard", middleware.Auth()(adminController.Dashboard))

// Or apply middleware to a group of routes
admin := r.Group("/admin")
admin.Use(middleware.Auth())
admin.Use(middleware.Role("admin"))

admin.Get("/dashboard", adminController.Dashboard)
admin.Get("/users", adminController.Users)`}</code>
      </pre>

      <h2>Web Controllers</h2>
      <p>Web controllers are used to render views:</p>

      <pre>
        <code>{`// controllers/web_controller.go
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
}`}</code>
      </pre>

      <DocPagination
        prev={{
          title: "Routing",
          href: "/docs/core-concepts/routing",
        }}
        next={{
          title: "Models",
          href: "/docs/core-concepts/models",
        }}
      />
    </div>
  )
}

