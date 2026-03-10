package promptrails

import (
	"context"
	"fmt"
)

// WebhookTriggersService handles webhook trigger API calls.
type WebhookTriggersService struct {
	http *httpClient
}

// ListWebhookTriggersParams are parameters for listing webhook triggers.
type ListWebhookTriggersParams struct {
	Page    int
	Limit   int
	AgentID string
}

// CreateWebhookTriggerParams are parameters for creating a webhook trigger.
type CreateWebhookTriggerParams struct {
	Name      string `json:"name"`
	AgentID   string `json:"agent_id"`
	HasSecret bool   `json:"has_secret,omitempty"`
}

// UpdateWebhookTriggerParams are parameters for updating a webhook trigger.
type UpdateWebhookTriggerParams struct {
	Name     *string `json:"name,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// List returns a paginated list of webhook triggers.
func (s *WebhookTriggersService) List(ctx context.Context, params *ListWebhookTriggersParams) (*PaginatedResponse[WebhookTrigger], error) {
	if params == nil {
		params = &ListWebhookTriggersParams{}
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
	var result PaginatedResponse[WebhookTrigger]
	err := s.http.get(ctx, "/api/v1/webhook-triggers", qp, &result)
	return &result, err
}

// Get returns a single webhook trigger by ID.
func (s *WebhookTriggersService) Get(ctx context.Context, id string) (*WebhookTrigger, error) {
	var result WebhookTrigger
	err := s.http.get(ctx, "/api/v1/webhook-triggers/"+id, nil, &result)
	return &result, err
}

// Create creates a new webhook trigger.
func (s *WebhookTriggersService) Create(ctx context.Context, params *CreateWebhookTriggerParams) (*WebhookTriggerCreateResponse, error) {
	var result WebhookTriggerCreateResponse
	err := s.http.post(ctx, "/api/v1/webhook-triggers", params, &result)
	return &result, err
}

// Update updates an existing webhook trigger.
func (s *WebhookTriggersService) Update(ctx context.Context, id string, params *UpdateWebhookTriggerParams) (*WebhookTrigger, error) {
	var result WebhookTrigger
	err := s.http.patch(ctx, "/api/v1/webhook-triggers/"+id, params, &result)
	return &result, err
}

// Delete removes a webhook trigger.
func (s *WebhookTriggersService) Delete(ctx context.Context, id string) error {
	return s.http.del(ctx, "/api/v1/webhook-triggers/"+id)
}
