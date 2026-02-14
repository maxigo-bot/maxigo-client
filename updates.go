package maxigo

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	// defaultPollingTimeout is the server-side default long-polling timeout (seconds).
	defaultPollingTimeout = 30
	// pollingBuffer is extra time added on top of the server-side polling timeout
	// to prevent the HTTP client from timing out before the server responds.
	pollingBuffer = 5 * time.Second
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
//
// The client automatically adjusts the HTTP timeout to accommodate the
// server-side long-polling duration, preventing spurious timeout errors.
func (c *Client) GetUpdates(ctx context.Context, opts GetUpdatesOpts) (*UpdateList, error) {
	q := make(url.Values)
	if opts.Limit > 0 {
		q.Set("limit", strconv.Itoa(opts.Limit))
	}

	serverTimeout := opts.Timeout
	if serverTimeout <= 0 {
		serverTimeout = defaultPollingTimeout
	}
	q.Set("timeout", strconv.Itoa(serverTimeout))

	if opts.Marker > 0 {
		q.Set("marker", strconv.FormatInt(opts.Marker, 10))
	}
	if len(opts.Types) > 0 {
		q.Set("types", strings.Join(opts.Types, ","))
	}

	pollingDuration := time.Duration(serverTimeout)*time.Second + pollingBuffer

	if deadline, ok := ctx.Deadline(); ok {
		if time.Until(deadline) < pollingDuration {
			return nil, ErrPollDeadline
		}
	} else {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, pollingDuration)
		defer cancel()
	}

	var result UpdateList
	if err := c.do(ctx, "GetUpdates", http.MethodGet, "/updates", q, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
