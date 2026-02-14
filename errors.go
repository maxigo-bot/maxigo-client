package maxigo

import (
	"errors"
	"fmt"
)

// ErrEmptyToken is returned when an empty token is passed to New.
var ErrEmptyToken = errors.New("bot token is empty")

// ErrorKind classifies the category of an error returned by the client.
type ErrorKind int

const (
	// ErrAPI indicates the Max Bot API returned a non-200 HTTP response.
	ErrAPI ErrorKind = iota
	// ErrNetwork indicates an HTTP transport failure (connection refused, DNS, etc.).
	ErrNetwork
	// ErrTimeout indicates the request timed out or the context deadline was exceeded.
	ErrTimeout
	// ErrDecode indicates a JSON marshal or unmarshal failure.
	ErrDecode
)

// String returns a human-readable name for the error kind.
func (k ErrorKind) String() string {
	switch k {
	case ErrAPI:
		return "api"
	case ErrNetwork:
		return "network"
	case ErrTimeout:
		return "timeout"
	case ErrDecode:
		return "decode"
	default:
		return "unknown"
	}
}

// Error represents any error produced by the maxigo client.
//
// Use [errors.As] to extract the Error and inspect its fields:
//
//	var e *maxigo.Error
//	if errors.As(err, &e) {
//	    switch e.Kind {
//	    case maxigo.ErrAPI:
//	        log.Printf("API error %d: %s", e.StatusCode, e.Message)
//	    case maxigo.ErrTimeout:
//	        log.Printf("timeout in %s", e.Op)
//	    }
//	}
type Error struct {
	// Kind classifies the error category.
	Kind ErrorKind
	// StatusCode is the HTTP status code. Only set when Kind is ErrAPI.
	StatusCode int
	// Message is a human-readable error description.
	Message string
	// Op is the client operation that failed (e.g. "SendMessage", "GetChat").
	Op string
	// Err is the underlying error, if any.
	Err error
}

// Error returns a formatted error string including the operation, kind, and details.
func (e *Error) Error() string {
	if e.Kind == ErrAPI {
		return fmt.Sprintf("%s: %s error %d: %s", e.Op, e.Kind, e.StatusCode, e.Message)
	}
	if e.Message != "" {
		return fmt.Sprintf("%s: %s: %s", e.Op, e.Kind, e.Message)
	}
	return fmt.Sprintf("%s: %s", e.Op, e.Kind)
}

// Unwrap returns the underlying error for use with [errors.Is] and [errors.As].
func (e *Error) Unwrap() error {
	return e.Err
}

// Timeout reports whether the error represents a timeout.
// This implements the informal timeout interface used by [net.Error].
func (e *Error) Timeout() bool {
	return e.Kind == ErrTimeout
}

func apiError(op string, statusCode int, message string) *Error {
	return &Error{
		Kind:       ErrAPI,
		StatusCode: statusCode,
		Message:    message,
		Op:         op,
	}
}

func networkError(op string, err error) *Error {
	return &Error{
		Kind:    ErrNetwork,
		Message: err.Error(),
		Op:      op,
		Err:     err,
	}
}

func timeoutError(op string, err error) *Error {
	return &Error{
		Kind:    ErrTimeout,
		Message: err.Error(),
		Op:      op,
		Err:     err,
	}
}

func decodeError(op string, err error) *Error {
	return &Error{
		Kind:    ErrDecode,
		Message: err.Error(),
		Op:      op,
		Err:     err,
	}
}
