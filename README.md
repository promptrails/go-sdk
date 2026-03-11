# PromptRails Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/promptrails/go-sdk.svg)](https://pkg.go.dev/github.com/promptrails/go-sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Official Go SDK for [PromptRails](https://promptrails.ai) — the AI agent orchestration platform.

## Installation

```bash
go get github.com/promptrails/go-sdk
```

## Quick Start

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

	ctx := context.Background()

	// Execute an agent
	result, err := client.Agents.Execute(ctx, "agent-id", &promptrails.ExecuteAgentParams{
		Input: map[string]any{"query": "Summarise this week's sales"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Output)
}
```

## Error Handling

```go
import "errors"

result, err := client.Agents.Get(ctx, "missing-id")
if err != nil {
	var notFound *promptrails.NotFoundError
	var rateLimit *promptrails.RateLimitError
	var quota *promptrails.QuotaExceededError

	switch {
	case errors.As(err, &notFound):
		fmt.Println("Agent not found")
	case errors.As(err, &rateLimit):
		fmt.Println("Rate limited, back off and retry")
	case errors.As(err, &quota):
		fmt.Println("Execution limit reached")
	default:
		fmt.Printf("Error: %v\n", err)
	}
}
```

## Available Resources

| Resource                 | Methods                                                                  |
| ------------------------ | ------------------------------------------------------------------------ |
| `client.Agents`          | `List`, `Get`, `Create`, `Update`, `Delete`, `Execute`, `ListVersions`, `CreateVersion`, `PromoteVersion`, `ListGuardrails`, `CreateGuardrail`, `ListMemories`, `CreateMemory`, `SearchMemories`, `DeleteAllMemories` |
| `client.Prompts`         | `List`, `Get`, `Create`, `Update`, `Delete`, `ListVersions`, `CreateVersion`, `PromoteVersion`, `Run` |
| `client.Executions`      | `List`, `Get`                                                            |
| `client.Credentials`     | `List`, `Get`, `Create`, `Update`, `Delete`, `SetDefault`, `CheckConnection` |
| `client.DataSources`     | `List`, `Get`, `Create`, `Update`, `Delete`, `ListVersions`, `CreateVersion`, `TestConnection`, `Query` |
| `client.Chat`            | `ListSessions`, `GetSession`, `CreateSession`, `DeleteSession`, `ListMessages`, `SendMessage` |
| `client.Traces`          | `List`, `GetByTraceID`                                                   |
| `client.Costs`           | `GetSummary`, `GetAgentSummary`                                          |
| `client.Scores`          | `List`, `Get`, `Create`, `Update`, `Delete`, `Aggregates`, `ListConfigs`, `GetConfig`, `CreateConfig`, `UpdateConfig`, `DeleteConfig` |
| `client.MCPTools`        | `List`, `Get`, `Create`, `Update`, `Delete`                              |
| `client.Approvals`       | `List`, `Get`, `Decide`                                                  |
| `client.WebhookTriggers` | `List`, `Get`, `Create`, `Update`, `Delete`                              |
| `client.A2A`             | `GetAgentCard`, `SendMessage`, `GetTask`, `ListTasks`, `CancelTask`      |
| `client.Media`           | `Generate`                                                               |
| `client.MediaModels`     | `List`                                                                   |
| `client.Assets`          | `List`, `Get`, `Delete`, `GetSignedURL`                                  |

## Media Studio

```go
// List available media models
models, err := client.MediaModels.List(ctx, &promptrails.ListMediaModelsParams{
	Provider:  "fal",
	MediaType: "image",
})

// Generate an image
resp, err := client.Media.Generate(ctx, &promptrails.GenerateMediaParams{
	Provider:  "fal",
	MediaType: "image",
	Model:     "fal-ai/flux/schnell",
	Prompt:    "A sunset over mountains",
	Config:    map[string]any{"width": 1024, "height": 768},
})
fmt.Println(resp.AssetURL)

// List assets
assets, err := client.Assets.List(ctx, &promptrails.ListAssetsParams{
	MediaType: "image",
})

// Get a signed URL for an asset
signed, err := client.Assets.GetSignedURL(ctx, "asset-id")
fmt.Println(signed.URL)

// Delete an asset
err = client.Assets.Delete(ctx, "asset-id")
```

## Configuration

```go
client := promptrails.NewClient("pr_key_...",
	promptrails.WithBaseURL("http://localhost:8082"),
	promptrails.WithTimeout(10 * time.Second),
	promptrails.WithMaxRetries(5),
)
```

| Option          | Default                      | Description                       |
| --------------- | ---------------------------- | --------------------------------- |
| `WithBaseURL`   | `https://api.promptrails.ai` | API base URL                      |
| `WithTimeout`   | `30s`                        | Request timeout                   |
| `WithMaxRetries`| `3`                          | Max retries on network/5xx errors |

## Contributing

```bash
go test ./... -v -race
go vet ./...
gofmt -w .
```

## License

MIT
