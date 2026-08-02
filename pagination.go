package promptrails

import "encoding/json"

// PaginatedResponse wraps a paginated API response.
type PaginatedResponse[T any] struct {
	Data []T            `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// PaginationMeta holds pagination metadata. API v2 standardizes on "pages";
// older payloads used "total_pages", which UnmarshalJSON still accepts.
type PaginationMeta struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Pages int `json:"pages"`
}

// UnmarshalJSON reads the standardized "pages" field, falling back to the
// legacy "total_pages" key when "pages" is absent.
func (m *PaginationMeta) UnmarshalJSON(b []byte) error {
	type alias PaginationMeta
	aux := struct {
		*alias
		TotalPages *int `json:"total_pages"`
	}{alias: (*alias)(m)}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	if m.Pages == 0 && aux.TotalPages != nil {
		m.Pages = *aux.TotalPages
	}
	return nil
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
