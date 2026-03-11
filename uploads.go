package maxigo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

// GetVideoDetails returns detailed information about a video attachment.
// Corresponds to GET /videos/{videoToken}.
func (c *Client) GetVideoDetails(ctx context.Context, videoToken string) (*VideoAttachmentDetails, error) {
	var result VideoAttachmentDetails
	path := fmt.Sprintf("/videos/%s", url.PathEscape(videoToken))
	if err := c.do(ctx, "GetVideoDetails", http.MethodGet, path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetUploadURL returns a URL to upload a file of the given type.
// Corresponds to POST /uploads.
func (c *Client) GetUploadURL(ctx context.Context, uploadType UploadType) (*UploadEndpoint, error) {
	q := make(url.Values)
	q.Set("type", string(uploadType))

	var result UploadEndpoint
	if err := c.do(ctx, "GetUploadURL", http.MethodPost, "/uploads", q, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UploadPhoto uploads an image and returns photo tokens.
// This is a two-step operation: get upload URL, then upload the file.
func (c *Client) UploadPhoto(ctx context.Context, filename string, reader io.Reader) (*PhotoTokens, error) {
	endpoint, err := c.GetUploadURL(ctx, UploadImage)
	if err != nil {
		return nil, err
	}

	body, err := c.doUpload(ctx, "UploadPhoto", endpoint.URL, filename, reader)
	if err != nil {
		return nil, err
	}

	var result PhotoTokens
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, decodeError("UploadPhoto", fmt.Errorf("unmarshal upload response: %w", err))
	}
	return &result, nil
}

// UploadMedia uploads a video, audio, or file and returns the token.
// This is a two-step operation: get upload URL, then upload the file.
func (c *Client) UploadMedia(ctx context.Context, uploadType UploadType, filename string, reader io.Reader) (*UploadedInfo, error) {
	endpoint, err := c.GetUploadURL(ctx, uploadType)
	if err != nil {
		return nil, err
	}

	body, err := c.doUpload(ctx, "UploadMedia", endpoint.URL, filename, reader)
	if err != nil {
		return nil, err
	}

	var result UploadedInfo
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, decodeError("UploadMedia", fmt.Errorf("unmarshal upload response: %w", err))
	}
	return &result, nil
}

// UploadPhotoFromFile opens a local file and uploads it as a photo.
func (c *Client) UploadPhotoFromFile(ctx context.Context, filePath string) (*PhotoTokens, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, &Error{Kind: ErrFetch, Op: "UploadPhotoFromFile", Message: err.Error(), Err: err}
	}
	defer func() { _ = f.Close() }()

	return c.UploadPhoto(ctx, filepath.Base(filePath), f)
}

// UploadMediaFromFile opens a local file and uploads it as the given media type.
func (c *Client) UploadMediaFromFile(ctx context.Context, uploadType UploadType, filePath string) (*UploadedInfo, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, &Error{Kind: ErrFetch, Op: "UploadMediaFromFile", Message: err.Error(), Err: err}
	}
	defer func() { _ = f.Close() }()

	return c.UploadMedia(ctx, uploadType, filepath.Base(filePath), f)
}

// UploadPhotoFromURL fetches an image from a URL and uploads it as a photo.
// Only http and https schemes are allowed.
//
// Security: do not pass untrusted user input directly as imageURL
// without validation — this could allow SSRF attacks against internal networks.
func (c *Client) UploadPhotoFromURL(ctx context.Context, imageURL string) (*PhotoTokens, error) {
	ctx, cancel := c.ensureTimeout(ctx)
	defer cancel()

	body, filename, err := c.fetchURL(ctx, "UploadPhotoFromURL", imageURL)
	if err != nil {
		return nil, err
	}
	defer func() { _ = body.Close() }()

	if filename == "" {
		filename = "photo"
	}
	return c.UploadPhoto(ctx, filename, body)
}

// UploadMediaFromURL fetches a file from a URL and uploads it as the given media type.
// Only http and https schemes are allowed.
//
// Security: do not pass untrusted user input directly as fileURL
// without validation — this could allow SSRF attacks against internal networks.
func (c *Client) UploadMediaFromURL(ctx context.Context, uploadType UploadType, fileURL string) (*UploadedInfo, error) {
	ctx, cancel := c.ensureTimeout(ctx)
	defer cancel()

	body, filename, err := c.fetchURL(ctx, "UploadMediaFromURL", fileURL)
	if err != nil {
		return nil, err
	}
	defer func() { _ = body.Close() }()

	if filename == "" {
		filename = "file"
	}
	return c.UploadMedia(ctx, uploadType, filename, body)
}

// maxFetchSize is the maximum number of bytes fetchURL will read (50 MB).
const maxFetchSize = 50 << 20

// fetchURL downloads content from the URL. Caller must close the body.
// Only http and https schemes are allowed.
func (c *Client) fetchURL(ctx context.Context, op, rawURL string) (io.ReadCloser, string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, "", networkError(op, fmt.Errorf("parse URL: %w", err))
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, "", networkError(op, fmt.Errorf("unsupported URL scheme %q: only http and https are allowed", u.Scheme))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, "", networkError(op, fmt.Errorf("create request: %w", err))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, "", timeoutError(op, ctx.Err())
		}
		if isTimeout(err) {
			return nil, "", timeoutError(op, err)
		}
		return nil, "", networkError(op, err)
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, "", fetchError(op, resp.StatusCode, fmt.Sprintf("fetch %s: %s", rawURL, http.StatusText(resp.StatusCode)))
	}

	limited := io.NopCloser(io.LimitReader(resp.Body, maxFetchSize))
	body := readCloser{limited, resp.Body}
	return body, extractFilename(resp, rawURL), nil
}

// readCloser combines a limited reader with the original closer.
type readCloser struct {
	io.ReadCloser        // limited reader (reads up to maxFetchSize)
	underlying    io.Closer // original resp.Body
}

func (rc readCloser) Close() error {
	return errors.Join(rc.ReadCloser.Close(), rc.underlying.Close())
}

// extractFilename gets name from Content-Disposition header or URL path.
// Filenames are sanitized with filepath.Base to prevent path traversal.
func extractFilename(resp *http.Response, rawURL string) string {
	if cd := resp.Header.Get("Content-Disposition"); cd != "" {
		if _, params, err := mime.ParseMediaType(cd); err == nil {
			if name := filepath.Base(params["filename"]); name != "" && name != "." {
				return name
			}
		}
	}

	if u, err := url.Parse(rawURL); err == nil {
		if base := path.Base(u.Path); base != "." && base != "/" {
			return base
		}
	}

	return ""
}
