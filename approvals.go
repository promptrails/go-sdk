package promptrails

import (
	"context"
	"fmt"
)

// ApprovalsService handles approval-related API calls.
type ApprovalsService struct {
	http *httpClient
}

// ListApprovalsParams are parameters for listing approvals.
type ListApprovalsParams struct {
	Page   int
	Limit  int
	Status string
}

// DecideApprovalParams are parameters for deciding on an approval.
type DecideApprovalParams struct {
	Decision string `json:"decision"` // "approved" or "rejected"
}

// List returns a paginated list of approval requests.
func (s *ApprovalsService) List(ctx context.Context, params *ListApprovalsParams) (*PaginatedResponse[ApprovalRequest], error) {
	if params == nil {
		params = &ListApprovalsParams{}
	}
	p := &ListParams{Page: params.Page, Limit: params.Limit}
	p.defaults()
	qp := map[string]string{
		"page":  fmt.Sprintf("%d", p.Page),
		"limit": fmt.Sprintf("%d", p.Limit),
	}
	if params.Status != "" {
		qp["status"] = params.Status
	}
	var result PaginatedResponse[ApprovalRequest]
	err := s.http.get(ctx, "/api/v1/approvals", qp, &result)
	return &result, err
}

// Get returns a single approval request by ID.
func (s *ApprovalsService) Get(ctx context.Context, id string) (*ApprovalRequest, error) {
	var result ApprovalRequest
	err := s.http.get(ctx, "/api/v1/approvals/"+id, nil, &result)
	return &result, err
}

// Decide approves or rejects an approval request.
func (s *ApprovalsService) Decide(ctx context.Context, id string, params *DecideApprovalParams) (*ApprovalRequest, error) {
	var result ApprovalRequest
	err := s.http.post(ctx, "/api/v1/approvals/"+id+"/decide", params, &result)
	return &result, err
}
