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

// AgentTriggerSource identifies the inbound channel that fires the trigger.
// One of: "generic", "slack", "telegram", "whatsapp", "teams", "schedule".
type AgentTriggerSource string

const (
	AgentTriggerSourceGeneric  AgentTriggerSource = "generic"
	AgentTriggerSourceSlack    AgentTriggerSource = "slack"
	AgentTriggerSourceTelegram AgentTriggerSource = "telegram"
	AgentTriggerSourceWhatsApp AgentTriggerSource = "whatsapp"
	AgentTriggerSourceTeams    AgentTriggerSource = "teams"
	AgentTriggerSourceSchedule AgentTriggerSource = "schedule"
)

// AgentTrigger represents an agent trigger from any source.
type AgentTrigger struct {
	ID           string                 `json:"id"`
	WorkspaceID  string                 `json:"workspace_id"`
	AgentID      string                 `json:"agent_id"`
	Name         string                 `json:"name"`
	Token        string                 `json:"token"`
	TokenPrefix  string                 `json:"token_prefix"`
	Source       AgentTriggerSource     `json:"source"`
	SourceConfig map[string]interface{} `json:"source_config,omitempty"`
	ReplyConfig  map[string]interface{} `json:"reply_config,omitempty"`
	IsActive     bool                   `json:"is_active"`
	HasSecret    bool                   `json:"has_secret"`
	LastUsedAt   *time.Time             `json:"last_used_at,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// AgentTriggerCreateResponse includes the secret (shown only once).
type AgentTriggerCreateResponse struct {
	AgentTrigger
	Secret string `json:"secret,omitempty"`
}

// AgentVFSFile is a single file or directory entry in an agent's Virtual Filesystem.
type AgentVFSFile struct {
	ID             string                 `json:"id"`
	WorkspaceID    string                 `json:"workspace_id"`
	AgentID        string                 `json:"agent_id"`
	Path           string                 `json:"path"`
	ParentPath     string                 `json:"parent_path"`
	Name           string                 `json:"name"`
	IsDir          bool                   `json:"is_dir"`
	Content        string                 `json:"content,omitempty"`
	SizeBytes      int64                  `json:"size_bytes"`
	MimeType       string                 `json:"mime_type,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	LastWriterKind string                 `json:"last_writer_kind"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// AgentVFSGrepMatch is a single line match from a VFS grep call.
type AgentVFSGrepMatch struct {
	Path       string `json:"path"`
	LineNumber int    `json:"line_number"`
	Line       string `json:"line"`
}

// AgentVFSReadResult is the full response body of a VFS read.
type AgentVFSReadResult struct {
	File       AgentVFSFile `json:"file"`
	Content    string       `json:"content"`
	TotalLines int          `json:"total_lines"`
	Truncated  bool         `json:"truncated"`
}

// AgentVFSWriteMode controls VFS write behavior.
type AgentVFSWriteMode string

const (
	AgentVFSWriteOverwrite AgentVFSWriteMode = "overwrite"
	AgentVFSWriteAppend    AgentVFSWriteMode = "append"
)

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

// LLMModel represents an LLM model in the catalog.
type LLMModel struct {
	ID                    string     `json:"id"`
	Provider              string     `json:"provider"`
	ModelID               string     `json:"model_id"`
	DisplayName           string     `json:"display_name"`
	InputPrice            *float64   `json:"input_price"`
	OutputPrice           *float64   `json:"output_price"`
	CachedInputPrice      *float64   `json:"cached_input_price"`
	MaxTokens             *int       `json:"max_tokens"`
	SupportsVision        bool       `json:"supports_vision"`
	SupportsTools         bool       `json:"supports_tools"`
	SupportsJSON          bool       `json:"supports_json"`
	SupportsStreaming     bool       `json:"supports_streaming"`
	SupportsTemperature   bool       `json:"supports_temperature"`
	SupportsTopP          bool       `json:"supports_top_p"`
	SupportsTopK          bool       `json:"supports_top_k"`
	SupportsReasoning     bool       `json:"supports_reasoning"`
	SupportsWebSearch     bool       `json:"supports_web_search"`
	SupportsPromptCaching bool       `json:"supports_prompt_caching"`
	IsActive              bool       `json:"is_active"`
	IsDeprecated          bool       `json:"is_deprecated"`
	DeprecatedAt          *time.Time `json:"deprecated_at"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

// AvailableModelEntry is a model offered by the available-models endpoint,
// filtered to the providers the workspace holds credentials for.
type AvailableModelEntry struct {
	ID                    string   `json:"id"`
	ModelID               string   `json:"model_id"`
	DisplayName           string   `json:"display_name"`
	MaxTokens             *int     `json:"max_tokens"`
	SupportsVision        bool     `json:"supports_vision"`
	SupportsTools         bool     `json:"supports_tools"`
	SupportsJSON          bool     `json:"supports_json"`
	SupportsTemperature   bool     `json:"supports_temperature"`
	SupportsTopP          bool     `json:"supports_top_p"`
	SupportsTopK          bool     `json:"supports_top_k"`
	SupportsReasoning     bool     `json:"supports_reasoning"`
	SupportsWebSearch     bool     `json:"supports_web_search"`
	SupportsPromptCaching bool     `json:"supports_prompt_caching"`
	InputPrice            *float64 `json:"input_price"`
	OutputPrice           *float64 `json:"output_price"`
	IsDeprecated          bool     `json:"is_deprecated"`
}

// AvailableModelGroup groups available models by provider.
type AvailableModelGroup struct {
	Provider string                `json:"provider"`
	Models   []AvailableModelEntry `json:"models"`
}

// MediaModel represents an available media model.
type MediaModel struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Provider     string         `json:"provider"`
	MediaType    string         `json:"media_type"`
	IsActive     bool           `json:"is_active"`
	IsDeprecated bool           `json:"is_deprecated"`
	DeprecatedAt *time.Time     `json:"deprecated_at"`
	Config       map[string]any `json:"config,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
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
