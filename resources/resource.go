package resources

// Resource represents a generic API resource
type Resource interface {
	ToMap() map[string]interface{}
}

// Collection represents a collection of resources
type Collection interface {
	ToSlice() []map[string]interface{}
}

// Paginator represents a paginated collection of resources
type Paginator struct {
	Data       []map[string]interface{} `json:"data"`
	Total      int                      `json:"total"`
	PerPage    int                      `json:"per_page"`
	CurrentPage int                     `json:"current_page"`
	LastPage   int                      `json:"last_page"`
	From       int                      `json:"from"`
	To         int                      `json:"to"`
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

