package promptrails

import (
	"context"
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
		if _, err := c.Prompts.CreateVersion(ctx, "p1", &CreateVersionParams{}); err != nil {
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
	t.Run("Run", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "POST", "/api/v1/prompts/p1/run",
			`{"data":{"output":"hello"}}`))
		defer srv.Close()
		res, err := c.Prompts.Run(ctx, "p1", &RunPromptParams{})
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if res.Output != "hello" {
			t.Errorf("run result output = %q, want hello", res.Output)
		}
	})
}
