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
	Costs         *CostsService
	Scores        *ScoresService
	MCPTools      *MCPToolsService
	Approvals     *ApprovalsService
	AgentTriggers *AgentTriggersService
	AgentVFS      *AgentVFSService
	A2A           *A2AService
	Media         *MediaService
	MediaModels   *MediaModelsService
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
		Costs:         &CostsService{http: h},
		Scores:        &ScoresService{http: h},
		MCPTools:      &MCPToolsService{http: h},
		Approvals:     &ApprovalsService{http: h},
		AgentTriggers: &AgentTriggersService{http: h},
		AgentVFS:      &AgentVFSService{http: h},
		A2A:           &A2AService{http: h},
		Media:         &MediaService{http: h},
		MediaModels:   &MediaModelsService{http: h},
		Assets:        &AssetsService{http: h},
		http:          h,
	}
}
