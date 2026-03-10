package promptrails

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testServer(handler http.HandlerFunc) (*httptest.Server, *Client) {
	srv := httptest.NewServer(handler)
	client := NewClient("test-key", WithBaseURL(srv.URL), WithMaxRetries(0))
	return srv, client
}

func TestAgentsList(t *testing.T) {
	srv, client := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.Header.Get("X-API-Key") != "test-key" {
			t.Error("missing API key header")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": "a1", "name": "Agent 1", "type": "simple", "status": "active", "created_at": "2025-01-01T00:00:00Z", "updated_at": "2025-01-01T00:00:00Z"},
				{"id": "a2", "name": "Agent 2", "type": "chain", "status": "draft", "created_at": "2025-01-01T00:00:00Z", "updated_at": "2025-01-01T00:00:00Z"},
			},
			"meta": map[string]any{"total": 2, "page": 1, "limit": 20, "total_pages": 1},
		})
	})
	defer srv.Close()

	result, err := client.Agents.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data) != 2 {
		t.Errorf("expected 2 agents, got %d", len(result.Data))
	}
	if result.Meta.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Meta.Total)
	}
	if result.Data[0].Name != "Agent 1" {
		t.Errorf("expected Agent 1, got %s", result.Data[0].Name)
	}
}

func TestAgentsGet(t *testing.T) {
	srv, client := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/agents/a1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id": "a1", "name": "Test Agent", "type": "simple", "status": "active",
				"created_at": "2025-01-01T00:00:00Z", "updated_at": "2025-01-01T00:00:00Z",
			},
		})
	})
	defer srv.Close()

	agent, err := client.Agents.Get(context.Background(), "a1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if agent.ID != "a1" {
		t.Errorf("expected id a1, got %s", agent.ID)
	}
	if agent.Name != "Test Agent" {
		t.Errorf("expected Test Agent, got %s", agent.Name)
	}
}

func TestAgentsCreate(t *testing.T) {
	srv, client := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id": "a3", "name": "New Agent", "type": "simple", "status": "draft",
				"created_at": "2025-01-01T00:00:00Z", "updated_at": "2025-01-01T00:00:00Z",
			},
		})
	})
	defer srv.Close()

	agent, err := client.Agents.Create(context.Background(), &CreateAgentParams{Name: "New Agent", Type: "simple"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if agent.Name != "New Agent" {
		t.Errorf("expected New Agent, got %s", agent.Name)
	}
}

func TestAgentsDelete(t *testing.T) {
	srv, client := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(200)
	})
	defer srv.Close()

	err := client.Agents.Delete(context.Background(), "a1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAgentsExecute(t *testing.T) {
	srv, client := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/agents/a1/execute" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"output":       map[string]any{"result": "hello"},
				"trace_id":     "t1",
				"execution_id": "e1",
				"status":       "completed",
				"cost":         0.001,
			},
		})
	})
	defer srv.Close()

	result, err := client.Agents.Execute(context.Background(), "a1", &ExecuteAgentParams{
		Input: map[string]any{"prompt": "hi"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TraceID != "t1" {
		t.Errorf("expected trace_id t1, got %s", result.TraceID)
	}
	if result.Status != "completed" {
		t.Errorf("expected status completed, got %s", result.Status)
	}
}

func TestAgentsNotFound(t *testing.T) {
	srv, client := testServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]any{"error": map[string]any{"message": "agent not found"}})
	})
	defer srv.Close()

	_, err := client.Agents.Get(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error")
	}
	var notFound *NotFoundError
	if !isNotFoundError(err, &notFound) {
		t.Errorf("expected NotFoundError, got %T", err)
	}
}

func isNotFoundError(err error, target **NotFoundError) bool {
	e, ok := err.(*NotFoundError)
	if ok {
		*target = e
	}
	return ok
}
