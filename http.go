package promptrails

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"
)

type httpClient struct {
	baseURL    string
	apiKey     string
	maxRetries int
	client     *http.Client
}

func newHTTPClient(cfg config) *httpClient {
	return &httpClient{
		baseURL:    strings.TrimRight(cfg.baseURL, "/"),
		apiKey:     cfg.apiKey,
		maxRetries: cfg.maxRetries,
		client:     &http.Client{Timeout: cfg.timeout},
	}
}

// getRaw does GET and returns raw bytes for caller to handle.
func (h *httpClient) getRaw(ctx context.Context, path string, params map[string]string) ([]byte, error) {
	var result json.RawMessage
	if err := h.request(ctx, http.MethodGet, path, params, nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (h *httpClient) get(ctx context.Context, path string, params map[string]string, result any) error {
	return h.request(ctx, http.MethodGet, path, params, nil, result)
}

func (h *httpClient) post(ctx context.Context, path string, body any, result any) error {
	return h.request(ctx, http.MethodPost, path, nil, body, result)
}

func (h *httpClient) patch(ctx context.Context, path string, body any, result any) error {
	return h.request(ctx, http.MethodPatch, path, nil, body, result)
}

func (h *httpClient) put(ctx context.Context, path string, body any, result any) error {
	return h.request(ctx, http.MethodPut, path, nil, body, result)
}

func (h *httpClient) del(ctx context.Context, path string) error {
	return h.request(ctx, http.MethodDelete, path, nil, nil, nil)
}

func (h *httpClient) request(ctx context.Context, method, path string, params map[string]string, body any, result any) error {
	var lastErr error
	attempts := h.maxRetries + 1

	for i := 0; i < attempts; i++ {
		if i > 0 {
			wait := time.Duration(math.Pow(2, float64(i-1))) * 500 * time.Millisecond
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(wait):
			}
		}

		err := h.doRequest(ctx, method, path, params, body, result)
		if err == nil {
			return nil
		}
		lastErr = err

		if _, ok := err.(*ServerError); ok {
			continue
		}
		return err
	}
	return lastErr
}

func (h *httpClient) doRequest(ctx context.Context, method, path string, params map[string]string, body any, result any) error {
	url := h.baseURL + path

	if len(params) > 0 {
		parts := make([]string, 0, len(params))
		for k, v := range params {
			if v != "" {
				parts = append(parts, k+"="+v)
			}
		}
		if len(parts) > 0 {
			url += "?" + strings.Join(parts, "&")
		}
	}

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-Key", h.apiKey)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return h.parseError(resp.StatusCode, respBody)
	}

	if result != nil && len(respBody) > 0 {
		return unmarshalResponse(respBody, result)
	}

	return nil
}

// unmarshalResponse handles the API response envelope.
// The API returns responses in two formats:
// 1. Single object: {"data": {...}} — unwrap "data"
// 2. Paginated: {"data": [...], "meta": {...}} — direct unmarshal
// 3. Array: {"data": [...]} — unwrap "data"
func unmarshalResponse(body []byte, result any) error {
	// Check if the response has a "meta" field (paginated response)
	var probe struct {
		Meta *json.RawMessage `json:"meta"`
		Data json.RawMessage  `json:"data"`
	}
	if err := json.Unmarshal(body, &probe); err == nil {
		if probe.Meta != nil {
			// Paginated response — unmarshal the full body directly
			return json.Unmarshal(body, result)
		}
		if probe.Data != nil {
			// Single/array response wrapped in envelope — unwrap "data"
			return json.Unmarshal(probe.Data, result)
		}
	}

	// Fallback: direct unmarshal
	return json.Unmarshal(body, result)
}

func (h *httpClient) parseError(statusCode int, body []byte) error {
	var raw struct {
		Error   any    `json:"error"`
		Message string `json:"message"`
	}
	msg := fmt.Sprintf("request failed (status %d)", statusCode)
	code := ""
	var details map[string]any

	if json.Unmarshal(body, &raw) == nil {
		switch e := raw.Error.(type) {
		case string:
			msg = e
		case map[string]any:
			if m, ok := e["message"].(string); ok {
				msg = m
			}
			if c, ok := e["code"].(string); ok {
				code = c
			}
			if d, ok := e["details"].(map[string]any); ok {
				details = d
			}
		}
		if msg == fmt.Sprintf("request failed (status %d)", statusCode) && raw.Message != "" {
			msg = raw.Message
		}
	}

	return newErrorForStatus(statusCode, msg, code, details)
}
