package tracing

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// captureServer records every span posted to the ingest endpoint.
type captureServer struct {
	*httptest.Server
	mu    sync.Mutex
	spans []map[string]any
}

func newCaptureServer(t *testing.T) *captureServer {
	t.Helper()
	cs := &captureServer{}
	cs.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != ingestPath {
			http.NotFound(w, r)
			return
		}
		body, _ := io.ReadAll(r.Body)
		var payload struct {
			Spans []map[string]any `json:"spans"`
		}
		_ = json.Unmarshal(body, &payload)
		cs.mu.Lock()
		cs.spans = append(cs.spans, payload.Spans...)
		cs.mu.Unlock()
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"ingested":1}}`))
	}))
	t.Cleanup(cs.Close)
	return cs
}

func (cs *captureServer) byName() map[string]map[string]any {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	out := make(map[string]map[string]any, len(cs.spans))
	for _, s := range cs.spans {
		out[s["name"].(string)] = s
	}
	return out
}

func newTestTracer(url string) *Tracer {
	// Long flush interval so only the explicit Flush drives delivery.
	return NewTracer("key", WithBaseURL(url), WithFlushInterval(time.Hour))
}

func TestNestedSpansShareTraceAndLinkParent(t *testing.T) {
	cs := newCaptureServer(t)
	tracer := newTestTracer(cs.URL)

	err := tracer.Span(context.Background(), "agent-run", KindAgent, func(ctx context.Context, root *Span) error {
		root.SetInput(map[string]any{"q": "weather?"})
		return tracer.Span(ctx, "llm-call", KindLLM, func(ctx context.Context, llm *Span) error {
			llm.SetModel("gpt-4o").SetUsage(120, 30, -1).SetOutput(map[string]any{"text": "rainy"})
			return nil
		})
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tracer.Flush()

	spans := cs.byName()
	root, llm := spans["agent-run"], spans["llm-call"]
	if root == nil || llm == nil {
		t.Fatalf("expected both spans, got %v", spans)
	}
	if root["trace_id"] != llm["trace_id"] {
		t.Errorf("trace ids differ: %v vs %v", root["trace_id"], llm["trace_id"])
	}
	if _, ok := root["parent_span_id"]; ok {
		t.Errorf("root should have no parent")
	}
	if llm["parent_span_id"] != root["span_id"] {
		t.Errorf("llm parent %v != root span %v", llm["parent_span_id"], root["span_id"])
	}
	if llm["kind"] != "llm" || llm["model_name"] != "gpt-4o" {
		t.Errorf("unexpected llm fields: %v", llm)
	}
	if llm["total_tokens"].(float64) != 150 {
		t.Errorf("expected total_tokens 150, got %v", llm["total_tokens"])
	}
}

func TestSpanRecordsError(t *testing.T) {
	cs := newCaptureServer(t)
	tracer := newTestTracer(cs.URL)

	wantErr := errors.New("boom")
	err := tracer.Span(context.Background(), "tool", KindTool, func(ctx context.Context, span *Span) error {
		return wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected error to propagate, got %v", err)
	}
	tracer.Flush()

	span := cs.byName()["tool"]
	if span["status"] != "error" || span["error_message"] != "boom" {
		t.Errorf("expected error status/message, got %v", span)
	}
}

func TestTagsSentAsArray(t *testing.T) {
	cs := newCaptureServer(t)
	tracer := newTestTracer(cs.URL)

	_ = tracer.Span(context.Background(), "tagged", KindChain, func(ctx context.Context, span *Span) error {
		span.SetTags("prod", "checkout")
		return nil
	})
	tracer.Flush()

	tags, ok := cs.byName()["tagged"]["tags"].([]any)
	if !ok || len(tags) != 2 || tags[0] != "prod" || tags[1] != "checkout" {
		t.Errorf("expected [prod checkout], got %v", cs.byName()["tagged"]["tags"])
	}
}
