package promptrails

import (
	"context"
	"fmt"
)

// PromptsService handles prompt-related API calls.
type PromptsService struct {
	http *httpClient
}

// ListPromptsParams are parameters for listing prompts.
type ListPromptsParams struct {
	Page   int
	Limit  int
	Search string
}

// CreatePromptParams are parameters for creating a prompt.
type CreatePromptParams struct {
	Name         string         `json:"name"`
	Description  string         `json:"description,omitempty"`
	SystemPrompt string         `json:"system_prompt,omitempty"`
	UserPrompt   string         `json:"user_prompt,omitempty"`
	Config       map[string]any `json:"config,omitempty"`
}

// UpdatePromptParams are parameters for updating a prompt.
type UpdatePromptParams struct {
	Name         *string        `json:"name,omitempty"`
	Description  *string        `json:"description,omitempty"`
	SystemPrompt *string        `json:"system_prompt,omitempty"`
	UserPrompt   *string        `json:"user_prompt,omitempty"`
	Config       map[string]any `json:"config,omitempty"`
}

// Reasoning effort levels accepted by RunPromptParams.ReasoningEffort. They
// mirror langrails' provider-agnostic ReasoningEffort and only take effect on
// models whose SupportsReasoning capability is set.
const (
	ReasoningEffortMinimal = "minimal"
	ReasoningEffortLow     = "low"
	ReasoningEffortMedium  = "medium"
	ReasoningEffortHigh    = "high"
)

// RunMessage is a role/content message seeded into a prompt run.
type RunMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// RunPromptParams are parameters for running a prompt. All fields except
// UserPrompt are optional; omitted fields fall back to the prompt version's
// saved configuration.
type RunPromptParams struct {
	SystemPrompt       string         `json:"system_prompt,omitempty"`
	UserPrompt         string         `json:"user_prompt"`
	LLMModelID         string         `json:"llm_model_id,omitempty"`
	FallbackLLMModelID string         `json:"fallback_llm_model_id,omitempty"`
	Temperature        *float64       `json:"temperature,omitempty"`
	MaxTokens          *int           `json:"max_tokens,omitempty"`
	TopP               *float64       `json:"top_p,omitempty"`
	TopK               *int           `json:"top_k,omitempty"`
	Input              map[string]any `json:"input,omitempty"`
	OutputSchema       map[string]any `json:"output_schema,omitempty"`
	Tools              []string       `json:"tools,omitempty"`
	InitialMessages    []RunMessage   `json:"initial_messages,omitempty"`
	CredentialID       string         `json:"credential_id,omitempty"`
	CacheTimeout       int            `json:"cache_timeout,omitempty"`
	// Feature toggles, gated by the model's capability flags.
	ReasoningEffort string `json:"reasoning_effort,omitempty"`
	WebSearch       *bool  `json:"web_search,omitempty"`
	PromptCaching   *bool  `json:"prompt_caching,omitempty"`
}

// List returns a paginated list of prompts.
func (s *PromptsService) List(ctx context.Context, params *ListPromptsParams) (*PaginatedResponse[Prompt], error) {
	if params == nil {
		params = &ListPromptsParams{}
	}
	p := &ListParams{Page: params.Page, Limit: params.Limit}
	p.defaults()
	qp := map[string]string{
		"page":  fmt.Sprintf("%d", p.Page),
		"limit": fmt.Sprintf("%d", p.Limit),
	}
	if params.Search != "" {
		qp["search"] = params.Search
	}
	var result PaginatedResponse[Prompt]
	err := s.http.get(ctx, "/api/v1/prompts", qp, &result)
	return &result, err
}

// Get returns a single prompt by ID.
func (s *PromptsService) Get(ctx context.Context, id string) (*Prompt, error) {
	var result Prompt
	err := s.http.get(ctx, "/api/v1/prompts/"+id, nil, &result)
	return &result, err
}

// Create creates a new prompt.
func (s *PromptsService) Create(ctx context.Context, params *CreatePromptParams) (*Prompt, error) {
	var result Prompt
	err := s.http.post(ctx, "/api/v1/prompts", params, &result)
	return &result, err
}

// Update updates an existing prompt.
func (s *PromptsService) Update(ctx context.Context, id string, params *UpdatePromptParams) (*Prompt, error) {
	var result Prompt
	err := s.http.patch(ctx, "/api/v1/prompts/"+id, params, &result)
	return &result, err
}

// Delete removes a prompt.
func (s *PromptsService) Delete(ctx context.Context, id string) error {
	return s.http.del(ctx, "/api/v1/prompts/"+id)
}

// ListVersions returns versions for a prompt.
func (s *PromptsService) ListVersions(ctx context.Context, id string, params *ListParams) (*PaginatedResponse[PromptVersion], error) {
	if params == nil {
		params = &ListParams{}
	}
	params.defaults()
	qp := map[string]string{
		"page":  fmt.Sprintf("%d", params.Page),
		"limit": fmt.Sprintf("%d", params.Limit),
	}
	var result PaginatedResponse[PromptVersion]
	err := s.http.get(ctx, "/api/v1/prompts/"+id+"/versions", qp, &result)
	return &result, err
}

// CreateVersion creates a new version for a prompt.
func (s *PromptsService) CreateVersion(ctx context.Context, id string, params *CreateVersionParams) (*PromptVersion, error) {
	var result PromptVersion
	err := s.http.post(ctx, "/api/v1/prompts/"+id+"/versions", params, &result)
	return &result, err
}

// PromoteVersion sets a version as the active version.
func (s *PromptsService) PromoteVersion(ctx context.Context, id, versionID string) error {
	return s.http.post(ctx, "/api/v1/prompts/"+id+"/versions/"+versionID+"/promote", nil, nil)
}

// Run executes a prompt with the given parameters.
func (s *PromptsService) Run(ctx context.Context, id string, params *RunPromptParams) (*RunPromptResponse, error) {
	var result RunPromptResponse
	err := s.http.post(ctx, "/api/v1/prompts/"+id+"/run", params, &result)
	return &result, err
}
