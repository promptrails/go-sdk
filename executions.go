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

// DecideParams are parameters for approving or denying a parked execution.
type DecideParams struct {
	Reason string `json:"reason,omitempty"`
}

// Tree fetches the execution with its full Children tree populated.
func (s *ExecutionsService) Tree(ctx context.Context, id string) (*Execution, error) {
	var result Execution
	err := s.http.get(ctx, "/api/v1/executions/"+id+"/tree", nil, &result)
	return &result, err
}

// Cancel requests cooperative cancellation of a running execution.
func (s *ExecutionsService) Cancel(ctx context.Context, id string) (*Execution, error) {
	var result Execution
	err := s.http.post(ctx, "/api/v1/executions/"+id+"/cancel", map[string]any{}, &result)
	return &result, err
}

// ApprovalInbox lists executions parked at "waiting_approval".
func (s *ExecutionsService) ApprovalInbox(ctx context.Context, params *ListParams) (*PaginatedResponse[Execution], error) {
	if params == nil {
		params = &ListParams{}
	}
	params.defaults()
	qp := map[string]string{
		"page":  fmt.Sprintf("%d", params.Page),
		"limit": fmt.Sprintf("%d", params.Limit),
	}
	var result PaginatedResponse[Execution]
	err := s.http.get(ctx, "/api/v1/executions/approval-inbox", qp, &result)
	return &result, err
}

// Approve approves a run parked at "waiting_approval" and resumes it.
func (s *ExecutionsService) Approve(ctx context.Context, id string, params *DecideParams) (*Execution, error) {
	if params == nil {
		params = &DecideParams{}
	}
	var result Execution
	err := s.http.post(ctx, "/api/v1/executions/"+id+"/approve", params, &result)
	return &result, err
}

// Deny denies a run parked at "waiting_approval" and resumes with a denial.
func (s *ExecutionsService) Deny(ctx context.Context, id string, params *DecideParams) (*Execution, error) {
	if params == nil {
		params = &DecideParams{}
	}
	var result Execution
	err := s.http.post(ctx, "/api/v1/executions/"+id+"/deny", params, &result)
	return &result, err
}

// Stream subscribes to the live SSE event stream for an execution. Useful
// when the execution was started outside a chat (e.g. Agents.Execute) and
// the caller wants progressive updates. Always Close() the returned stream.
func (s *ExecutionsService) Stream(ctx context.Context, executionID string) (*ChatStream, error) {
	body, err := s.http.stream(ctx, "GET", "/api/v1/executions/"+executionID+"/stream", nil)
	if err != nil {
		return nil, err
	}
	return newChatStream(body), nil
}
