package maxigo

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
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

func TestUploadPhotoFromFile(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "test-image.jpg")
	if err := os.WriteFile(tmp, []byte("fake image data"), 0o644); err != nil {
		t.Fatal(err)
	}

	var requestCount atomic.Int32
	var uploadedFilename string
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		count := requestCount.Add(1)
		if count == 1 {
			writeJSON(t, w, UploadEndpoint{URL: "http://" + r.Host + "/do-upload"})
			return
		}
		_ = r.ParseMultipartForm(1 << 20)
		if f, fh, err := r.FormFile("data"); err == nil {
			uploadedFilename = fh.Filename
			_ = f.Close()
		}
		writeJSON(t, w, PhotoTokens{
			Photos: map[string]PhotoToken{"default": {Token: "tok-from-file"}},
		})
	})

	result, err := c.UploadPhotoFromFile(context.Background(), tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Photos["default"].Token != "tok-from-file" {
		t.Errorf("token = %q, want %q", result.Photos["default"].Token, "tok-from-file")
	}
	if uploadedFilename != "test-image.jpg" {
		t.Errorf("filename = %q, want %q", uploadedFilename, "test-image.jpg")
	}
}

func TestUploadPhotoFromFileNotFound(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, _ *http.Request) {
		t.Fatal("no HTTP request should be made")
	})

	_, err := c.UploadPhotoFromFile(context.Background(), "/nonexistent/photo.jpg")
	if err == nil {
		t.Fatal("expected error")
	}
	var e *Error
	if !errors.As(err, &e) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if e.Kind != ErrFetch {
		t.Errorf("Kind = %v, want ErrFetch", e.Kind)
	}
}

func TestUploadMediaFromFile(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "clip.mp4")
	if err := os.WriteFile(tmp, []byte("fake video data"), 0o644); err != nil {
		t.Fatal(err)
	}

	var requestCount atomic.Int32
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		count := requestCount.Add(1)
		if count == 1 {
			writeJSON(t, w, UploadEndpoint{URL: "http://" + r.Host + "/do-upload"})
			return
		}
		writeJSON(t, w, UploadedInfo{Token: "media-from-file"})
	})

	result, err := c.UploadMediaFromFile(context.Background(), UploadVideo, tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Token != "media-from-file" {
		t.Errorf("Token = %q, want %q", result.Token, "media-from-file")
	}
}

func TestUploadPhotoFromURL(t *testing.T) {
	var requestCount atomic.Int32
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		count := requestCount.Add(1)
		switch {
		case r.URL.Path == "/external/photo.png":
			// Serve the external image
			w.Header().Set("Content-Disposition", `attachment; filename="served.png"`)
			_, _ = w.Write([]byte("image bytes"))
		case count == 2:
			// GetUploadURL
			writeJSON(t, w, UploadEndpoint{URL: "http://" + r.Host + "/do-upload"})
		default:
			// Actual upload
			writeJSON(t, w, PhotoTokens{
				Photos: map[string]PhotoToken{"default": {Token: "tok-from-url"}},
			})
		}
	})

	result, err := c.UploadPhotoFromURL(context.Background(), srv.URL+"/external/photo.png")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Photos["default"].Token != "tok-from-url" {
		t.Errorf("token = %q, want %q", result.Photos["default"].Token, "tok-from-url")
	}
}

func TestUploadPhotoFromURLFetchError(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/missing.png" {
			http.NotFound(w, r)
			return
		}
		t.Fatalf("unexpected request to %s", r.URL.Path)
	})

	_, err := c.UploadPhotoFromURL(context.Background(), srv.URL+"/missing.png")
	if err == nil {
		t.Fatal("expected error")
	}
	var e *Error
	if !errors.As(err, &e) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if e.Kind != ErrFetch {
		t.Errorf("Kind = %v, want ErrFetch", e.Kind)
	}
	if e.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d, want 404", e.StatusCode)
	}
}

func TestUploadMediaFromURL(t *testing.T) {
	var requestCount atomic.Int32
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		count := requestCount.Add(1)
		switch {
		case r.URL.Path == "/external/video.mp4":
			_, _ = w.Write([]byte("video bytes"))
		case count == 2:
			writeJSON(t, w, UploadEndpoint{URL: "http://" + r.Host + "/do-upload"})
		default:
			writeJSON(t, w, UploadedInfo{Token: "media-from-url"})
		}
	})

	result, err := c.UploadMediaFromURL(context.Background(), UploadVideo, srv.URL+"/external/video.mp4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Token != "media-from-url" {
		t.Errorf("Token = %q, want %q", result.Token, "media-from-url")
	}
}
