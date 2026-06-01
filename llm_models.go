package promptrails

import (
	"context"
	"fmt"
)

// LLMModelsService handles LLM model catalog API calls.
type LLMModelsService struct {
	http *httpClient
}

// ListLLMModelsParams are parameters for listing LLM models.
type ListLLMModelsParams struct {
	Page     int
	Limit    int
	Provider string
	IsActive string
}

// List returns a paginated list of LLM models in the catalog.
func (s *LLMModelsService) List(ctx context.Context, params *ListLLMModelsParams) (*PaginatedResponse[LLMModel], error) {
	if params == nil {
		params = &ListLLMModelsParams{}
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
	if params.IsActive != "" {
		qp["is_active"] = params.IsActive
	}
	var result PaginatedResponse[LLMModel]
	err := s.http.get(ctx, "/api/v1/llm-models", qp, &result)
	return &result, err
}

// ListAvailable returns models grouped by provider, restricted to the
// providers the current workspace holds credentials for. Deprecated models
// are included with IsDeprecated set so callers can warn rather than hide.
func (s *LLMModelsService) ListAvailable(ctx context.Context) ([]AvailableModelGroup, error) {
	var result struct {
		Groups []AvailableModelGroup `json:"groups"`
	}
	err := s.http.get(ctx, "/api/v1/llm-models/available", nil, &result)
	return result.Groups, err
}
