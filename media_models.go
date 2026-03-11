package promptrails

import (
	"context"
	"fmt"
)

// MediaModelsService handles media model-related API calls.
type MediaModelsService struct {
	http *httpClient
}

// ListMediaModelsParams are parameters for listing media models.
type ListMediaModelsParams struct {
	Page      int
	Limit     int
	Provider  string
	MediaType string
	IsActive  string
}

// List returns a paginated list of media models.
func (s *MediaModelsService) List(ctx context.Context, params *ListMediaModelsParams) (*PaginatedResponse[MediaModel], error) {
	if params == nil {
		params = &ListMediaModelsParams{}
	}
	p := &ListParams{Page: params.Page, Limit: params.Limit}
	p.defaults()
	qp := map[string]string{
		"page":  fmt.Sprintf("%d", p.Page),
		"limit": fmt.Sprintf("%d", p.Limit),
	}
	if params.Provider != "" {
		qp["provider"] = params.Provider
	}
	if params.MediaType != "" {
		qp["media_type"] = params.MediaType
	}
	if params.IsActive != "" {
		qp["is_active"] = params.IsActive
	}
	var result PaginatedResponse[MediaModel]
	err := s.http.get(ctx, "/api/v1/media-models", qp, &result)
	return &result, err
}
