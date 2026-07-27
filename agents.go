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

// CreateAgentParams are parameters for creating an agent. Type is one of
// "agent" or "workflow" — API v2 has exactly two agent kinds.
type CreateAgentParams struct {
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Description string         `json:"description,omitempty"`
	Config      map[string]any `json:"config,omitempty"`
	TemplateID  string         `json:"template_id,omitempty"`
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
//
// The version owns the model + sampling (ModelConfig), the execution-tree
// budget (RunBudget), the approval policy, cache TTL, version-scoped VFS /
// masking overrides, and the attached Tools / SubAgents / Guardrails — all
// siblings of Config, not part of it. Config carries only the agent-kind
// discriminator payload (PromptAgentConfig or WorkflowAgentConfig).
type CreateVersionParams struct {
	Version        string               `json:"version,omitempty"`
	Message        string               `json:"message,omitempty"`
	Config         AgentConfig          `json:"config,omitempty"`
	InputSchema    map[string]any       `json:"input_schema,omitempty"`
	OutputSchema   map[string]any       `json:"output_schema,omitempty"`
	SetCurrent     bool                 `json:"set_current,omitempty"`
	ModelConfig    *ModelConfig         `json:"model_config,omitempty"`
	RunBudget      *RunBudget           `json:"run_budget,omitempty"`
	ApprovalPolicy *ApprovalPolicy      `json:"approval_policy,omitempty"`
	CacheTimeout   *int                 `json:"cache_timeout,omitempty"`
	VFSEnabled     *bool                `json:"vfs_enabled,omitempty"`
	MaskingEnabled *bool                `json:"masking_enabled,omitempty"`
	Tools          []ToolAttachment     `json:"tools,omitempty"`
	SubAgents      []SubAgentAttachment `json:"sub_agents,omitempty"`
	Guardrails     []GuardrailSpec      `json:"guardrails,omitempty"`
}

// AgentConfig is implemented by every concrete agent-version config type. The
// SDK injects a "type" discriminator on marshal so the backend can route.
// API v2 has exactly two kinds: PromptAgentConfig ("agent") and
// WorkflowAgentConfig ("workflow").
type AgentConfig interface {
	agentConfig()
}

// WorkflowNode is one node in a "workflow" agent's DAG.
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

// ToolAttachment is an MCP tool attached to an agent version with per-tool
// policy. Pass as an element of CreateVersionParams.Tools.
type ToolAttachment struct {
	MCPToolID        string `json:"mcp_tool_id"`
	RequiresApproval bool   `json:"requires_approval,omitempty"`
	NoRetry          bool   `json:"no_retry,omitempty"`
	SortOrder        *int   `json:"sort_order,omitempty"`
}

// SubAgentAttachment is a delegate sub-agent attached to an agent version
// (agents-as-tools). Pass as an element of CreateVersionParams.SubAgents.
type SubAgentAttachment struct {
	AgentID          string `json:"agent_id"`
	Alias            string `json:"alias"`
	Description      string `json:"description,omitempty"`
	Mode             string `json:"mode,omitempty"`         // "delegate" | "handoff"
	ContextMode      string `json:"context_mode,omitempty"` // "task" | "window"
	RequiresApproval bool   `json:"requires_approval,omitempty"`
	SortOrder        *int   `json:"sort_order,omitempty"`
}

// PromptAgentConfig is the Config for an "agent" — a single prompt, optionally
// extended with tools and sub-agents (a supervisor when it has sub-agents).
type PromptAgentConfig struct {
	PromptID string `json:"prompt_id"`
}

// WorkflowAgentConfig is the Config for a "workflow" — a deterministic DAG.
type WorkflowAgentConfig struct {
	Nodes []WorkflowNode `json:"nodes"`
}

func (PromptAgentConfig) agentConfig()   {}
func (WorkflowAgentConfig) agentConfig() {}

// MarshalJSON injects the "type" discriminator so the backend (and other
// SDKs) can route to the correct shape without separate API surface.

func (c PromptAgentConfig) MarshalJSON() ([]byte, error) {
	type alias PromptAgentConfig
	return marshalConfigWithType("agent", alias(c))
}
func (c WorkflowAgentConfig) MarshalJSON() ([]byte, error) {
	type alias WorkflowAgentConfig
	return marshalConfigWithType("workflow", alias(c))
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

// CreateGuardrailParams are parameters for attaching a guardrail to an agent.
type CreateGuardrailParams struct {
	Type        string         `json:"type"`         // "input" | "output"
	ScannerType string         `json:"scanner_type"` // e.g. "pii", "prompt_injection"
	Action      string         `json:"action,omitempty"`
	Config      map[string]any `json:"config,omitempty"`
	IsActive    *bool          `json:"is_active,omitempty"`
	SortOrder   *int           `json:"sort_order,omitempty"`
}

// PlaygroundParams are parameters for an ad-hoc agent run with a prompt
// override, without saving a version.
type PlaygroundParams struct {
	Input          map[string]any `json:"input"`
	PromptOverride map[string]any `json:"prompt_override"`
	VersionID      string         `json:"version_id,omitempty"`
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

// Playground runs the agent with an ad-hoc prompt override without saving a
// version. PromptOverride may carry "system_prompt", "user_prompt" and
// "input_schema"; VersionID selects whose runtime behavior is used (defaults
// to the current version).
func (s *AgentsService) Playground(ctx context.Context, id string, params *PlaygroundParams) (map[string]any, error) {
	var result map[string]any
	err := s.http.post(ctx, "/api/v1/agents/"+id+"/playground", params, &result)
	return result, err
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
