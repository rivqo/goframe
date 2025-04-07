package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MakeModel creates a new model file
func MakeModel(name string, migration bool) {
	// Create models directory if it doesn't exist
	if err := os.MkdirAll("models", 0755); err != nil {
		fmt.Printf("Failed to create models directory: %v\n", err)
		os.Exit(1)
	}

	// Create model file
	modelName := strings.Title(name)
	tableName := toSnakeCase(name) + "s"
	path := filepath.Join("models", toSnakeCase(name)+".go")

	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Failed to create model file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Write model template
	template := fmt.Sprintf(`package models

import (
	"github.com/example/goframe/db"
)

// %s represents the %s model
type %s struct {
	db.Entity
	// Add your fields here
	Name string ` + "`db:\"name\" json:\"name\"`" + `
}

// TableName returns the table name for the model
func (%s) TableName() string {
	return "%s"
}

// %sRepository provides methods to interact with %s
type %sRepository struct {
	repo *db.Repository[%s]
}

// New%sRepository creates a new %s repository
func New%sRepository(database *db.Database) *%sRepository {
	return &%sRepository{
		repo: db.NewRepository[%s](database),
	}
}

// Create creates a new %s
func (r *%sRepository) Create(model *%s) error {
	return r.repo.Create(model)
}

// FindByID finds a %s by ID
func (r *%sRepository) FindByID(id uint) (*%s, error) {
	var model %s
	if err := r.repo.FindByID(id, &model); err != nil {
		return nil, err
	}
	return &model, nil
}

// FindAll finds all %ss
func (r *%sRepository) FindAll() ([]*%s, error) {
	var models []*%s
	if err := r.repo.FindAll(&models); err != nil {
		return nil, err
	}
	return models, nil
}

// Update updates a %s
func (r *%sRepository) Update(model *%s) error {
	return r.repo.Update(model)
}

// Delete deletes a %s
func (r *%sRepository) Delete(model *%s) error {
	return r.repo.Delete(model)
}
`, modelName, modelName, modelName, modelName, tableName, 
   modelName, modelName, modelName, modelName, 
   modelName, modelName, modelName, modelName, modelName, modelName,
   modelName, modelName, modelName, 
   modelName, modelName, modelName, modelName,
   modelName, modelName, modelName, modelName,
   modelName, modelName, modelName,
   modelName, modelName, modelName)

	if _, err := file.WriteString(template); err != nil {
		fmt.Printf("Failed to write model file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created model: %s\n", path)

	// Create migration if requested
	if migration {
		migrationName := fmt.Sprintf("create_%s_table", tableName)
		MakeMigration(migrationName)
	}
}

// MakeController creates a new controller file
func MakeController(name string, resource bool) {
	// Create controllers directory if it doesn't exist
	if err := os.MkdirAll("controllers", 0755); err != nil {
		fmt.Printf("Failed to create controllers directory: %v\n", err)
		os.Exit(1)
	}

	// Create controller file
	controllerName := strings.Title(name) + "Controller"
	path := filepath.Join("controllers", toSnakeCase(name)+"_controller.go")

	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Failed to create controller file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Write controller template
	var template string
	if resource {
		modelName := strings.TrimSuffix(strings.Title(name), "Controller")
		template = fmt.Sprintf(`package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/example/goframe/models"
	"github.com/example/goframe/resources"
)

// %s handles %s-related HTTP requests
type %s struct {
	repo *models.%sRepository
}

// New%s creates a new %s controller
func New%s(repo *models.%sRepository) *%s {
	return &%s{
		repo: repo,
	}
}

// Index returns all %ss
func (c *%s) Index(w http.ResponseWriter, r *http.Request) {
	models, err := c.repo.FindAll()
	if err != nil {
		http.Error(w, "Failed to fetch %ss", http.StatusInternalServerError)
		return
	}

	resource := resources.New%sCollection(models)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resource)
}

// Show returns a single %s
func (c *%s) Show(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	model, err := c.repo.FindByID(uint(id))
	if err != nil {
		http.Error(w, "%s not found", http.StatusNotFound)
		return
	}

	resource := resources.New%sResource(model)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resource)
}

// Store creates a new %s
func (c *%s) Store(w http.ResponseWriter, r *http.Request) {
	var model models.%s
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := c.repo.Create(&model); err != nil {
		http.Error(w, "Failed to create %s", http.StatusInternalServerError)
		return
	}

	resource := resources.New%sResource(&model)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resource)
}

// Update updates a %s
func (c *%s) Update(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	model, err := c.repo.FindByID(uint(id))
	if err != nil {
		http.Error(w, "%s not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(model); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := c.repo.Update(model); err != nil {
		http.Error(w, "Failed to update %s", http.StatusInternalServerError)
		return
	}

	resource := resources.New%sResource(model)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resource)
}

// Destroy deletes a %s
func (c *%s) Destroy(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	model, err := c.repo.FindByID(uint(id))
	if err != nil {
		http.Error(w, "%s not found", http.StatusNotFound)
		return
	}

	if err := c.repo.Delete(model); err != nil {
		http.Error(w, "Failed to delete %s", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
`, controllerName, strings.ToLower(modelName), controllerName, modelName, 
   controllerName, controllerName, controllerName, modelName, controllerName, controllerName,
   strings.ToLower(modelName), controllerName, strings.ToLower(modelName), modelName,
   modelName, controllerName, modelName, modelName,
   modelName, controllerName, modelName,
   modelName, controllerName, strings.ToLower(modelName), modelName,
   modelName, controllerName, modelName, strings.ToLower(modelName),
   modelName, controllerName, strings.ToLower(modelName))
	} else {
		template = fmt.Sprintf(`package controllers

import (
	"net/http"
)

// %s handles HTTP requests
type %s struct {
	// Add your dependencies here
}

// New%s creates a new %s controller
func New%s() *%s {
	return &%s{
		// Initialize your dependencies here
	}
}

// Index handles the index route
func (c *%s) Index(w http.ResponseWriter, r *http.Request) {
	// Implement your logic here
	w.Write([]byte("Hello from %s"))
}
`, controllerName, controllerName, controllerName, controllerName, controllerName, controllerName, controllerName, controllerName, controllerName)
	}

	if _, err := file.WriteString(template); err != nil {
		fmt.Printf("Failed to write controller file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created controller: %s\n", path)
}

// MakeResource creates a new resource file
func MakeResource(name string) {
	// Create resources directory if it doesn't exist
	if err := os.MkdirAll("resources", 0755); err != nil {
		fmt.Printf("Failed to create resources directory: %v\n", err)
		os.Exit(1)
	}

	// Create resource file
	resourceName := strings.Title(name) + "Resource"
	path := filepath.Join("resources", toSnakeCase(name)+"_resource.go")

	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Failed to create resource file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Write resource template
	modelName := strings.Title(name)
	template := fmt.Sprintf(`package resources

import (
	"github.com/example/goframe/models"
)

// %s represents a %s resource
type %s struct {
	ID   uint   ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
	// Add more fields as needed
}

// %sCollection represents a collection of %s resources
type %sCollection struct {
	Data []*%s ` + "`json:\"data\"`" + `
}

// New%sResource creates a new %s resource from a model
func New%sResource(model *models.%s) *%s {
	return &%s{
		ID:   model.ID,
		Name: model.Name,
		// Map more fields as needed
	}
}

// New%sCollection creates a new collection of %s resources
func New%sCollection(models []*models.%s) *%sCollection {
	resources := make([]*%s, len(models))
	for i, model := range models {
		resources[i] = New%sResource(model)
	}
	return &%sCollection{
		Data: resources,
	}
}
`, resourceName, modelName, resourceName, 
   resourceName, modelName, resourceName, resourceName,
   resourceName, modelName, resourceName, modelName, resourceName, resourceName,
   resourceName, modelName, resourceName, modelName, resourceName, resourceName, resourceName, resourceName)

	if _, err := file.WriteString(template); err != nil {
		fmt.Printf("Failed to write resource file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created resource: %s\n", path)
}

// toSnakeCase converts a string to snake_case
func toSnakeCase(s string) string {
	s = strings.TrimSuffix(s, "Controller")
	var result strings.Builder
	for i, r := range s {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

