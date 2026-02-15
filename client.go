package maxigo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://botapi.max.ru"
	defaultTimeout = 30 * time.Second
)

// Client is an HTTP client for the Max Bot API.
// Create one with [New] and configure it with [Option] functions.
//
// All methods are safe for concurrent use.
type Client struct {
	httpClient *http.Client
	baseURL    string
	token      string
	timeout    time.Duration
}

// New creates a new Max Bot API client with the given token.
//
// Use functional options to customize the client:
//
//	client, err := maxigo.New("your-token",
//	    maxigo.WithTimeout(10 * time.Second),
//	)
func New(token string, opts ...Option) (*Client, error) {
	if token == "" {
		return nil, ErrEmptyToken
	}

	c := &Client{
		httpClient: &http.Client{},
		baseURL:    defaultBaseURL,
		token:      token,
		timeout:    defaultTimeout,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// do performs an HTTP request and decodes the JSON response into result.
// If result is nil, the response body is discarded.
// If the context has no deadline, the client's default timeout is applied.
func (c *Client) do(ctx context.Context, op, method, path string, query url.Values, body any, result any) error {
	if _, ok := ctx.Deadline(); !ok && c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	u, err := c.buildURL(path, query)
	if err != nil {
		return networkError(op, err)
	}

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return decodeError(op, fmt.Errorf("marshal request: %w", err))
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return networkError(op, fmt.Errorf("create request: %w", err))
	}

	req.Header.Set("Authorization", c.token)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return timeoutError(op, ctx.Err())
		}
		if isTimeout(err) {
			return timeoutError(op, err)
		}
		return networkError(op, err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return networkError(op, fmt.Errorf("read response: %w", err))
	}

	if resp.StatusCode != http.StatusOK {
		return c.parseAPIError(op, resp.StatusCode, respBody)
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return decodeError(op, fmt.Errorf("unmarshal response: %w", err))
		}
	}

	return nil
}

// doUpload performs a multipart file upload to the given URL.
func (c *Client) doUpload(ctx context.Context, op, uploadURL, filename string, reader io.Reader) ([]byte, error) {
	if _, ok := ctx.Deadline(); !ok && c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("data", filename)
	if err != nil {
		return nil, networkError(op, fmt.Errorf("create form file: %w", err))
	}
	if _, err := io.Copy(part, reader); err != nil {
		return nil, networkError(op, fmt.Errorf("copy file data: %w", err))
	}
	if err := writer.Close(); err != nil {
		return nil, networkError(op, fmt.Errorf("close multipart writer: %w", err))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadURL, &buf)
	if err != nil {
		return nil, networkError(op, fmt.Errorf("create upload request: %w", err))
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, timeoutError(op, ctx.Err())
		}
		return nil, networkError(op, err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, networkError(op, fmt.Errorf("read upload response: %w", err))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(op, resp.StatusCode, string(body))
	}

	return body, nil
}

func (c *Client) buildURL(path string, query url.Values) (string, error) {
	base := strings.TrimRight(c.baseURL, "/")
	u, err := url.Parse(base + path)
	if err != nil {
		return "", fmt.Errorf("parse URL: %w", err)
	}

	if query != nil {
		u.RawQuery = query.Encode()
	}

	return u.String(), nil
}

func (c *Client) parseAPIError(op string, statusCode int, body []byte) *Error {
	var apiResp apiErrorResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return apiError(op, statusCode, http.StatusText(statusCode))
	}

	msg := apiResp.Message
	if msg == "" {
		msg = apiResp.Error
	}
	if msg == "" {
		msg = http.StatusText(statusCode)
	}

	return apiError(op, statusCode, msg)
}

func isTimeout(err error) bool {
	var t interface{ Timeout() bool }
	if errors.As(err, &t) {
		return t.Timeout()
	}
	return false
}
