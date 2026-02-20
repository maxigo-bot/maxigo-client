# maxigo-client

[![Go Reference](https://pkg.go.dev/badge/github.com/maxigo-bot/maxigo-client.svg)](https://pkg.go.dev/github.com/maxigo-bot/maxigo-client)
[![Go Report Card](https://goreportcard.com/badge/github.com/maxigo-bot/maxigo-client)](https://goreportcard.com/report/github.com/maxigo-bot/maxigo-client)
[![CI](https://github.com/maxigo-bot/maxigo-client/actions/workflows/ci.yml/badge.svg)](https://github.com/maxigo-bot/maxigo-client/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/maxigo-bot/maxigo-client/branch/main/graph/badge.svg)](https://codecov.io/gh/maxigo-bot/maxigo-client)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/maxigo-bot/maxigo-client)](https://github.com/maxigo-bot/maxigo-client)

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
package main

import (
	"context"
	"fmt"
	"log"

	maxigo "github.com/maxigo-bot/maxigo-client"
)

func main() {
	client, err := maxigo.New("YOUR_BOT_TOKEN")
	if err != nil {
		log.Fatal(err)
	}

	bot, err := client.GetBot(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Bot: %s (ID: %d)\n", bot.FirstName, bot.UserID)
}
```


## Features

- Zero external dependencies — only stdlib
- All methods take `context.Context`
- Structured error handling with `errors.As`
- Full Max Bot API coverage: messages, chats, uploads, webhooks, long polling
- Type-safe constructors for all button types and attachments
- `Optional[T]` generics for three-state fields (unset / zero / value)
- Testable without real API via `WithBaseURL`

## Type-Safe Constructors

No need to remember string constants — use constructors for buttons and attachments:

**Buttons:**

```go
maxigo.NewCallbackButton("Click", "payload")                        // callback
maxigo.NewCallbackButtonWithIntent("Yes", "yes", maxigo.IntentPositive) // callback with intent
maxigo.NewLinkButton("Open", "https://example.com")                 // link
maxigo.NewRequestContactButton("Share contact")                     // request contact
maxigo.NewRequestGeoLocationButton("Send location", true)           // request geo (quick=true)
maxigo.NewChatButton("Create chat", "Title")                        // create chat
maxigo.NewMessageButton("Send")                                     // message from user
maxigo.NewOpenAppButton("Open WebApp", "bot_username")               // open mini-app
```

**Attachments:**

```go
maxigo.NewInlineKeyboardAttachment(buttons) // inline keyboard
maxigo.NewPhotoAttachment(payload)          // image
maxigo.NewVideoAttachment(payload)          // video
maxigo.NewAudioAttachment(payload)          // audio
maxigo.NewFileAttachment(payload)           // file
maxigo.NewStickerAttachment(payload)        // sticker
maxigo.NewContactAttachment(payload)        // contact card
maxigo.NewShareAttachment(payload)          // share link
maxigo.NewLocationAttachment(lat, lng)      // location
```

**Example — inline keyboard with contact and geo buttons:**

```go
msg, err := client.SendMessage(ctx, chatID, &maxigo.NewMessageBody{
    Text: maxigo.Some("Choose an option:"),
    Attachments: []maxigo.AttachmentRequest{
        maxigo.NewInlineKeyboardAttachment([][]maxigo.Button{
            {
                maxigo.NewRequestContactButton("Share contact"),
                maxigo.NewRequestGeoLocationButton("Send location", false),
            },
            {
                maxigo.NewCallbackButtonWithIntent("Cancel", "cancel", maxigo.IntentNegative),
            },
        }),
    },
})
```

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

| Package | Description |
|---------|-------------|
| [maxigo-client](https://github.com/maxigo-bot/maxigo-client) | Idiomatic Go HTTP client for Max Bot API (zero external deps) |
| [maxigo-bot](https://github.com/maxigo-bot/maxigo-bot) | Bot framework with router, middleware, and context |

## License

MIT
