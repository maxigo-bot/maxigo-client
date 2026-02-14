package maxigo

import (
	"context"
	"net/http"
	"testing"
)

func TestSubscribe(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/subscriptions" {
			t.Errorf("path = %q, want /subscriptions", r.URL.Path)
		}

		var body SubscriptionRequestBody
		readJSON(t, r, &body)
		if body.URL != "https://example.com/webhook" {
			t.Errorf("URL = %q, want %q", body.URL, "https://example.com/webhook")
		}
		if len(body.UpdateTypes) != 2 {
			t.Errorf("len(UpdateTypes) = %d, want 2", len(body.UpdateTypes))
		}
		if body.Secret != "my-secret" {
			t.Errorf("Secret = %q, want %q", body.Secret, "my-secret")
		}

		writeJSON(t, w, SimpleQueryResult{Success: true})
	})

	result, err := c.Subscribe(context.Background(), "https://example.com/webhook",
		[]string{"message_created", "bot_started"}, "my-secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Error("Success should be true")
	}
}

func TestUnsubscribe(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Query().Get("url") != "https://example.com/webhook" {
			t.Errorf("url = %q, want %q", r.URL.Query().Get("url"), "https://example.com/webhook")
		}
		writeJSON(t, w, SimpleQueryResult{Success: true})
	})

	result, err := c.Unsubscribe(context.Background(), "https://example.com/webhook")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Error("Success should be true")
	}
}

func TestGetSubscriptions(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/subscriptions" {
			t.Errorf("path = %q, want /subscriptions", r.URL.Path)
		}
		writeJSON(t, w, getSubscriptionsResult{
			Subscriptions: []Subscription{
				{
					URL:         "https://example.com/webhook",
					Time:        1234567890,
					UpdateTypes: []string{"message_created"},
				},
			},
		})
	})

	subs, err := c.GetSubscriptions(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(subs) != 1 {
		t.Fatalf("len(subs) = %d, want 1", len(subs))
	}
	if subs[0].URL != "https://example.com/webhook" {
		t.Errorf("URL = %q, want %q", subs[0].URL, "https://example.com/webhook")
	}
}
