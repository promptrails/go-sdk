# API Client

The `promptrails.Client` wraps the PromptRails REST API for managing agents,
prompts, executions, media, and more.

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

    result, err := client.Agents.Execute(ctx, "agent-id", &promptrails.ExecuteAgentParams{
        Input: map[string]any{"query": "Summarise this week's sales"},
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result.Output)
}
```

## Error handling

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

## Available resources

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
| `client.AgentTriggers`   | `List`, `Get`, `Create` (with `Source` + `SourceConfig`), `Update`, `Delete` |
| `client.AgentVFS`        | `List`, `Read`, `Write`, `Stat`, `Mkdir`, `Move`, `Copy`, `Delete`, `Grep`, `Glob`, `Usage` |
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

// Browse and manage assets
assets, err := client.Assets.List(ctx, &promptrails.ListAssetsParams{MediaType: "image"})
signed, err := client.Assets.GetSignedURL(ctx, "asset-id")
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

| Option           | Default                      | Description                       |
| ---------------- | ---------------------------- | --------------------------------- |
| `WithBaseURL`    | `https://api.promptrails.ai` | API base URL                      |
| `WithTimeout`    | `30s`                        | Request timeout                   |
| `WithMaxRetries` | `3`                          | Max retries on network/5xx errors |
