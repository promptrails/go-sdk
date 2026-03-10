package promptrails

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestPromptsList(t *testing.T) {
	srv, client := testServer(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": "p1", "name": "Prompt 1", "status": "active", "created_at": "2025-01-01T00:00:00Z", "updated_at": "2025-01-01T00:00:00Z"},
			},
			"meta": map[string]any{"total": 1, "page": 1, "limit": 20, "total_pages": 1},
		})
	})
	defer srv.Close()

	result, err := client.Prompts.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data) != 1 {
		t.Errorf("expected 1 prompt, got %d", len(result.Data))
	}
}

func TestPromptsGet(t *testing.T) {
	srv, client := testServer(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id": "p1", "name": "Test Prompt", "status": "active",
				"created_at": "2025-01-01T00:00:00Z", "updated_at": "2025-01-01T00:00:00Z",
			},
		})
	})
	defer srv.Close()

	prompt, err := client.Prompts.Get(context.Background(), "p1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prompt.Name != "Test Prompt" {
		t.Errorf("expected Test Prompt, got %s", prompt.Name)
	}
}

func TestPromptsCreate(t *testing.T) {
	srv, client := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id": "p2", "name": "New Prompt", "status": "draft",
				"created_at": "2025-01-01T00:00:00Z", "updated_at": "2025-01-01T00:00:00Z",
			},
		})
	})
	defer srv.Close()

	prompt, err := client.Prompts.Create(context.Background(), &CreatePromptParams{Name: "New Prompt"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prompt.Name != "New Prompt" {
		t.Errorf("expected New Prompt, got %s", prompt.Name)
	}
}

func TestPromptsDelete(t *testing.T) {
	srv, client := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(200)
	})
	defer srv.Close()

	err := client.Prompts.Delete(context.Background(), "p1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
