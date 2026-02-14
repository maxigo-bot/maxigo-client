package maxigo

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		c, err := New("test-token")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c == nil {
			t.Fatal("client should not be nil")
		}
		if c.token != "test-token" {
			t.Errorf("token = %q, want %q", c.token, "test-token")
		}
		if c.baseURL != defaultBaseURL {
			t.Errorf("baseURL = %q, want %q", c.baseURL, defaultBaseURL)
		}
	})

	t.Run("empty token", func(t *testing.T) {
		_, err := New("")
		if !errors.Is(err, ErrEmptyToken) {
			t.Errorf("err = %v, want ErrEmptyToken", err)
		}
	})

	t.Run("with base URL", func(t *testing.T) {
		c, err := New("token", WithBaseURL("http://localhost:8080"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.baseURL != "http://localhost:8080" {
			t.Errorf("baseURL = %q, want %q", c.baseURL, "http://localhost:8080")
		}
	})

	t.Run("with timeout", func(t *testing.T) {
		c, err := New("token", WithTimeout(5*time.Second))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.httpClient.Timeout != 5*time.Second {
			t.Errorf("timeout = %v, want %v", c.httpClient.Timeout, 5*time.Second)
		}
	})

	t.Run("with custom HTTP client", func(t *testing.T) {
		custom := &http.Client{Timeout: 99 * time.Second}
		c, err := New("token", WithHTTPClient(custom))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.httpClient.Timeout != 99*time.Second {
			t.Errorf("timeout = %v, want %v", c.httpClient.Timeout, 99*time.Second)
		}
		if c.httpClient == custom {
			t.Error("httpClient should be a copy, not the original")
		}
	})
}

func testClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	c, err := New("test-token", WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return c, srv
}

// strPtr returns a pointer to s. Test helper for *string fields.
func strPtr(s string) *string { return &s }

// writeJSON writes v as JSON to the response writer. Test helper.
func writeJSON(t *testing.T, w http.ResponseWriter, v any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		t.Fatalf("failed to encode JSON: %v", err)
	}
}

// readJSON decodes the request body as JSON. Test helper.
func readJSON(t *testing.T, r *http.Request, v any) {
	t.Helper()
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
}

// writeError writes an error response. Test helper.
func writeError(t *testing.T, w http.ResponseWriter, statusCode int, body string) {
	t.Helper()
	w.WriteHeader(statusCode)
	if _, err := w.Write([]byte(body)); err != nil {
		t.Fatalf("failed to write: %v", err)
	}
}

func TestDoSuccess(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "test-token" {
			t.Errorf("Authorization = %q, want %q", r.Header.Get("Authorization"), "test-token")
		}
		if r.URL.Query().Get("access_token") != "" {
			t.Errorf("access_token should not be in query, got %q", r.URL.Query().Get("access_token"))
		}
		w.Header().Set("Content-Type", "application/json")
		writeJSON(t, w, map[string]string{"status": "ok"})
	})

	var result map[string]string
	err := c.do(context.Background(), "Test", http.MethodGet, "/test", nil, nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["status"] != "ok" {
		t.Errorf("status = %q, want %q", result["status"], "ok")
	}
}

func TestDoWithBody(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %q, want %q", r.Header.Get("Content-Type"), "application/json")
		}
		var body map[string]string
		readJSON(t, r, &body)
		if body["text"] != "hello" {
			t.Errorf("text = %q, want %q", body["text"], "hello")
		}
		writeJSON(t, w, map[string]bool{"success": true})
	})

	body := map[string]string{"text": "hello"}
	var result map[string]bool
	err := c.do(context.Background(), "Test", http.MethodPost, "/test", nil, body, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result["success"] {
		t.Error("success should be true")
	}
}

func TestDoAPIError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantMsg    string
	}{
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			body:       `{"code":"verify.token","message":"Invalid access_token"}`,
			wantMsg:    "Invalid access_token",
		},
		{
			name:       "forbidden",
			statusCode: http.StatusForbidden,
			body:       `{"code":"access.denied","message":"you don't have permission"}`,
			wantMsg:    "you don't have permission",
		},
		{
			name:       "not found",
			statusCode: http.StatusNotFound,
			body:       `{"code":"not.found","message":"chat not found"}`,
			wantMsg:    "chat not found",
		},
		{
			name:       "rate limited",
			statusCode: http.StatusTooManyRequests,
			body:       `{"code":"rate.limit","message":"too many requests"}`,
			wantMsg:    "too many requests",
		},
		{
			name:       "server error",
			statusCode: http.StatusInternalServerError,
			body:       `{"code":"internal","message":"internal server error"}`,
			wantMsg:    "internal server error",
		},
		{
			name:       "invalid JSON error body",
			statusCode: http.StatusBadGateway,
			body:       `not json`,
			wantMsg:    "Bad Gateway",
		},
		{
			name:       "error with error field only",
			statusCode: http.StatusBadRequest,
			body:       `{"code":"bad","error":"something broke"}`,
			wantMsg:    "something broke",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
				writeError(t, w, tt.statusCode, tt.body)
			})

			err := c.do(context.Background(), "TestOp", http.MethodGet, "/test", nil, nil, nil)
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
			if e.StatusCode != tt.statusCode {
				t.Errorf("StatusCode = %d, want %d", e.StatusCode, tt.statusCode)
			}
			if e.Message != tt.wantMsg {
				t.Errorf("Message = %q, want %q", e.Message, tt.wantMsg)
			}
			if e.Op != "TestOp" {
				t.Errorf("Op = %q, want %q", e.Op, "TestOp")
			}
		})
	}
}

func TestDoContextCanceled(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	err := c.do(ctx, "TestOp", http.MethodGet, "/test", nil, nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}

	var e *Error
	if !errors.As(err, &e) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if e.Kind != ErrTimeout {
		t.Errorf("Kind = %v, want ErrTimeout", e.Kind)
	}
	if !e.Timeout() {
		t.Error("Timeout() should return true")
	}
}

func TestDoDecodeError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{invalid json`))
	})

	var result map[string]string
	err := c.do(context.Background(), "TestOp", http.MethodGet, "/test", nil, nil, &result)
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
}

func TestDoNilResult(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ignored":"data"}`))
	})

	err := c.do(context.Background(), "TestOp", http.MethodDelete, "/test", nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDoQueryParams(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("chat_id") != "123" {
			t.Errorf("chat_id = %q, want %q", r.URL.Query().Get("chat_id"), "123")
		}
		if r.URL.Query().Get("count") != "50" {
			t.Errorf("count = %q, want %q", r.URL.Query().Get("count"), "50")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})

	q := make(url.Values)
	q.Set("chat_id", "123")
	q.Set("count", "50")
	err := c.do(context.Background(), "TestOp", http.MethodGet, "/test", q, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
