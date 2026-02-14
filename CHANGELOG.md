# Changelog

## [v0.1.0] - 2026-02-14

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