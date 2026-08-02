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

// TraceFilterParams are query filters for trace summary / PII-report queries.
// All fields are optional; empty values are omitted from the query.
type TraceFilterParams struct {
	DateFrom    string
	DateTo      string
	Status      string
	Level       string
	Kind        string
	ModelName   string
	AgentID     string
	SessionID   string
	ExecutionID string
}

func (p *TraceFilterParams) query() map[string]string {
	qp := map[string]string{}
	if p == nil {
		return qp
	}
	for k, v := range map[string]string{
		"date_from":    p.DateFrom,
		"date_to":      p.DateTo,
		"status":       p.Status,
		"level":        p.Level,
		"kind":         p.Kind,
		"model_name":   p.ModelName,
		"agent_id":     p.AgentID,
		"session_id":   p.SessionID,
		"execution_id": p.ExecutionID,
	} {
		if v != "" {
			qp[k] = v
		}
	}
	return qp
}

// GetSummary returns aggregate statistics over a filtered set of traces.
func (s *TracesService) GetSummary(ctx context.Context, params *TraceFilterParams) (*TraceSummary, error) {
	var result TraceSummary
	err := s.http.get(ctx, "/api/v1/traces/summary", params.query(), &result)
	return &result, err
}

// PIIReport returns the PII-masking report over a filtered set of traces.
func (s *TracesService) PIIReport(ctx context.Context, params *TraceFilterParams) (map[string]any, error) {
	var result map[string]any
	err := s.http.get(ctx, "/api/v1/traces/pii-report", params.query(), &result)
	return result, err
}

// Ingest ingests up to 1000 raw spans in one request.
func (s *TracesService) Ingest(ctx context.Context, spans []map[string]any) (map[string]any, error) {
	var result map[string]any
	err := s.http.post(ctx, "/api/v1/traces/ingest", map[string]any{"spans": spans}, &result)
	return result, err
}
