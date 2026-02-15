package maxigo

import (
	"encoding/json"
	"fmt"
)

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
	UpdateBotStopped         UpdateType = "bot_stopped"
	UpdateBotAdded           UpdateType = "bot_added"
	UpdateBotRemoved         UpdateType = "bot_removed"
	UpdateUserAdded          UpdateType = "user_added"
	UpdateUserRemoved        UpdateType = "user_removed"
	UpdateChatTitleChanged   UpdateType = "chat_title_changed"
	UpdateMessageChatCreated UpdateType = "message_chat_created"
	UpdateDialogMuted        UpdateType = "dialog_muted"
	UpdateDialogUnmuted      UpdateType = "dialog_unmuted"
	UpdateDialogCleared      UpdateType = "dialog_cleared"
	UpdateDialogRemoved      UpdateType = "dialog_removed"
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

// User represents a Max user or bot.
type User struct {
	// Unique identifier of the user or bot.
	UserID int64 `json:"user_id"`
	// Display name of the user or bot.
	FirstName string `json:"first_name"`
	// Display last name. Not returned for bots.
	LastName *string `json:"last_name,omitempty"`
	// Bot username or unique public name. May be null for users.
	Username *string `json:"username,omitempty"`
	// True if this is a bot.
	IsBot bool `json:"is_bot"`
	// Last activity time in MAX (Unix time in milliseconds).
	// May be absent if the user disabled online status in settings.
	LastActivityTime int64 `json:"last_activity_time"`
}

// UserWithPhoto extends User with avatar and description.
type UserWithPhoto struct {
	User
	// User or bot description (up to 16000 characters).
	Description *string `json:"description,omitempty"`
	// Small avatar URL.
	AvatarURL string `json:"avatar_url,omitempty"`
	// Full-size avatar URL.
	FullAvatarURL string `json:"full_avatar_url,omitempty"`
}

// BotInfo represents the bot's info returned by GET /me.
type BotInfo struct {
	UserWithPhoto
	// Commands supported by the bot (up to 32).
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
	Name        OptString    `json:"name,omitzero"`
	FirstName   OptString    `json:"first_name,omitzero"`
	Description OptString    `json:"description,omitzero"`
	Commands    []BotCommand `json:"commands,omitempty"`
	Photo       *PhotoAttachmentRequestPayload `json:"photo,omitempty"`
}

// Image represents a generic image object.
type Image struct {
	URL string `json:"url"`
}

// Chat represents a Max chat.
type Chat struct {
	// Chat identifier.
	ChatID int64 `json:"chat_id"`
	// Chat type: "chat" (group), "dialog" (direct), or "channel".
	Type ChatType `json:"type"`
	// Bot's status in the chat: "active", "removed", "left", "closed".
	Status ChatStatus `json:"status"`
	// Display title. May be null for dialogs.
	Title *string `json:"title"`
	// Chat icon.
	Icon *Image `json:"icon"`
	// Last event time in the chat (Unix time).
	LastEventTime int64 `json:"last_event_time"`
	// Number of participants. Always 2 for dialogs.
	ParticipantsCount int `json:"participants_count"`
	// Chat owner ID.
	OwnerID *int64 `json:"owner_id,omitempty"`
	// Participants with last activity time. May be null for chat lists.
	Participants map[string]int64 `json:"participants,omitempty"`
	// Whether the chat is publicly accessible (always false for dialogs).
	IsPublic bool `json:"is_public"`
	// Chat invite link.
	Link *string `json:"link,omitempty"`
	// Chat description.
	Description *string `json:"description"`
	// User info for dialog chats (type "dialog" only).
	DialogWithUser *UserWithPhoto `json:"dialog_with_user,omitempty"`
	// Message count. Only for group chats and channels, not dialogs.
	MessagesCount *int `json:"messages_count,omitempty"`
	// ID of the message containing the button that initiated this chat.
	ChatMessageID *string `json:"chat_message_id,omitempty"`
	// Pinned message. Only returned when requesting a specific chat.
	PinnedMessage *Message `json:"pinned_message,omitempty"`
}

// ChatList represents a paginated list of chats.
type ChatList struct {
	Chats  []Chat `json:"chats"`
	Marker *int64 `json:"marker"`
}

// ChatPatch represents the request body for PATCH /chats/{chatId}.
type ChatPatch struct {
	Icon   *PhotoAttachmentRequestPayload `json:"icon,omitempty"`
	Title  OptString                      `json:"title,omitzero"`
	Pin    OptString                      `json:"pin,omitzero"`
	Notify OptBool                        `json:"notify,omitzero"`
}

// ChatMember represents a member of a chat.
type ChatMember struct {
	UserWithPhoto
	// Last activity time in the chat. May be stale for superchats.
	LastAccessTime int64 `json:"last_access_time"`
	// Whether the user is the chat owner.
	IsOwner bool `json:"is_owner"`
	// Whether the user is a chat administrator.
	IsAdmin bool `json:"is_admin"`
	// Time when the user joined the chat (Unix time).
	JoinTime int64 `json:"join_time"`
	// Admin permissions. Null if the member is not an admin.
	Permissions []ChatAdminPermission `json:"permissions"`
	// Custom admin title shown in chat.
	Alias *string `json:"alias"`
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
	// User who sent the message.
	Sender *User `json:"sender,omitempty"`
	// Recipient — can be a user or a chat.
	Recipient Recipient `json:"recipient"`
	// Message creation time (Unix time).
	Timestamp int64 `json:"timestamp"`
	// Forwarded or replied-to message.
	Link *LinkedMessage `json:"link,omitempty"`
	// Message content: text and attachments.
	Body MessageBody `json:"body"`
	// Message statistics.
	Stat *MessageStat `json:"stat,omitempty"`
	// Public link to a channel post. Absent for dialogs and group chats.
	URL *string `json:"url,omitempty"`
}

// MessageBody represents the body of a message.
type MessageBody struct {
	MID         string            `json:"mid"`
	Seq         int64             `json:"seq"`
	Text        *string           `json:"text"`
	Attachments []json.RawMessage `json:"attachments"`
	Markup      []MarkupElement   `json:"markup,omitempty"`
}

// attachmentFactories maps JSON "type" values to factory functions
// that return a pointer to the corresponding Go struct.
var attachmentFactories = map[string]func() Attachment{
	"image":          func() Attachment { return new(PhotoAttachment) },
	"video":          func() Attachment { return new(VideoAttachment) },
	"audio":          func() Attachment { return new(AudioAttachment) },
	"file":           func() Attachment { return new(FileAttachment) },
	"sticker":        func() Attachment { return new(StickerAttachment) },
	"contact":        func() Attachment { return new(ContactAttachment) },
	"share":          func() Attachment { return new(ShareAttachment) },
	"location":       func() Attachment { return new(LocationAttachment) },
	"data":           func() Attachment { return new(DataAttachment) },
	"inline_keyboard": func() Attachment { return new(InlineKeyboardAttachment) },
	"reply_keyboard": func() Attachment { return new(ReplyKeyboardAttachment) },
}

// ParseAttachments unmarshals raw JSON attachments into typed structs.
// Each returned element is a pointer to one of the attachment types
// (e.g. *PhotoAttachment, *ContactAttachment). Use a type switch to
// inspect individual attachments.
//
// Unknown attachment types are silently skipped for forward compatibility.
// Returns nil, nil when there are no attachments.
func (mb *MessageBody) ParseAttachments() ([]Attachment, error) {
	if len(mb.Attachments) == 0 {
		return nil, nil
	}

	var header struct {
		Type string `json:"type"`
	}

	result := make([]Attachment, 0, len(mb.Attachments))
	for _, raw := range mb.Attachments {
		if err := json.Unmarshal(raw, &header); err != nil {
			return nil, fmt.Errorf("parse attachment type: %w", err)
		}

		factory, ok := attachmentFactories[header.Type]
		if !ok {
			continue // unknown type — skip for forward compat
		}

		att := factory()
		if err := json.Unmarshal(raw, att); err != nil {
			return nil, fmt.Errorf("parse %s attachment: %w", header.Type, err)
		}
		result = append(result, att)
	}

	return result, nil
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
// Use [Some] to set optional fields. Unset fields are omitted from JSON,
// which tells the server to keep the existing value when editing a message.
type NewMessageBody struct {
	// Message text (up to 4000 characters).
	Text OptString `json:"text,omitzero"`
	// Message attachments. If empty, all existing attachments will be removed.
	Attachments []AttachmentRequest `json:"attachments,omitzero"`
	// Link to another message (for reply or forward).
	Link *NewMessageLink `json:"link,omitempty"`
	// If false, chat members will not be notified (default true).
	Notify OptBool `json:"notify,omitzero"`
	// Text formatting mode: "markdown" or "html".
	Format Optional[TextFormat] `json:"format,omitzero"`

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

// Attachment is implemented by all attachment response types.
// Use [MessageBody.ParseAttachments] to convert raw JSON attachments
// into typed structs, then type-switch on the result.
type Attachment interface {
	// GetType returns the attachment type string (e.g. "image", "video", "contact").
	GetType() string
}

// AttachmentType is embedded in all attachment responses.
type AttachmentType struct {
	Type string `json:"type"`
}

// GetType implements the [Attachment] interface.
func (a AttachmentType) GetType() string {
	return a.Type
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
	URL   OptString `json:"url,omitzero"`
	Token OptString `json:"token,omitzero"`
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
// "request_contact", "request_geo_location", "chat", "message", "open_app".
type Button struct {
	Type             string    `json:"type"`
	Text             string    `json:"text"`
	Payload          string    `json:"payload,omitempty"`
	URL              string    `json:"url,omitempty"`
	Intent           Intent    `json:"intent,omitempty"`
	Quick            bool      `json:"quick,omitempty"`
	ChatTitle        string    `json:"chat_title,omitempty"`
	ChatDescription  OptString `json:"chat_description,omitzero"`
	StartPayload     OptString `json:"start_payload,omitzero"`
	UUID             OptInt64  `json:"uuid,omitzero"`
	WebApp           string    `json:"web_app,omitempty"`
}

// NewCallbackButton creates a callback button that sends payload to the bot.
func NewCallbackButton(text, payload string) Button {
	return Button{Type: "callback", Text: text, Payload: payload}
}

// NewCallbackButtonWithIntent creates a callback button with a visual intent.
func NewCallbackButtonWithIntent(text, payload string, intent Intent) Button {
	return Button{Type: "callback", Text: text, Payload: payload, Intent: intent}
}

// NewLinkButton creates a button that opens a URL when pressed.
func NewLinkButton(text, url string) Button {
	return Button{Type: "link", Text: text, URL: url}
}

// NewRequestContactButton creates a button that requests the user's contact information.
func NewRequestContactButton(text string) Button {
	return Button{Type: "request_contact", Text: text}
}

// NewRequestGeoLocationButton creates a button that requests the user's location.
// If quick is true, the location is sent without asking user's confirmation.
func NewRequestGeoLocationButton(text string, quick bool) Button {
	return Button{Type: "request_geo_location", Text: text, Quick: quick}
}

// NewChatButton creates a button that creates a new chat when pressed.
// The bot will be added as administrator and the message author will own the chat.
func NewChatButton(text, chatTitle string) Button {
	return Button{Type: "chat", Text: text, ChatTitle: chatTitle}
}

// NewMessageButton creates a button that sends a message from the user in chat.
func NewMessageButton(text string) Button {
	return Button{Type: "message", Text: text}
}

// NewOpenAppButton creates a button that opens a mini app inside the messenger.
// The webApp parameter is the bot username whose mini app to launch.
func NewOpenAppButton(text, webApp string) Button {
	return Button{Type: "open_app", Text: text, WebApp: webApp}
}

// ReplyButton represents a button in a reply keyboard.
type ReplyButton struct {
	Type    string    `json:"type,omitempty"`
	Text    string    `json:"text"`
	Payload OptString `json:"payload,omitzero"`
	Intent  Intent    `json:"intent,omitempty"`
	Quick   bool      `json:"quick,omitempty"`
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
	DirectUserID OptInt64         `json:"direct_user_id,omitzero"`
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
	URL    OptString             `json:"url,omitzero"`
	Token  OptString             `json:"token,omitzero"`
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
	Name      OptString `json:"name,omitzero"`
	ContactID OptInt64  `json:"contact_id,omitzero"`
	VCFInfo   OptString `json:"vcf_info,omitzero"`
	VCFPhone  OptString `json:"vcf_phone,omitzero"`
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
	Notification OptString       `json:"notification,omitzero"`
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
	MessageID string  `json:"message_id"`
	Notify    OptBool `json:"notify,omitzero"`
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
	// Discriminator that determines the update type.
	UpdateType UpdateType `json:"update_type"`
	// Unix time when the event occurred.
	Timestamp int64 `json:"timestamp"`
}

// MessageCreatedUpdate is received when a new message is created.
type MessageCreatedUpdate struct {
	Update
	// The newly created message.
	Message Message `json:"message"`
	// User's current locale (IETF BCP 47). Only available in dialogs.
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

// BotStoppedUpdate is received when a user stops the bot.
type BotStoppedUpdate struct {
	Update
	ChatID int64 `json:"chat_id"`
	User   User  `json:"user"`
}

// DialogMutedUpdate is received when a user mutes the dialog with the bot.
type DialogMutedUpdate struct {
	Update
	ChatID int64 `json:"chat_id"`
	User   User  `json:"user"`
}

// DialogUnmutedUpdate is received when a user unmutes the dialog with the bot.
type DialogUnmutedUpdate struct {
	Update
	ChatID int64 `json:"chat_id"`
	User   User  `json:"user"`
}

// DialogClearedUpdate is received when a user clears the dialog history.
type DialogClearedUpdate struct {
	Update
	ChatID int64 `json:"chat_id"`
	User   User  `json:"user"`
}

// DialogRemovedUpdate is received when a user removes the dialog with the bot.
type DialogRemovedUpdate struct {
	Update
	ChatID int64 `json:"chat_id"`
	User   User  `json:"user"`
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
