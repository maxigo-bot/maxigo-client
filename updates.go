package maxigo

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GetUpdatesOpts holds optional parameters for [Client.GetUpdates].
type GetUpdatesOpts struct {
	// Limit sets the maximum number of updates to return. 0 uses the server default (100).
	Limit int
	// Timeout sets the long-polling timeout in seconds. 0 uses the server default (30).
	Timeout int
	// Marker is the pagination cursor. 0 returns uncommitted updates.
	Marker int64
	// Types filters updates by type (e.g. "message_created", "message_callback").
	Types []string
}

// GetUpdates fetches updates using long polling.
// Corresponds to GET /updates.
func (c *Client) GetUpdates(ctx context.Context, opts GetUpdatesOpts) (*UpdateList, error) {
	q := make(url.Values)
	if opts.Limit > 0 {
		q.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Timeout > 0 {
		q.Set("timeout", strconv.Itoa(opts.Timeout))
	}
	if opts.Marker > 0 {
		q.Set("marker", strconv.FormatInt(opts.Marker, 10))
	}
	if len(opts.Types) > 0 {
		q.Set("types", strings.Join(opts.Types, ","))
	}

	var result UpdateList
	if err := c.do(ctx, "GetUpdates", http.MethodGet, "/updates", q, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
