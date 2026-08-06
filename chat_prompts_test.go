package promptrails

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

// ---------- Chat (non-streaming) ----------

func TestChat(t *testing.T) {
	ctx := context.Background()

	t.Run("ListSessions", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/chat/sessions", listEnvelope))
		defer srv.Close()
		if _, err := c.Chat.ListSessions(ctx, nil); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("GetSession", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/chat/sessions/s1", `{"data":{"id":"s1"}}`))
		defer srv.Close()
		if _, err := c.Chat.GetSession(ctx, "s1"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("CreateSession", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "POST", "/api/v1/chat/sessions", `{"data":{"id":"s1"}}`))
		defer srv.Close()
		if _, err := c.Chat.CreateSession(ctx, &CreateSessionParams{AgentID: "a1"}); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("DeleteSession", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "DELETE", "/api/v1/chat/sessions/s1", ""))
		defer srv.Close()
		if err := c.Chat.DeleteSession(ctx, "s1"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("ListMessages", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/chat/sessions/s1/messages", listEnvelope))
		defer srv.Close()
		if _, err := c.Chat.ListMessages(ctx, "s1", nil); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("SendMessage", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "POST", "/api/v1/chat/sessions/s1/messages", `{"data":{"id":"m1"}}`))
		defer srv.Close()
		if _, err := c.Chat.SendMessage(ctx, "s1", &SendMessageParams{Content: "hi"}); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("SubmitFeedback", func(t *testing.T) {
		srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost || r.URL.Path != "/api/v1/chat/sessions/s1/feedback" {
				t.Fatalf("got %s %s", r.Method, r.URL.Path)
			}
			var payload SubmitFeedbackParams
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if payload.ExecutionID != "exec1" || payload.Value != -1 {
				t.Fatalf("unexpected payload: %#v", payload)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":{"submitted":true}}`))
		})
		defer srv.Close()

		result, err := c.Chat.SubmitFeedback(ctx, "s1", &SubmitFeedbackParams{
			ExecutionID: "exec1",
			Value:       -1,
		})
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if !result.Submitted {
			t.Fatal("expected submitted response")
		}
	})
}

// ---------- Prompts (extended methods) ----------

func TestPromptsExtended(t *testing.T) {
	ctx := context.Background()

	t.Run("Update", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "PATCH", "/api/v1/prompts/p1", `{"data":{"id":"p1"}}`))
		defer srv.Close()
		if _, err := c.Prompts.Update(ctx, "p1", &UpdatePromptParams{}); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("ListVersions", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/prompts/p1/versions", listEnvelope))
		defer srv.Close()
		if _, err := c.Prompts.ListVersions(ctx, "p1", nil); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("CreateVersion", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "POST", "/api/v1/prompts/p1/versions", `{"data":{"id":"v1"}}`))
		defer srv.Close()
		if _, err := c.Prompts.CreateVersion(ctx, "p1", &CreatePromptVersionParams{UserPrompt: "hi"}); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("PromoteVersion", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "POST", "/api/v1/prompts/p1/versions/v1/promote", ""))
		defer srv.Close()
		if err := c.Prompts.PromoteVersion(ctx, "p1", "v1"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
}
