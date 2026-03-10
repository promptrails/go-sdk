package promptrails

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestExecutionsList(t *testing.T) {
	srv, client := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": "e1", "agent_id": "a1", "status": "completed", "created_at": "2025-01-01T00:00:00Z", "updated_at": "2025-01-01T00:00:00Z"},
			},
			"meta": map[string]any{"total": 1, "page": 1, "limit": 20, "total_pages": 1},
		})
	})
	defer srv.Close()

	result, err := client.Executions.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data) != 1 {
		t.Errorf("expected 1 execution, got %d", len(result.Data))
	}
	if result.Data[0].Status != "completed" {
		t.Errorf("expected completed, got %s", result.Data[0].Status)
	}
}

func TestExecutionsGet(t *testing.T) {
	srv, client := testServer(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id": "e1", "agent_id": "a1", "status": "completed",
				"input":      map[string]any{"query": "test"},
				"output":     map[string]any{"result": "ok"},
				"created_at": "2025-01-01T00:00:00Z", "updated_at": "2025-01-01T00:00:00Z",
			},
		})
	})
	defer srv.Close()

	exec, err := client.Executions.Get(context.Background(), "e1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exec.ID != "e1" {
		t.Errorf("expected e1, got %s", exec.ID)
	}
	if exec.Status != "completed" {
		t.Errorf("expected completed, got %s", exec.Status)
	}
}
