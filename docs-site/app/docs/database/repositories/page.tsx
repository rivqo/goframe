import { DocPagination } from "@/components/doc-pagination"

export default function RepositoriesPage() {
  return (
    <div className="space-y-6">
      <h1>Repositories</h1>
      <p>
        GoFrame uses the repository pattern to encapsulate database operations. Repositories provide a clean,
        object-oriented interface for interacting with your database.
      </p>

      <h2>Basic Repository</h2>
      <p>
        The <code>db.Repository</code> struct provides generic methods for common database operations. It uses Go's
        generics to provide type-safe operations.
      </p>

      <pre>
        <code>{`// db/orm.go
// Repository provides a generic way to interact with entities
type Repository[T any] struct {
    db *Database
}

// NewRepository creates a new repository for the given entity type
func NewRepository[T any](db *Database) *Repository[T] {
    return &Repository[T]{db: db}
}

// Create creates a new entity
func (r *Repository[T]) Create(entity *T) error {
    // Set created_at and updated_at
    setTimestamps(entity)

    return r.db.Create(entity)
}

// FindByID finds an entity by ID
func (r *Repository[T]) FindByID(id uint, dest *T) error {
    return r.db.Table(dest).Where("id = ?", id).First(dest)
}

// FindAll finds all entities matching the query
func (r *Repository[T]) FindAll(dest *[]T, conditions ...interface{}) error {
    query := r.db.Table(*new(T))

    if len(conditions) > 0 {
        if condition, ok := conditions[0].(string); ok {
            query = query.Where(condition, conditions[1:]...)
        }
    }

    return query.Find(dest)
}

// Update updates an entity
func (r *Repository[T]) Update(entity *T) error {
    // Set updated_at
    setUpdatedAt(entity)

    values := toMap(entity)
    delete(values, "id")       // Don't update ID
    delete(values, "created_at") // Don't update created_at

    return r.db.Table(entity).Where("id = ?", getID(entity)).Update(values)
}

// Delete deletes an entity
func (r *Repository[T]) Delete(entity *T) error {
    return r.db.Table(entity).Where("id = ?", getID(entity)).Delete()
}`}</code>
      </pre>

      <h2>Creating a Repository</h2>
      <p>
        To create a repository for a specific model, you typically define a struct that wraps the generic repository and
        provides model-specific methods:
      </p>

      <pre>
        <code>{`// models/user.go
package models

import (
    "github.com/example/goframe/db"
)

// User represents a user in the database
type User struct {
    db.Entity
    Name     string \`db:"name" json:"name"\`
    Email    string \`db:"email" json:"email"\`
    Password string \`db:"password" json:"-"\`
}

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

      <h2>Using Repositories</h2>
      <p>Repositories are typically used in controllers or services:</p>

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

      <h2>Advanced Repository Patterns</h2>

      <h3>Relationships</h3>
      <p>You can implement methods to handle relationships between models:</p>

      <pre>
        <code>{`// models/post.go
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
}

// FindByUser finds all posts by a user
func (r *PostRepository) FindByUser(userID uint) ([]*Post, error) {
    var posts []*Post
    err := r.repo.db.Table("posts").
        Where("user_id", "=", userID).
        OrderBy("created_at", "desc").
        Get(&posts)
    return posts, err
}`}</code>
      </pre>

      <h3>Pagination</h3>
      <p>You can implement pagination methods to handle large result sets:</p>

      <pre>
        <code>{`// models/user.go
// Paginate returns a paginated list of users
func (r *UserRepository) Paginate(page, perPage int) ([]*User, error) {
    var users []*User
    err := r.repo.db.Table("users").
        Limit(perPage).
        Offset((page - 1) * perPage).
        OrderBy("created_at", "desc").
        Get(&users)
    return users, err
}

// Count returns the total number of users
func (r *UserRepository) Count() (int, error) {
    return r.repo.db.Table("users").Count()
}`}</code>
      </pre>

      <h3>Custom Queries</h3>
      <p>You can implement methods for custom queries:</p>

      <pre>
        <code>{`// models/user.go
// FindActive finds all active users
func (r *UserRepository) FindActive() ([]*User, error) {
    var users []*User
    err := r.repo.db.Table("users").
        Where("active", "=", true).
        OrderBy("name", "asc").
        Get(&users)
    return users, err  "=", true).
        OrderBy("name", "asc").
        Get(&users)
    return users, err
}

// FindByRole finds all users with a specific role
func (r *UserRepository) FindByRole(role string) ([]*User, error) {
    var users []*User
    err := r.repo.db.Table("users").
        Where("role", "=", role).
        OrderBy("name", "asc").
        Get(&users)
    return users, err
}

// Search searches for users by name or email
func (r *UserRepository) Search(query string) ([]*User, error) {
    var users []*User
    searchQuery := "%" + query + "%"
    err := r.repo.db.Table("users").
        WhereRaw("(name LIKE ? OR email LIKE ?)", searchQuery, searchQuery).
        OrderBy("name", "asc").
        Get(&users)
    return users, err
}`}</code>
      </pre>

      <DocPagination
        prev={{
          title: "Query Builder",
          href: "/docs/database/query-builder",
        }}
        next={{
          title: "Features",
          href: "/docs/features",
        }}
      />
    </div>
  )
}

