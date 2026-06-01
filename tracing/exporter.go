package tracing

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const ingestPath = "/api/v1/traces/ingest"

// exporter buffers finished spans and ships them to the ingest API in batches.
// Delivery is best-effort: failures are logged and dropped rather than returned
// to the caller, and the buffer is bounded.
type exporter struct {
	apiKey       string
	baseURL      string
	client       *http.Client
	maxBatchSize int
	maxQueueSize int

	mu     sync.Mutex
	buffer []spanPayload

	flushCh chan struct{}
	stopCh  chan struct{}
	doneCh  chan struct{}
}

func newExporter(apiKey, baseURL string, batchSize int, interval time.Duration, timeout time.Duration) *exporter {
	e := &exporter{
		apiKey:       apiKey,
		baseURL:      strings.TrimRight(baseURL, "/"),
		client:       &http.Client{Timeout: timeout},
		maxBatchSize: batchSize,
		maxQueueSize: 10_000,
		flushCh:      make(chan struct{}, 1),
		stopCh:       make(chan struct{}),
		doneCh:       make(chan struct{}),
	}
	go e.run(interval)
	return e
}

func (e *exporter) submit(p spanPayload) {
	e.mu.Lock()
	if len(e.buffer) >= e.maxQueueSize {
		e.mu.Unlock()
		log.Println("promptrails tracing buffer full; dropping span")
		return
	}
	e.buffer = append(e.buffer, p)
	full := len(e.buffer) >= e.maxBatchSize
	e.mu.Unlock()
	if full {
		select {
		case e.flushCh <- struct{}{}:
		default:
		}
	}
}

func (e *exporter) run(interval time.Duration) {
	defer close(e.doneCh)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-e.stopCh:
			e.drain()
			return
		case <-ticker.C:
			e.drain()
		case <-e.flushCh:
			e.drain()
		}
	}
}

func (e *exporter) drain() {
	for {
		e.mu.Lock()
		if len(e.buffer) == 0 {
			e.mu.Unlock()
			return
		}
		n := len(e.buffer)
		if n > e.maxBatchSize {
			n = e.maxBatchSize
		}
		batch := make([]spanPayload, n)
		copy(batch, e.buffer[:n])
		e.buffer = e.buffer[n:]
		e.mu.Unlock()
		e.send(batch)
	}
}

func (e *exporter) send(spans []spanPayload) {
	body, err := json.Marshal(map[string]any{"spans": spans})
	if err != nil {
		log.Printf("promptrails trace export: marshal failed: %v", err)
		return
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, e.baseURL+ingestPath, bytes.NewReader(body))
	if err != nil {
		log.Printf("promptrails trace export: request failed: %v", err)
		return
	}
	req.Header.Set("X-API-Key", e.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := e.client.Do(req)
	if err != nil {
		log.Printf("promptrails trace export failed (%d spans): %v", len(spans), err)
		return
	}
	_ = resp.Body.Close()
}

// flush sends buffered spans synchronously.
func (e *exporter) flush() {
	e.drain()
}

// shutdown stops the worker after flushing remaining spans.
func (e *exporter) shutdown() {
	select {
	case <-e.stopCh:
		// already stopped
	default:
		close(e.stopCh)
	}
	<-e.doneCh
}
