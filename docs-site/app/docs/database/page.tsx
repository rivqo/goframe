import { DocPagination } from "@/components/doc-pagination"

export default function DatabasePage() {
  return (
    <div className="space-y-6">
      <h1>Database</h1>
      <p>
        GoFrame provides a simple and intuitive database layer that makes it easy to interact with your database. The
        database layer is built on top of Go's standard database/sql package and provides a fluent query builder,
        migrations, and an ORM-like repository pattern.
      </p>

      <h2>Configuration</h2>
      <p>
        Database configuration is stored in the <code>config.yaml</code> file:
      </p>
      <pre>
        <code>{`database:
  driver: postgres
  host: localhost
  port: 5432
  name: goframe
  user: postgres
  password: postgres`}</code>
      </pre>

      <h2>Connecting to the Database</h2>
      <p>GoFrame handles database connections for you. The connection is established when your application starts:</p>
      <pre>
        <code>{`// main.go
// Initialize database
database, err := db.Connect(cfg.Database)
if err != nil {
    log.Fatalf("Failed to connect to database: %v", err)
}
defer database.Close()`}</code>
      </pre>

      <h2>Query Builder</h2>
      <p>GoFrame provides a fluent query builder that makes it easy to build SQL queries:</p>
      <pre>
        <code>{`// Find a user by ID
user := &models.User{}
err := db.Table(user).Where("id = ?", id).First(user)

// Find all active users
var users []models.User
err := db.Table(models.User{}).Where("active = ?", true).Find(&users)

// Count users
var count int
err := db.Table(models.User{}).Count(&count)

// Insert a new user
user := &models.User{
    Name:  "John Doe",
    Email: "john@example.com",
}
err := db.Create(user)

// Update a user
user.Name = "Jane Doe"
err := db.Table(user).Where("id = ?", user.ID).Update(map[string]interface{}{
    "name": user.Name,
})

// Delete a user
err := db.Table(user).Where("id = ?", user.ID).Delete()`}</code>
      </pre>

      <h2>Repositories</h2>
      <p>GoFrame encourages the use of the repository pattern to encapsulate database logic:</p>
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
    if err := r.repo.db.Table(user).Where("email = ?", email).First(&user); err != nil {
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

      <h2>Transactions</h2>
      <p>GoFrame supports database transactions:</p>
      <pre>
        <code>{`// Start a transaction
tx, err := db.Begin()
if err != nil {
    return err
}

// Defer rollback in case of error
defer tx.Rollback()

// Perform operations within the transaction
if err := tx.Create(user); err != nil {
    return err
}

if err := tx.Create(post); err != nil {
    return err
}

// Commit the transaction
return tx.Commit()`}</code>
      </pre>

      <h2>Migrations</h2>
      <p>GoFrame provides a powerful migration system for managing your database schema:</p>
      <pre>
        <code>{`// Create a migration
goframe make:migration create_users_table

// Run migrations
goframe migrate

// Rollback the last batch of migrations
goframe migrate --rollback

// Rollback all migrations
goframe migrate --reset

// Rollback and re-run all migrations
goframe migrate --refresh`}</code>
      </pre>

      <p>
        Learn more about migrations in the <a href="/docs/database/migrations">Migrations</a> section.
      </p>

      <DocPagination
        prev={{
          title: "Project Structure",
          href: "/docs/project-structure",
        }}
        next={{
          title: "Migrations",
          href: "/docs/database/migrations",
        }}
      />
    </div>
  )
}

