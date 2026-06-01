# PromptRails Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/promptrails/go-sdk.svg)](https://pkg.go.dev/github.com/promptrails/go-sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Official Go SDK for [PromptRails](https://promptrails.ai) — the AI agent orchestration platform.

The SDK has two independent parts:

- **API client** (`github.com/promptrails/go-sdk`) — manage agents, prompts, executions, and more.
- **Tracing** (`github.com/promptrails/go-sdk/tracing`) — send spans to PromptRails
  from any code, without managing your prompts/agents on the platform. Standard
  library only.

## Installation

```bash
go get github.com/promptrails/go-sdk
```

## Quick Start

### API client

```go
package main

import (
    "context"
    "fmt"
    "log"

    promptrails "github.com/promptrails/go-sdk"
)

func main() {
    client := promptrails.NewClient("pr_key_...")
    result, err := client.Agents.Execute(context.Background(), "agent-id", &promptrails.ExecuteAgentParams{
        Input: map[string]any{"query": "Summarise this week's sales"},
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result.Output)
}
```

See the [API client guide](docs/api-client.md) for resources, error handling,
media studio, and configuration.

### Tracing

```go
import "github.com/promptrails/go-sdk/tracing"

tracer := tracing.NewTracer("pr_...")
defer tracer.Shutdown()

_ = tracer.Span(ctx, "agent-run", tracing.KindAgent, func(ctx context.Context, root *tracing.Span) error {
    root.SetInput(map[string]any{"q": "weather?"})
    return tracer.Span(ctx, "llm-call", tracing.KindLLM, func(ctx context.Context, llm *tracing.Span) error {
        llm.SetModel("gpt-4o").SetUsage(120, 30, -1)
        return nil
    })
})
```

See the [tracing guide](docs/tracing.md) for manual spans, span kinds,
configuration, and the OpenTelemetry bridge.

## Documentation

- [API client](docs/api-client.md) — resources, error handling, media studio, configuration
- [Tracing](docs/tracing.md) — spans, batching, configuration, OpenTelemetry

## Contributing

```bash
go test ./... -v -race
go vet ./...
gofmt -w .
```

## License

MIT
