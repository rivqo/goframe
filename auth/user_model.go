package auth

import (
	// "time"

	"github.com/example/goframe/db"
	"golang.org/x/crypto/bcrypt"
)

// UserModel represents a user in the database
type UserModel struct {
	db.Entity
	Name         string `db:"name" json:"name"`
	Email        string `db:"email" json:"email"`
	PasswordHash string `db:"password_hash" json:"-"`
}

// UserRepository provides methods to interact with users
type UserRepository struct {
	repo *db.Repository[UserModel]
}

// NewUserRepository creates a new user repository
func NewUserRepository(database *db.Database) *UserRepository {
	return &UserRepository{
		repo: db.NewRepository[UserModel](database),
	}
}

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

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*UserModel, error) {
    var user UserModel
    err := r.repo.FindByString("email", &user, email)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id uint) (*UserModel, error) {
	var user UserModel
	if err := r.repo.FindByID(id, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// CheckPassword checks if a password is valid for a user
func (r *UserRepository) CheckPassword(user *UserModel, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}

