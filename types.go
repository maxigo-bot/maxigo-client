package maxigo

import "encoding/json"

// ChatType represents the type of chat.
type ChatType string

const (
	ChatDialog  ChatType = "dialog"
	ChatGroup   ChatType = "chat"
	ChatChannel ChatType = "channel"
)

// ChatStatus represents the bot's status in a chat.
type ChatStatus string

const (
	ChatStatusActive    ChatStatus = "active"
	ChatStatusRemoved   ChatStatus = "removed"
	ChatStatusLeft      ChatStatus = "left"
	ChatStatusClosed    ChatStatus = "closed"
	ChatStatusSuspended ChatStatus = "suspended"
)

// SenderAction represents an action to send to chat members.
type SenderAction string

const (
	ActionTypingOn  SenderAction = "typing_on"
	ActionSendPhoto SenderAction = "sending_photo"
	ActionSendVideo SenderAction = "sending_video"
	ActionSendAudio SenderAction = "sending_audio"
	ActionSendFile  SenderAction = "sending_file"
	ActionMarkSeen  SenderAction = "mark_seen"
)

// UploadType represents the type of file being uploaded.
type UploadType string

const (
	UploadImage UploadType = "image"
	UploadVideo UploadType = "video"
	UploadAudio UploadType = "audio"
	UploadFile  UploadType = "file"
)

// Intent represents the visual intent of a button.
type Intent string

const (
	IntentPositive Intent = "positive"
	IntentNegative Intent = "negative"
	IntentDefault  Intent = "default"
)

// MessageLinkType represents the type of linked message.
type MessageLinkType string

const (
	LinkForward MessageLinkType = "forward"
	LinkReply   MessageLinkType = "reply"
)

// TextFormat represents the text formatting mode.
type TextFormat string

const (
	FormatMarkdown TextFormat = "markdown"
	FormatHTML     TextFormat = "html"
)

// UpdateType represents the type of an update event.
type UpdateType string

const (
	UpdateMessageCreated     UpdateType = "message_created"
	UpdateMessageCallback    UpdateType = "message_callback"
	UpdateMessageEdited      UpdateType = "message_edited"
	UpdateMessageRemoved     UpdateType = "message_removed"
	UpdateBotStarted         UpdateType = "bot_started"
	UpdateBotAdded           UpdateType = "bot_added"
	UpdateBotRemoved         UpdateType = "bot_removed"
	UpdateUserAdded          UpdateType = "user_added"
	UpdateUserRemoved        UpdateType = "user_removed"
	UpdateChatTitleChanged   UpdateType = "chat_title_changed"
	UpdateMessageChatCreated UpdateType = "message_chat_created"
)

// ChatAdminPermission represents a permission granted to a chat admin.
type ChatAdminPermission string

const (
	PermReadAllMessages  ChatAdminPermission = "read_all_messages"
	PermAddRemoveMembers ChatAdminPermission = "add_remove_members"
	PermAddAdmins        ChatAdminPermission = "add_admins"
	PermChangeChatInfo   ChatAdminPermission = "change_chat_info"
	PermPinMessage       ChatAdminPermission = "pin_message"
	PermWrite            ChatAdminPermission = "write"
)

// User represents a Max user.
type User struct {
	UserID           int64  `json:"user_id"`
	FirstName        string `json:"first_name"`
	LastName         *string `json:"last_name,omitempty"`
	Username         *string `json:"username,omitempty"`
	IsBot            bool   `json:"is_bot"`
	LastActivityTime int64  `json:"last_activity_time"`
}

// UserWithPhoto extends User with avatar and description.
type UserWithPhoto struct {
	User
	Description   *string `json:"description,omitempty"`
	AvatarURL     string  `json:"avatar_url,omitempty"`
	FullAvatarURL string  `json:"full_avatar_url,omitempty"`
}

// BotInfo represents the bot's info returned by GET /me.
type BotInfo struct {
	UserWithPhoto
	Commands []BotCommand `json:"commands,omitempty"`
}

// BotCommand represents a command supported by the bot.
type BotCommand struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// BotPatch represents the request body for PATCH /me.
type BotPatch struct {
	// Deprecated: use FirstName instead. Will be removed in a future API version.
	Name        *string      `json:"name,omitempty"`
	FirstName   *string      `json:"first_name,omitempty"`
	Description *string      `json:"description,omitempty"`
	Commands    []BotCommand `json:"commands,omitempty"`
	Photo       *PhotoAttachmentRequestPayload `json:"photo,omitempty"`
}

// Image represents a generic image object.
type Image struct {
	URL string `json:"url"`
}

// Chat represents a Max chat.
type Chat struct {
	ChatID            int64          `json:"chat_id"`
	Type              ChatType       `json:"type"`
	Status            ChatStatus     `json:"status"`
	Title             *string        `json:"title"`
	Icon              *Image         `json:"icon"`
	LastEventTime     int64          `json:"last_event_time"`
	ParticipantsCount int            `json:"participants_count"`
	OwnerID           *int64         `json:"owner_id,omitempty"`
	Participants      map[string]int64 `json:"participants,omitempty"`
	IsPublic          bool           `json:"is_public"`
	Link              *string        `json:"link,omitempty"`
	Description       *string        `json:"description"`
	DialogWithUser    *UserWithPhoto `json:"dialog_with_user,omitempty"`
	MessagesCount     *int           `json:"messages_count,omitempty"`
	ChatMessageID     *string        `json:"chat_message_id,omitempty"`
	PinnedMessage     *Message       `json:"pinned_message,omitempty"`
}

// ChatList represents a paginated list of chats.
type ChatList struct {
	Chats  []Chat `json:"chats"`
	Marker *int64 `json:"marker"`
}

// ChatPatch represents the request body for PATCH /chats/{chatId}.
type ChatPatch struct {
	Icon   *PhotoAttachmentRequestPayload `json:"icon,omitempty"`
	Title  *string                        `json:"title,omitempty"`
	Pin    *string                        `json:"pin,omitempty"`
	Notify *bool                          `json:"notify,omitempty"`
}

// ChatMember represents a member of a chat.
type ChatMember struct {
	UserWithPhoto
	LastAccessTime int64                 `json:"last_access_time"`
	IsOwner        bool                  `json:"is_owner"`
	IsAdmin        bool                  `json:"is_admin"`
	JoinTime       int64                 `json:"join_time"`
	Permissions    []ChatAdminPermission `json:"permissions"`
	Alias          *string               `json:"alias"`
}

// ChatMembersList represents a paginated list of chat members.
type ChatMembersList struct {
	Members []ChatMember `json:"members"`
	Marker  *int64       `json:"marker,omitempty"`
}

// ChatAdmin represents an administrator with permissions.
type ChatAdmin struct {
	UserID      int64                 `json:"user_id"`
	Permissions []ChatAdminPermission `json:"permissions"`
	Alias       *string               `json:"alias,omitempty"`
}

// ChatAdminsList represents a list of chat admins.
type ChatAdminsList struct {
	Admins []ChatAdmin `json:"admins"`
}

// Recipient represents a message recipient (user or chat).
type Recipient struct {
	ChatID   *int64   `json:"chat_id"`
	ChatType ChatType `json:"chat_type"`
	UserID   *int64   `json:"user_id"`
}

// MessageStat represents message statistics.
type MessageStat struct {
	Views int `json:"views"`
}

// Message represents a message in a chat.
type Message struct {
	Sender    *User          `json:"sender,omitempty"`
	Recipient Recipient      `json:"recipient"`
	Timestamp int64          `json:"timestamp"`
	Link      *LinkedMessage `json:"link,omitempty"`
	Body      MessageBody    `json:"body"`
	Stat      *MessageStat   `json:"stat,omitempty"`
	URL       *string        `json:"url,omitempty"`
}

// MessageBody represents the body of a message.
type MessageBody struct {
	MID         string            `json:"mid"`
	Seq         int64             `json:"seq"`
	Text        *string           `json:"text"`
	Attachments []json.RawMessage `json:"attachments"`
	Markup      []MarkupElement   `json:"markup,omitempty"`
}

// LinkedMessage represents a forwarded or replied message.
type LinkedMessage struct {
	Type    MessageLinkType `json:"type"`
	Sender  *User           `json:"sender,omitempty"`
	ChatID  int64           `json:"chat_id,omitempty"`
	Message MessageBody     `json:"message"`
}

// MessageList represents a paginated list of messages.
type MessageList struct {
	Messages []Message `json:"messages"`
}

// NewMessageBody represents the body for sending or editing a message.
//
// Text is a pointer to distinguish between "not set" (nil, field omitted)
// and "set to empty string" (clears message text). When editing a message,
// nil means "keep existing text", while a pointer to "" means "clear text".
type NewMessageBody struct {
	Text        *string            `json:"text"`
	Attachments []AttachmentRequest `json:"attachments,omitempty"`
	Link        *NewMessageLink    `json:"link,omitempty"`
	Notify      *bool              `json:"notify,omitempty"`
	Format      *TextFormat        `json:"format,omitempty"`

	// DisableLinkPreview prevents the server from generating link previews.
	// Sent as a query parameter, not in the JSON body.
	DisableLinkPreview bool `json:"-"`
}

// NewMessageLink represents a link to another message (for reply/forward).
type NewMessageLink struct {
	Type MessageLinkType `json:"type"`
	MID  string          `json:"mid"`
}

// sendMessageResult represents the response from POST /messages.
type sendMessageResult struct {
	Message Message `json:"message"`
}

// Attachment types

// AttachmentType is embedded in all attachment responses.
type AttachmentType struct {
	Type string `json:"type"`
}

// PhotoAttachmentPayload represents the payload of a photo attachment.
type PhotoAttachmentPayload struct {
	PhotoID int64  `json:"photo_id"`
	Token   string `json:"token"`
	URL     string `json:"url"`
}

// PhotoAttachment represents an image attachment in a message.
type PhotoAttachment struct {
	AttachmentType
	Payload PhotoAttachmentPayload `json:"payload"`
}

// MediaAttachmentPayload represents the payload of video/audio attachments.
type MediaAttachmentPayload struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

// VideoThumbnail represents a video thumbnail.
type VideoThumbnail struct {
	URL string `json:"url"`
}

// VideoAttachment represents a video attachment in a message.
type VideoAttachment struct {
	AttachmentType
	Payload   MediaAttachmentPayload `json:"payload"`
	Thumbnail *VideoThumbnail        `json:"thumbnail,omitempty"`
	Width     *int                   `json:"width,omitempty"`
	Height    *int                   `json:"height,omitempty"`
	Duration  *int                   `json:"duration,omitempty"`
}

// AudioAttachment represents an audio attachment in a message.
type AudioAttachment struct {
	AttachmentType
	Payload       MediaAttachmentPayload `json:"payload"`
	Transcription *string                `json:"transcription,omitempty"`
}

// FileAttachmentPayload represents the payload of a file attachment.
type FileAttachmentPayload struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

// FileAttachment represents a file attachment in a message.
type FileAttachment struct {
	AttachmentType
	Payload  FileAttachmentPayload `json:"payload"`
	Filename string                `json:"filename"`
	Size     int64                 `json:"size"`
}

// StickerAttachmentPayload represents the payload of a sticker attachment.
type StickerAttachmentPayload struct {
	URL  string `json:"url"`
	Code string `json:"code"`
}

// StickerAttachment represents a sticker attachment in a message.
type StickerAttachment struct {
	AttachmentType
	Payload StickerAttachmentPayload `json:"payload"`
	Width   int                      `json:"width"`
	Height  int                      `json:"height"`
}

// ContactAttachmentPayload represents the payload of a contact attachment.
type ContactAttachmentPayload struct {
	VCFInfo *string `json:"vcf_info,omitempty"`
	MaxInfo *User   `json:"max_info,omitempty"`
}

// ContactAttachment represents a contact attachment in a message.
type ContactAttachment struct {
	AttachmentType
	Payload ContactAttachmentPayload `json:"payload"`
}

// ShareAttachmentPayload represents the payload of a share attachment.
type ShareAttachmentPayload struct {
	URL   *string `json:"url,omitempty"`
	Token *string `json:"token,omitempty"`
}

// ShareAttachment represents a link share attachment in a message.
type ShareAttachment struct {
	AttachmentType
	Payload     ShareAttachmentPayload `json:"payload"`
	Title       *string                `json:"title,omitempty"`
	Description *string                `json:"description,omitempty"`
	ImageURL    *string                `json:"image_url,omitempty"`
}

// LocationAttachment represents a location attachment in a message.
type LocationAttachment struct {
	AttachmentType
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// DataAttachment contains a payload sent through a SendMessageButton.
type DataAttachment struct {
	AttachmentType
	Data string `json:"data"`
}

// Keyboard represents a two-dimensional array of buttons.
type Keyboard struct {
	Buttons [][]Button `json:"buttons"`
}

// InlineKeyboardAttachment represents an inline keyboard in a message.
type InlineKeyboardAttachment struct {
	AttachmentType
	Payload Keyboard `json:"payload"`
}

// Button represents a button in an inline keyboard.
// Use the Type field to determine the button kind: "callback", "link",
// "request_contact", "request_geo_location", "chat", "message".
type Button struct {
	Type             string  `json:"type"`
	Text             string  `json:"text"`
	Payload          string  `json:"payload,omitempty"`
	URL              string  `json:"url,omitempty"`
	Intent           Intent  `json:"intent,omitempty"`
	Quick            bool    `json:"quick,omitempty"`
	ChatTitle        string  `json:"chat_title,omitempty"`
	ChatDescription  *string `json:"chat_description,omitempty"`
	StartPayload     *string `json:"start_payload,omitempty"`
	UUID             *int64  `json:"uuid,omitempty"`
}

// ReplyButton represents a button in a reply keyboard.
type ReplyButton struct {
	Type    string  `json:"type,omitempty"`
	Text    string  `json:"text"`
	Payload *string `json:"payload,omitempty"`
	Intent  Intent  `json:"intent,omitempty"`
	Quick   bool    `json:"quick,omitempty"`
}

// ReplyKeyboardAttachment represents a reply keyboard in a message.
type ReplyKeyboardAttachment struct {
	AttachmentType
	Buttons [][]ReplyButton `json:"buttons"`
}

// Attachment request types (for sending messages)

// AttachmentRequest represents a request to attach something to a message.
//
// Use the constructor functions (NewPhotoAttachment, NewVideoAttachment, etc.)
// for type-safe attachment creation. The Payload field accepts specific payload
// types depending on Type — see each constructor for details.
type AttachmentRequest struct {
	Type    string      `json:"type"`
	Payload any `json:"payload,omitempty"`

	// Fields for specific attachment types (set based on Type):

	// Location
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`

	// Reply keyboard
	Buttons      [][]ReplyButton `json:"buttons,omitempty"`
	Direct       bool            `json:"direct,omitempty"`
	DirectUserID *int64          `json:"direct_user_id,omitempty"`
}

// NewPhotoAttachment creates an image attachment request.
func NewPhotoAttachment(payload PhotoAttachmentRequestPayload) AttachmentRequest {
	return AttachmentRequest{Type: "image", Payload: payload}
}

// NewVideoAttachment creates a video attachment request from an upload token.
func NewVideoAttachment(payload UploadedInfo) AttachmentRequest {
	return AttachmentRequest{Type: "video", Payload: payload}
}

// NewAudioAttachment creates an audio attachment request from an upload token.
func NewAudioAttachment(payload UploadedInfo) AttachmentRequest {
	return AttachmentRequest{Type: "audio", Payload: payload}
}

// NewFileAttachment creates a file attachment request from an upload token.
func NewFileAttachment(payload UploadedInfo) AttachmentRequest {
	return AttachmentRequest{Type: "file", Payload: payload}
}

// NewStickerAttachment creates a sticker attachment request.
func NewStickerAttachment(payload StickerAttachmentRequestPayload) AttachmentRequest {
	return AttachmentRequest{Type: "sticker", Payload: payload}
}

// NewContactAttachment creates a contact attachment request.
func NewContactAttachment(payload ContactAttachmentRequestPayload) AttachmentRequest {
	return AttachmentRequest{Type: "contact", Payload: payload}
}

// NewShareAttachment creates a share (link) attachment request.
func NewShareAttachment(payload ShareAttachmentPayload) AttachmentRequest {
	return AttachmentRequest{Type: "share", Payload: payload}
}

// NewInlineKeyboardAttachment creates an inline keyboard attachment request.
func NewInlineKeyboardAttachment(buttons [][]Button) AttachmentRequest {
	return AttachmentRequest{Type: "inline_keyboard", Payload: Keyboard{Buttons: buttons}}
}

// NewLocationAttachment creates a location attachment request.
func NewLocationAttachment(latitude, longitude float64) AttachmentRequest {
	return AttachmentRequest{Type: "location", Latitude: latitude, Longitude: longitude}
}

// PhotoAttachmentRequestPayload is the payload for attaching an image.
// Fields are mutually exclusive.
type PhotoAttachmentRequestPayload struct {
	URL    *string               `json:"url,omitempty"`
	Token  *string               `json:"token,omitempty"`
	Photos map[string]PhotoToken `json:"photos,omitempty"`
}

// PhotoToken represents an uploaded image token.
type PhotoToken struct {
	Token string `json:"token"`
}

// PhotoTokens is the response after uploading images.
type PhotoTokens struct {
	Photos map[string]PhotoToken `json:"photos"`
}

// UploadedInfo contains a token received after uploading media.
type UploadedInfo struct {
	Token string `json:"token"`
}

// ContactAttachmentRequestPayload is the payload for attaching a contact.
type ContactAttachmentRequestPayload struct {
	Name      *string `json:"name"`
	ContactID *int64  `json:"contact_id,omitempty"`
	VCFInfo   *string `json:"vcf_info,omitempty"`
	VCFPhone  *string `json:"vcf_phone,omitempty"`
}

// StickerAttachmentRequestPayload is the payload for attaching a sticker.
type StickerAttachmentRequestPayload struct {
	Code string `json:"code"`
}

// MarkupElement represents a text formatting element in a message.
type MarkupElement struct {
	Type     string  `json:"type"`
	From     int     `json:"from"`
	Length   int     `json:"length"`
	URL      string  `json:"url,omitempty"`
	UserLink *string `json:"user_link,omitempty"`
	UserID   *int64  `json:"user_id,omitempty"`
}

// Callback represents the data received when a user presses an inline button.
type Callback struct {
	Timestamp  int64  `json:"timestamp"`
	CallbackID string `json:"callback_id"`
	Payload    string `json:"payload,omitempty"`
	User       User   `json:"user"`
}

// CallbackAnswer represents the response to a callback.
type CallbackAnswer struct {
	Message      *NewMessageBody `json:"message,omitempty"`
	Notification *string         `json:"notification,omitempty"`
}

// Subscription represents a WebHook subscription.
type Subscription struct {
	URL         string   `json:"url"`
	Time        int64    `json:"time"`
	UpdateTypes []string `json:"update_types"`
	Version     *string  `json:"version"`
}

// getSubscriptionsResult is the response from GET /subscriptions.
type getSubscriptionsResult struct {
	Subscriptions []Subscription `json:"subscriptions"`
}

// SubscriptionRequestBody represents the request body for POST /subscriptions.
type SubscriptionRequestBody struct {
	URL         string   `json:"url"`
	Secret      string   `json:"secret,omitempty"`
	UpdateTypes []string `json:"update_types,omitempty"`
	Version     string   `json:"version,omitempty"`
}

// SimpleQueryResult represents a simple success/failure response.
type SimpleQueryResult struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// UploadEndpoint is returned by POST /uploads.
type UploadEndpoint struct {
	URL   string  `json:"url"`
	Token *string `json:"token,omitempty"`
}

// UserIDsList is the request body for POST /chats/{chatId}/members.
type UserIDsList struct {
	UserIDs []int64 `json:"user_ids"`
}

// ActionRequestBody is the request body for POST /chats/{chatId}/actions.
type ActionRequestBody struct {
	Action SenderAction `json:"action"`
}

// PinMessageBody is the request body for PUT /chats/{chatId}/pin.
type PinMessageBody struct {
	MessageID string `json:"message_id"`
	Notify    *bool  `json:"notify,omitempty"`
}

// GetPinnedMessageResult is the response from GET /chats/{chatId}/pin.
type GetPinnedMessageResult struct {
	Message *Message `json:"message"`
}

// VideoURLs contains video playback URLs in various resolutions.
type VideoURLs struct {
	MP41080 *string `json:"mp4_1080,omitempty"`
	MP4720  *string `json:"mp4_720,omitempty"`
	MP4480  *string `json:"mp4_480,omitempty"`
	MP4360  *string `json:"mp4_360,omitempty"`
	MP4240  *string `json:"mp4_240,omitempty"`
	MP4144  *string `json:"mp4_144,omitempty"`
	HLS      *string `json:"hls,omitempty"`
}

// VideoAttachmentDetails is the response from GET /videos/{videoToken}.
type VideoAttachmentDetails struct {
	Token     string                  `json:"token"`
	URLs      *VideoURLs              `json:"urls,omitempty"`
	Thumbnail *PhotoAttachmentPayload `json:"thumbnail,omitempty"`
	Width     int                     `json:"width"`
	Height    int                     `json:"height"`
	Duration  int                     `json:"duration"`
}

// Update types

// Update is the base for all update events.
type Update struct {
	UpdateType UpdateType `json:"update_type"`
	Timestamp  int64      `json:"timestamp"`
}

// MessageCreatedUpdate is received when a new message is created.
type MessageCreatedUpdate struct {
	Update
	Message    Message `json:"message"`
	UserLocale *string `json:"user_locale,omitempty"`
}

// MessageCallbackUpdate is received when a user presses an inline button.
type MessageCallbackUpdate struct {
	Update
	Callback   Callback `json:"callback"`
	Message    *Message `json:"message"`
	UserLocale *string  `json:"user_locale,omitempty"`
}

// MessageEditedUpdate is received when a message is edited.
type MessageEditedUpdate struct {
	Update
	Message Message `json:"message"`
}

// MessageRemovedUpdate is received when a message is removed.
type MessageRemovedUpdate struct {
	Update
	MessageID string `json:"message_id"`
	ChatID    int64  `json:"chat_id"`
	UserID    int64  `json:"user_id"`
}

// BotStartedUpdate is received when a user presses the Start button.
type BotStartedUpdate struct {
	Update
	ChatID     int64   `json:"chat_id"`
	User       User    `json:"user"`
	Payload    *string `json:"payload,omitempty"`
	UserLocale *string `json:"user_locale,omitempty"`
}

// BotAddedUpdate is received when the bot is added to a chat.
type BotAddedUpdate struct {
	Update
	ChatID    int64 `json:"chat_id"`
	User      User  `json:"user"`
	IsChannel bool  `json:"is_channel"`
}

// BotRemovedUpdate is received when the bot is removed from a chat.
type BotRemovedUpdate struct {
	Update
	ChatID    int64 `json:"chat_id"`
	User      User  `json:"user"`
	IsChannel bool  `json:"is_channel"`
}

// UserAddedUpdate is received when a user is added to a chat.
type UserAddedUpdate struct {
	Update
	ChatID    int64  `json:"chat_id"`
	User      User   `json:"user"`
	InviterID *int64 `json:"inviter_id,omitempty"`
	IsChannel bool   `json:"is_channel"`
}

// UserRemovedUpdate is received when a user is removed from a chat.
type UserRemovedUpdate struct {
	Update
	ChatID    int64  `json:"chat_id"`
	User      User   `json:"user"`
	AdminID   *int64 `json:"admin_id,omitempty"`
	IsChannel bool   `json:"is_channel"`
}

// ChatTitleChangedUpdate is received when a chat title is changed.
type ChatTitleChangedUpdate struct {
	Update
	ChatID int64  `json:"chat_id"`
	User   User   `json:"user"`
	Title  string `json:"title"`
}

// MessageChatCreatedUpdate is received when a chat is created via a chat button.
type MessageChatCreatedUpdate struct {
	Update
	Chat         Chat    `json:"chat"`
	MessageID    string  `json:"message_id"`
	StartPayload *string `json:"start_payload,omitempty"`
}

// UpdateList is the response from GET /updates.
type UpdateList struct {
	Updates []json.RawMessage `json:"updates"`
	Marker  *int64            `json:"marker"`
}

// apiErrorResponse is the error response from the API (internal use).
type apiErrorResponse struct {
	Error   string `json:"error,omitempty"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
