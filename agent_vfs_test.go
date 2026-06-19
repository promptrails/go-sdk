package promptrails

import (
	"context"
	"io"
	"net/http"
	"testing"
)

func TestAgentVFS(t *testing.T) {
	ctx := context.Background()
	const base = "/api/v1/agents/ag1/vfs"

	t.Run("List", func(t *testing.T) {
		srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != base {
				t.Errorf("path=%s", r.URL.Path)
			}
			if r.URL.Query().Get("recursive") != "true" || r.URL.Query().Get("path") != "/docs" {
				t.Errorf("query not propagated: %v", r.URL.Query())
			}
			_, _ = io.WriteString(w, `{"data":{"path":"/docs","items":[{"path":"/docs/a"}],"total":1}}`)
		})
		defer srv.Close()
		res, err := c.AgentVFS.List(ctx, "ag1", &ListAgentVFSParams{Path: "/docs", Recursive: true})
		if err != nil || res.Total != 1 || len(res.Items) != 1 {
			t.Fatalf("res=%+v err=%v", res, err)
		}
	})

	t.Run("Read", func(t *testing.T) {
		srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != base+"/file" || r.URL.Query().Get("path") != "/a.txt" {
				t.Errorf("path/query=%s ?%s", r.URL.Path, r.URL.RawQuery)
			}
			_, _ = io.WriteString(w, `{"data":{"path":"/a.txt","content":"hi"}}`)
		})
		defer srv.Close()
		if _, err := c.AgentVFS.Read(ctx, "ag1", "/a.txt", nil); err != nil {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("Read with line range", func(t *testing.T) {
		srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			if q.Get("line_offset") != "10" || q.Get("line_limit") != "5" {
				t.Errorf("line range params not propagated: %v", q)
			}
			_, _ = io.WriteString(w, `{"data":{"path":"/a.txt","content":"x"}}`)
		})
		defer srv.Close()
		_, err := c.AgentVFS.Read(ctx, "ag1", "/a.txt", &ReadAgentVFSParams{LineOffset: 10, LineLimit: 5})
		if err != nil {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("Stat", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", base+"/stat", `{"data":{"path":"/a.txt"}}`))
		defer srv.Close()
		if _, err := c.AgentVFS.Stat(ctx, "ag1", "/a.txt"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("Write", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "PUT", base+"/file", `{"data":{"path":"/a.txt"}}`))
		defer srv.Close()
		params := &WriteAgentVFSParams{Path: "/a.txt", Content: "hi"}
		if _, err := c.AgentVFS.Write(ctx, "ag1", params); err != nil {
			t.Fatalf("err=%v", err)
		}
		// Default write mode should be filled in.
		if params.Mode != AgentVFSWriteOverwrite {
			t.Errorf("Mode = %q, want default overwrite", params.Mode)
		}
	})

	t.Run("Write nil params", func(t *testing.T) {
		_, c := testServer(jsonHandler(t, "PUT", base+"/file", ""))
		if _, err := c.AgentVFS.Write(ctx, "ag1", nil); err == nil {
			t.Error("expected error for nil params")
		}
	})

	t.Run("Mkdir", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "POST", base+"/mkdir", `{"data":{"path":"/d"}}`))
		defer srv.Close()
		if _, err := c.AgentVFS.Mkdir(ctx, "ag1", "/d"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("Move", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "POST", base+"/move", ""))
		defer srv.Close()
		if err := c.AgentVFS.Move(ctx, "ag1", "/a", "/b"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("Copy", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "POST", base+"/copy", ""))
		defer srv.Close()
		if err := c.AgentVFS.Copy(ctx, "ag1", "/a", "/b"); err != nil {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("Delete recursive", func(t *testing.T) {
		srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "DELETE" || r.URL.Path != base {
				t.Errorf("method/path=%s %s", r.Method, r.URL.Path)
			}
			if r.URL.Query().Get("recursive") != "true" || r.URL.Query().Get("path") != "/d" {
				t.Errorf("query not propagated: %v", r.URL.Query())
			}
		})
		defer srv.Close()
		if err := c.AgentVFS.Delete(ctx, "ag1", "/d", true); err != nil {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("Grep", func(t *testing.T) {
		srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != base+"/grep" || r.URL.Query().Get("q") != "todo" {
				t.Errorf("path/query=%s ?%s", r.URL.Path, r.URL.RawQuery)
			}
			_, _ = io.WriteString(w, `{"data":{"query":"todo","matches":[]}}`)
		})
		defer srv.Close()
		if _, err := c.AgentVFS.Grep(ctx, "ag1", &GrepAgentVFSParams{Query: "todo"}); err != nil {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("Grep empty query", func(t *testing.T) {
		_, c := testServer(jsonHandler(t, "GET", base+"/grep", ""))
		if _, err := c.AgentVFS.Grep(ctx, "ag1", &GrepAgentVFSParams{}); err == nil {
			t.Error("expected error for empty query")
		}
	})

	t.Run("Glob", func(t *testing.T) {
		srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != base+"/glob" || r.URL.Query().Get("pattern") != "*.go" {
				t.Errorf("path/query=%s ?%s", r.URL.Path, r.URL.RawQuery)
			}
			_, _ = io.WriteString(w, `{"data":{"pattern":"*.go","items":[]}}`)
		})
		defer srv.Close()
		if _, err := c.AgentVFS.Glob(ctx, "ag1", &GlobAgentVFSParams{Pattern: "*.go"}); err != nil {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("Glob empty pattern", func(t *testing.T) {
		_, c := testServer(jsonHandler(t, "GET", base+"/glob", ""))
		if _, err := c.AgentVFS.Glob(ctx, "ag1", &GlobAgentVFSParams{}); err == nil {
			t.Error("expected error for empty pattern")
		}
	})

	t.Run("Usage", func(t *testing.T) {
		srv, c := testServer(jsonHandler(t, "GET", base+"/usage", `{"data":{"bytes_used":42}}`))
		defer srv.Close()
		used, err := c.AgentVFS.Usage(ctx, "ag1")
		if err != nil || used != 42 {
			t.Fatalf("used=%d err=%v", used, err)
		}
	})
}
