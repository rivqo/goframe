import { DocPagination } from "@/components/doc-pagination"

export default function ResourcesPage() {
  return (
    <div className="space-y-6">
      <h1>Resources</h1>
      <p>
        Resources in GoFrame provide a way to transform your models into API responses. They are similar to Laravel's
        API Resources and help you build a consistent API layer.
      </p>
      <h2>Creating Resources</h2>
      <p>
        To create a new resource, use the <code>make:resource</code> command:
      </p>
      <pre>
        <code>{`goframe make:resource User`}</code>
      </pre>
      <p>
        This will create a new resource file in the <code>resources</code> directory:
      </p>
      <pre>
        <code>{`// resources/user_resource.go
package resources

import (
	"github.com/example/goframe/models"
)

// UserResource represents a User resource
type UserResource struct {
	ID   uint   \`json:"id"\`
	Name string \`json:"name"\`
	Email string \`json:"email"\`
	// Add more fields as needed
}

// UserCollection represents a collection of User resources
type UserCollection struct {
	Data []*UserResource \`json:"data"\`
}

// NewUserResource creates a new User resource from a model
func NewUserResource(model *models.User) *UserResource {
	return &UserResource{
		ID:   model.ID,
		Name: model.Name,
		Email: model.Email,
		// Map more fields as needed
	}
}

// NewUserCollection creates a new collection of User resources
func NewUserCollection(models []*models.User) *UserCollection {
	resources := make([]*UserResource, len(models))
	for i, model := range models {
		resources[i] = NewUserResource(model)
	}
	return &UserCollection{
		Data: resources,
	}
}
`}</code>
      </pre>
      <h2>Using Resources</h2>
      <p>You can use resources in your controllers to transform your models into API responses:</p>
      <pre>
        <code>{`// controllers/user_controller.go
package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/example/goframe/models"
	"github.com/example/goframe/resources"
)

// UserController handles user-related HTTP requests
type UserController struct {
	repo *models.UserRepository
}

// Index returns all users
func (c *UserController) Index(w http.ResponseWriter, r *http.Request) {
	users, err := c.repo.FindAll()
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	resource := resources.NewUserCollection(users)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resource)
}

// Show returns a single user
func (c *UserController) Show(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, err := c.repo.FindByID(uint(id))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	resource := resources.NewUserResource(user)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resource)
}
`}</code>
      </pre>
      <h2>Resource Collections</h2>
      <p>Resource collections allow you to transform a collection of models into a consistent API response:</p>
      <pre>
        <code>{`users, _ := userRepo.FindAll()
collection := resources.NewUserCollection(users)
json.NewEncoder(w).Encode(collection)

// Output:
// {
//   "data": [
//     {"id": 1, "name": "John Doe", "email": "john@example.com"},
//     {"id": 2, "name": "Jane Smith", "email": "jane@example.com"}
//   ]
// }
`}</code>
      </pre>
      <h2>Pagination</h2>
      <p>
        GoFrame also provides a <code>Paginator</code> for paginated collections:
      </p>
      <pre>
        <code>{`// resources/resource.go
package resources

// Paginator represents a paginated collection of resources
type Paginator struct {
	Data       []map[string]interface{} \`json:"data"\`
	Total      int                      \`json:"total"\`
	PerPage    int                      \`json:"per_page"\`
	CurrentPage int                     \`json:"current_page"\`
	LastPage   int                      \`json:"last_page"\`
	From       int                      \`json:"from"\`
	To         int                      \`json:"to"\`
}

// NewPaginator creates a new paginator
func NewPaginator(data []map[string]interface{}, total, perPage, currentPage int) *Paginator {
	lastPage := total / perPage
	if total % perPage > 0 {
		lastPage++
	}

	from := (currentPage - 1) * perPage + 1
	to := from + len(data) - 1

	if from > total {
		from = 0
	}

	if to > total {
		to = total
	}

	return &Paginator{
		Data:       data,
		Total:      total,
		PerPage:    perPage,
		CurrentPage: currentPage,
		LastPage:   lastPage,
		From:       from,
		To:         to,
	}
}
`}</code>
      </pre>
      <p>You can use the paginator in your controllers:</p>
      <pre>
        <code>{`// Get pagination parameters
page, _ := strconv.Atoi(r.URL.Query().Get("page"))
if page < 1 {
	page = 1
}
perPage := 10

// Get total count
total, _ := c.repo.Count()

// Get paginated users
users, _ := c.repo.Paginate(page, perPage)

// Transform users to resources
resources := make([]map[string]interface{}, len(users))
for i, user := range users {
	resource := resources.NewUserResource(user)
	resources[i] = resource.ToMap()
}

// Create paginator
paginator := resources.NewPaginator(resources, total, perPage, page)

// Return response
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(paginator)

// Output:
// {
//   "data": [
//     {"id": 1, "name": "John Doe", "email": "john@example.com"},
//     {"id": 2, "name": "Jane Smith", "email": "jane@example.com"}
//   ],
//   "total": 25,
//   "per_page": 10,
//   "current_page": 1,
//   "last_page": 3,
//   "from": 1,
//   "to": 10
// }
`}</code>
      </pre>
      <DocPagination
        prev={{
          title: "Authentication",
          href: "/docs/features/authentication",
        }}
        next={{
          title: "Rate Limiting",
          href: "/docs/features/rate-limiting",
        }}
      />
    </div>
  )
}

