# maxigo-client Guide

Go HTTP client for [Max Bot API](https://dev.max.ru). Zero external dependencies.

> **[Документация на русском](guide-ru.md)** | **[README](../README.md)**

## Why This Project?

The official [`max-bot-api-client-go`](https://github.com/max-messenger/max-bot-api-client-go) has systemic issues that make it unsuitable for production:

**Errors are logged instead of returned** — 30+ places use `log.Println` or `slog.Error` in library code. Users cannot control or suppress this output. Some errors are silently swallowed (`json.Decode` failure returns `nil`).

**Cannot test without real API** — no simple `WithBaseURL()` option. Upload methods use `http.DefaultClient` directly, bypassing all client configuration. Testing requires implementing a 7-method `ConfigInterface`.

**6 external dependencies** — zerolog, yaml parser, env parser, gomock — all unnecessary for an HTTP client.

**Broken methods** — `GetChatID()` returns 0 for callbacks (chat ID is available but ignored). `GetCommand()` returns the full message text. `schemes.Error` used as response struct is always non-nil, so `Check()` always returns an error.

**Wrong types** — `time.Duration` for Unix timestamps (interprets as nanoseconds). `int64→int` casts in 10+ places (truncates on 32-bit). `[]interface{}` for attachments (no type safety).

**Non-idiomatic Go** — builder pattern, `SCREAMING_CASE` constants, `Api` instead of `API`, no `context.Context` in uploads, no functional options.

| Problem                      | Official client                                    | maxigo-client                                                 |
|------------------------------|----------------------------------------------------|---------------------------------------------------------------|
| Error handling               | `log.Println` in 30+ places                       | All errors returned as `*Error` with Kind/StatusCode/Op       |
| Testability                  | Need full `ConfigInterface` mock                   | `maxigo.New("token", WithBaseURL(srv.URL))`                   |
| Dependencies                 | 6 transitive (zerolog, yaml, env...)               | 0 — only stdlib                                               |
| `GetChatID()` for callbacks  | Returns 0                                          | Extract from `Message.Recipient.ChatId`                       |
| Types                        | `time.Duration` for timestamps, `int→int64` casts  | Correct `int64` everywhere                                    |
| Uploads                      | `http.Get()` without context/timeout               | All requests through configured client with `context.Context` |
| API style                    | `NewMessage().SetChat().SetText()`                 | `SendMessage(ctx, chatID, &NewMessageBody{Text: &text})`     |
| Constants                    | `TYPING_ON`, `CALLBACK`, `POSITIVE`                | `ActionTypingOn`, `IntentPositive`                            |
| Configuration                | YAML files + env parser                            | Functional options: `WithTimeout`, `WithHTTPClient`           |

**maxigo-client** fixes all of these.

## Installation

```bash
go get github.com/maxigo-bot/maxigo-client
```

Requires Go 1.25+.

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

    ctx := context.Background()

    // Get bot info
    bot, err := client.GetBot(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Bot: %s (ID: %d)\n", bot.FirstName, bot.UserID)

    // Send a message
    text := "Hello from maxigo!"
    msg, err := client.SendMessage(ctx, 123456, &maxigo.NewMessageBody{
        Text: &text,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Sent message: %s\n", msg.Body.MID)
}
```

## Configuration

The client is configured using functional options:

```go
client, err := maxigo.New("token",
    maxigo.WithTimeout(10 * time.Second),      // HTTP request timeout (default: 30s)
    maxigo.WithHTTPClient(customHTTPClient),    // custom *http.Client
    maxigo.WithBaseURL("https://custom.api"),   // custom base URL
)
```

`WithBaseURL` is useful for testing — point the client at an `httptest.Server`.

## Messages

### Sending

```go
// To a chat
text := "Hello!"
msg, err := client.SendMessage(ctx, chatID, &maxigo.NewMessageBody{
    Text: &text,
})

// To a specific user
text = "Direct message"
msg, err := client.SendMessageToUser(ctx, userID, &maxigo.NewMessageBody{
    Text: &text,
})

// With formatting
text = "**Bold** and _italic_"
format := maxigo.FormatMarkdown
msg, err := client.SendMessage(ctx, chatID, &maxigo.NewMessageBody{
    Text:   &text,
    Format: &format,
})

// With inline keyboard (type-safe constructor)
text = "Choose an action:"
msg, err := client.SendMessage(ctx, chatID, &maxigo.NewMessageBody{
    Text: &text,
    Attachments: []maxigo.AttachmentRequest{
        maxigo.NewInlineKeyboardAttachment([][]maxigo.Button{
            {
                {Type: "callback", Text: "Yes", Payload: "yes", Intent: maxigo.IntentPositive},
                {Type: "callback", Text: "No", Payload: "no", Intent: maxigo.IntentNegative},
            },
        }),
    },
})

// With link button
text = "Visit our website:"
msg, err := client.SendMessage(ctx, chatID, &maxigo.NewMessageBody{
    Text: &text,
    Attachments: []maxigo.AttachmentRequest{
        maxigo.NewInlineKeyboardAttachment([][]maxigo.Button{
            {
                {Type: "link", Text: "Open", URL: "https://example.com"},
            },
        }),
    },
})

// Reply to a message
text = "This is a reply"
msg, err := client.SendMessage(ctx, chatID, &maxigo.NewMessageBody{
    Text: &text,
    Link: &maxigo.NewMessageLink{
        Type: maxigo.LinkReply,
        MID:  "mid-original",
    },
})

// Forward a message
msg, err := client.SendMessage(ctx, chatID, &maxigo.NewMessageBody{
    Link: &maxigo.NewMessageLink{
        Type: maxigo.LinkForward,
        MID:  "mid-to-forward",
    },
})
```

### Editing and Deleting

```go
// Edit a message
text := "Updated text"
result, err := client.EditMessage(ctx, "mid-123", &maxigo.NewMessageBody{
    Text: &text,
})

// Delete a message
result, err := client.DeleteMessage(ctx, "mid-123")
```

### Retrieving Messages

```go
// List messages from a chat
messages, err := client.GetMessages(ctx, maxigo.GetMessagesOpts{ChatID: chatID, Count: 50})

// Get a specific message by ID
msg, err := client.GetMessageByID(ctx, "mid-123")
```

### Answering Callbacks

When a user presses an inline button:

```go
notif := "Done!"
result, err := client.AnswerCallback(ctx, callbackID, &maxigo.CallbackAnswer{
    Notification: &notif,
})

// Or replace the message:
text := "Button was pressed!"
result, err := client.AnswerCallback(ctx, callbackID, &maxigo.CallbackAnswer{
    Message: &maxigo.NewMessageBody{
        Text: &text,
    },
})
```

## Chats

```go
// Get a chat
chat, err := client.GetChat(ctx, chatID)

// Get a chat by invite link
chat, err := client.GetChatByLink(ctx, "https://max.ru/join/abc123")

// List chats (paginated)
list, err := client.GetChats(ctx, maxigo.GetChatsOpts{Count: 50})
// Next page:
list2, err := client.GetChats(ctx, maxigo.GetChatsOpts{Count: 50, Marker: *list.Marker})

// Edit a chat
title := "New Title"
chat, err := client.EditChat(ctx, chatID, &maxigo.ChatPatch{
    Title: &title,
})

// Delete a chat
result, err := client.DeleteChat(ctx, chatID)

// Members
members, err := client.GetMembers(ctx, chatID, maxigo.GetMembersOpts{Count: 100})
admins, err := client.GetAdmins(ctx, chatID)

// Add / remove members
result, err := client.AddMembers(ctx, chatID, []int64{userID1, userID2})
result, err := client.RemoveMember(ctx, chatID, userID, false) // block=false

// Send typing action
result, err := client.SendAction(ctx, chatID, maxigo.ActionTypingOn)

// Pin / unpin messages
result, err := client.PinMessage(ctx, chatID, &maxigo.PinMessageBody{MessageID: "mid-1"})
result, err := client.UnpinMessage(ctx, chatID)
pinned, err := client.GetPinnedMessage(ctx, chatID)

// Bot's own membership
membership, err := client.GetMembership(ctx, chatID)

// Leave a chat
result, err := client.LeaveChat(ctx, chatID)
```

### Sender Actions

| Constant           | Description              |
|--------------------|--------------------------|
| `ActionTypingOn`   | Bot is typing            |
| `ActionSendPhoto`  | Bot is sending a photo   |
| `ActionSendVideo`  | Bot is sending a video   |
| `ActionSendAudio`  | Bot is sending audio     |
| `ActionSendFile`   | Bot is sending a file    |
| `ActionMarkSeen`   | Mark messages as read    |

## File Uploads

Uploading is a two-step process: get an upload URL, then upload the file.

```go
// Photo (simplified)
file, _ := os.Open("photo.jpg")
tokens, err := client.UploadPhoto(ctx, "photo.jpg", file)

// Then send with the token:
client.SendMessage(ctx, chatID, &maxigo.NewMessageBody{
    Attachments: []maxigo.AttachmentRequest{
        {
            Type: "image",
            Payload: maxigo.PhotoAttachmentRequestPayload{
                Photos: tokens.Photos,
            },
        },
    },
})

// Video / audio / file
file, _ := os.Open("video.mp4")
info, err := client.UploadMedia(ctx, maxigo.UploadVideo, "video.mp4", file)

// Manual two-step (if you need more control):
endpoint, err := client.GetUploadURL(ctx, maxigo.UploadFile)
// Then POST the file to endpoint.URL
```

### Upload Types

| Constant      | Description                 |
|---------------|-----------------------------|
| `UploadImage` | Image files (jpg, png, gif) |
| `UploadVideo` | Video files                 |
| `UploadAudio` | Audio files                 |
| `UploadFile`  | Any file                    |

## Subscriptions (Webhooks)

```go
// Subscribe to updates
result, err := client.Subscribe(ctx,
    "https://example.com/webhook",
    []string{"message_created", "message_callback"},
    "my-secret",
)

// Unsubscribe
result, err := client.Unsubscribe(ctx, "https://example.com/webhook")

// List active subscriptions
subs, err := client.GetSubscriptions(ctx)
for _, s := range subs {
    fmt.Printf("Webhook: %s, types: %v\n", s.URL, s.UpdateTypes)
}
```

## Long Polling

```go
var marker int64

for {
    result, err := client.GetUpdates(ctx, maxigo.GetUpdatesOpts{Limit: 100, Timeout: 30, Marker: marker})
    if err != nil {
        log.Println("error:", err)
        time.Sleep(time.Second)
        continue
    }

    for _, raw := range result.Updates {
        var base maxigo.Update
        json.Unmarshal(raw, &base)

        switch base.UpdateType {
        case maxigo.UpdateMessageCreated:
            var upd maxigo.MessageCreatedUpdate
            json.Unmarshal(raw, &upd)
            fmt.Println("New message:", *upd.Message.Body.Text)

        case maxigo.UpdateMessageCallback:
            var upd maxigo.MessageCallbackUpdate
            json.Unmarshal(raw, &upd)
            fmt.Println("Callback:", upd.Callback.Payload)

        case maxigo.UpdateBotStarted:
            var upd maxigo.BotStartedUpdate
            json.Unmarshal(raw, &upd)
            fmt.Printf("User %d pressed Start\n", upd.User.UserID)

        case maxigo.UpdateBotAdded:
            var upd maxigo.BotAddedUpdate
            json.Unmarshal(raw, &upd)
            fmt.Printf("Bot added to chat %d\n", upd.ChatID)

        case maxigo.UpdateUserAdded:
            var upd maxigo.UserAddedUpdate
            json.Unmarshal(raw, &upd)
            fmt.Printf("User %d added to chat %d\n", upd.User.UserID, upd.ChatID)
        }
    }

    if result.Marker != nil {
        marker = *result.Marker
    }
}
```

### Update Types

| Constant                     | Struct                       | Description              |
|------------------------------|------------------------------|--------------------------|
| `UpdateMessageCreated`       | `MessageCreatedUpdate`       | New message              |
| `UpdateMessageCallback`      | `MessageCallbackUpdate`      | Inline button pressed    |
| `UpdateMessageEdited`        | `MessageEditedUpdate`        | Message edited           |
| `UpdateMessageRemoved`       | `MessageRemovedUpdate`       | Message deleted          |
| `UpdateBotStarted`           | `BotStartedUpdate`           | User pressed Start       |
| `UpdateBotAdded`             | `BotAddedUpdate`             | Bot added to chat        |
| `UpdateBotRemoved`           | `BotRemovedUpdate`           | Bot removed from chat    |
| `UpdateUserAdded`            | `UserAddedUpdate`            | User added to chat       |
| `UpdateUserRemoved`          | `UserRemovedUpdate`          | User removed from chat   |
| `UpdateChatTitleChanged`     | `ChatTitleChangedUpdate`     | Chat title changed       |
| `UpdateMessageChatCreated`   | `MessageChatCreatedUpdate`   | Chat created via button  |

## Error Handling

All errors are returned as `*maxigo.Error` with structured fields:

```go
msg, err := client.SendMessage(ctx, chatID, body)
if err != nil {
    var e *maxigo.Error
    if errors.As(err, &e) {
        switch e.Kind {
        case maxigo.ErrAPI:
            // API returned non-200: e.StatusCode (401, 403, 404, 429, 500...)
            fmt.Printf("API error %d: %s\n", e.StatusCode, e.Message)
        case maxigo.ErrNetwork:
            // Connection or DNS failure
            fmt.Println("Network:", e.Message)
        case maxigo.ErrTimeout:
            // Request timeout or context cancellation
            fmt.Println("Timeout")
        case maxigo.ErrDecode:
            // JSON marshal/unmarshal failure
            fmt.Println("Decode error:", e.Message)
        }
        // e.Op — operation name ("SendMessage", "GetChat", ...)
        // e.Err — underlying error (for Unwrap)
    }
}
```

### Error Kinds

| Kind         | Description                                |
|--------------|--------------------------------------------|
| `ErrAPI`     | HTTP response with status != 200           |
| `ErrNetwork` | Connection, DNS, or transport failure      |
| `ErrTimeout` | Request timeout or `context` cancellation  |
| `ErrDecode`  | JSON marshal/unmarshal failure             |

### Error Methods

- `e.Error() string` — formatted error message including Op, Kind, StatusCode
- `e.Timeout() bool` — returns `true` for `ErrTimeout`
- `e.Unwrap() error` — returns the underlying error for `errors.Is/As` chains

### Handling Specific HTTP Status Codes

```go
var e *maxigo.Error
if errors.As(err, &e) && e.Kind == maxigo.ErrAPI {
    switch e.StatusCode {
    case 401:
        // Invalid token
    case 403:
        // No permission
    case 404:
        // Chat/message not found
    case 429:
        // Rate limited — back off and retry
    }
}
```

## Max Bot API Quirks

- Commands use `:` as separator (not space like Telegram): `/start:payload`. No `@botname`.
- `MessageCallbackUpdate` has no direct ChatID — extract from `Message.Recipient.ChatId`.

## Testing

The client is fully testable without hitting the real API:

```go
func TestMyBot(t *testing.T) {
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(maxigo.BotInfo{
            UserWithPhoto: maxigo.UserWithPhoto{
                User: maxigo.User{UserID: 1, FirstName: "TestBot", IsBot: true},
            },
        })
    }))
    defer srv.Close()

    client, _ := maxigo.New("test-token", maxigo.WithBaseURL(srv.URL))
    bot, err := client.GetBot(context.Background())
    // assert...
}
```

## Ecosystem

| Package                          | Description                                  |
|----------------------------------|----------------------------------------------|
| `github.com/maxigo-bot/maxigo-client` | HTTP client (this package)                   |
| `github.com/maxigo/maxigo`       | Bot framework with router/middleware/context |

> *Both packages are currently in development and not yet published.*

## License

MIT
