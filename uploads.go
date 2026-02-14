package maxigo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
