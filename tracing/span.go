package tracing

import (
	"sync"
	"time"
)

// Kind identifies the component a span represents. Any of the PromptRails span
// kinds is valid; the constants below cover the common ones.
type Kind = string

const (
	KindAgent      Kind = "agent"
	KindLLM        Kind = "llm"
	KindTool       Kind = "tool"
	KindDatasource Kind = "datasource"
	KindChain      Kind = "chain"
	KindWorkflow   Kind = "workflow"
	KindAgentStep  Kind = "agent_step"
	KindEmbedding  Kind = "embedding"
	KindGuardrail  Kind = "guardrail"
	KindMemory     Kind = "memory"
)

// spanPayload is the wire shape sent to POST /api/v1/traces/ingest.
type spanPayload struct {
	TraceID          string         `json:"trace_id"`
	SpanID           string         `json:"span_id"`
	ParentSpanID     string         `json:"parent_span_id,omitempty"`
	Name             string         `json:"name"`
	Kind             string         `json:"kind"`
	Status           string         `json:"status"`
	Level            string         `json:"level"`
	Input            any            `json:"input,omitempty"`
	Output           any            `json:"output,omitempty"`
	Attributes       map[string]any `json:"attributes,omitempty"`
	Tags             []string       `json:"tags,omitempty"`
	ModelName        string         `json:"model_name,omitempty"`
	PromptTokens     *int           `json:"prompt_tokens,omitempty"`
	CompletionTokens *int           `json:"completion_tokens,omitempty"`
	TotalTokens      *int           `json:"total_tokens,omitempty"`
	Cost             *float64       `json:"cost,omitempty"`
	SessionID        string         `json:"session_id,omitempty"`
	ErrorMessage     string         `json:"error_message,omitempty"`
	ErrorType        string         `json:"error_type,omitempty"`
	StartedAt        time.Time      `json:"started_at"`
	EndedAt          *time.Time     `json:"ended_at,omitempty"`
}

// Span is a single unit of work within a trace. Build it with the Set* methods,
// then call End (or use Tracer.Span, which ends it for you). Methods are safe
// for concurrent use.
type Span struct {
	mu      sync.Mutex
	payload spanPayload
	onEnd   func(spanPayload)
	ended   bool
}

// TraceID returns the span's trace identifier.
func (s *Span) TraceID() string { return s.payload.TraceID }

// SpanID returns the span's identifier.
func (s *Span) SpanID() string { return s.payload.SpanID }

// SessionID returns the span's session identifier, if set.
func (s *Span) SessionID() string { return s.payload.SessionID }

func (s *Span) SetInput(v any) *Span {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.payload.Input = v
	return s
}

func (s *Span) SetOutput(v any) *Span {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.payload.Output = v
	return s
}

func (s *Span) SetAttributes(attrs map[string]any) *Span {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.payload.Attributes == nil {
		s.payload.Attributes = map[string]any{}
	}
	for k, v := range attrs {
		s.payload.Attributes[k] = v
	}
	return s
}

func (s *Span) SetTags(tags ...string) *Span {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.payload.Tags = append(s.payload.Tags, tags...)
	return s
}

func (s *Span) SetModel(model string) *Span {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.payload.ModelName = model
	return s
}

// SetUsage records token usage. Pass -1 for an unknown count; total is computed
// from prompt+completion when total is negative and both are known.
func (s *Span) SetUsage(promptTokens, completionTokens, totalTokens int) *Span {
	s.mu.Lock()
	defer s.mu.Unlock()
	if promptTokens >= 0 {
		s.payload.PromptTokens = &promptTokens
	}
	if completionTokens >= 0 {
		s.payload.CompletionTokens = &completionTokens
	}
	if totalTokens < 0 && promptTokens >= 0 && completionTokens >= 0 {
		totalTokens = promptTokens + completionTokens
	}
	if totalTokens >= 0 {
		s.payload.TotalTokens = &totalTokens
	}
	return s
}

func (s *Span) SetCost(cost float64) *Span {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.payload.Cost = &cost
	return s
}

func (s *Span) SetSession(sessionID string) *Span {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.payload.SessionID = sessionID
	return s
}

// SetError marks the span as failed and records the error message/type.
func (s *Span) SetError(err error) *Span {
	if err == nil {
		return s
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.payload.Status = "error"
	s.payload.Level = "error"
	s.payload.ErrorMessage = err.Error()
	return s
}

// End finalizes the span and hands it to the exporter. It is idempotent.
func (s *Span) End() {
	s.mu.Lock()
	if s.ended {
		s.mu.Unlock()
		return
	}
	s.ended = true
	now := time.Now().UTC()
	s.payload.EndedAt = &now
	payload := s.payload
	onEnd := s.onEnd
	s.mu.Unlock()

	if onEnd != nil {
		onEnd(payload)
	}
}
