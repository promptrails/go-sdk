package promptrails

import (
	"encoding/json"
	"io"
	"strings"
	"testing"
)

// nopCloser wraps an io.Reader so ChatStream can Close it.
type nopCloser struct{ io.Reader }

func (nopCloser) Close() error { return nil }

func collectStream(t *testing.T, body string) []ChatStreamEvent {
	t.Helper()
	stream := newChatStream(nopCloser{strings.NewReader(body)})
	defer stream.Close()
	var events []ChatStreamEvent
	for stream.Next() {
		events = append(events, stream.Event())
	}
	if err := stream.Err(); err != nil {
		t.Fatalf("stream error: %v", err)
	}
	return events
}

func TestStreamParsesAllEventTypes(t *testing.T) {
	body := "event: execution\n" +
		`data: {"execution_id":"e1","user_message_id":"m1"}` + "\n\n" +
		"event: thinking\n" +
		`data: {"content":"Let me check"}` + "\n\n" +
		"event: tool_start\n" +
		`data: {"id":"t1","name":"search"}` + "\n\n" +
		"event: tool_end\n" +
		`data: {"id":"t1","name":"search","summary":"3 hits"}` + "\n\n" +
		"event: content\n" +
		`data: {"content":"Hello"}` + "\n\n" +
		"event: done\n" +
		`data: {"output":{"content":"Hello world"},"token_usage":{"total_tokens":42}}` + "\n\n"

	events := collectStream(t, body)
	if len(events) != 6 {
		t.Fatalf("expected 6 events, got %d", len(events))
	}

	if e, ok := events[0].(*ExecutionEvent); !ok || e.ExecutionID != "e1" || e.UserMessageID != "m1" {
		t.Errorf("bad execution event: %#v", events[0])
	}
	if e, ok := events[1].(*ThinkingEvent); !ok || e.Content != "Let me check" {
		t.Errorf("bad thinking event: %#v", events[1])
	}
	if e, ok := events[2].(*ToolStartEvent); !ok || e.ID != "t1" || e.Name != "search" {
		t.Errorf("bad tool_start event: %#v", events[2])
	}
	if e, ok := events[3].(*ToolEndEvent); !ok || e.Summary != "3 hits" {
		t.Errorf("bad tool_end event: %#v", events[3])
	}
	if e, ok := events[4].(*ContentEvent); !ok || e.Content != "Hello" {
		t.Errorf("bad content event: %#v", events[4])
	}
	if e, ok := events[5].(*DoneEvent); !ok ||
		e.Output["content"] != "Hello world" ||
		e.TokenUsage.TotalTokens != 42 {
		t.Errorf("bad done event: %#v", events[5])
	}
}

func TestStreamSkipsUnknownEventTypes(t *testing.T) {
	body := "event: ping\ndata: {\"ok\":true}\n\n" +
		"event: content\ndata: {\"content\":\"hi\"}\n\n"
	events := collectStream(t, body)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if _, ok := events[0].(*ContentEvent); !ok {
		t.Errorf("expected ContentEvent, got %T", events[0])
	}
}

func TestStreamErrorEvent(t *testing.T) {
	events := collectStream(t, "event: error\ndata: {\"message\":\"quota exceeded\"}\n\n")
	e, ok := events[0].(*ErrorEvent)
	if !ok {
		t.Fatalf("expected ErrorEvent, got %T", events[0])
	}
	if e.Message != "quota exceeded" {
		t.Errorf("bad error message: %q", e.Message)
	}
}

func TestAgentConfigMarshalInjectsTypeTag(t *testing.T) {
	cases := []struct {
		cfg     AgentConfig
		wantTyp string
	}{
		{SimpleAgentConfig{PromptID: "p1"}, "simple"},
		{ChainAgentConfig{PromptIDs: []PromptLink{{PromptID: "p1", Role: "main", SortOrder: 0}}}, "chain"},
		{MultiAgentConfig{PromptIDs: []PromptLink{{PromptID: "p1", Role: "a", SortOrder: 0}}}, "multi_agent"},
		{WorkflowAgentConfig{Nodes: []WorkflowNode{{ID: "n1", DependsOn: []string{}}}}, "workflow"},
		{CompositeAgentConfig{Steps: []CompositeStep{{ID: "s1", AgentID: "a1"}}}, "composite"},
	}
	for _, tc := range cases {
		b, err := marshalJSONViaParams(tc.cfg)
		if err != nil {
			t.Fatalf("marshal %T: %v", tc.cfg, err)
		}
		if !strings.Contains(b, `"type":"`+tc.wantTyp+`"`) {
			t.Errorf("%T did not emit type=%q — got: %s", tc.cfg, tc.wantTyp, b)
		}
	}
}

func marshalJSONViaParams(cfg AgentConfig) (string, error) {
	params := &CreateVersionParams{Version: "1.0.0", Config: cfg}
	b, err := json.Marshal(params)
	return string(b), err
}
