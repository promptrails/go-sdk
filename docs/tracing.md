# Tracing

Send spans to PromptRails from any Go code, without managing your prompts or
agents on the platform — PromptRails as a standalone LLM observability backend.

The tracing package (`github.com/promptrails/go-sdk/tracing`) depends only on the
standard library and is independent of the API client. It needs an API key with
the `traces:write` scope.

## Quick start

```go
import (
    "context"

    "github.com/promptrails/go-sdk/tracing"
)

tracer := tracing.NewTracer("pr_...")
defer tracer.Shutdown()

ctx := context.Background()
_ = tracer.Span(ctx, "agent-run", tracing.KindAgent, func(ctx context.Context, root *tracing.Span) error {
    root.SetInput(map[string]any{"q": "What's the weather?"})

    return tracer.Span(ctx, "llm-call", tracing.KindLLM, func(ctx context.Context, llm *tracing.Span) error {
        llm.SetModel("gpt-4o").SetUsage(120, 30, -1)
        llm.SetOutput(map[string]any{"text": "18°C and rainy."})
        return nil
    })
})
```

A span created from a child `context.Context` automatically shares the parent's
trace ID and links to it, so the tree renders correctly in the trace viewer. If
the callback returns an error, the span is marked `error`.

## Manual spans

When the callback form doesn't fit, use `StartSpan` and end the span yourself:

```go
ctx, span := tracer.StartSpan(ctx, "step", tracing.KindAgentStep)
defer span.End()
// ... use ctx for child spans ...
```

## Span data

The `Set*` methods chain:

```go
span.
    SetInput(map[string]any{"query": "weather istanbul"}).
    SetOutput(map[string]any{"rows": 3}).
    SetAttributes(map[string]any{"index": "weather"}).
    SetTags("prod", "search").
    SetSession("session-abc")
```

For LLM spans:

```go
span.SetModel("gpt-4o")
span.SetUsage(120, 30, -1) // prompt, completion, total (-1 computes total)
span.SetCost(0.0023)
span.SetError(err) // marks the span failed
```

Pass `-1` to `SetUsage` for any unknown count; total is computed from
prompt+completion when total is `-1` and both are known.

## Span kinds

`tracing.Kind` is a string; the package exports constants for the common kinds
(`KindAgent`, `KindLLM`, `KindTool`, `KindDatasource`, `KindChain`,
`KindWorkflow`, `KindAgentStep`, `KindEmbedding`, `KindGuardrail`). Any
PromptRails span kind string is accepted.

## Lifecycle & flushing

Spans are buffered and shipped in batches by a background goroutine — on an
interval, when the buffer fills, and on `Shutdown`. Export is best-effort:
failures are logged (via the standard `log` package) and dropped, never returned
to your code.

- `tracer.Flush()` sends buffered spans synchronously.
- `tracer.Shutdown()` flushes and stops the worker — call it with `defer` in
  `main`.

## Configuration

```go
tracer := tracing.NewTracer("pr_...",
    tracing.WithBaseURL("https://api.promptrails.ai"),
    tracing.WithBatchSize(100),
    tracing.WithFlushInterval(time.Second),
    tracing.WithTimeout(30*time.Second),
)
```

Each span is POSTed to `POST /api/v1/traces/ingest`; the workspace comes from the
API key and the `source` attribute is set to `sdk`.

## OpenTelemetry

If you already instrument with OpenTelemetry, you don't need this package to
export. Point a standard OTLP/HTTP exporter at the PromptRails OTLP endpoint and
your spans flow in — `gen_ai.*` semantic-convention attributes are mapped onto
the trace model:

```go
import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"

exporter, _ := otlptracehttp.New(ctx,
    otlptracehttp.WithEndpointURL("https://api.promptrails.ai/api/v1/otel/v1/traces"),
    otlptracehttp.WithHeaders(map[string]string{"X-API-Key": "pr_..."}),
)
```
