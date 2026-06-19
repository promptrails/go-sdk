package promptrails

import (
	"context"
	"io"
	"net/http"
	"testing"
)

// jsonHandler asserts method + path (+ API key) and writes a fixed body.
func jsonHandler(t *testing.T, wantMethod, wantPath, respBody string) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != wantMethod {
			t.Errorf("method = %s, want %s", r.Method, wantMethod)
		}
		if r.URL.Path != wantPath {
			t.Errorf("path = %s, want %s", r.URL.Path, wantPath)
		}
		if r.Header.Get("X-API-Key") != "test-key" {
			t.Error("missing API key header")
		}
		if respBody != "" {
			_, _ = io.WriteString(w, respBody)
		}
	}
}

const listEnvelope = `{"data":[{"id":"x1"}],"meta":{"total":1,"page":1,"limit":20,"total_pages":1}}`

// ---------- A2A ----------

func TestA2A(t *testing.T) {
	ctx := context.Background()

	t.Run("GetAgentCard", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/a2a/agents/ag1/card", `{"data":{"name":"Card","version":"1.0"}}`))
		defer srv.Close()
		card, err := c.A2A.GetAgentCard(ctx, "ag1")
		if err != nil || card.Name != "Card" {
			t.Fatalf("card=%+v err=%v", card, err)
		}
	})

	t.Run("SendMessage", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "POST", "/api/v1/a2a/messages", `{"data":{"id":"t1","status":"queued"}}`))
		defer srv.Close()
		task, err := c.A2A.SendMessage(ctx, &A2ASendMessageParams{AgentID: "ag1", Message: map[string]any{"text": "hi"}})
		if err != nil || task.ID != "t1" {
			t.Fatalf("task=%+v err=%v", task, err)
		}
	})

	t.Run("GetTask", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/a2a/tasks/t1", `{"data":{"id":"t1","status":"completed"}}`))
		defer srv.Close()
		task, err := c.A2A.GetTask(ctx, "t1")
		if err != nil || task.Status != "completed" {
			t.Fatalf("task=%+v err=%v", task, err)
		}
	})

	t.Run("ListTasks", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/a2a/tasks", listEnvelope))
		defer srv.Close()
		res, err := c.A2A.ListTasks(ctx, nil)
		if err != nil || res.Meta.Total != 1 {
			t.Fatalf("res=%+v err=%v", res, err)
		}
	})

	t.Run("CancelTask", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "POST", "/api/v1/a2a/tasks/t1/cancel", ""))
		defer srv.Close()
		if err := c.A2A.CancelTask(ctx, "t1"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
}

// ---------- Assets ----------

func TestAssets(t *testing.T) {
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("media_type") != "image" {
				t.Errorf("media_type filter not propagated: %q", r.URL.Query().Get("media_type"))
			}
			_, _ = io.WriteString(w, listEnvelope)
		})
		defer srv.Close()
		res, err := c.Assets.List(ctx, &ListAssetsParams{MediaType: "image"})
		if err != nil || len(res.Data) != 1 {
			t.Fatalf("res=%+v err=%v", res, err)
		}
	})

	t.Run("Get", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/assets/a1", `{"data":{"id":"a1","file_name":"x.png"}}`))
		defer srv.Close()
		a, err := c.Assets.Get(ctx, "a1")
		if err != nil || a.FileName != "x.png" {
			t.Fatalf("a=%+v err=%v", a, err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "DELETE", "/api/v1/assets/a1", ""))
		defer srv.Close()
		if err := c.Assets.Delete(ctx, "a1"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("GetSignedURL", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/assets/a1/signed-url", `{"data":{"url":"https://signed"}}`))
		defer srv.Close()
		s, err := c.Assets.GetSignedURL(ctx, "a1")
		if err != nil || s.URL != "https://signed" {
			t.Fatalf("s=%+v err=%v", s, err)
		}
	})
}

// ---------- Costs ----------

func TestCosts(t *testing.T) {
	ctx := context.Background()

	t.Run("GetSummary", func(t *testing.T) {
		srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/v1/costs/summary" {
				t.Errorf("path=%s", r.URL.Path)
			}
			if r.URL.Query().Get("from") != "2025-01-01" {
				t.Errorf("from filter not propagated: %q", r.URL.Query().Get("from"))
			}
			_, _ = io.WriteString(w, `{"data":{"total_cost":12.5,"total_tokens":100}}`)
		})
		defer srv.Close()
		s, err := c.Costs.GetSummary(ctx, &CostParams{From: "2025-01-01"})
		if err != nil || s.TotalCost != 12.5 {
			t.Fatalf("s=%+v err=%v", s, err)
		}
	})

	t.Run("GetAgentSummary", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/costs/agents/ag1/summary", `{"data":{"total_cost":3.0}}`))
		defer srv.Close()
		s, err := c.Costs.GetAgentSummary(ctx, "ag1", nil)
		if err != nil || s.TotalCost != 3.0 {
			t.Fatalf("s=%+v err=%v", s, err)
		}
	})
}

// ---------- Credentials ----------

func TestCredentials(t *testing.T) {
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/credentials", listEnvelope))
		defer srv.Close()
		res, err := c.Credentials.List(ctx, nil)
		if err != nil || res.Meta.Total != 1 {
			t.Fatalf("res=%+v err=%v", res, err)
		}
	})

	t.Run("Get", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/credentials/c1", `{"data":{"id":"c1","provider":"openai"}}`))
		defer srv.Close()
		cr, err := c.Credentials.Get(ctx, "c1")
		if err != nil || cr.Provider != "openai" {
			t.Fatalf("cr=%+v err=%v", cr, err)
		}
	})

	t.Run("Create", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "POST", "/api/v1/credentials", `{"data":{"id":"c1","name":"key"}}`))
		defer srv.Close()
		cr, err := c.Credentials.Create(ctx, &CreateCredentialParams{Name: "key", Provider: "openai", APIKey: "sk"})
		if err != nil || cr.Name != "key" {
			t.Fatalf("cr=%+v err=%v", cr, err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "PATCH", "/api/v1/credentials/c1", `{"data":{"id":"c1","name":"new"}}`))
		defer srv.Close()
		name := "new"
		cr, err := c.Credentials.Update(ctx, "c1", &UpdateCredentialParams{Name: &name})
		if err != nil || cr.Name != "new" {
			t.Fatalf("cr=%+v err=%v", cr, err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "DELETE", "/api/v1/credentials/c1", ""))
		defer srv.Close()
		if err := c.Credentials.Delete(ctx, "c1"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("SetDefault", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "POST", "/api/v1/credentials/c1/default", ""))
		defer srv.Close()
		if err := c.Credentials.SetDefault(ctx, "c1"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("CheckConnection", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "POST", "/api/v1/credentials/c1/check", `{"data":{"id":"c1"}}`))
		defer srv.Close()
		if _, err := c.Credentials.CheckConnection(ctx, "c1"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
}

// ---------- Traces ----------

func TestTraces(t *testing.T) {
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("kind") != "llm" {
				t.Errorf("kind filter not propagated: %q", r.URL.Query().Get("kind"))
			}
			_, _ = io.WriteString(w, listEnvelope)
		})
		defer srv.Close()
		res, err := c.Traces.List(ctx, &ListTracesParams{Kind: "llm"})
		if err != nil || len(res.Data) != 1 {
			t.Fatalf("res=%+v err=%v", res, err)
		}
	})

	t.Run("GetByTraceID", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/traces/tr1", `{"data":[{"id":"s1","trace_id":"tr1"}]}`))
		defer srv.Close()
		traces, err := c.Traces.GetByTraceID(ctx, "tr1")
		if err != nil || len(traces) != 1 || traces[0].TraceID != "tr1" {
			t.Fatalf("traces=%+v err=%v", traces, err)
		}
	})
}

// ---------- Scores ----------

func TestScores(t *testing.T) {
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("execution_id") != "e1" {
				t.Errorf("execution_id filter not propagated")
			}
			_, _ = io.WriteString(w, listEnvelope)
		})
		defer srv.Close()
		res, err := c.Scores.List(ctx, &ListScoresParams{ExecutionID: "e1"})
		if err != nil || res.Meta.Total != 1 {
			t.Fatalf("res=%+v err=%v", res, err)
		}
	})

	t.Run("Get", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/scores/s1", `{"data":{"id":"s1","value":0.9}}`))
		defer srv.Close()
		s, err := c.Scores.Get(ctx, "s1")
		if err != nil || s.Value != 0.9 {
			t.Fatalf("s=%+v err=%v", s, err)
		}
	})

	t.Run("Create", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "POST", "/api/v1/scores", `{"data":{"id":"s1","value":1}}`))
		defer srv.Close()
		if _, err := c.Scores.Create(ctx, &CreateScoreParams{ExecutionID: "e1", Value: 1}); err != nil {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "PATCH", "/api/v1/scores/s1", `{"data":{"id":"s1"}}`))
		defer srv.Close()
		v := 0.5
		if _, err := c.Scores.Update(ctx, "s1", &UpdateScoreParams{Value: &v}); err != nil {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "DELETE", "/api/v1/scores/s1", ""))
		defer srv.Close()
		if err := c.Scores.Delete(ctx, "s1"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("Aggregates", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/scores/aggregates", `{"data":[{"config_id":"cfg1","count":3,"average":0.8}]}`))
		defer srv.Close()
		aggs, err := c.Scores.Aggregates(ctx, nil)
		if err != nil || len(aggs) != 1 || aggs[0].ConfigID != "cfg1" {
			t.Fatalf("aggs=%+v err=%v", aggs, err)
		}
	})

	t.Run("ListConfigs", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/scores/configs", listEnvelope))
		defer srv.Close()
		if _, err := c.Scores.ListConfigs(ctx, nil); err != nil {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("GetConfig", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", "/api/v1/scores/configs/cfg1", `{"data":{"id":"cfg1"}}`))
		defer srv.Close()
		if _, err := c.Scores.GetConfig(ctx, "cfg1"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("CreateConfig", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "POST", "/api/v1/scores/configs", `{"data":{"id":"cfg1"}}`))
		defer srv.Close()
		if _, err := c.Scores.CreateConfig(ctx, &CreateScoreConfigParams{Name: "acc", DataType: "numeric"}); err != nil {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("UpdateConfig", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "PATCH", "/api/v1/scores/configs/cfg1", `{"data":{"id":"cfg1"}}`))
		defer srv.Close()
		name := "renamed"
		if _, err := c.Scores.UpdateConfig(ctx, "cfg1", &UpdateScoreConfigParams{Name: &name}); err != nil {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("DeleteConfig", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "DELETE", "/api/v1/scores/configs/cfg1", ""))
		defer srv.Close()
		if err := c.Scores.DeleteConfig(ctx, "cfg1"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})
}
