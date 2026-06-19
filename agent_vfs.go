package promptrails

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// AgentVFSService handles the per-agent Virtual Filesystem API.
//
// The VFS is a per-agent persistent file tree the agent reads and writes
// through builtin tools. Files survive across executions. Every operation
// is also reachable from this service for tooling and the studio.
type AgentVFSService struct {
	http *httpClient
}

// AgentVFSListResult is the response shape of a list call.
type AgentVFSListResult struct {
	Path  string         `json:"path"`
	Items []AgentVFSFile `json:"items"`
	Total int64          `json:"total"`
}

// AgentVFSGrepResult is the response shape of a grep call.
type AgentVFSGrepResult struct {
	Query   string              `json:"query"`
	Matches []AgentVFSGrepMatch `json:"matches"`
}

// AgentVFSGlobResult is the response shape of a glob call.
type AgentVFSGlobResult struct {
	Pattern string         `json:"pattern"`
	Items   []AgentVFSFile `json:"items"`
}

// ListAgentVFSParams scopes a List call.
type ListAgentVFSParams struct {
	Path      string
	Recursive bool
	Offset    int
	Limit     int
}

// ReadAgentVFSParams optionally scopes a Read call to a line range.
type ReadAgentVFSParams struct {
	LineOffset int
	LineLimit  int
}

// WriteAgentVFSParams is the body of a write call.
type WriteAgentVFSParams struct {
	Path     string            `json:"path"`
	Content  string            `json:"content"`
	Mode     AgentVFSWriteMode `json:"mode,omitempty"`
	MimeType string            `json:"mime_type,omitempty"`
}

// GrepAgentVFSParams scopes a Grep call.
type GrepAgentVFSParams struct {
	Query string
	Path  string
	Limit int
}

// GlobAgentVFSParams scopes a Glob call.
type GlobAgentVFSParams struct {
	Pattern string
	Path    string
	Limit   int
}

// List returns directory entries.
func (s *AgentVFSService) List(ctx context.Context, agentID string, params *ListAgentVFSParams) (*AgentVFSListResult, error) {
	qp := map[string]string{}
	if params != nil {
		if params.Path != "" {
			qp["path"] = params.Path
		}
		if params.Recursive {
			qp["recursive"] = "true"
		}
		if params.Offset > 0 {
			qp["offset"] = strconv.Itoa(params.Offset)
		}
		if params.Limit > 0 {
			qp["limit"] = strconv.Itoa(params.Limit)
		}
	}
	var result AgentVFSListResult
	err := s.http.get(ctx, fmt.Sprintf("/api/v1/agents/%s/vfs", agentID), qp, &result)
	return &result, err
}

// Read returns a file's content (optionally a line range).
func (s *AgentVFSService) Read(ctx context.Context, agentID, path string, params *ReadAgentVFSParams) (*AgentVFSReadResult, error) {
	qp := map[string]string{"path": path}
	if params != nil {
		if params.LineOffset > 0 {
			qp["line_offset"] = strconv.Itoa(params.LineOffset)
		}
		if params.LineLimit > 0 {
			qp["line_limit"] = strconv.Itoa(params.LineLimit)
		}
	}
	var result AgentVFSReadResult
	err := s.http.get(ctx, fmt.Sprintf("/api/v1/agents/%s/vfs/file", agentID), qp, &result)
	return &result, err
}

// Stat returns metadata for a single path.
func (s *AgentVFSService) Stat(ctx context.Context, agentID, path string) (*AgentVFSFile, error) {
	var result AgentVFSFile
	err := s.http.get(ctx, fmt.Sprintf("/api/v1/agents/%s/vfs/stat", agentID), map[string]string{"path": path}, &result)
	return &result, err
}

// Write creates or updates a file.
func (s *AgentVFSService) Write(ctx context.Context, agentID string, params *WriteAgentVFSParams) (*AgentVFSFile, error) {
	if params == nil {
		return nil, fmt.Errorf("params is required")
	}
	if params.Mode == "" {
		params.Mode = AgentVFSWriteOverwrite
	}
	var result AgentVFSFile
	err := s.http.put(ctx, fmt.Sprintf("/api/v1/agents/%s/vfs/file", agentID), params, &result)
	return &result, err
}

// Mkdir creates a directory (parents auto-created).
func (s *AgentVFSService) Mkdir(ctx context.Context, agentID, path string) (*AgentVFSFile, error) {
	var result AgentVFSFile
	err := s.http.post(ctx, fmt.Sprintf("/api/v1/agents/%s/vfs/mkdir", agentID), map[string]string{"path": path}, &result)
	return &result, err
}

// Move renames or moves a file or directory.
func (s *AgentVFSService) Move(ctx context.Context, agentID, from, to string) error {
	return s.http.post(ctx, fmt.Sprintf("/api/v1/agents/%s/vfs/move", agentID), map[string]string{"from": from, "to": to}, nil)
}

// Copy duplicates a file or subtree.
func (s *AgentVFSService) Copy(ctx context.Context, agentID, from, to string) error {
	return s.http.post(ctx, fmt.Sprintf("/api/v1/agents/%s/vfs/copy", agentID), map[string]string{"from": from, "to": to}, nil)
}

// Delete removes a file. Set recursive=true to remove a directory subtree.
func (s *AgentVFSService) Delete(ctx context.Context, agentID, path string, recursive bool) error {
	q := url.Values{}
	q.Set("path", path)
	if recursive {
		q.Set("recursive", "true")
	}
	return s.http.del(ctx, fmt.Sprintf("/api/v1/agents/%s/vfs?%s", agentID, q.Encode()))
}

// Grep searches file contents.
func (s *AgentVFSService) Grep(ctx context.Context, agentID string, params *GrepAgentVFSParams) (*AgentVFSGrepResult, error) {
	if params == nil || params.Query == "" {
		return nil, fmt.Errorf("query is required")
	}
	qp := map[string]string{"q": params.Query}
	if params.Path != "" {
		qp["path"] = params.Path
	}
	if params.Limit > 0 {
		qp["limit"] = strconv.Itoa(params.Limit)
	}
	var result AgentVFSGrepResult
	err := s.http.get(ctx, fmt.Sprintf("/api/v1/agents/%s/vfs/grep", agentID), qp, &result)
	return &result, err
}

// Glob finds files by name or path pattern.
func (s *AgentVFSService) Glob(ctx context.Context, agentID string, params *GlobAgentVFSParams) (*AgentVFSGlobResult, error) {
	if params == nil || params.Pattern == "" {
		return nil, fmt.Errorf("pattern is required")
	}
	qp := map[string]string{"pattern": params.Pattern}
	if params.Path != "" {
		qp["path"] = params.Path
	}
	if params.Limit > 0 {
		qp["limit"] = strconv.Itoa(params.Limit)
	}
	var result AgentVFSGlobResult
	err := s.http.get(ctx, fmt.Sprintf("/api/v1/agents/%s/vfs/glob", agentID), qp, &result)
	return &result, err
}

// Usage returns total bytes used by the agent's VFS.
func (s *AgentVFSService) Usage(ctx context.Context, agentID string) (int64, error) {
	var result struct {
		BytesUsed int64 `json:"bytes_used"`
	}
	err := s.http.get(ctx, fmt.Sprintf("/api/v1/agents/%s/vfs/usage", agentID), nil, &result)
	return result.BytesUsed, err
}
