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

// AgentVersion represents a version of an agent. In API v2 the version owns
// the model + runtime configuration (ModelConfig, RunBudget, ApprovalPolicy,
// cache TTL, VFS/masking overrides) and the attached tools, sub-agents and
// guardrails — siblings of Config.
type AgentVersion struct {
	ID             string           `json:"id"`
	AgentID        string           `json:"agent_id"`
	Version        string           `json:"version"`
	IsActive       bool             `json:"is_active"`
	Message        string           `json:"message"`
	Config         map[string]any   `json:"config,omitempty"`
	InputSchema    map[string]any   `json:"input_schema,omitempty"`
	OutputSchema   map[string]any   `json:"output_schema,omitempty"`
	ModelConfig    *ModelConfig     `json:"model_config,omitempty"`
	RunBudget      *RunBudget       `json:"run_budget,omitempty"`
	ApprovalPolicy *ApprovalPolicy  `json:"approval_policy,omitempty"`
	CacheTimeout   *int             `json:"cache_timeout,omitempty"`
	VFSEnabled     *bool            `json:"vfs_enabled,omitempty"`
	MaskingEnabled *bool            `json:"masking_enabled,omitempty"`
	Tools          []map[string]any `json:"tools,omitempty"`
	SubAgents      []map[string]any `json:"sub_agents,omitempty"`
	Guardrails     []map[string]any `json:"guardrails,omitempty"`
	CreatedAt      time.Time        `json:"created_at"`
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

// PromptVersion represents a content-only version of a prompt (API v2).
// Model, sampling, tools, output schema and cache TTL live on the agent
// version, not on the prompt.
type PromptVersion struct {
	ID           string         `json:"id"`
	PromptID     string         `json:"prompt_id"`
	Version      string         `json:"version"`
	SystemPrompt string         `json:"system_prompt,omitempty"`
	UserPrompt   string         `json:"user_prompt,omitempty"`
	InputSchema  map[string]any `json:"input_schema,omitempty"`
	IsActive     bool           `json:"is_active"`
	Config       map[string]any `json:"config,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
}

// Execution represents an agent execution record.
//
// API v2 executions form a tree: a sub-agent delegation, handoff continuation
// or workflow agent-node run has ParentExecutionID set and appears in its
// parent's Children. Root executions have an empty ParentExecutionID. Status
// gains "waiting_approval" (parked at an approval-gated tool call) and
// "cancel_requested" (cooperative cancel observed before finalizing).
type Execution struct {
	ID                string         `json:"id"`
	AgentID           string         `json:"agent_id"`
	AgentName         string         `json:"agent_name"`
	WorkspaceID       string         `json:"workspace_id"`
	SessionID         string         `json:"session_id,omitempty"`
	ParentExecutionID *string        `json:"parent_execution_id,omitempty"`
	Status            string         `json:"status"`
	Input             map[string]any `json:"input,omitempty"`
	Output            map[string]any `json:"output,omitempty"`
	Error             *string        `json:"error,omitempty"`
	TraceID           *string        `json:"trace_id,omitempty"`
	DurationMs        int            `json:"duration_ms"`
	TotalTokens       int            `json:"total_tokens"`
	Cost              float64        `json:"cost"`
	ApprovalExpiresAt *time.Time     `json:"approval_expires_at,omitempty"`
	Children          []Execution    `json:"children,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
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

// Trace represents an observability span (API v2). TokenUsage carries raw
// provider counts and, when reported, extends beyond prompt/completion/total
// with "cached_tokens" (prompt-cache read hits), "cache_creation_tokens" and
// "reasoning_tokens".
type Trace struct {
	ID                  string         `json:"id"`
	TraceID             string         `json:"trace_id"`
	ExecutionID         *string        `json:"execution_id,omitempty"`
	WorkspaceID         string         `json:"workspace_id"`
	ParentID            *string        `json:"parent_id,omitempty"`
	Name                string         `json:"name"`
	Kind                string         `json:"kind"`
	Status              string         `json:"status"`
	Level               string         `json:"level,omitempty"`
	Input               any            `json:"input,omitempty"`
	Output              any            `json:"output,omitempty"`
	Attributes          map[string]any `json:"attributes,omitempty"`
	Tags                []string       `json:"tags,omitempty"`
	TokenUsage          map[string]any `json:"token_usage,omitempty"`
	Cost                *float64       `json:"cost,omitempty"`
	DurationMs          int            `json:"duration_ms"`
	CompletionStartTime *string        `json:"completion_start_time,omitempty"`
	ErrorMessage        string         `json:"error_message,omitempty"`
	ErrorType           string         `json:"error_type,omitempty"`
	PromptTokens        *int           `json:"prompt_tokens,omitempty"`
	CompletionTokens    *int           `json:"completion_tokens,omitempty"`
	TotalTokens         *int           `json:"total_tokens,omitempty"`
	ModelName           string         `json:"model_name,omitempty"`
	AgentID             *string        `json:"agent_id,omitempty"`
	AgentName           string         `json:"agent_name,omitempty"`
	AgentType           string         `json:"agent_type,omitempty"`
	UserID              *string        `json:"user_id,omitempty"`
	SessionID           string         `json:"session_id,omitempty"`
	DataSourceID        *string        `json:"data_source_id,omitempty"`
	DataSourceName      string         `json:"data_source_name,omitempty"`
	PromptName          string         `json:"prompt_name,omitempty"`
	MCPToolName         string         `json:"mcp_tool_name,omitempty"`
	MCPToolType         string         `json:"mcp_tool_type,omitempty"`
	ServiceName         string         `json:"service_name,omitempty"`
	StartedAt           string         `json:"started_at,omitempty"`
	EndedAt             *string        `json:"ended_at,omitempty"`
	CreatedAt           time.Time      `json:"created_at"`
}

// TraceSummary holds aggregate statistics over a filtered set of traces
// (GET /traces/summary).
type TraceSummary struct {
	TotalTraces    int     `json:"total_traces"`
	TotalTokens    int     `json:"total_tokens"`
	TotalCost      float64 `json:"total_cost"`
	AvgDurationMs  float64 `json:"avg_duration_ms"`
	ErrorCount     int     `json:"error_count"`
	UniqueModels   int     `json:"unique_models"`
	UniqueSessions int     `json:"unique_sessions"`
}

// ModelConfig is the version-scoped model + sampling ownership (API v2). Every
// field is optional; unset sampling inherits the provider/model default.
type ModelConfig struct {
	ModelID         string   `json:"model_id,omitempty"`
	FallbackModelID string   `json:"fallback_model_id,omitempty"`
	Temperature     *float64 `json:"temperature,omitempty"`
	TopP            *float64 `json:"top_p,omitempty"`
	TopK            *int     `json:"top_k,omitempty"`
	MaxTokens       *int     `json:"max_tokens,omitempty"`
}

// RunBudget bounds the whole execution tree, enforced at the root. Every field
// is optional.
type RunBudget struct {
	MaxCost        *float64 `json:"max_cost,omitempty"`
	MaxTotalTokens *int     `json:"max_total_tokens,omitempty"`
	MaxToolCalls   *int     `json:"max_tool_calls,omitempty"`
	MaxChildren    *int     `json:"max_children,omitempty"`
	MaxDepth       *int     `json:"max_depth,omitempty"`
}

// ApprovalPolicy configures who may approve/deny a parked, approval-gated call.
type ApprovalPolicy struct {
	Mode      string   `json:"mode,omitempty"` // "admins" (default) | "assigned" | "any_member"
	MemberIDs []string `json:"member_ids,omitempty"`
}

// GuardrailSpec is a guardrail attachment spec — the input shape used to
// create/attach a guardrail on an agent version. ID is present on responses
// only; omit it on create.
type GuardrailSpec struct {
	ID          string         `json:"id,omitempty"`
	Type        string         `json:"type"` // "input" | "output"
	ScannerType string         `json:"scanner_type"`
	Action      string         `json:"action,omitempty"`
	Config      map[string]any `json:"config,omitempty"`
	IsActive    *bool          `json:"is_active,omitempty"`
	SortOrder   *int           `json:"sort_order,omitempty"`
}

// ScannerMeta is metadata for an available guardrail scanner
// (GET /guardrails/scanners).
type ScannerMeta struct {
	Type           string   `json:"type"`
	Label          string   `json:"label"`
	Description    string   `json:"description"`
	Category       string   `json:"category"` // "local" | "llm_guard"
	Enabled        bool     `json:"enabled"`
	ConfigFields   []string `json:"config_fields,omitempty"`
	DisabledReason string   `json:"disabled_reason,omitempty"`
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
