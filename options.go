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

// WithTimeout sets the HTTP client timeout. Default is 30 seconds.
// Note: if used together with WithHTTPClient, apply WithTimeout after it,
// otherwise WithHTTPClient will replace the client and discard the timeout.
func WithTimeout(d time.Duration) Option {
	return func(cl *Client) {
		cl.httpClient.Timeout = d
	}
}
