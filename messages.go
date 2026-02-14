package maxigo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// SendMessage sends a message to a chat.
// Set body.DisableLinkPreview to true to prevent the server from generating link previews.
// Corresponds to POST /messages.
func (c *Client) SendMessage(ctx context.Context, chatID int64, body *NewMessageBody) (*Message, error) {
	q := make(url.Values)
	if chatID != 0 {
		q.Set("chat_id", strconv.FormatInt(chatID, 10))
	}
	if body != nil && body.DisableLinkPreview {
		q.Set("disable_link_preview", "true")
	}

	var result sendMessageResult
	if err := c.do(ctx, "SendMessage", http.MethodPost, "/messages", q, body, &result); err != nil {
		return nil, err
	}
	return &result.Message, nil
}

// SendMessageToUser sends a message directly to a user.
// Set body.DisableLinkPreview to true to prevent the server from generating link previews.
// Corresponds to POST /messages with user_id query parameter.
func (c *Client) SendMessageToUser(ctx context.Context, userID int64, body *NewMessageBody) (*Message, error) {
	q := make(url.Values)
	q.Set("user_id", strconv.FormatInt(userID, 10))
	if body != nil && body.DisableLinkPreview {
		q.Set("disable_link_preview", "true")
	}

	var result sendMessageResult
	if err := c.do(ctx, "SendMessageToUser", http.MethodPost, "/messages", q, body, &result); err != nil {
		return nil, err
	}
	return &result.Message, nil
}

// EditMessage edits an existing message.
// Corresponds to PUT /messages.
func (c *Client) EditMessage(ctx context.Context, messageID string, body *NewMessageBody) (*SimpleQueryResult, error) {
	q := make(url.Values)
	q.Set("message_id", messageID)

	var result SimpleQueryResult
	if err := c.do(ctx, "EditMessage", http.MethodPut, "/messages", q, body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteMessage deletes a message.
// Corresponds to DELETE /messages.
func (c *Client) DeleteMessage(ctx context.Context, messageID string) (*SimpleQueryResult, error) {
	q := make(url.Values)
	q.Set("message_id", messageID)

	var result SimpleQueryResult
	if err := c.do(ctx, "DeleteMessage", http.MethodDelete, "/messages", q, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMessagesOpts holds optional parameters for [Client.GetMessages].
type GetMessagesOpts struct {
	// ChatID filters messages by chat. Required unless MessageIDs is set.
	ChatID int64
	// Count limits the number of returned messages. 0 uses the server default.
	Count int
	// MessageIDs requests specific messages by their IDs.
	MessageIDs []string
	// From is the start timestamp (inclusive). Messages are returned in reverse
	// chronological order, so From should be greater than To.
	From int64
	// To is the end timestamp (inclusive).
	To int64
}

// GetMessages returns messages from a chat.
// Corresponds to GET /messages.
func (c *Client) GetMessages(ctx context.Context, opts GetMessagesOpts) (*MessageList, error) {
	q := make(url.Values)
	if opts.ChatID != 0 {
		q.Set("chat_id", strconv.FormatInt(opts.ChatID, 10))
	}
	if opts.Count > 0 {
		q.Set("count", strconv.Itoa(opts.Count))
	}
	if len(opts.MessageIDs) > 0 {
		q.Set("message_ids", strings.Join(opts.MessageIDs, ","))
	}
	if opts.From != 0 {
		q.Set("from", strconv.FormatInt(opts.From, 10))
	}
	if opts.To != 0 {
		q.Set("to", strconv.FormatInt(opts.To, 10))
	}

	var result MessageList
	if err := c.do(ctx, "GetMessages", http.MethodGet, "/messages", q, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMessageByID returns a single message by its identifier.
// Corresponds to GET /messages/{messageId}.
func (c *Client) GetMessageByID(ctx context.Context, messageID string) (*Message, error) {
	var result Message
	path := fmt.Sprintf("/messages/%s", url.PathEscape(messageID))
	if err := c.do(ctx, "GetMessageByID", http.MethodGet, path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AnswerCallback sends a response to a callback button press.
// Corresponds to POST /answers.
func (c *Client) AnswerCallback(ctx context.Context, callbackID string, answer *CallbackAnswer) (*SimpleQueryResult, error) {
	q := make(url.Values)
	q.Set("callback_id", callbackID)

	var result SimpleQueryResult
	if err := c.do(ctx, "AnswerCallback", http.MethodPost, "/answers", q, answer, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
