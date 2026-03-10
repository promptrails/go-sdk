package promptrails

import (
	"context"
	"fmt"
)

// ExecutionsService handles execution-related API calls.
type ExecutionsService struct {
	http *httpClient
}

// ListExecutionsParams are parameters for listing executions.
type ListExecutionsParams struct {
	Page      int
	Limit     int
	AgentID   string
	SessionID string
	Status    string
}

// List returns a paginated list of executions.
func (s *ExecutionsService) List(ctx context.Context, params *ListExecutionsParams) (*PaginatedResponse[Execution], error) {
	if params == nil {
		params = &ListExecutionsParams{}
	}
	p := &ListParams{Page: params.Page, Limit: params.Limit}
	p.defaults()
	qp := map[string]string{
		"page":  fmt.Sprintf("%d", p.Page),
		"limit": fmt.Sprintf("%d", p.Limit),
	}
	if params.AgentID != "" {
		qp["agent_id"] = params.AgentID
	}
	if params.SessionID != "" {
		qp["session_id"] = params.SessionID
	}
	if params.Status != "" {
		qp["status"] = params.Status
	}
	var result PaginatedResponse[Execution]
	err := s.http.get(ctx, "/api/v1/executions", qp, &result)
	return &result, err
}

// Get returns a single execution by ID.
func (s *ExecutionsService) Get(ctx context.Context, id string) (*Execution, error) {
	var result Execution
	err := s.http.get(ctx, "/api/v1/executions/"+id, nil, &result)
	return &result, err
}
