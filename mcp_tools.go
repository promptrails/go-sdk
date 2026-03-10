package promptrails

import (
	"context"
	"fmt"
)

// MCPToolsService handles MCP tool API calls.
type MCPToolsService struct {
	http *httpClient
}

// ListMCPToolsParams are parameters for listing MCP tools.
type ListMCPToolsParams struct {
	Page  int
	Limit int
}

// CreateMCPToolParams are parameters for creating an MCP tool.
type CreateMCPToolParams struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	ServerURL   string         `json:"server_url"`
	ToolSchema  map[string]any `json:"tool_schema,omitempty"`
}

// UpdateMCPToolParams are parameters for updating an MCP tool.
type UpdateMCPToolParams struct {
	Name        *string        `json:"name,omitempty"`
	Description *string        `json:"description,omitempty"`
	ServerURL   *string        `json:"server_url,omitempty"`
	ToolSchema  map[string]any `json:"tool_schema,omitempty"`
}

// List returns a paginated list of MCP tools.
func (s *MCPToolsService) List(ctx context.Context, params *ListMCPToolsParams) (*PaginatedResponse[MCPTool], error) {
	if params == nil {
		params = &ListMCPToolsParams{}
	}
	p := &ListParams{Page: params.Page, Limit: params.Limit}
	p.defaults()
	qp := map[string]string{
		"page":  fmt.Sprintf("%d", p.Page),
		"limit": fmt.Sprintf("%d", p.Limit),
	}
	var result PaginatedResponse[MCPTool]
	err := s.http.get(ctx, "/api/v1/mcp-tools", qp, &result)
	return &result, err
}

// Get returns a single MCP tool by ID.
func (s *MCPToolsService) Get(ctx context.Context, id string) (*MCPTool, error) {
	var result MCPTool
	err := s.http.get(ctx, "/api/v1/mcp-tools/"+id, nil, &result)
	return &result, err
}

// Create creates a new MCP tool.
func (s *MCPToolsService) Create(ctx context.Context, params *CreateMCPToolParams) (*MCPTool, error) {
	var result MCPTool
	err := s.http.post(ctx, "/api/v1/mcp-tools", params, &result)
	return &result, err
}

// Update updates an existing MCP tool.
func (s *MCPToolsService) Update(ctx context.Context, id string, params *UpdateMCPToolParams) (*MCPTool, error) {
	var result MCPTool
	err := s.http.patch(ctx, "/api/v1/mcp-tools/"+id, params, &result)
	return &result, err
}

// Delete removes an MCP tool.
func (s *MCPToolsService) Delete(ctx context.Context, id string) error {
	return s.http.del(ctx, "/api/v1/mcp-tools/"+id)
}
