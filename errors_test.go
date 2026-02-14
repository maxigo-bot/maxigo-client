package maxigo

import (
	"errors"
	"fmt"
	"io"
	"testing"
)

func TestErrorKindString(t *testing.T) {
	tests := []struct {
		kind ErrorKind
		want string
	}{
		{ErrAPI, "api"},
		{ErrNetwork, "network"},
		{ErrTimeout, "timeout"},
		{ErrDecode, "decode"},
		{ErrorKind(99), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.kind.String(); got != tt.want {
			t.Errorf("ErrorKind(%d).String() = %q, want %q", tt.kind, got, tt.want)
		}
	}
}

func TestErrorError(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{
			name: "api error with status code",
			err:  apiError("SendMessage", 403, "you don't have permission"),
			want: "SendMessage: api error 403: you don't have permission",
		},
		{
			name: "network error",
			err:  networkError("GetChat", fmt.Errorf("connection refused")),
			want: "GetChat: network: connection refused",
		},
		{
			name: "timeout error",
			err:  timeoutError("GetUpdates", fmt.Errorf("context deadline exceeded")),
			want: "GetUpdates: timeout: context deadline exceeded",
		},
		{
			name: "decode error",
			err:  decodeError("SendMessage", fmt.Errorf("unexpected EOF")),
			want: "SendMessage: decode: unexpected EOF",
		},
		{
			name: "error without message",
			err:  &Error{Kind: ErrNetwork, Op: "GetBot"},
			want: "GetBot: network",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestErrorUnwrap(t *testing.T) {
	underlying := io.ErrUnexpectedEOF
	err := networkError("GetChat", underlying)

	if !errors.Is(err, underlying) {
		t.Error("errors.Is should find the underlying error")
	}

	var target *Error
	if !errors.As(err, &target) {
		t.Fatal("errors.As should extract *Error")
	}
	if target.Kind != ErrNetwork {
		t.Errorf("Kind = %v, want ErrNetwork", target.Kind)
	}
	if target.Op != "GetChat" {
		t.Errorf("Op = %q, want %q", target.Op, "GetChat")
	}
}

func TestErrorUnwrapNil(t *testing.T) {
	err := apiError("SendMessage", 403, "forbidden")
	if err.Unwrap() != nil {
		t.Error("API error should have nil Unwrap()")
	}
}

func TestErrorTimeout(t *testing.T) {
	tests := []struct {
		kind ErrorKind
		want bool
	}{
		{ErrTimeout, true},
		{ErrAPI, false},
		{ErrNetwork, false},
		{ErrDecode, false},
	}
	for _, tt := range tests {
		err := &Error{Kind: tt.kind, Op: "test"}
		if got := err.Timeout(); got != tt.want {
			t.Errorf("Error{Kind: %v}.Timeout() = %v, want %v", tt.kind, got, tt.want)
		}
	}
}

func TestErrorWrappedInFmtErrorf(t *testing.T) {
	original := apiError("SendMessage", 429, "rate limited")
	wrapped := fmt.Errorf("handler failed: %w", original)

	var target *Error
	if !errors.As(wrapped, &target) {
		t.Fatal("errors.As should extract *Error from wrapped error")
	}
	if target.StatusCode != 429 {
		t.Errorf("StatusCode = %d, want 429", target.StatusCode)
	}
}

func TestErrEmptyToken(t *testing.T) {
	if ErrEmptyToken == nil {
		t.Fatal("ErrEmptyToken should not be nil")
	}
	if ErrEmptyToken.Error() != "bot token is empty" {
		t.Errorf("ErrEmptyToken.Error() = %q, want %q", ErrEmptyToken.Error(), "bot token is empty")
	}
}

func TestConstructors(t *testing.T) {
	t.Run("apiError", func(t *testing.T) {
		err := apiError("SendMessage", 500, "internal server error")
		if err.Kind != ErrAPI {
			t.Errorf("Kind = %v, want ErrAPI", err.Kind)
		}
		if err.StatusCode != 500 {
			t.Errorf("StatusCode = %d, want 500", err.StatusCode)
		}
		if err.Message != "internal server error" {
			t.Errorf("Message = %q, want %q", err.Message, "internal server error")
		}
		if err.Op != "SendMessage" {
			t.Errorf("Op = %q, want %q", err.Op, "SendMessage")
		}
		if err.Err != nil {
			t.Error("Err should be nil for API errors")
		}
	})

	t.Run("networkError", func(t *testing.T) {
		underlying := io.EOF
		err := networkError("GetChat", underlying)
		if err.Kind != ErrNetwork {
			t.Errorf("Kind = %v, want ErrNetwork", err.Kind)
		}
		if err.StatusCode != 0 {
			t.Errorf("StatusCode = %d, want 0", err.StatusCode)
		}
		if err.Err != underlying {
			t.Error("Err should be the underlying error")
		}
	})

	t.Run("timeoutError", func(t *testing.T) {
		underlying := fmt.Errorf("context deadline exceeded")
		err := timeoutError("GetUpdates", underlying)
		if err.Kind != ErrTimeout {
			t.Errorf("Kind = %v, want ErrTimeout", err.Kind)
		}
		if !err.Timeout() {
			t.Error("Timeout() should return true")
		}
		if err.Err != underlying {
			t.Error("Err should be the underlying error")
		}
	})

	t.Run("decodeError", func(t *testing.T) {
		underlying := fmt.Errorf("invalid character")
		err := decodeError("SendMessage", underlying)
		if err.Kind != ErrDecode {
			t.Errorf("Kind = %v, want ErrDecode", err.Kind)
		}
		if err.Err != underlying {
			t.Error("Err should be the underlying error")
		}
	})
}

func TestErrorImplementsErrorInterface(t *testing.T) {
	var _ error = (*Error)(nil)
}
