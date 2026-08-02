// Package promptrails provides a Go client for the PromptRails API.
//
// Usage:
//
//	client := promptrails.NewClient("pr_key_...")
//	agents, err := client.Agents.List(ctx, nil)
package promptrails

// Client is the PromptRails API client.
type Client struct {
	Agents        *AgentsService
	Prompts       *PromptsService
	Executions    *ExecutionsService
	Credentials   *CredentialsService
	DataSources   *DataSourcesService
	Chat          *ChatService
	Traces        *TracesService
	MCPTools      *MCPToolsService
	Guardrails    *GuardrailsService
	AgentTriggers *AgentTriggersService
	AgentVFS      *AgentVFSService
	A2A           *A2AService
	LLMModels     *LLMModelsService
	Assets        *AssetsService

	http *httpClient
}

// NewClient creates a new PromptRails API client.
func NewClient(apiKey string, opts ...Option) *Client {
	cfg := defaultConfig()
	cfg.apiKey = apiKey
	for _, o := range opts {
		o(&cfg)
	}
	h := newHTTPClient(cfg)
	return &Client{
		Agents:        &AgentsService{http: h},
		Prompts:       &PromptsService{http: h},
		Executions:    &ExecutionsService{http: h},
		Credentials:   &CredentialsService{http: h},
		DataSources:   &DataSourcesService{http: h},
		Chat:          &ChatService{http: h},
		Traces:        &TracesService{http: h},
		MCPTools:      &MCPToolsService{http: h},
		Guardrails:    &GuardrailsService{http: h},
		AgentTriggers: &AgentTriggersService{http: h},
		AgentVFS:      &AgentVFSService{http: h},
		A2A:           &A2AService{http: h},
		LLMModels:     &LLMModelsService{http: h},
		Assets:        &AssetsService{http: h},
		http:          h,
	}
}
