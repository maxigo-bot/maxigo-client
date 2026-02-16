package maxigo

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
)

func TestGetVideoDetails(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/videos/abc-123" {
			t.Errorf("path = %q, want /videos/abc-123", r.URL.Path)
		}
		writeJSON(t, w, VideoAttachmentDetails{
			Token:    "abc-123",
			Width:    1920,
			Height:   1080,
			Duration: 120,
		})
	})

	result, err := c.GetVideoDetails(context.Background(), "abc-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Token != "abc-123" {
		t.Errorf("Token = %q, want %q", result.Token, "abc-123")
	}
	if result.Width != 1920 {
		t.Errorf("Width = %d, want 1920", result.Width)
	}
}

func TestGetUploadURL(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/uploads" {
			t.Errorf("path = %q, want /uploads", r.URL.Path)
		}
		if r.URL.Query().Get("type") != "image" {
			t.Errorf("type = %q, want image", r.URL.Query().Get("type"))
		}
		writeJSON(t, w, UploadEndpoint{
			URL: "https://upload.example.com/upload",
		})
	})

	result, err := c.GetUploadURL(context.Background(), UploadImage)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.URL != "https://upload.example.com/upload" {
		t.Errorf("URL = %q, want %q", result.URL, "https://upload.example.com/upload")
	}
}

func TestUploadPhoto(t *testing.T) {
	var requestCount atomic.Int32
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		count := requestCount.Add(1)
		if count == 1 {
			// First request: get upload URL — return the same server URL
			writeJSON(t, w, UploadEndpoint{
				URL: "http://" + r.Host + "/do-upload",
			})
			return
		}
		// Second request: the actual upload
		if r.URL.Path != "/do-upload" {
			t.Errorf("upload path = %q, want /do-upload", r.URL.Path)
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Errorf("Content-Type = %q, want multipart/form-data", r.Header.Get("Content-Type"))
		}
		writeJSON(t, w, PhotoTokens{
			Photos: map[string]PhotoToken{
				"default": {Token: "photo-token-123"},
			},
		})
	})

	result, err := c.UploadPhoto(context.Background(), "test.jpg", strings.NewReader("fake image data"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Photos["default"].Token != "photo-token-123" {
		t.Errorf("token = %q, want %q", result.Photos["default"].Token, "photo-token-123")
	}
}

func TestUploadPhotoGetURLError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeError(t, w, http.StatusForbidden, `{"code":"forbidden","message":"access denied"}`)
	})

	_, err := c.UploadPhoto(context.Background(), "test.jpg", strings.NewReader("data"))
	if err == nil {
		t.Fatal("expected error")
	}

	var e *Error
	if !errors.As(err, &e) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if e.Kind != ErrAPI {
		t.Errorf("Kind = %v, want ErrAPI", e.Kind)
	}
	if e.Op != "GetUploadURL" {
		t.Errorf("Op = %q, want GetUploadURL", e.Op)
	}
}

func TestUploadPhotoUploadError(t *testing.T) {
	var requestCount atomic.Int32
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		count := requestCount.Add(1)
		if count == 1 {
			writeJSON(t, w, UploadEndpoint{
				URL: "http://" + r.Host + "/do-upload",
			})
			return
		}
		writeError(t, w, http.StatusInternalServerError, "upload failed")
	})

	_, err := c.UploadPhoto(context.Background(), "test.jpg", strings.NewReader("data"))
	if err == nil {
		t.Fatal("expected error")
	}

	var e *Error
	if !errors.As(err, &e) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if e.Kind != ErrAPI {
		t.Errorf("Kind = %v, want ErrAPI", e.Kind)
	}
}

func TestUploadPhotoUnmarshalError(t *testing.T) {
	var requestCount atomic.Int32
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		count := requestCount.Add(1)
		if count == 1 {
			writeJSON(t, w, UploadEndpoint{
				URL: "http://" + r.Host + "/do-upload",
			})
			return
		}
		_, _ = w.Write([]byte(`{invalid json`))
	})

	_, err := c.UploadPhoto(context.Background(), "test.jpg", strings.NewReader("data"))
	if err == nil {
		t.Fatal("expected error")
	}

	var e *Error
	if !errors.As(err, &e) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if e.Kind != ErrDecode {
		t.Errorf("Kind = %v, want ErrDecode", e.Kind)
	}
	if e.Op != "UploadPhoto" {
		t.Errorf("Op = %q, want UploadPhoto", e.Op)
	}
}

func TestUploadMedia(t *testing.T) {
	var requestCount atomic.Int32
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		count := requestCount.Add(1)
		if count == 1 {
			writeJSON(t, w, UploadEndpoint{
				URL: "http://" + r.Host + "/do-upload",
			})
			return
		}
		writeJSON(t, w, UploadedInfo{Token: "video-token-456"})
	})

	result, err := c.UploadMedia(context.Background(), UploadVideo, "test.mp4", strings.NewReader("fake video data"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Token != "video-token-456" {
		t.Errorf("Token = %q, want %q", result.Token, "video-token-456")
	}
}

func TestUploadMediaGetURLError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeError(t, w, http.StatusForbidden, `{"code":"forbidden","message":"access denied"}`)
	})

	_, err := c.UploadMedia(context.Background(), UploadVideo, "test.mp4", strings.NewReader("data"))
	if err == nil {
		t.Fatal("expected error")
	}

	var e *Error
	if !errors.As(err, &e) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if e.Kind != ErrAPI {
		t.Errorf("Kind = %v, want ErrAPI", e.Kind)
	}
	if e.Op != "GetUploadURL" {
		t.Errorf("Op = %q, want GetUploadURL", e.Op)
	}
}

func TestUploadMediaUploadError(t *testing.T) {
	var requestCount atomic.Int32
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		count := requestCount.Add(1)
		if count == 1 {
			writeJSON(t, w, UploadEndpoint{
				URL: "http://" + r.Host + "/do-upload",
			})
			return
		}
		writeError(t, w, http.StatusInternalServerError, "upload failed")
	})

	_, err := c.UploadMedia(context.Background(), UploadVideo, "test.mp4", strings.NewReader("data"))
	if err == nil {
		t.Fatal("expected error")
	}

	var e *Error
	if !errors.As(err, &e) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if e.Kind != ErrAPI {
		t.Errorf("Kind = %v, want ErrAPI", e.Kind)
	}
}

func TestUploadMediaUnmarshalError(t *testing.T) {
	var requestCount atomic.Int32
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		count := requestCount.Add(1)
		if count == 1 {
			writeJSON(t, w, UploadEndpoint{
				URL: "http://" + r.Host + "/do-upload",
			})
			return
		}
		_, _ = w.Write([]byte(`{invalid json`))
	})

	_, err := c.UploadMedia(context.Background(), UploadVideo, "test.mp4", strings.NewReader("data"))
	if err == nil {
		t.Fatal("expected error")
	}

	var e *Error
	if !errors.As(err, &e) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if e.Kind != ErrDecode {
		t.Errorf("Kind = %v, want ErrDecode", e.Kind)
	}
	if e.Op != "UploadMedia" {
		t.Errorf("Op = %q, want UploadMedia", e.Op)
	}
}
