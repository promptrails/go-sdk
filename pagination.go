package promptrails

// PaginatedResponse wraps a paginated API response.
type PaginatedResponse[T any] struct {
	Data []T            `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// PaginationMeta holds pagination metadata.
type PaginationMeta struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
}

// ListParams are common pagination parameters.
type ListParams struct {
	Page  int
	Limit int
}

func (p *ListParams) defaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 20
	}
}
