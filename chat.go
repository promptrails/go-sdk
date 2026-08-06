package promptrails

import (
	"context"
	"fmt"
)

// ChatService handles chat-related API calls.
type ChatService struct {
	http *httpClient
}

// CreateSessionParams are parameters for creating a chat session.
type CreateSessionParams struct {
	AgentID string `json:"agent_id"`
	Title   string `json:"title,omitempty"`
}

// SendMessageParams are parameters for sending a chat message.
type SendMessageParams struct {
	Content string `json:"content"`
}

// SubmitFeedbackParams are parameters for thumbs-up/down chat feedback.
type SubmitFeedbackParams struct {
	ExecutionID string `json:"execution_id"`
	Value       int    `json:"value"`
}

// ChatFeedbackResponse confirms that feedback was created or updated.
type ChatFeedbackResponse struct {
	Submitted bool `json:"submitted"`
}

// ListSessions returns a paginated list of chat sessions.
func (s *ChatService) ListSessions(ctx context.Context, params *ListParams) (*PaginatedResponse[ChatSession], error) {
	if params == nil {
		params = &ListParams{}
	}
	params.defaults()
	qp := map[string]string{
		"page":  fmt.Sprintf("%d", params.Page),
		"limit": fmt.Sprintf("%d", params.Limit),
	}
	var result PaginatedResponse[ChatSession]
	err := s.http.get(ctx, "/api/v1/chat/sessions", qp, &result)
	return &result, err
}

// GetSession returns a single chat session.
func (s *ChatService) GetSession(ctx context.Context, id string) (*ChatSession, error) {
	var result ChatSession
	err := s.http.get(ctx, "/api/v1/chat/sessions/"+id, nil, &result)
	return &result, err
}

// CreateSession creates a new chat session.
func (s *ChatService) CreateSession(ctx context.Context, params *CreateSessionParams) (*ChatSession, error) {
	var result ChatSession
	err := s.http.post(ctx, "/api/v1/chat/sessions", params, &result)
	return &result, err
}

// DeleteSession removes a chat session.
func (s *ChatService) DeleteSession(ctx context.Context, id string) error {
	return s.http.del(ctx, "/api/v1/chat/sessions/"+id)
}

// ListMessages returns messages in a session.
func (s *ChatService) ListMessages(ctx context.Context, sessionID string, params *ListParams) (*PaginatedResponse[ChatMessage], error) {
	if params == nil {
		params = &ListParams{}
	}
	params.defaults()
	qp := map[string]string{
		"page":  fmt.Sprintf("%d", params.Page),
		"limit": fmt.Sprintf("%d", params.Limit),
	}
	var result PaginatedResponse[ChatMessage]
	err := s.http.get(ctx, "/api/v1/chat/sessions/"+sessionID+"/messages", qp, &result)
	return &result, err
}

// SendMessage sends a message to a chat session.
func (s *ChatService) SendMessage(ctx context.Context, sessionID string, params *SendMessageParams) (*ChatMessage, error) {
	var result ChatMessage
	err := s.http.post(ctx, "/api/v1/chat/sessions/"+sessionID+"/messages", params, &result)
	return &result, err
}

// SubmitFeedback submits or updates feedback for an execution in a chat session.
// Value must be 1 for thumbs-up or -1 for thumbs-down.
func (s *ChatService) SubmitFeedback(ctx context.Context, sessionID string, params *SubmitFeedbackParams) (*ChatFeedbackResponse, error) {
	var result ChatFeedbackResponse
	err := s.http.post(ctx, "/api/v1/chat/sessions/"+sessionID+"/feedback", params, &result)
	return &result, err
}

// SendMessageStream posts a user message and returns a ChatStream that
// yields the agent's execution events on the same HTTP connection. Cancel
// by cancelling ctx, and always Close() the returned stream.
func (s *ChatService) SendMessageStream(ctx context.Context, sessionID string, params *SendMessageParams) (*ChatStream, error) {
	body, err := s.http.stream(ctx, "POST", "/api/v1/chat/sessions/"+sessionID+"/messages/stream", params)
	if err != nil {
		return nil, err
	}
	return newChatStream(body), nil
}
