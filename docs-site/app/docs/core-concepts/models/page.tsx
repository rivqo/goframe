import { DocPagination } from "@/components/doc-pagination"

export default function ModelsPage() {
  return (
    <div className="space-y-6">
      <h1>Models</h1>
      <p>
        Models in GoFrame represent the data structures of your application and provide methods for interacting with the
        database. They are the central component of the ORM system.
      </p>

      <h2>Defining Models</h2>
      <p>
        Models are defined as Go structs that embed the <code>db.Entity</code> struct. The <code>db.Entity</code>
        struct provides common fields like <code>ID</code>, <code>CreatedAt</code>, and <code>UpdatedAt</code>.
      </p>

      <pre>
        <code>{`// models/user.go
package models

import (
    "time"
    
    "github.com/example/goframe/db"
)

// User represents a user in the database
type User struct {
    db.Entity
    Name     string \`db:"name" json:"name"\`
    Email    string \`db:"email" json:"email"\`
    Password string \`db:"password" json:"-"\`
}

// TableName returns the table name for the model
func (User) TableName() string {
    return "users"
}`}</code>
      </pre>

      <h3>Field Tags</h3>
      <p>Field tags are used to specify how fields are mapped to database columns and JSON properties:</p>
      <ul>
        <li>
          <code>db:"name"</code> - Specifies the database column name
        </li>
        <li>
          <code>json:"name"</code> - Specifies the JSON property name
        </li>
        <li>
          <code>json:"-"</code> - Excludes the field from JSON serialization
        </li>
      </ul>

      <h2>Creating Models</h2>
      <p>
        You can create a new model using the <code>make:model</code> command:
      </p>

      <pre>
        <code>{`goframe make:model User`}</code>
      </pre>

      <p>
        This will create a new model file in the <code>models</code> directory. You can also create a migration for the
        model at the same time:
      </p>

      <pre>
        <code>{`goframe make:model User -m`}</code>
      </pre>

      <h2>Repositories</h2>
      <p>
        GoFrame uses the repository pattern to encapsulate database operations. Each model typically has a corresponding
        repository that provides methods for interacting with the database.
      </p>

      <pre>
        <code>{`// models/user.go
// UserRepository provides methods to interact with users
type UserRepository struct {
    repo *db.Repository[User]
}

// NewUserRepository creates a new user repository
func NewUserRepository(database *db.Database) *UserRepository {
    return &UserRepository{
        repo: db.NewRepository[User](database),
    }
}

// Create creates a new user
func (r *UserRepository) Create(user *User) error {
    return r.repo.Create(user)
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id uint) (*User, error) {
    var user User
    if err := r.repo.FindByID(id, &user); err != nil {
        return nil, err
    }
    return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*User, error) {
    var user User
    if err := r.repo.db.Table("users").Where("email", "=", email).First(&user); err != nil {
        return nil, err
    }
    return &user, nil
}

// FindAll finds all users
func (r *UserRepository) FindAll() ([]*User, error) {
    var users []*User
    if err := r.repo.FindAll(&users); err != nil {
        return nil, err
    }
    return users, nil
}

// Update updates a user
func (r *UserRepository) Update(user *User) error {
    return r.repo.Update(user)
}

// Delete deletes a user
func (r *UserRepository) Delete(user *User) error {
    return r.repo.Delete(user)
}`}</code>
      </pre>

      <h2>Using Models</h2>
      <p>Here's an example of how to use models and repositories in a controller:</p>

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
}`}</code>
      </pre>

      <h2>Relationships</h2>
      <p>You can define relationships between models using struct fields:</p>

      <pre>
        <code>{`// models/post.go
package models

import (
    "time"
    
    "github.com/example/goframe/db"
)

// Post represents a blog post
type Post struct {
    db.Entity
    Title    string \`db:"title" json:"title"\`
    Content  string \`db:"content" json:"content"\`
    UserID   uint   \`db:"user_id" json:"user_id"\`
    User     *User  \`db:"-" json:"user,omitempty"\`
}

// PostRepository provides methods to interact with posts
type PostRepository struct {
    repo *db.Repository[Post]
    userRepo *UserRepository
}

// NewPostRepository creates a new post repository
func NewPostRepository(database *db.Database, userRepo *UserRepository) *PostRepository {
    return &PostRepository{
        repo: db.NewRepository[Post](database),
        userRepo: userRepo,
    }
}

// FindWithUser finds a post by ID and loads the associated user
func (r *PostRepository) FindWithUser(id uint) (*Post, error) {
    post, err := r.repo.FindByID(id, nil)
    if err != nil {
        return nil, err
    }
    
    // Load the user
    user, err := r.userRepo.FindByID(post.UserID)
    if err != nil {
        return post, nil // Return post even if user can't be loaded
    }
    
    post.User = user
    return post, nil
}`}</code>
      </pre>

      <DocPagination
        prev={{
          title: "Controllers",
          href: "/docs/core-concepts/controllers",
        }}
        next={{
          title: "Middleware",
          href: "/docs/core-concepts/middleware",
        }}
      />
    </div>
  )
}

