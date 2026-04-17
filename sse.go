package promptrails

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ChatStreamEvent is one frame of a PromptRails SSE stream. Concrete types
// (ExecutionEvent, ThinkingEvent, ToolStartEvent, ToolEndEvent,
// ContentEvent, DoneEvent, ErrorEvent) implement it via type switch.
type ChatStreamEvent interface {
	chatStreamEvent()
}

// ExecutionEvent carries the execution identifier emitted at the start of
// a chat stream so the caller can correlate downstream state.
type ExecutionEvent struct {
	ExecutionID   string
	UserMessageID string
}

// ThinkingEvent surfaces intermediate reasoning text the agent produces
// between tool rounds.
type ThinkingEvent struct {
	Content string
}

// ToolStartEvent announces that the agent is about to execute a tool.
type ToolStartEvent struct {
	ID   string
	Name string
}

// ToolEndEvent reports a finished tool call. Summary is a short string
// suitable for UI display (full result lives in the execution trace).
type ToolEndEvent struct {
	ID      string
	Name    string
	Summary string
}

// ContentEvent is a delta of the final assistant response.
type ContentEvent struct {
	Content string
}

// TokenUsage is the prompt/completion token accounting for a run.
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
}

// DoneEvent closes the stream. Output may carry the full response when the
// execution did not emit progressive content deltas (e.g. non-tool agents).
type DoneEvent struct {
	Output     map[string]any
	TokenUsage TokenUsage
}

// ErrorEvent carries a terminal error from the server.
type ErrorEvent struct {
	Message string
}

func (*ExecutionEvent) chatStreamEvent() {}
func (*ThinkingEvent) chatStreamEvent()  {}
func (*ToolStartEvent) chatStreamEvent() {}
func (*ToolEndEvent) chatStreamEvent()   {}
func (*ContentEvent) chatStreamEvent()   {}
func (*DoneEvent) chatStreamEvent()      {}
func (*ErrorEvent) chatStreamEvent()     {}

// ChatStream iterates events from an SSE response. Typical use:
//
//	stream, err := client.Chat.SendMessageStream(ctx, sid, params)
//	if err != nil { return err }
//	defer stream.Close()
//	for stream.Next() {
//	    switch e := stream.Event().(type) {
//	    case *ContentEvent: io.WriteString(os.Stdout, e.Content)
//	    case *ToolStartEvent: log.Printf("tool: %s", e.Name)
//	    }
//	}
//	if err := stream.Err(); err != nil { return err }
//
// Cancel by cancelling the ctx passed to the opener.
type ChatStream struct {
	body    io.ReadCloser
	scanner *bufio.Scanner
	current ChatStreamEvent
	err     error
	closed  bool
}

func newChatStream(body io.ReadCloser) *ChatStream {
	scanner := bufio.NewScanner(body)
	// SSE frames can legitimately be larger than bufio's default 64KB buffer
	// (long tool summaries, schema payloads, etc). Give the scanner headroom.
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	return &ChatStream{body: body, scanner: scanner}
}

// Next advances to the next event. Returns false at EOF or on error —
// check Err() to distinguish.
func (s *ChatStream) Next() bool {
	if s.closed || s.err != nil {
		return false
	}

	var eventName, dataLine string
	for s.scanner.Scan() {
		line := s.scanner.Text()
		switch {
		case strings.HasPrefix(line, "event:"):
			eventName = strings.TrimSpace(line[6:])
		case strings.HasPrefix(line, "data:"):
			dataLine = strings.TrimSpace(line[5:])
		case line == "":
			if dataLine == "" {
				continue
			}
			evt, err := parseSSEEvent(eventName, dataLine)
			if err != nil {
				s.err = err
				return false
			}
			if evt == nil {
				// Unknown/invalid frame — skip silently.
				eventName, dataLine = "", ""
				continue
			}
			s.current = evt
			return true
		}
	}
	if err := s.scanner.Err(); err != nil {
		s.err = err
	}
	return false
}

// Event returns the current event (valid after Next returns true).
func (s *ChatStream) Event() ChatStreamEvent { return s.current }

// Err returns any error that terminated the stream, excluding EOF.
func (s *ChatStream) Err() error { return s.err }

// Close releases the underlying response body. Always call in a deferred
// context after opening the stream.
func (s *ChatStream) Close() error {
	if s.closed {
		return nil
	}
	s.closed = true
	return s.body.Close()
}

func parseSSEEvent(name, data string) (ChatStreamEvent, error) {
	var raw map[string]any
	if err := json.Unmarshal([]byte(data), &raw); err != nil {
		return nil, fmt.Errorf("parse sse data: %w", err)
	}

	str := func(k string) string {
		if v, ok := raw[k].(string); ok {
			return v
		}
		return ""
	}

	switch name {
	case "execution":
		if raw["execution_id"] == nil {
			return nil, nil
		}
		return &ExecutionEvent{
			ExecutionID:   str("execution_id"),
			UserMessageID: str("user_message_id"),
		}, nil
	case "thinking":
		return &ThinkingEvent{Content: str("content")}, nil
	case "tool_start":
		if raw["id"] == nil {
			return nil, nil
		}
		return &ToolStartEvent{ID: str("id"), Name: str("name")}, nil
	case "tool_end":
		if raw["id"] == nil {
			return nil, nil
		}
		return &ToolEndEvent{ID: str("id"), Name: str("name"), Summary: str("summary")}, nil
	case "content":
		return &ContentEvent{Content: str("content")}, nil
	case "done":
		evt := &DoneEvent{}
		if out, ok := raw["output"].(map[string]any); ok {
			evt.Output = out
		}
		if usage, ok := raw["token_usage"].(map[string]any); ok {
			if b, err := json.Marshal(usage); err == nil {
				_ = json.Unmarshal(b, &evt.TokenUsage)
			}
		}
		return evt, nil
	case "error":
		msg := str("message")
		if msg == "" {
			msg = str("error")
		}
		if msg == "" {
			msg = "stream error"
		}
		return &ErrorEvent{Message: msg}, nil
	default:
		return nil, nil
	}
}

// stream opens a streaming request and returns the raw response body. The
// caller owns closing the body; normally this is wrapped in a ChatStream.
// No retries / no timeout wrapping — streams are long-lived by design.
func (h *httpClient) stream(ctx context.Context, method, path string, body any) (io.ReadCloser, error) {
	url := h.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
		bodyReader = strings.NewReader(string(data))
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("X-API-Key", h.apiKey)
	req.Header.Set("User-Agent", "promptrails-go/"+Version)

	// Use a client without the short request timeout — long-lived SSE.
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	if resp.StatusCode >= 400 {
		buf, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		return nil, h.parseError(resp.StatusCode, buf)
	}
	return resp.Body, nil
}
