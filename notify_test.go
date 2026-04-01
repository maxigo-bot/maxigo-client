package maxigo

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestCheckPhoneNumbers(t *testing.T) {
	t.Run("existing numbers", func(t *testing.T) {
		c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("method = %q, want GET", r.Method)
			}
			if r.URL.Path != "/notify/exists" {
				t.Errorf("path = %q, want /notify/exists", r.URL.Path)
			}
			phones := r.URL.Query().Get("phone_numbers")
			if phones != "79001234567,79007654321" {
				t.Errorf("phone_numbers = %q, want %q", phones, "79001234567,79007654321")
			}
			writeJSON(t, w, map[string][]string{
				"existing_phone_numbers": {"79001234567"},
			})
		})

		result, err := c.CheckPhoneNumbers(context.Background(), []string{"79001234567", "79007654321"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 1 || result[0] != "79001234567" {
			t.Errorf("result = %v, want [79001234567]", result)
		}
	})

	t.Run("no numbers exist", func(t *testing.T) {
		c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(t, w, map[string][]string{
				"existing_phone_numbers": {},
			})
		})

		result, err := c.CheckPhoneNumbers(context.Background(), []string{"79009999999"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("result = %v, want empty", result)
		}
	})

	t.Run("single number", func(t *testing.T) {
		c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
			phones := r.URL.Query().Get("phone_numbers")
			if phones != "79001234567" {
				t.Errorf("phone_numbers = %q, want single number", phones)
			}
			writeJSON(t, w, map[string][]string{
				"existing_phone_numbers": {"79001234567"},
			})
		})

		result, err := c.CheckPhoneNumbers(context.Background(), []string{"79001234567"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 1 {
			t.Errorf("len(result) = %d, want 1", len(result))
		}
	})

	t.Run("API error", func(t *testing.T) {
		c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
			writeError(t, w, http.StatusUnauthorized, `{"code":"verify.token","message":"Invalid access_token"}`)
		})

		_, err := c.CheckPhoneNumbers(context.Background(), []string{"79001234567"})
		var e *Error
		if !errors.As(err, &e) {
			t.Fatalf("expected *Error, got %T", err)
		}
		if e.Kind != ErrAPI {
			t.Errorf("Kind = %v, want ErrAPI", e.Kind)
		}
		if e.Op != "CheckPhoneNumbers" {
			t.Errorf("Op = %q, want CheckPhoneNumbers", e.Op)
		}
	})
}
