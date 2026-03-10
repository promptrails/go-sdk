package promptrails

import (
	"context"
	"fmt"
)

// AgentsService handles agent-related API calls.
type AgentsService struct {
	http *httpClient
}

// ListAgentsParams are parameters for listing agents.
type ListAgentsParams struct {
	Page   int
	Limit  int
	Type   string
	Status string
	Search string
}

// CreateAgentParams are parameters for creating an agent.
type CreateAgentParams struct {
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Description string         `json:"description,omitempty"`
	Config      map[string]any `json:"config,omitempty"`
}

// UpdateAgentParams are parameters for updating an agent.
type UpdateAgentParams struct {
	Name        *string        `json:"name,omitempty"`
	Description *string        `json:"description,omitempty"`
	Status      *string        `json:"status,omitempty"`
	Config      map[string]any `json:"config,omitempty"`
}

// ExecuteAgentParams are parameters for executing an agent.
type ExecuteAgentParams struct {
	Input     map[string]any `json:"input"`
	SessionID string         `json:"session_id,omitempty"`
	Sync      bool           `json:"sync,omitempty"`
}

// CreateVersionParams are parameters for creating an agent version.
type CreateVersionParams struct {
	Message string         `json:"message,omitempty"`
	Config  map[string]any `json:"config,omitempty"`
}

// CreateGuardrailParams are parameters for creating a guardrail.
type CreateGuardrailParams struct {
	Name   string         `json:"name"`
	Type   string         `json:"type"`
	Config map[string]any `json:"config,omitempty"`
}

// CreateMemoryParams are parameters for creating a memory.
type CreateMemoryParams struct {
	Content  string         `json:"content"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// SearchMemoriesParams are parameters for searching memories.
type SearchMemoriesParams struct {
	Query string `json:"query"`
	Limit int    `json:"limit,omitempty"`
}

// List returns a paginated list of agents.
func (s *AgentsService) List(ctx context.Context, params *ListAgentsParams) (*PaginatedResponse[Agent], error) {
	if params == nil {
		params = &ListAgentsParams{}
	}
	p := &ListParams{Page: params.Page, Limit: params.Limit}
	p.defaults()
	qp := map[string]string{
		"page":  fmt.Sprintf("%d", p.Page),
		"limit": fmt.Sprintf("%d", p.Limit),
	}
	if params.Type != "" {
		qp["type"] = params.Type
	}
	if params.Status != "" {
		qp["status"] = params.Status
	}
	if params.Search != "" {
		qp["search"] = params.Search
	}
	var result PaginatedResponse[Agent]
	err := s.http.get(ctx, "/api/v1/agents", qp, &result)
	return &result, err
}

// Get returns a single agent by ID.
func (s *AgentsService) Get(ctx context.Context, id string) (*Agent, error) {
	var result Agent
	err := s.http.get(ctx, "/api/v1/agents/"+id, nil, &result)
	return &result, err
}

// Create creates a new agent.
func (s *AgentsService) Create(ctx context.Context, params *CreateAgentParams) (*Agent, error) {
	var result Agent
	err := s.http.post(ctx, "/api/v1/agents", params, &result)
	return &result, err
}

// Update updates an existing agent.
func (s *AgentsService) Update(ctx context.Context, id string, params *UpdateAgentParams) (*Agent, error) {
	var result Agent
	err := s.http.patch(ctx, "/api/v1/agents/"+id, params, &result)
	return &result, err
}

// Delete removes an agent.
func (s *AgentsService) Delete(ctx context.Context, id string) error {
	return s.http.del(ctx, "/api/v1/agents/"+id)
}

// Execute runs an agent.
func (s *AgentsService) Execute(ctx context.Context, id string, params *ExecuteAgentParams) (*ExecuteResult, error) {
	var result ExecuteResult
	err := s.http.post(ctx, "/api/v1/agents/"+id+"/execute", params, &result)
	return &result, err
}

// ListVersions returns versions for an agent.
func (s *AgentsService) ListVersions(ctx context.Context, id string, params *ListParams) (*PaginatedResponse[AgentVersion], error) {
	if params == nil {
		params = &ListParams{}
	}
	params.defaults()
	qp := map[string]string{
		"page":  fmt.Sprintf("%d", params.Page),
		"limit": fmt.Sprintf("%d", params.Limit),
	}
	var result PaginatedResponse[AgentVersion]
	err := s.http.get(ctx, "/api/v1/agents/"+id+"/versions", qp, &result)
	return &result, err
}

// CreateVersion creates a new version for an agent.
func (s *AgentsService) CreateVersion(ctx context.Context, id string, params *CreateVersionParams) (*AgentVersion, error) {
	var result AgentVersion
	err := s.http.post(ctx, "/api/v1/agents/"+id+"/versions", params, &result)
	return &result, err
}

// PromoteVersion sets a version as the active version.
func (s *AgentsService) PromoteVersion(ctx context.Context, id, versionID string) error {
	return s.http.post(ctx, "/api/v1/agents/"+id+"/versions/"+versionID+"/promote", nil, nil)
}

// ListGuardrails returns guardrails for an agent.
func (s *AgentsService) ListGuardrails(ctx context.Context, id string) ([]Guardrail, error) {
	var result []Guardrail
	err := s.http.get(ctx, "/api/v1/agents/"+id+"/guardrails", nil, &result)
	return result, err
}

// CreateGuardrail adds a guardrail to an agent.
func (s *AgentsService) CreateGuardrail(ctx context.Context, id string, params *CreateGuardrailParams) (*Guardrail, error) {
	var result Guardrail
	err := s.http.post(ctx, "/api/v1/agents/"+id+"/guardrails", params, &result)
	return &result, err
}

// ListMemories returns memories for an agent.
func (s *AgentsService) ListMemories(ctx context.Context, id string, params *ListParams) (*PaginatedResponse[AgentMemory], error) {
	if params == nil {
		params = &ListParams{}
	}
	params.defaults()
	qp := map[string]string{
		"page":  fmt.Sprintf("%d", params.Page),
		"limit": fmt.Sprintf("%d", params.Limit),
	}
	var result PaginatedResponse[AgentMemory]
	err := s.http.get(ctx, "/api/v1/agents/"+id+"/memories", qp, &result)
	return &result, err
}

// CreateMemory adds a memory entry to an agent.
func (s *AgentsService) CreateMemory(ctx context.Context, id string, params *CreateMemoryParams) (*AgentMemory, error) {
	var result AgentMemory
	err := s.http.post(ctx, "/api/v1/agents/"+id+"/memories", params, &result)
	return &result, err
}

// SearchMemories searches agent memories.
func (s *AgentsService) SearchMemories(ctx context.Context, id string, params *SearchMemoriesParams) ([]AgentMemory, error) {
	var result []AgentMemory
	err := s.http.post(ctx, "/api/v1/agents/"+id+"/memories/search", params, &result)
	return result, err
}

// DeleteAllMemories removes all memories for an agent.
func (s *AgentsService) DeleteAllMemories(ctx context.Context, id string) error {
	return s.http.del(ctx, "/api/v1/agents/"+id+"/memories")
}
