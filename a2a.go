package promptrails

import (
	"context"
	"fmt"
)

// A2AService handles A2A (Agent-to-Agent) API calls.
type A2AService struct {
	http *httpClient
}

// A2ASendMessageParams are parameters for sending an A2A message.
type A2ASendMessageParams struct {
	AgentID string         `json:"agent_id"`
	Message map[string]any `json:"message"`
}

// GetAgentCard returns the A2A agent card for an agent.
func (s *A2AService) GetAgentCard(ctx context.Context, agentID string) (*A2AAgentCard, error) {
	var result A2AAgentCard
	err := s.http.get(ctx, "/api/v1/a2a/agents/"+agentID+"/card", nil, &result)
	return &result, err
}

// SendMessage sends an A2A message.
func (s *A2AService) SendMessage(ctx context.Context, params *A2ASendMessageParams) (*A2ATask, error) {
	var result A2ATask
	err := s.http.post(ctx, "/api/v1/a2a/messages", params, &result)
	return &result, err
}

// GetTask returns an A2A task by ID.
func (s *A2AService) GetTask(ctx context.Context, taskID string) (*A2ATask, error) {
	var result A2ATask
	err := s.http.get(ctx, "/api/v1/a2a/tasks/"+taskID, nil, &result)
	return &result, err
}

// ListTasks returns a paginated list of A2A tasks.
func (s *A2AService) ListTasks(ctx context.Context, params *ListParams) (*PaginatedResponse[A2ATask], error) {
	if params == nil {
		params = &ListParams{}
	}
	params.defaults()
	qp := map[string]string{
		"page":  fmt.Sprintf("%d", params.Page),
		"limit": fmt.Sprintf("%d", params.Limit),
	}
	var result PaginatedResponse[A2ATask]
	err := s.http.get(ctx, "/api/v1/a2a/tasks", qp, &result)
	return &result, err
}

// CancelTask cancels an A2A task.
func (s *A2AService) CancelTask(ctx context.Context, taskID string) error {
	return s.http.post(ctx, "/api/v1/a2a/tasks/"+taskID+"/cancel", nil, nil)
}
