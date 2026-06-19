package promptrails

import (
	"context"
	"encoding/json"
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
// Config is a typed discriminated union by agent type.
type CreateVersionParams struct {
	Version      string         `json:"version,omitempty"`
	Message      string         `json:"message,omitempty"`
	Config       AgentConfig    `json:"config,omitempty"`
	InputSchema  map[string]any `json:"input_schema,omitempty"`
	OutputSchema map[string]any `json:"output_schema,omitempty"`
	SetCurrent   bool           `json:"set_current,omitempty"`
}

// AgentConfig is implemented by every concrete agent config type. The SDK
// injects a "type" discriminator on marshal so the backend can route.
type AgentConfig interface {
	agentConfig()
}

// PromptLink pins a prompt into a chain/multi-agent step at a role.
type PromptLink struct {
	PromptID  string `json:"prompt_id"`
	Role      string `json:"role"`
	SortOrder int    `json:"sort_order"`
}

// WorkflowNode is one step in a workflow DAG.
type WorkflowNode struct {
	ID            string         `json:"id"`
	PromptID      string         `json:"prompt_id,omitempty"`
	DependsOn     []string       `json:"depends_on"`
	NodeType      string         `json:"node_type,omitempty"`
	MediaProvider string         `json:"media_provider,omitempty"`
	MediaType     string         `json:"media_type,omitempty"`
	MediaModel    string         `json:"media_model,omitempty"`
	MediaConfig   map[string]any `json:"media_config,omitempty"`
}

// CompositeStep references another agent inside a composite agent.
type CompositeStep struct {
	ID           string         `json:"id"`
	AgentID      string         `json:"agent_id"`
	DependsOn    []string       `json:"depends_on,omitempty"`
	InputMapping map[string]any `json:"input_mapping,omitempty"`
}

// SimpleAgentConfig is the config for a simple (single-prompt) agent.
type SimpleAgentConfig struct {
	PromptID               string   `json:"prompt_id"`
	ApprovalRequired       bool     `json:"approval_required,omitempty"`
	ApprovalCheckpointName string   `json:"approval_checkpoint_name,omitempty"`
	MaxTokens              int      `json:"max_tokens,omitempty"`
	Temperature            *float64 `json:"temperature,omitempty"`
	LLMModelID             string   `json:"llm_model_id,omitempty"`
}

// ChainAgentConfig is the config for a sequential chain agent.
type ChainAgentConfig struct {
	PromptIDs              []PromptLink `json:"prompt_ids"`
	ApprovalRequired       bool         `json:"approval_required,omitempty"`
	ApprovalCheckpointName string       `json:"approval_checkpoint_name,omitempty"`
}

// MultiAgentConfig is the config for a parallel multi-agent run.
type MultiAgentConfig struct {
	PromptIDs []PromptLink `json:"prompt_ids"`
}

// WorkflowAgentConfig is the config for a DAG-style workflow agent.
type WorkflowAgentConfig struct {
	Nodes []WorkflowNode `json:"nodes"`
}

// CompositeAgentConfig is the config for an agent composed of sub-agents.
type CompositeAgentConfig struct {
	Steps []CompositeStep `json:"steps"`
}

func (SimpleAgentConfig) agentConfig()    {}
func (ChainAgentConfig) agentConfig()     {}
func (MultiAgentConfig) agentConfig()     {}
func (WorkflowAgentConfig) agentConfig()  {}
func (CompositeAgentConfig) agentConfig() {}

// MarshalJSON injects the "type" discriminator so the backend (and other
// SDKs) can route to the correct shape without separate API surface.

func (c SimpleAgentConfig) MarshalJSON() ([]byte, error) {
	type alias SimpleAgentConfig
	return marshalConfigWithType("simple", alias(c))
}
func (c ChainAgentConfig) MarshalJSON() ([]byte, error) {
	type alias ChainAgentConfig
	return marshalConfigWithType("chain", alias(c))
}
func (c MultiAgentConfig) MarshalJSON() ([]byte, error) {
	type alias MultiAgentConfig
	return marshalConfigWithType("multi_agent", alias(c))
}
func (c WorkflowAgentConfig) MarshalJSON() ([]byte, error) {
	type alias WorkflowAgentConfig
	return marshalConfigWithType("workflow", alias(c))
}
func (c CompositeAgentConfig) MarshalJSON() ([]byte, error) {
	type alias CompositeAgentConfig
	return marshalConfigWithType("composite", alias(c))
}

func marshalConfigWithType(typ string, inner any) ([]byte, error) {
	// Encode inner as a map first so the "type" key merges with the
	// concrete fields in a single flat JSON object.
	b, err := json.Marshal(inner)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	if m == nil {
		m = map[string]any{}
	}
	m["type"] = typ
	return json.Marshal(m)
}

// CreateGuardrailParams are parameters for creating a guardrail.
type CreateGuardrailParams struct {
	Name   string         `json:"name"`
	Type   string         `json:"type"`
	Config map[string]any `json:"config,omitempty"`
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
