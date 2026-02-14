# maxigo-client

Go HTTP-клиент для [Max Bot API](https://dev.max.ru). Без внешних зависимостей.

> **[English Guide](guide.md)** | **[README](../README.md)**

## Зачем этот проект?

Официальный [`max-bot-api-client-go`](https://github.com/max-messenger/max-bot-api-client-go) имеет системные проблемы, делающие его непригодным для продакшена:

**Ошибки логируются вместо возврата** — 30+ мест с `log.Println` и `slog.Error` прямо в библиотеке. Пользователь не может подавить или перенаправить эти логи. Некоторые ошибки молча проглатываются (`json.Decode` падает — возвращается `nil`).

**Невозможно тестировать без реального API** — нет простого `WithBaseURL()`. Загрузка файлов идёт через `http.DefaultClient` напрямую, минуя настройки клиента. Для тестирования нужно реализовать `ConfigInterface` из 7 методов.

**6 внешних зависимостей** — zerolog, YAML-парсер, парсер env-переменных, gomock — всё это не нужно HTTP-клиенту.

**Сломанные методы** — `GetChatID()` возвращает 0 для callback (хотя chat ID есть, но игнорируется). `GetCommand()` возвращает весь текст сообщения. `schemes.Error` используется как структура ответа и всегда non-nil, поэтому `Check()` всегда возвращает ошибку.

**Неправильные типы** — `time.Duration` для Unix-таймстампов (интерпретирует как наносекунды). `int64→int` в 10+ местах (обрезка на 32-бит). `[]interface{}` для вложений (никакой типобезопасности).

**Неидиоматичный Go** — builder-паттерн, `SCREAMING_CASE` константы, `Api` вместо `API`, нет `context.Context` в загрузках, нет функциональных опций.

| Проблема                   | Официальный клиент                                   | maxigo-client                                             |
|----------------------------|------------------------------------------------------|-----------------------------------------------------------|
| Обработка ошибок           | `log.Println` в 30+ местах                           | Все ошибки возвращаются как `*Error` с Kind/StatusCode/Op |
| Тестируемость              | Нужен мок `ConfigInterface` из 7 методов             | `maxigo.New("token", WithBaseURL(srv.URL))`               |
| Зависимости                | 6 транзитивных (zerolog, yaml, env...)               | 0 — только stdlib                                         |
| `GetChatID()` для callback | Возвращает 0                                         | Извлекаем из `Message.Recipient.ChatId`                   |
| Типы                       | `time.Duration` для таймстампов, `int→int64` кастинг | Корректный `int64` везде                                  |
| Загрузки файлов            | `http.Get()` без context/timeout                     | Все запросы через настроенный клиент с `context.Context`  |
| Стиль API                  | `NewMessage().SetChat().SetText()`                   | `SendMessage(ctx, chatID, &NewMessageBody{Text: &text})`  |
| Константы                  | `TYPING_ON`, `CALLBACK`, `POSITIVE`                  | `ActionTypingOn`, `IntentPositive`                        |
| Конфигурация               | YAML-файлы + парсер env                              | Функциональные опции: `WithTimeout`, `WithHTTPClient`     |

**maxigo-client** исправляет все эти проблемы.

## Установка

```bash
go get github.com/maxigo-bot/maxigo-client
```

Требуется Go 1.25+.

## Быстрый старт

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

    // Получаем информацию о боте
    bot, err := client.GetBot(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Бот: %s (ID: %d)\n", bot.FirstName, bot.UserID)

    // Отправляем сообщение
    text := "Привет из maxigo!"
    msg, err := client.SendMessage(ctx, 123456, &maxigo.NewMessageBody{
        Text: &text,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Отправлено: %s\n", msg.Body.MID)
}
```

## Конфигурация

Клиент настраивается через функциональные опции:

```go
client, err := maxigo.New("token",
    maxigo.WithTimeout(10 * time.Second),      // таймаут HTTP-запросов (по умолчанию 30с)
    maxigo.WithHTTPClient(customHTTPClient),    // свой *http.Client
    maxigo.WithBaseURL("https://custom.api"),   // другой базовый URL
)
```

`WithBaseURL` полезен для тестирования — можно направить клиент на `httptest.Server`.

## Работа с сообщениями

### Отправка

```go
// В чат
text := "Привет!"
msg, err := client.SendMessage(ctx, chatID, &maxigo.NewMessageBody{
    Text: &text,
})

// Конкретному пользователю
text := "Личное сообщение"
msg, err := client.SendMessageToUser(ctx, userID, &maxigo.NewMessageBody{
    Text: &text,
})

// С форматированием
text := "**Жирный** и _курсив_"
format := maxigo.FormatMarkdown
msg, err := client.SendMessage(ctx, chatID, &maxigo.NewMessageBody{
    Text:   &text,
    Format: &format,
})

// С инлайн-клавиатурой (типобезопасный конструктор)
text := "Выберите действие:"
msg, err := client.SendMessage(ctx, chatID, &maxigo.NewMessageBody{
    Text: &text,
    Attachments: []maxigo.AttachmentRequest{
        maxigo.NewInlineKeyboardAttachment([][]maxigo.Button{
            {
                {Type: "callback", Text: "Да", Payload: "yes", Intent: maxigo.IntentPositive},
                {Type: "callback", Text: "Нет", Payload: "no", Intent: maxigo.IntentNegative},
            },
        }),
    },
})
```

### Редактирование и удаление

```go
// Редактировать сообщение
text := "Обновлённый текст"
result, err := client.EditMessage(ctx, "mid-123", &maxigo.NewMessageBody{
    Text: &text,
})

// Удалить сообщение
result, err := client.DeleteMessage(ctx, "mid-123")
```

### Получение сообщений

```go
// Список сообщений из чата
messages, err := client.GetMessages(ctx, maxigo.GetMessagesOpts{ChatID: chatID, Count: 50})

// Конкретное сообщение по ID
msg, err := client.GetMessageByID(ctx, "mid-123")
```

### Ответ на callback (нажатие кнопки)

```go
notif := "Готово!"
result, err := client.AnswerCallback(ctx, callbackID, &maxigo.CallbackAnswer{
    Notification: &notif,
})
```

## Работа с чатами

```go
// Получить чат
chat, err := client.GetChat(ctx, chatID)

// Список чатов (с пагинацией)
list, err := client.GetChats(ctx, maxigo.GetChatsOpts{Count: 50})
// Следующая страница:
list2, err := client.GetChats(ctx, maxigo.GetChatsOpts{Count: 50, Marker: *list.Marker})

// Редактировать чат
title := "Новое название"
chat, err := client.EditChat(ctx, chatID, &maxigo.ChatPatch{
    Title: &title,
})

// Удалить чат
result, err := client.DeleteChat(ctx, chatID)

// Участники
members, err := client.GetMembers(ctx, chatID, maxigo.GetMembersOpts{Count: 100})
admins, err := client.GetAdmins(ctx, chatID)

// Добавить/удалить участников
result, err := client.AddMembers(ctx, chatID, []int64{userID1, userID2})
result, err := client.RemoveMember(ctx, chatID, userID, false) // block=false

// Отправить действие (набирает текст...)
result, err := client.SendAction(ctx, chatID, maxigo.ActionTypingOn)

// Закреплённое сообщение
result, err := client.PinMessage(ctx, chatID, &maxigo.PinMessageBody{MessageID: "mid-1"})
result, err := client.UnpinMessage(ctx, chatID)
pinned, err := client.GetPinnedMessage(ctx, chatID)

// Покинуть чат
result, err := client.LeaveChat(ctx, chatID)
```

## Загрузка файлов

Загрузка выполняется в два шага: получение URL для загрузки, затем сама загрузка.

```go
// Фото (упрощённый метод)
file, _ := os.Open("photo.jpg")
tokens, err := client.UploadPhoto(ctx, "photo.jpg", file)

// Затем отправить с токеном:
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

// Видео / аудио / файл
file, _ := os.Open("video.mp4")
info, err := client.UploadMedia(ctx, maxigo.UploadVideo, "video.mp4", file)
```

## Подписки (Webhooks)

```go
// Подписаться
result, err := client.Subscribe(ctx,
    "https://example.com/webhook",
    []string{"message_created", "message_callback"},
    "my-secret",
)

// Отписаться
result, err := client.Unsubscribe(ctx, "https://example.com/webhook")

// Список подписок
subs, err := client.GetSubscriptions(ctx)
```

## Получение обновлений (Long Polling)

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
            fmt.Println("Новое сообщение:", *upd.Message.Body.Text)

        case maxigo.UpdateMessageCallback:
            var upd maxigo.MessageCallbackUpdate
            json.Unmarshal(raw, &upd)
            fmt.Println("Callback:", upd.Callback.Payload)

        case maxigo.UpdateBotStarted:
            var upd maxigo.BotStartedUpdate
            json.Unmarshal(raw, &upd)
            fmt.Printf("Пользователь %d нажал Start\n", upd.User.UserID)
        }
    }

    if result.Marker != nil {
        marker = *result.Marker
    }
}
```

### Типы обновлений

| Константа                  | Тип структуры              | Описание                    |
|----------------------------|----------------------------|-----------------------------|
| `UpdateMessageCreated`     | `MessageCreatedUpdate`     | Новое сообщение             |
| `UpdateMessageCallback`    | `MessageCallbackUpdate`    | Нажатие инлайн-кнопки       |
| `UpdateMessageEdited`      | `MessageEditedUpdate`      | Сообщение отредактировано   |
| `UpdateMessageRemoved`     | `MessageRemovedUpdate`     | Сообщение удалено           |
| `UpdateBotStarted`         | `BotStartedUpdate`         | Пользователь нажал Start    |
| `UpdateBotAdded`           | `BotAddedUpdate`           | Бот добавлен в чат          |
| `UpdateBotRemoved`         | `BotRemovedUpdate`         | Бот удалён из чата          |
| `UpdateUserAdded`          | `UserAddedUpdate`          | Пользователь добавлен в чат |
| `UpdateUserRemoved`        | `UserRemovedUpdate`        | Пользователь удалён из чата |
| `UpdateChatTitleChanged`   | `ChatTitleChangedUpdate`   | Название чата изменено      |
| `UpdateMessageChatCreated` | `MessageChatCreatedUpdate` | Чат создан через кнопку     |

## Обработка ошибок

Все ошибки возвращаются как `*maxigo.Error` со структурированными полями:

```go
msg, err := client.SendMessage(ctx, chatID, body)
if err != nil {
    var e *maxigo.Error
    if errors.As(err, &e) {
        switch e.Kind {
        case maxigo.ErrAPI:
            // Ошибка от API: e.StatusCode (401, 403, 404, 429, 500...)
            fmt.Printf("Ошибка API %d: %s\n", e.StatusCode, e.Message)
        case maxigo.ErrNetwork:
            // Проблемы с сетью
            fmt.Println("Сеть:", e.Message)
        case maxigo.ErrTimeout:
            // Таймаут или отмена контекста
            fmt.Println("Таймаут")
        case maxigo.ErrDecode:
            // Ошибка парсинга JSON
            fmt.Println("Ошибка декодирования:", e.Message)
        }
        // e.Op — название операции ("SendMessage", "GetChat", ...)
        // e.Err — оригинальная ошибка (для Unwrap)
    }
}
```

| ErrorKind    | Описание                                  |
|--------------|-------------------------------------------|
| `ErrAPI`     | HTTP-ответ с кодом != 200                 |
| `ErrNetwork` | Ошибка соединения, DNS                    |
| `ErrTimeout` | Таймаут запроса или отмена `context`      |
| `ErrDecode`  | Ошибка сериализации/десериализации JSON   |

Дополнительные методы:
- `e.Timeout() bool` — `true` для ErrTimeout
- `e.Unwrap() error` — оригинальная ошибка для цепочки `errors.Is/As`

## Особенности Max Bot API

- Команды используют `:` как разделитель (не пробел как в Telegram): `/start:payload`
- `MessageCallbackUpdate` не содержит прямого ChatID — извлекайте из `Message.Recipient.ChatId`

## Экосистема

```
github.com/maxigo-bot/maxigo-client  — HTTP-клиент (этот пакет)
github.com/maxigo/maxigo         — фреймворк для ботов с роутером/middleware/контекстом
```

> *Оба пакета находятся в разработке и ещё не опубликованы.*

## Лицензия

MIT
