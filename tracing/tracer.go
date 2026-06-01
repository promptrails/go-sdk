package tracing

import (
	"context"
	"time"
)

type spanCtxKey struct{}

// SpanFromContext returns the active span carried by ctx, or nil.
func SpanFromContext(ctx context.Context) *Span {
	span, _ := ctx.Value(spanCtxKey{}).(*Span)
	return span
}

// Tracer creates spans and ships them to PromptRails. It is safe for concurrent
// use. Create one per process and reuse it.
type Tracer struct {
	exporter *exporter
}

// Option configures a Tracer.
type Option func(*tracerConfig)

type tracerConfig struct {
	baseURL       string
	batchSize     int
	flushInterval time.Duration
	timeout       time.Duration
}

// WithBaseURL points the tracer at a specific PromptRails deployment.
func WithBaseURL(url string) Option { return func(c *tracerConfig) { c.baseURL = url } }

// WithBatchSize sets how many spans are sent per request.
func WithBatchSize(n int) Option { return func(c *tracerConfig) { c.batchSize = n } }

// WithFlushInterval sets how often buffered spans are flushed.
func WithFlushInterval(d time.Duration) Option { return func(c *tracerConfig) { c.flushInterval = d } }

// WithTimeout sets the HTTP timeout for export requests.
func WithTimeout(d time.Duration) Option { return func(c *tracerConfig) { c.timeout = d } }

// NewTracer creates a tracer authenticated with the given API key (which must
// have the traces:write scope). Remember to call Shutdown before exit to flush
// any buffered spans.
func NewTracer(apiKey string, opts ...Option) *Tracer {
	cfg := tracerConfig{
		baseURL:       "https://api.promptrails.ai",
		batchSize:     100,
		flushInterval: time.Second,
		timeout:       30 * time.Second,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return &Tracer{
		exporter: newExporter(apiKey, cfg.baseURL, cfg.batchSize, cfg.flushInterval, cfg.timeout),
	}
}

// StartSpan creates a span and returns a context carrying it. The span inherits
// its trace ID and parent from any span already in ctx. The caller must call
// span.End (typically via defer).
//
//	ctx, span := tracer.StartSpan(ctx, "llm-call", tracing.KindLLM)
//	defer span.End()
func (t *Tracer) StartSpan(ctx context.Context, name string, kind Kind) (context.Context, *Span) {
	parent := SpanFromContext(ctx)
	now := time.Now().UTC()

	span := &Span{
		onEnd: t.exporter.submit,
		payload: spanPayload{
			TraceID:   generateTraceID(),
			SpanID:    generateSpanID(),
			Name:      name,
			Kind:      kind,
			Status:    "ok",
			Level:     "default",
			StartedAt: now,
		},
	}
	if parent != nil {
		span.payload.TraceID = parent.payload.TraceID
		span.payload.ParentSpanID = parent.payload.SpanID
		span.payload.SessionID = parent.payload.SessionID
	}
	return context.WithValue(ctx, spanCtxKey{}, span), span
}

// Span runs fn inside a span: the span is the active parent for fn, ended when
// fn returns, and marked as an error if fn returns one.
func (t *Tracer) Span(ctx context.Context, name string, kind Kind, fn func(context.Context, *Span) error) error {
	ctx, span := t.StartSpan(ctx, name, kind)
	defer span.End()
	if err := fn(ctx, span); err != nil {
		span.SetError(err)
		return err
	}
	return nil
}

// Flush sends any buffered spans synchronously.
func (t *Tracer) Flush() { t.exporter.flush() }

// Shutdown flushes remaining spans and stops the background worker.
func (t *Tracer) Shutdown() { t.exporter.shutdown() }
