package maxigo

import (
	"net/http"
	"time"
)

// Option configures the Client.
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client. Useful for testing or custom transports.
// The client is shallow-copied to prevent mutation of the original.
func WithHTTPClient(c *http.Client) Option {
	return func(cl *Client) {
		copied := *c
		cl.httpClient = &copied
	}
}

// WithBaseURL overrides the default API base URL (https://botapi.max.ru).
// This is primarily useful for testing with httptest.Server.
func WithBaseURL(url string) Option {
	return func(cl *Client) {
		cl.baseURL = url
	}
}

// WithTimeout sets the request timeout. Default is 30 seconds.
// This timeout is applied via context to each request that has no deadline set.
// For long-polling requests ([Client.GetUpdates]), the timeout is automatically
// extended to accommodate the server-side polling.
func WithTimeout(d time.Duration) Option {
	return func(cl *Client) {
		cl.timeout = d
	}
}

// DefaultRetryIntervals are the default retry intervals used by [WithRetry]
// when no custom intervals are provided.
var DefaultRetryIntervals = []time.Duration{
	500 * time.Millisecond,
	1 * time.Second,
	2 * time.Second,
	5 * time.Second,
}

// WithRetry enables automatic retry for retryable API errors.
// Retryable errors are HTTP 429 (Too Many Requests) and API errors
// with messages containing "not.ready" or "not.processed" (attachment processing).
//
// With no arguments, [DefaultRetryIntervals] are used (500ms, 1s, 2s, 5s).
// Custom intervals can be provided:
//
//	client, err := maxigo.New("token", maxigo.WithRetry())                           // defaults
//	client, err := maxigo.New("token", maxigo.WithRetry(time.Second, 3*time.Second)) // custom
//
// Retry respects context cancellation between attempts.
func WithRetry(intervals ...time.Duration) Option {
	return func(cl *Client) {
		if len(intervals) == 0 {
			cl.retryIntervals = append([]time.Duration(nil), DefaultRetryIntervals...)
		} else {
			cl.retryIntervals = append([]time.Duration(nil), intervals...)
		}
	}
}
