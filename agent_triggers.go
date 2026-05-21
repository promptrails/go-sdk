package promptrails

import (
	"context"
	"fmt"
)

// AgentTriggersService handles agent trigger API calls.
//
// One trigger record, six sources: generic webhook, Slack, Telegram, Microsoft
// Teams, WhatsApp Business, or a cron schedule. Pass the matching SourceConfig
// fields when calling Create.
type AgentTriggersService struct {
	http *httpClient
}

// ListAgentTriggersParams are parameters for listing agent triggers.
type ListAgentTriggersParams struct {
	Page    int
	Limit   int
	AgentID string
}

// CreateAgentTriggerParams are parameters for creating an agent trigger.
type CreateAgentTriggerParams struct {
	Name           string                 `json:"name"`
	AgentID        string                 `json:"agent_id"`
	Source         AgentTriggerSource     `json:"source,omitempty"`
	SourceConfig   map[string]interface{} `json:"source_config,omitempty"`
	ReplyConfig    map[string]interface{} `json:"reply_config,omitempty"`
	GenerateSecret bool                   `json:"generate_secret,omitempty"`
}

// UpdateAgentTriggerParams are parameters for updating an agent trigger.
type UpdateAgentTriggerParams struct {
	Name         *string                `json:"name,omitempty"`
	IsActive     *bool                  `json:"is_active,omitempty"`
	Source       *AgentTriggerSource    `json:"source,omitempty"`
	SourceConfig map[string]interface{} `json:"source_config,omitempty"`
	ReplyConfig  map[string]interface{} `json:"reply_config,omitempty"`
}

// List returns a paginated list of agent triggers.
func (s *AgentTriggersService) List(ctx context.Context, params *ListAgentTriggersParams) (*PaginatedResponse[AgentTrigger], error) {
	if params == nil {
		params = &ListAgentTriggersParams{}
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
	var result PaginatedResponse[AgentTrigger]
	err := s.http.get(ctx, "/api/v1/triggers", qp, &result)
	return &result, err
}

// Get returns a single agent trigger by ID.
func (s *AgentTriggersService) Get(ctx context.Context, id string) (*AgentTrigger, error) {
	var result AgentTrigger
	err := s.http.get(ctx, "/api/v1/triggers/"+id, nil, &result)
	return &result, err
}

// Create creates a new agent trigger.
func (s *AgentTriggersService) Create(ctx context.Context, params *CreateAgentTriggerParams) (*AgentTriggerCreateResponse, error) {
	var result AgentTriggerCreateResponse
	err := s.http.post(ctx, "/api/v1/triggers", params, &result)
	return &result, err
}

// Update updates an existing agent trigger.
func (s *AgentTriggersService) Update(ctx context.Context, id string, params *UpdateAgentTriggerParams) (*AgentTrigger, error) {
	var result AgentTrigger
	err := s.http.patch(ctx, "/api/v1/triggers/"+id, params, &result)
	return &result, err
}

// Delete removes an agent trigger.
func (s *AgentTriggersService) Delete(ctx context.Context, id string) error {
	return s.http.del(ctx, "/api/v1/triggers/"+id)
}
