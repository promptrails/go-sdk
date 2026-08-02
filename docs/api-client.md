# API Client

The `promptrails.Client` wraps the PromptRails REST API for managing agents,
prompts, executions, traces, and more.

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
| `client.Agents`          | `List`, `Get`, `Create`, `Update`, `Delete`, `Execute`, `Playground`, `ListVersions`, `CreateVersion`, `PromoteVersion`, `ListGuardrails`, `CreateGuardrail` |
| `client.Prompts`         | `List`, `Get`, `Create`, `Update`, `Delete`, `ListVersions`, `CreateVersion`, `PromoteVersion` |
| `client.Executions`      | `List`, `Get`, `Tree`, `Cancel`, `ApprovalInbox`, `Approve`, `Deny`, `Stream` |
| `client.Credentials`     | `List`, `Get`, `Create`, `Update`, `Delete`, `SetDefault`, `CheckConnection` |
| `client.DataSources`     | `List`, `Get`, `Create`, `Update`, `Delete`, `ListVersions`, `CreateVersion`, `TestConnection`, `Query` |
| `client.Chat`            | `ListSessions`, `GetSession`, `CreateSession`, `DeleteSession`, `ListMessages`, `SendMessage` |
| `client.Traces`          | `List`, `GetByTraceID`, `GetSummary`, `PIIReport`, `Ingest`             |
| `client.MCPTools`        | `List`, `Get`, `Create`, `Update`, `Delete`                              |
| `client.Guardrails`      | `ListScanners`, `Update`, `Delete`                                       |
| `client.AgentTriggers`   | `List`, `Get`, `Create` (with `Source` + `SourceConfig`), `Update`, `Delete` |
| `client.AgentVFS`        | `List`, `Read`, `Write`, `Stat`, `Mkdir`, `Move`, `Copy`, `Delete`, `Grep`, `Glob`, `Usage` |
| `client.A2A`             | `GetAgentCard`, `SendMessage`, `GetTask`, `ListTasks`, `CancelTask`      |
| `client.LLMModels`       | `List`, `ListAvailable`                                                  |
| `client.Assets`          | `List`, `Get`, `Delete`, `GetSignedURL`                                  |

## Agent versions

In API v2 an agent is one of two kinds — `agent` or `workflow` — and the
**version** owns the model + runtime configuration. Model, sampling, run
budget, approval policy, cache TTL, VFS/masking overrides and the attached
tools / sub-agents / guardrails are all passed to `CreateVersion` as siblings
of `Config` (a prompt carries no model configuration):

```go
temp := 0.2
maxCost := 1.0
_, err := client.Agents.CreateVersion(ctx, "agent-id", &promptrails.CreateVersionParams{
    Version:    "1.0.0",
    SetCurrent: true,
    Config:     promptrails.PromptAgentConfig{PromptID: "prompt-id"},
    ModelConfig: &promptrails.ModelConfig{
        ModelID:     "llm-model-id",
        Temperature: &temp,
    },
    RunBudget:  &promptrails.RunBudget{MaxCost: &maxCost},
    Tools:      []promptrails.ToolAttachment{{MCPToolID: "tool-id", RequiresApproval: true}},
    SubAgents:  []promptrails.SubAgentAttachment{{AgentID: "sub-id", Alias: "researcher"}},
    Guardrails: []promptrails.GuardrailSpec{{Type: "input", ScannerType: "pii"}},
})
```

## Approvals & execution trees

Approval-gated tool calls park an execution at `waiting_approval`. Reach them
through the execution-scoped inbox and resume them in place:

```go
inbox, err := client.Executions.ApprovalInbox(ctx, nil)
_, err = client.Executions.Approve(ctx, "execution-id", &promptrails.DecideParams{Reason: "looks good"})
_, err = client.Executions.Deny(ctx, "execution-id", nil)

// Inspect the full sub-agent / workflow tree, or cancel a running root.
tree, err := client.Executions.Tree(ctx, "execution-id")
_, err = client.Executions.Cancel(ctx, "execution-id")
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
