# maxigo-client

  Idiomatic Go HTTP client for [Max Bot API](https://dev.max.ru) (OpenAPI schema v0.0.10). Zero external dependencies.

## Why This Project?

  The official [`max-bot-api-client-go`](https://github.com/max-messenger/max-bot-api-client-go) has systemic issues that make it unsuitable for production use:

  | Category                                          | Issues                                                   | Severity                                    |
  |---------------------------------------------------|----------------------------------------------------------|---------------------------------------------|
  | Errors swallowed via `log.Println` / `slog.Error` | 30+ places                                               | Critical — library must not write to stdout |
  | Untestable without real API                       | No `WithBaseURL`, uploads use `http.DefaultClient`       | Critical — impossible to unit test          |
  | Unnecessary dependencies                          | zerolog, yaml, env — 6 transitive deps                   | Critical — for an HTTP client               |
  | Broken API methods                                | `GetChatID()` returns 0 for callbacks                    | Critical — callbacks unusable               |
  | Wrong types                                       | `time.Duration` for Unix timestamps, `int64→int` casts   | Critical — data corruption on 32-bit        |
  | No `context.Context` in uploads                   | `http.Get()` without timeout or cancellation             | Critical — can hang forever                 |
  | Non-idiomatic API                                 | Builder pattern, `SCREAMING_CASE`, no functional options | Major                                       |
  | `schemes.Error` always non-nil                    | `Check()` returns error even on success                  | Critical — always errors                    |

  **maxigo-client** fixes all of these.

## Documentation

- **[English Guide](docs/guide.md)** — full API reference with examples
- **[Документация на русском](docs/guide-ru.md)** — полное описание API с примерами
  Requires Go 1.25+.

## Installation

```bash
go get github.com/maxigo-bot/maxigo-client
```


## Quick Start

```go
import (
"context"
"fmt"
"log"

maxigo "github.com/maxigo-bot/maxigo-client"
)

client, err := maxigo.New("YOUR_BOT_TOKEN")
if err != nil {
log.Fatal(err)
}

bot, err := client.GetBot(context.Background())
if err != nil {
log.Fatal(err)
}
fmt.Printf("Bot: %s (ID: %d)\n", bot.FirstName, bot.UserID)
```


## Features

- Zero external dependencies — only stdlib
- All methods take `context.Context`
- Structured error handling with `errors.As`
- Full Max Bot API coverage: messages, chats, uploads, webhooks, long polling
- Testable without real API via `WithBaseURL`

## API Overview

```go
// Bot
client.GetBot(ctx)
client.EditBot(ctx, patch)

// Messages
client.SendMessage(ctx, chatID, body)
client.EditMessage(ctx, messageID, body)
client.DeleteMessage(ctx, messageID)
client.AnswerCallback(ctx, callbackID, answer)

// Chats
client.GetChat(ctx, chatID)
client.GetChats(ctx, opts)
client.EditChat(ctx, chatID, patch)
client.GetMembers(ctx, chatID, opts)
client.SendAction(ctx, chatID, action)

// Uploads
client.UploadPhoto(ctx, filename, reader)
client.UploadMedia(ctx, uploadType, filename, reader)

// Webhooks
client.Subscribe(ctx, url, types, secret)
client.Unsubscribe(ctx, url)

// Long Polling
client.GetUpdates(ctx, opts)
```

## Error Handling

```go
var e *maxigo.Error
if errors.As(err, &e) {
    fmt.Printf("[%s] %s %d: %s\n", e.Op, e.Kind, e.StatusCode, e.Message)
}
```

Error kinds: `ErrAPI`, `ErrNetwork`, `ErrTimeout`, `ErrDecode`. See [guide](docs/guide.md#error-handling) for details.

## Ecosystem

| Package                           | Description                                  |
|-----------------------------------|----------------------------------------------|
| `github.com/maxigo-bot/maxigo-client` | HTTP client                                  |
| `github.com/maxigo/maxigo`        | Bot framework with router/middleware/context |

> *Both packages are currently in development and not yet published.*

## License

MIT
