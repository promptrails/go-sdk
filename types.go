package promptrails

import "time"

// Agent represents an agent.
type Agent struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	Type             string         `json:"type"`
	Description      string         `json:"description"`
	Status           string         `json:"status"`
	WorkspaceID      string         `json:"workspace_id"`
	Config           map[string]any `json:"config,omitempty"`
	CurrentVersionID *string        `json:"current_version_id,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// AgentVersion represents a version of an agent.
type AgentVersion struct {
	ID        string         `json:"id"`
	AgentID   string         `json:"agent_id"`
	Version   string         `json:"version"`
	IsActive  bool           `json:"is_active"`
	Message   string         `json:"message"`
	Config    map[string]any `json:"config,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

// Prompt represents a prompt template.
type Prompt struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	Description      string         `json:"description"`
	Status           string         `json:"status"`
	WorkspaceID      string         `json:"workspace_id"`
	SystemPrompt     *string        `json:"system_prompt,omitempty"`
	UserPrompt       *string        `json:"user_prompt,omitempty"`
	Config           map[string]any `json:"config,omitempty"`
	CurrentVersionID *string        `json:"current_version_id,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// PromptVersion represents a version of a prompt.
type PromptVersion struct {
	ID        string         `json:"id"`
	PromptID  string         `json:"prompt_id"`
	Version   string         `json:"version"`
	IsActive  bool           `json:"is_active"`
	Config    map[string]any `json:"config,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

// RunPromptResponse holds the result of running a prompt.
type RunPromptResponse struct {
	Output string `json:"output"`
}

// Execution represents an agent execution record.
type Execution struct {
	ID          string         `json:"id"`
	AgentID     string         `json:"agent_id"`
	AgentName   string         `json:"agent_name"`
	WorkspaceID string         `json:"workspace_id"`
	Status      string         `json:"status"`
	Input       map[string]any `json:"input,omitempty"`
	Output      map[string]any `json:"output,omitempty"`
	Error       *string        `json:"error,omitempty"`
	TraceID     *string        `json:"trace_id,omitempty"`
	DurationMs  int            `json:"duration_ms"`
	TotalTokens int            `json:"total_tokens"`
	Cost        float64        `json:"cost"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// ExecuteResult is the response from executing an agent.
type ExecuteResult struct {
	Output      map[string]any `json:"output,omitempty"`
	TraceID     string         `json:"trace_id"`
	ExecutionID string         `json:"execution_id"`
	Status      string         `json:"status"`
	Cost        float64        `json:"cost"`
}

// Credential represents a provider credential.
type Credential struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Provider    string    `json:"provider"`
	WorkspaceID string    `json:"workspace_id"`
	IsDefault   bool      `json:"is_default"`
	IsValid     bool      `json:"is_valid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DataSource represents a data source.
type DataSource struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	Type             string         `json:"type"`
	Description      string         `json:"description"`
	Status           string         `json:"status"`
	WorkspaceID      string         `json:"workspace_id"`
	Config           map[string]any `json:"config,omitempty"`
	CurrentVersionID *string        `json:"current_version_id,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// DataSourceVersion represents a version of a data source.
type DataSourceVersion struct {
	ID           string         `json:"id"`
	DataSourceID string         `json:"data_source_id"`
	Version      string         `json:"version"`
	IsActive     bool           `json:"is_active"`
	Config       map[string]any `json:"config,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
}

// ChatSession represents a chat session.
type ChatSession struct {
	ID          string     `json:"id"`
	AgentID     string     `json:"agent_id"`
	WorkspaceID string     `json:"workspace_id"`
	Title       string     `json:"title"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// ChatMessage represents a message in a chat session.
type ChatMessage struct {
	ID        string         `json:"id"`
	SessionID string         `json:"session_id"`
	Role      string         `json:"role"`
	Content   string         `json:"content"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

// Trace represents an execution trace.
type Trace struct {
	ID          string    `json:"id"`
	TraceID     string    `json:"trace_id"`
	ExecutionID string    `json:"execution_id"`
	WorkspaceID string    `json:"workspace_id"`
	Kind        string    `json:"kind"`
	Name        string    `json:"name"`
	Input       any       `json:"input,omitempty"`
	Output      any       `json:"output,omitempty"`
	DurationMs  int       `json:"duration_ms"`
	ParentID    *string   `json:"parent_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// CostSummary represents cost summary data.
type CostSummary struct {
	TotalCost       float64 `json:"total_cost"`
	TotalTokens     int     `json:"total_tokens"`
	TotalExecutions int     `json:"total_executions"`
	ByModel         any     `json:"by_model,omitempty"`
	ByDay           any     `json:"by_day,omitempty"`
}

// Score represents an evaluation score.
type Score struct {
	ID            string         `json:"id"`
	ExecutionID   string         `json:"execution_id"`
	ScoreConfigID string         `json:"score_config_id"`
	Value         float64        `json:"value"`
	Comment       string         `json:"comment"`
	Metadata      map[string]any `json:"metadata,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// ScoreConfig represents a score configuration.
type ScoreConfig struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	WorkspaceID string    `json:"workspace_id"`
	DataType    string    `json:"data_type"`
	MinValue    *float64  `json:"min_value,omitempty"`
	MaxValue    *float64  `json:"max_value,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ScoreAggregate holds aggregated score data.
type ScoreAggregate struct {
	ConfigID   string  `json:"config_id"`
	ConfigName string  `json:"config_name"`
	Count      int     `json:"count"`
	Average    float64 `json:"average"`
	Min        float64 `json:"min"`
	Max        float64 `json:"max"`
}

// MCPTool represents an MCP tool.
type MCPTool struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	WorkspaceID string         `json:"workspace_id"`
	ServerURL   string         `json:"server_url"`
	ToolSchema  map[string]any `json:"tool_schema,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// ApprovalRequest represents a human-in-the-loop approval request.
type ApprovalRequest struct {
	ID             string         `json:"id"`
	ExecutionID    string         `json:"execution_id"`
	AgentID        string         `json:"agent_id"`
	WorkspaceID    string         `json:"workspace_id"`
	CheckpointName string         `json:"checkpoint_name"`
	Payload        map[string]any `json:"payload,omitempty"`
	Reason         string         `json:"reason"`
	Decision       *string        `json:"decision,omitempty"`
	DecidedBy      *string        `json:"decided_by,omitempty"`
	DecidedAt      *time.Time     `json:"decided_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
}

// WebhookTrigger represents a webhook trigger.
type WebhookTrigger struct {
	ID          string     `json:"id"`
	AgentID     string     `json:"agent_id"`
	Name        string     `json:"name"`
	Token       string     `json:"token"`
	TokenPrefix string     `json:"token_prefix"`
	IsActive    bool       `json:"is_active"`
	HasSecret   bool       `json:"has_secret"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// WebhookTriggerCreateResponse includes the secret (shown only once).
type WebhookTriggerCreateResponse struct {
	WebhookTrigger
	Secret string `json:"secret,omitempty"`
}

// AgentMemory represents a memory entry for an agent.
type AgentMemory struct {
	ID        string         `json:"id"`
	AgentID   string         `json:"agent_id"`
	Content   string         `json:"content"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

// Guardrail represents a guardrail configuration for an agent.
type Guardrail struct {
	ID        string         `json:"id"`
	AgentID   string         `json:"agent_id"`
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	Config    map[string]any `json:"config,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

// A2AAgentCard represents an A2A agent card.
type A2AAgentCard struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	URL          string `json:"url"`
	Version      string `json:"version"`
	Capabilities any    `json:"capabilities,omitempty"`
	Skills       any    `json:"skills,omitempty"`
}

// A2ATask represents an A2A task.
type A2ATask struct {
	ID        string         `json:"id"`
	Status    string         `json:"status"`
	Messages  any            `json:"messages,omitempty"`
	Artifacts any            `json:"artifacts,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

// LLMModel represents an available LLM model.
type LLMModel struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
}

// MediaModel represents an available media model.
type MediaModel struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Provider  string         `json:"provider"`
	MediaType string         `json:"media_type"`
	IsActive  bool           `json:"is_active"`
	Config    map[string]any `json:"config,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// Asset represents a media asset.
type Asset struct {
	ID          string         `json:"id"`
	WorkspaceID string         `json:"workspace_id"`
	FileName    string         `json:"file_name"`
	ContentType string         `json:"content_type"`
	MediaType   string         `json:"media_type"`
	Provider    string         `json:"provider"`
	Size        int64          `json:"size"`
	URL         string         `json:"url,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// SignedURLResponse holds a signed URL for accessing an asset.
type SignedURLResponse struct {
	URL       string `json:"url"`
	ExpiresAt string `json:"expires_at,omitempty"`
}
