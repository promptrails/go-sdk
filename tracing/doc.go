// Package tracing sends spans to PromptRails from any Go code, without managing
// your prompts or agents on the platform — PromptRails as a standalone LLM
// observability backend.
//
// It depends only on the standard library and needs an API key with the
// traces:write scope.
//
//	tracer := tracing.NewTracer(os.Getenv("PROMPTRAILS_API_KEY"))
//	defer tracer.Shutdown()
//
//	err := tracer.Span(ctx, "agent-run", tracing.KindAgent, func(ctx context.Context, span *tracing.Span) error {
//	    span.SetInput(map[string]any{"q": "weather?"})
//	    return tracer.Span(ctx, "llm-call", tracing.KindLLM, func(ctx context.Context, llm *tracing.Span) error {
//	        llm.SetModel("gpt-4o").SetUsage(120, 30, -1)
//	        return nil
//	    })
//	})
//
// Spans created from a child context automatically share the parent's trace ID
// and link to it, so the tree renders correctly in the trace viewer.
//
// # OpenTelemetry
//
// If you already instrument with OpenTelemetry, you don't need this package for
// export: point a standard OTLP/HTTP exporter at the PromptRails OTLP endpoint
// (<base>/api/v1/otel) with the X-API-Key header. The backend maps gen_ai.*
// semantic conventions onto the trace model.
package tracing
