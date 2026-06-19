package promptrails

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
)

// sseServer streams the given SSE body and a client pointed at it.
func sseServer(t *testing.T, wantMethod, wantPath, sseBody string) (*Client, func()) {
	t.Helper()
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != wantMethod {
			t.Errorf("method = %s, want %s", r.Method, wantMethod)
		}
		if r.URL.Path != wantPath {
			t.Errorf("path = %s, want %s", r.URL.Path, wantPath)
		}
		if r.Header.Get("Accept") != "text/event-stream" {
			t.Errorf("Accept = %q, want text/event-stream", r.Header.Get("Accept"))
		}
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = io.WriteString(w, sseBody)
	})
	return c, srv.Close
}

func TestChatSendMessageStream(t *testing.T) {
	body := "event: execution\n" +
		`data: {"execution_id":"e1","user_message_id":"m1"}` + "\n\n" +
		"event: content\n" +
		`data: {"content":"Hello"}` + "\n\n" +
		"event: done\n" +
		`data: {"output":{"text":"Hello"},"token_usage":{"total_tokens":5}}` + "\n\n"

	c, closeFn := sseServer(t, "POST", "/api/v1/chat/sessions/s1/messages/stream", body)
	defer closeFn()

	stream, err := c.Chat.SendMessageStream(context.Background(), "s1", &SendMessageParams{Content: "hi"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer stream.Close()

	var gotExec, gotContent, gotDone bool
	for stream.Next() {
		switch e := stream.Event().(type) {
		case *ExecutionEvent:
			gotExec = e.ExecutionID == "e1"
		case *ContentEvent:
			gotContent = e.Content == "Hello"
		case *DoneEvent:
			gotDone = e.TokenUsage.TotalTokens == 5
		}
	}
	if err := stream.Err(); err != nil {
		t.Fatalf("stream error: %v", err)
	}
	if !gotExec || !gotContent || !gotDone {
		t.Errorf("missing events: exec=%v content=%v done=%v", gotExec, gotContent, gotDone)
	}
}

func TestExecutionsStream(t *testing.T) {
	body := "event: content\n" + `data: {"content":"tick"}` + "\n\n"

	c, closeFn := sseServer(t, "GET", "/api/v1/executions/ex1/stream", body)
	defer closeFn()

	stream, err := c.Executions.Stream(context.Background(), "ex1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer stream.Close()

	var content string
	for stream.Next() {
		if e, ok := stream.Event().(*ContentEvent); ok {
			content += e.Content
		}
	}
	if content != "tick" {
		t.Errorf("content = %q, want tick", content)
	}
}

// A non-2xx status on a stream open must surface as an APIError, not a stream.
func TestSendMessageStream_HTTPError(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = io.WriteString(w, `{"error":{"message":"forbidden","code":"forbidden"}}`)
	})
	defer srv.Close()

	_, err := c.Chat.SendMessageStream(context.Background(), "s1", &SendMessageParams{Content: "hi"})
	if err == nil {
		t.Fatal("expected error on 403 stream open")
	}
	var forbidden *ForbiddenError
	if !errors.As(err, &forbidden) || forbidden.StatusCode != http.StatusForbidden {
		t.Errorf("err = %v, want *ForbiddenError with status 403", err)
	}
}

// A terminal error frame is delivered as an ErrorEvent through the stream.
func TestStream_ServerErrorFrame(t *testing.T) {
	body := "event: error\n" + `data: {"message":"quota exceeded"}` + "\n\n"
	c, closeFn := sseServer(t, "POST", "/api/v1/chat/sessions/s1/messages/stream", body)
	defer closeFn()

	stream, err := c.Chat.SendMessageStream(context.Background(), "s1", &SendMessageParams{Content: "hi"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer stream.Close()

	var msg string
	for stream.Next() {
		if e, ok := stream.Event().(*ErrorEvent); ok {
			msg = e.Message
		}
	}
	if msg != "quota exceeded" {
		t.Errorf("error event message = %q, want 'quota exceeded'", msg)
	}
}
