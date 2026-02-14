# Changelog

## [v0.2.0] - 2026-02-14

### Fixed
- **Optional null handling**: `Optional[T].UnmarshalJSON` now treats JSON `null` as unset (`Set=false`) instead of `Set=true` with zero value — prevents accidental field overwrites when round-tripping API responses
- **Nil body panic**: `SendMessage` and `SendMessageToUser` no longer panic when `body` is nil
- **Upload goroutine leak**: `doUpload` now closes the pipe reader on all exit paths, preventing goroutine leaks on network errors
- **Upload timeout**: `doUpload` now applies the client's default timeout when the context has no deadline, matching `do()` behavior
- **Long polling timeout**: `GetUpdates` no longer races `http.Client.Timeout` against the server-side polling timeout — timeouts are now managed via `context` with an automatic 5s buffer on top of the server-side duration
- **JSON serialization**: `NewMessageBody.Attachments` changed from `omitempty` to `omitzero` — empty slice `[]AttachmentRequest{}` is now serialized as `"attachments":[]` instead of being omitted, allowing inline keyboard removal

### Added
- `Attachment` interface and `MessageBody.ParseAttachments()` — type-safe parsing of response attachments via discriminator map (all 11 types: image, video, audio, file, sticker, contact, share, location, data, inline_keyboard, reply_keyboard); unknown types are skipped for forward compatibility
- Type-safe button constructors: `NewCallbackButton`, `NewCallbackButtonWithIntent`, `NewLinkButton`, `NewRequestContactButton`, `NewRequestGeoLocationButton`, `NewChatButton`, `NewMessageButton`
- `ErrPollDeadline` — returned when the context deadline is too short for the requested polling timeout
- `Optional[T]`, `Some[T]()` — generic optional type with three-state JSON semantics (unset / zero value / value), replaces `*string`/`*bool` pointers in request types
- Type aliases: `OptString`, `OptBool`, `OptInt64`

### Changed
- `NewMessageBody.Text`: `*string` → `OptString`
- `NewMessageBody.Notify`: `*bool` → `OptBool`
- `NewMessageBody.Format`: `*TextFormat` → `Optional[TextFormat]`
- `BotPatch.Name`, `.FirstName`, `.Description`: `*string` → `OptString`
- `ChatPatch.Title`, `.Pin`: `*string` → `OptString`
- `ChatPatch.Notify`: `*bool` → `OptBool`
- `CallbackAnswer.Notification`: `*string` → `OptString`
- `PinMessageBody.Notify`: `*bool` → `OptBool`
- `Button.ChatDescription`, `.StartPayload`: `*string` → `OptString`
- `Button.UUID`: `*int64` → `OptInt64`
- `ReplyButton.Payload`: `*string` → `OptString`
- `AttachmentRequest.DirectUserID`: `*int64` → `OptInt64`
- `PhotoAttachmentRequestPayload.URL`, `.Token`: `*string` → `OptString`
- `ContactAttachmentRequestPayload.Name`, `.VCFInfo`, `.VCFPhone`: `*string` → `OptString`
- `ContactAttachmentRequestPayload.ContactID`: `*int64` → `OptInt64`
- `ShareAttachmentPayload.URL`, `.Token`: `*string` → `OptString`

## [v0.1.1] - 2026-02-14

### Fixed
- **Query params**: `message_ids` and `types` are now comma-separated (`style: simple`) as required by the OpenAPI schema
- **URL-encoding**: `GetChatByLink` and `GetMessageByID` now escape special characters in path via `url.PathEscape`
- **Security**: token is sent via `Authorization` header instead of `access_token` query parameter — prevents token leaking into proxy and CDN logs
- **Tests**: `requestCount` in upload tests replaced with `atomic.Int32` to eliminate race condition
- **Options**: `WithHTTPClient` now makes a shallow copy — `WithTimeout` no longer mutates the external `*http.Client`
- **Types**: `Button.UUID` changed from `*int` to `*int64` for consistency with other ID fields
- **Leak fix**: `doUpload` — added `pr.Close()` on request creation error
- **Types**: `interface{}` replaced with `any` throughout (`AttachmentRequest.Payload`, `do()` signature)
- **Types**: `VideoURLs` fields renamed to idiomatic Go (`MP4_1080` → `MP41080`, etc.); JSON tags unchanged
- **Types**: `UserIdsList` → `UserIDsList` (Go naming convention for acronyms)
- **Types**: `User.LastName` and `User.Username` changed from `string` to `*string` (nullable in API)
- **Types**: `NewMessageBody.Text` changed from `string` to `*string` — `nil` means "keep existing text" on edit, `""` means "clear text"
- **Types**: `GetAdmins` now returns `*ChatAdminsList` instead of `*ChatMembersList`
- **Types**: `BotCommand.Description` changed from `string` to `*string` (nullable in schema)
- **Types**: `ContactAttachmentRequestPayload.Name` changed from `string` to `*string` (nullable in schema)
- **Code quality**: `isTimeout` uses `errors.As` instead of direct type assertion (handles wrapped errors)

### Added
- **Endpoints**: `SetAdmins` (POST /chats/{chatId}/members/admins), `RemoveAdmin` (DELETE /chats/{chatId}/members/admins/{userId}), `GetVideoDetails` (GET /videos/{videoToken})
- **Params**: `DisableLinkPreview` field on `NewMessageBody` (sent as query param, not JSON)
- **Params**: `from`/`to` (time range filter) in `GetMessages` via new `GetMessagesOpts`
- **Params**: `user_ids` in `GetMembers` via new `GetMembersOpts`
- **Opts**: `GetChatsOpts` struct for `GetChats` pagination parameters
- **Opts**: `GetUpdatesOpts` struct for `GetUpdates` long-polling parameters
- Type-safe attachment constructors: `NewPhotoAttachment`, `NewVideoAttachment`, `NewAudioAttachment`, `NewFileAttachment`, `NewStickerAttachment`, `NewContactAttachment`, `NewShareAttachment`, `NewInlineKeyboardAttachment`, `NewLocationAttachment`
- Tests: `TestGetMessagesWithIDs`, `TestGetChatByLink`, `TestGetMembership`, `TestUnpinMessage`, `TestGetPinnedMessage`, `types_test.go` (all 9 attachment constructors)