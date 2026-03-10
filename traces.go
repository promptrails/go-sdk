package promptrails

import (
	"context"
	"fmt"
)

// TracesService handles trace-related API calls.
type TracesService struct {
	http *httpClient
}

// ListTracesParams are parameters for listing traces.
type ListTracesParams struct {
	Page    int
	Limit   int
	TraceID string
	Kind    string
}

// List returns a paginated list of traces.
func (s *TracesService) List(ctx context.Context, params *ListTracesParams) (*PaginatedResponse[Trace], error) {
	if params == nil {
		params = &ListTracesParams{}
	}
	p := &ListParams{Page: params.Page, Limit: params.Limit}
	p.defaults()
	qp := map[string]string{
		"page":  fmt.Sprintf("%d", p.Page),
		"limit": fmt.Sprintf("%d", p.Limit),
	}
	if params.TraceID != "" {
		qp["trace_id"] = params.TraceID
	}
	if params.Kind != "" {
		qp["kind"] = params.Kind
	}
	var result PaginatedResponse[Trace]
	err := s.http.get(ctx, "/api/v1/traces", qp, &result)
	return &result, err
}

// GetByTraceID returns traces grouped by trace ID.
func (s *TracesService) GetByTraceID(ctx context.Context, traceID string) ([]Trace, error) {
	var result []Trace
	err := s.http.get(ctx, "/api/v1/traces/"+traceID, nil, &result)
	return result, err
}
