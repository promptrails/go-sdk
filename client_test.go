package promptrails

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	c := NewClient("test-key")
	if c.Agents == nil {
		t.Error("Agents service should not be nil")
	}
	if c.Prompts == nil {
		t.Error("Prompts service should not be nil")
	}
	if c.Executions == nil {
		t.Error("Executions service should not be nil")
	}
	if c.Credentials == nil {
		t.Error("Credentials service should not be nil")
	}
	if c.DataSources == nil {
		t.Error("DataSources service should not be nil")
	}
	if c.Chat == nil {
		t.Error("Chat service should not be nil")
	}
	if c.Traces == nil {
		t.Error("Traces service should not be nil")
	}
	if c.Costs == nil {
		t.Error("Costs service should not be nil")
	}
	if c.Scores == nil {
		t.Error("Scores service should not be nil")
	}
	if c.MCPTools == nil {
		t.Error("MCPTools service should not be nil")
	}
	if c.Approvals == nil {
		t.Error("Approvals service should not be nil")
	}
	if c.WebhookTriggers == nil {
		t.Error("WebhookTriggers service should not be nil")
	}
	if c.A2A == nil {
		t.Error("A2A service should not be nil")
	}
}

func TestNewClientWithOptions(t *testing.T) {
	c := NewClient("test-key",
		WithBaseURL("http://localhost:8082"),
		WithTimeout(10*time.Second),
		WithMaxRetries(5),
	)
	if c.http.baseURL != "http://localhost:8082" {
		t.Errorf("expected base URL http://localhost:8082, got %s", c.http.baseURL)
	}
	if c.http.maxRetries != 5 {
		t.Errorf("expected max retries 5, got %d", c.http.maxRetries)
	}
	if c.http.apiKey != "test-key" {
		t.Errorf("expected api key test-key, got %s", c.http.apiKey)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()
	if cfg.baseURL != defaultBaseURL {
		t.Errorf("expected base URL %s, got %s", defaultBaseURL, cfg.baseURL)
	}
	if cfg.timeout != 30*time.Second {
		t.Errorf("expected timeout 30s, got %s", cfg.timeout)
	}
	if cfg.maxRetries != 3 {
		t.Errorf("expected max retries 3, got %d", cfg.maxRetries)
	}
}
