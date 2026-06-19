package promptrails

import (
	"context"
	"io"
	"net/http"
	"testing"
)

// ---------- Agents (extended methods) ----------

func TestAgentsExtended(t *testing.T) {
	ctx := context.Background()

	const (
		obj  = `{"data":{}}` // single object
		arr  = `{"data":[]}` // bare array
		page = listEnvelope  // paginated
		none = ""            // no body (delete/promote)
	)

	cases := []struct {
		name string
		path string
		resp string
		call func(c *Client) error
	}{
		{"Update", "/api/v1/agents/a1", obj, func(c *Client) error {
			_, err := c.Agents.Update(ctx, "a1", &UpdateAgentParams{})
			return err
		}},
		{"Delete", "/api/v1/agents/a1", none, func(c *Client) error {
			return c.Agents.Delete(ctx, "a1")
		}},
		{"Execute", "/api/v1/agents/a1/execute", obj, func(c *Client) error {
			_, err := c.Agents.Execute(ctx, "a1", &ExecuteAgentParams{})
			return err
		}},
		{"ListVersions", "/api/v1/agents/a1/versions", page, func(c *Client) error {
			_, err := c.Agents.ListVersions(ctx, "a1", nil)
			return err
		}},
		{"CreateVersion", "/api/v1/agents/a1/versions", obj, func(c *Client) error {
			_, err := c.Agents.CreateVersion(ctx, "a1", &CreateVersionParams{})
			return err
		}},
		{"PromoteVersion", "/api/v1/agents/a1/versions/v1/promote", none, func(c *Client) error {
			return c.Agents.PromoteVersion(ctx, "a1", "v1")
		}},
		{"ListGuardrails", "/api/v1/agents/a1/guardrails", arr, func(c *Client) error {
			_, err := c.Agents.ListGuardrails(ctx, "a1")
			return err
		}},
		{"CreateGuardrail", "/api/v1/agents/a1/guardrails", obj, func(c *Client) error {
			_, err := c.Agents.CreateGuardrail(ctx, "a1", &CreateGuardrailParams{})
			return err
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tc.path {
					t.Errorf("path = %s, want %s", r.URL.Path, tc.path)
				}
				if tc.resp != "" {
					_, _ = io.WriteString(w, tc.resp)
				}
			})
			defer srv.Close()
			if err := tc.call(c); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// ---------- AgentTriggers ----------

func TestAgentTriggers(t *testing.T) {
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/triggers", listEnvelope))
		defer srv.Close()
		if _, err := c.AgentTriggers.List(ctx, nil); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("Get", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/triggers/t1", `{"data":{"id":"t1"}}`))
		defer srv.Close()
		if _, err := c.AgentTriggers.Get(ctx, "t1"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("Create", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "POST", "/api/v1/triggers", `{"data":{"id":"t1"}}`))
		defer srv.Close()
		if _, err := c.AgentTriggers.Create(ctx, &CreateAgentTriggerParams{}); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("Update", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "PATCH", "/api/v1/triggers/t1", `{"data":{"id":"t1"}}`))
		defer srv.Close()
		if _, err := c.AgentTriggers.Update(ctx, "t1", &UpdateAgentTriggerParams{}); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("Delete", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "DELETE", "/api/v1/triggers/t1", ""))
		defer srv.Close()
		if err := c.AgentTriggers.Delete(ctx, "t1"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
}

// ---------- DataSources ----------

func TestDataSources(t *testing.T) {
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/data-sources", listEnvelope))
		defer srv.Close()
		if _, err := c.DataSources.List(ctx, nil); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("Get", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/data-sources/d1", `{"data":{"id":"d1"}}`))
		defer srv.Close()
		if _, err := c.DataSources.Get(ctx, "d1"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("Create", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "POST", "/api/v1/data-sources", `{"data":{"id":"d1"}}`))
		defer srv.Close()
		if _, err := c.DataSources.Create(ctx, &CreateDataSourceParams{}); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("Update", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "PATCH", "/api/v1/data-sources/d1", `{"data":{"id":"d1"}}`))
		defer srv.Close()
		if _, err := c.DataSources.Update(ctx, "d1", &UpdateDataSourceParams{}); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("Delete", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "DELETE", "/api/v1/data-sources/d1", ""))
		defer srv.Close()
		if err := c.DataSources.Delete(ctx, "d1"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("ListVersions", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/data-sources/d1/versions", listEnvelope))
		defer srv.Close()
		if _, err := c.DataSources.ListVersions(ctx, "d1", nil); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("CreateVersion", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "POST", "/api/v1/data-sources/d1/versions", `{"data":{"id":"v1"}}`))
		defer srv.Close()
		if _, err := c.DataSources.CreateVersion(ctx, "d1", &CreateVersionParams{}); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("TestConnection", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "POST", "/api/v1/data-sources/d1/test", ""))
		defer srv.Close()
		if err := c.DataSources.TestConnection(ctx, "d1"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("Query", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "POST", "/api/v1/data-sources/d1/query", `{"data":{"rows":[]}}`))
		defer srv.Close()
		if _, err := c.DataSources.Query(ctx, "d1", &QueryDataSourceParams{}); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
}

// ---------- MCPTools ----------

func TestMCPTools(t *testing.T) {
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/mcp-tools", listEnvelope))
		defer srv.Close()
		if _, err := c.MCPTools.List(ctx, nil); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("Get", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/mcp-tools/m1", `{"data":{"id":"m1"}}`))
		defer srv.Close()
		if _, err := c.MCPTools.Get(ctx, "m1"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("Create", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "POST", "/api/v1/mcp-tools", `{"data":{"id":"m1"}}`))
		defer srv.Close()
		if _, err := c.MCPTools.Create(ctx, &CreateMCPToolParams{}); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("Update", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "PATCH", "/api/v1/mcp-tools/m1", `{"data":{"id":"m1"}}`))
		defer srv.Close()
		if _, err := c.MCPTools.Update(ctx, "m1", &UpdateMCPToolParams{}); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("Delete", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "DELETE", "/api/v1/mcp-tools/m1", ""))
		defer srv.Close()
		if err := c.MCPTools.Delete(ctx, "m1"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
}

// ---------- Approvals ----------

func TestApprovals(t *testing.T) {
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/approvals", listEnvelope))
		defer srv.Close()
		if _, err := c.Approvals.List(ctx, nil); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("Get", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/approvals/ap1", `{"data":{"id":"ap1"}}`))
		defer srv.Close()
		if _, err := c.Approvals.Get(ctx, "ap1"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("Decide", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "POST", "/api/v1/approvals/ap1/decide", `{"data":{"id":"ap1"}}`))
		defer srv.Close()
		if _, err := c.Approvals.Decide(ctx, "ap1", &DecideApprovalParams{}); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
}

// ---------- LLMModels / MediaModels / Media ----------

func TestLLMModels(t *testing.T) {
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/llm-models", listEnvelope))
		defer srv.Close()
		if _, err := c.LLMModels.List(ctx, nil); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
	t.Run("ListAvailable", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/llm-models/available", `{"data":{"groups":[{"provider":"openai"}]}}`))
		defer srv.Close()
		if _, err := c.LLMModels.ListAvailable(ctx); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
}

func TestMediaModels(t *testing.T) {
	srv, c := testServer(jsonHandler(t, "GET", "/api/v1/media-models", listEnvelope))
	defer srv.Close()
	if _, err := c.MediaModels.List(context.Background(), nil); err != nil {
		t.Fatalf("err=%v", err)
	}
}

func TestMediaGenerate(t *testing.T) {
	srv, c := testServer(jsonHandler(t, "POST", "/api/v1/media/generate", `{"data":{"job_id":"j1"}}`))
	defer srv.Close()
	if _, err := c.Media.Generate(context.Background(), &GenerateMediaParams{}); err != nil {
		t.Fatalf("err=%v", err)
	}
}
