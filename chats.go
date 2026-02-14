package maxigo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GetChat returns info about a chat by its ID.
// Corresponds to GET /chats/{chatId}.
func (c *Client) GetChat(ctx context.Context, chatID int64) (*Chat, error) {
	var result Chat
	path := fmt.Sprintf("/chats/%d", chatID)
	if err := c.do(ctx, "GetChat", http.MethodGet, path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetChatByLink returns chat/channel information by its public link or username.
// Corresponds to GET /chats/{chatLink}.
func (c *Client) GetChatByLink(ctx context.Context, chatLink string) (*Chat, error) {
	var result Chat
	path := fmt.Sprintf("/chats/%s", url.PathEscape(chatLink))
	if err := c.do(ctx, "GetChatByLink", http.MethodGet, path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetChatsOpts holds optional parameters for [Client.GetChats].
type GetChatsOpts struct {
	// Count limits the number of returned chats. 0 uses the server default (50).
	Count int
	// Marker is the pagination cursor. 0 for the first page.
	Marker int64
}

// GetChats returns a paginated list of chats the bot participates in.
// Corresponds to GET /chats.
func (c *Client) GetChats(ctx context.Context, opts GetChatsOpts) (*ChatList, error) {
	q := make(url.Values)
	if opts.Count > 0 {
		q.Set("count", strconv.Itoa(opts.Count))
	}
	if opts.Marker > 0 {
		q.Set("marker", strconv.FormatInt(opts.Marker, 10))
	}

	var result ChatList
	if err := c.do(ctx, "GetChats", http.MethodGet, "/chats", q, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// EditChat edits a chat's info (title, icon, pin).
// Corresponds to PATCH /chats/{chatId}.
func (c *Client) EditChat(ctx context.Context, chatID int64, patch *ChatPatch) (*Chat, error) {
	var result Chat
	path := fmt.Sprintf("/chats/%d", chatID)
	if err := c.do(ctx, "EditChat", http.MethodPatch, path, nil, patch, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteChat deletes a chat for all participants.
// Corresponds to DELETE /chats/{chatId}.
func (c *Client) DeleteChat(ctx context.Context, chatID int64) (*SimpleQueryResult, error) {
	var result SimpleQueryResult
	path := fmt.Sprintf("/chats/%d", chatID)
	if err := c.do(ctx, "DeleteChat", http.MethodDelete, path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMembersOpts holds optional parameters for [Client.GetMembers].
type GetMembersOpts struct {
	// Count limits the number of returned members. 0 uses the server default (20).
	Count int
	// Marker is the pagination cursor. 0 for the first page.
	Marker int64
	// UserIDs filters members by specific user IDs.
	UserIDs []int64
}

// GetMembers returns a paginated list of chat members.
// Corresponds to GET /chats/{chatId}/members.
func (c *Client) GetMembers(ctx context.Context, chatID int64, opts GetMembersOpts) (*ChatMembersList, error) {
	q := make(url.Values)
	if opts.Count > 0 {
		q.Set("count", strconv.Itoa(opts.Count))
	}
	if opts.Marker > 0 {
		q.Set("marker", strconv.FormatInt(opts.Marker, 10))
	}
	if len(opts.UserIDs) > 0 {
		ids := make([]string, len(opts.UserIDs))
		for i, id := range opts.UserIDs {
			ids[i] = strconv.FormatInt(id, 10)
		}
		q.Set("user_ids", strings.Join(ids, ","))
	}

	var result ChatMembersList
	path := fmt.Sprintf("/chats/%d/members", chatID)
	if err := c.do(ctx, "GetMembers", http.MethodGet, path, q, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAdmins returns all administrators of a chat.
// The bot must be an administrator in the chat.
// Corresponds to GET /chats/{chatId}/members/admins.
func (c *Client) GetAdmins(ctx context.Context, chatID int64) (*ChatAdminsList, error) {
	var result ChatAdminsList
	path := fmt.Sprintf("/chats/%d/members/admins", chatID)
	if err := c.do(ctx, "GetAdmins", http.MethodGet, path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AddMembers adds members to a chat.
// Corresponds to POST /chats/{chatId}/members.
func (c *Client) AddMembers(ctx context.Context, chatID int64, userIDs []int64) (*SimpleQueryResult, error) {
	var result SimpleQueryResult
	path := fmt.Sprintf("/chats/%d/members", chatID)
	body := UserIDsList{UserIDs: userIDs}
	if err := c.do(ctx, "AddMembers", http.MethodPost, path, nil, body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RemoveMember removes a member from a chat.
// Set block=true to also block the user (only for chats with public/private links).
// Corresponds to DELETE /chats/{chatId}/members.
func (c *Client) RemoveMember(ctx context.Context, chatID int64, userID int64, block bool) (*SimpleQueryResult, error) {
	q := make(url.Values)
	q.Set("user_id", strconv.FormatInt(userID, 10))
	if block {
		q.Set("block", "true")
	}

	var result SimpleQueryResult
	path := fmt.Sprintf("/chats/%d/members", chatID)
	if err := c.do(ctx, "RemoveMember", http.MethodDelete, path, q, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetAdmins sets chat administrators with specified permissions.
// Corresponds to POST /chats/{chatId}/members/admins.
func (c *Client) SetAdmins(ctx context.Context, chatID int64, admins *ChatAdminsList) (*SimpleQueryResult, error) {
	var result SimpleQueryResult
	path := fmt.Sprintf("/chats/%d/members/admins", chatID)
	if err := c.do(ctx, "SetAdmins", http.MethodPost, path, nil, admins, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RemoveAdmin revokes admin rights from a user in the chat.
// Corresponds to DELETE /chats/{chatId}/members/admins/{userId}.
func (c *Client) RemoveAdmin(ctx context.Context, chatID int64, userID int64) (*SimpleQueryResult, error) {
	var result SimpleQueryResult
	path := fmt.Sprintf("/chats/%d/members/admins/%d", chatID, userID)
	if err := c.do(ctx, "RemoveAdmin", http.MethodDelete, path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SendAction sends a bot action (e.g. typing) to a chat.
// Corresponds to POST /chats/{chatId}/actions.
func (c *Client) SendAction(ctx context.Context, chatID int64, action SenderAction) (*SimpleQueryResult, error) {
	var result SimpleQueryResult
	path := fmt.Sprintf("/chats/%d/actions", chatID)
	body := ActionRequestBody{Action: action}
	if err := c.do(ctx, "SendAction", http.MethodPost, path, nil, body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMembership returns the bot's membership info for a chat.
// Corresponds to GET /chats/{chatId}/members/me.
func (c *Client) GetMembership(ctx context.Context, chatID int64) (*ChatMember, error) {
	var result ChatMember
	path := fmt.Sprintf("/chats/%d/members/me", chatID)
	if err := c.do(ctx, "GetMembership", http.MethodGet, path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// LeaveChat removes the bot from chat members.
// Corresponds to DELETE /chats/{chatId}/members/me.
func (c *Client) LeaveChat(ctx context.Context, chatID int64) (*SimpleQueryResult, error) {
	var result SimpleQueryResult
	path := fmt.Sprintf("/chats/%d/members/me", chatID)
	if err := c.do(ctx, "LeaveChat", http.MethodDelete, path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PinMessage pins a message in a chat or channel.
// Corresponds to PUT /chats/{chatId}/pin.
func (c *Client) PinMessage(ctx context.Context, chatID int64, body *PinMessageBody) (*SimpleQueryResult, error) {
	var result SimpleQueryResult
	path := fmt.Sprintf("/chats/%d/pin", chatID)
	if err := c.do(ctx, "PinMessage", http.MethodPut, path, nil, body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UnpinMessage unpins the pinned message in a chat or channel.
// Corresponds to DELETE /chats/{chatId}/pin.
func (c *Client) UnpinMessage(ctx context.Context, chatID int64) (*SimpleQueryResult, error) {
	var result SimpleQueryResult
	path := fmt.Sprintf("/chats/%d/pin", chatID)
	if err := c.do(ctx, "UnpinMessage", http.MethodDelete, path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPinnedMessage returns the pinned message in a chat or channel.
// Corresponds to GET /chats/{chatId}/pin.
func (c *Client) GetPinnedMessage(ctx context.Context, chatID int64) (*GetPinnedMessageResult, error) {
	var result GetPinnedMessageResult
	path := fmt.Sprintf("/chats/%d/pin", chatID)
	if err := c.do(ctx, "GetPinnedMessage", http.MethodGet, path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
