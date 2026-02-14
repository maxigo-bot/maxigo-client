package maxigo

import (
	"context"
	"net/http"
)

// GetBot returns info about the current bot.
// Corresponds to GET /me.
func (c *Client) GetBot(ctx context.Context) (*BotInfo, error) {
	var result BotInfo
	if err := c.do(ctx, "GetBot", http.MethodGet, "/me", nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// EditBot edits the current bot's info.
// Only filled fields will be updated; the rest remain unchanged.
// Corresponds to PATCH /me.
func (c *Client) EditBot(ctx context.Context, patch *BotPatch) (*BotInfo, error) {
	var result BotInfo
	if err := c.do(ctx, "EditBot", http.MethodPatch, "/me", nil, patch, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
