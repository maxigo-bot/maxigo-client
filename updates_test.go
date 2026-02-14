package maxigo

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetUpdates(t *testing.T) {
	t.Run("success with message_created", func(t *testing.T) {
		c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("method = %q, want GET", r.Method)
			}
			if r.URL.Path != "/updates" {
				t.Errorf("path = %q, want /updates", r.URL.Path)
			}
			if r.URL.Query().Get("limit") != "10" {
				t.Errorf("limit = %q, want 10", r.URL.Query().Get("limit"))
			}
			if r.URL.Query().Get("timeout") != "30" {
				t.Errorf("timeout = %q, want 30", r.URL.Query().Get("timeout"))
			}

			text := "Hello"
			marker := int64(12345)
			writeJSON(t, w, UpdateList{
				Updates: []json.RawMessage{
					mustMarshal(MessageCreatedUpdate{
						Update: Update{
							UpdateType: UpdateMessageCreated,
							Timestamp:  1000,
						},
						Message: Message{
							Timestamp: 1000,
							Body:      MessageBody{MID: "mid-1", Text: &text},
						},
					}),
				},
				Marker: &marker,
			})
		})

		result, err := c.GetUpdates(context.Background(), GetUpdatesOpts{Limit: 10, Timeout: 30})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Updates) != 1 {
			t.Fatalf("len(Updates) = %d, want 1", len(result.Updates))
		}
		if result.Marker == nil || *result.Marker != 12345 {
			t.Errorf("Marker = %v, want 12345", result.Marker)
		}

		// Parse the raw update
		var update MessageCreatedUpdate
		if err := json.Unmarshal(result.Updates[0], &update); err != nil {
			t.Fatalf("unmarshal update: %v", err)
		}
		if update.UpdateType != UpdateMessageCreated {
			t.Errorf("UpdateType = %q, want %q", update.UpdateType, UpdateMessageCreated)
		}
		if update.Message.Body.MID != "mid-1" {
			t.Errorf("MID = %q, want %q", update.Message.Body.MID, "mid-1")
		}
	})

	t.Run("empty updates", func(t *testing.T) {
		c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(t, w, UpdateList{
				Updates: []json.RawMessage{},
			})
		})

		result, err := c.GetUpdates(context.Background(), GetUpdatesOpts{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Updates) != 0 {
			t.Errorf("len(Updates) = %d, want 0", len(result.Updates))
		}
	})

	t.Run("with marker and types", func(t *testing.T) {
		c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("marker") != "5000" {
				t.Errorf("marker = %q, want 5000", r.URL.Query().Get("marker"))
			}
			types := r.URL.Query().Get("types")
			if types != "message_created,message_callback" {
				t.Errorf("types = %q, want %q", types, "message_created,message_callback")
			}
			writeJSON(t, w, UpdateList{Updates: []json.RawMessage{}})
		})

		_, err := c.GetUpdates(context.Background(), GetUpdatesOpts{
			Marker: 5000,
			Types:  []string{"message_created", "message_callback"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("multiple update types", func(t *testing.T) {
		c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(t, w, UpdateList{
				Updates: []json.RawMessage{
					mustMarshal(BotStartedUpdate{
						Update: Update{UpdateType: UpdateBotStarted, Timestamp: 1000},
						ChatID: 42,
						User:   User{UserID: 1, FirstName: "Test"},
					}),
					mustMarshal(MessageCallbackUpdate{
						Update: Update{UpdateType: UpdateMessageCallback, Timestamp: 2000},
						Callback: Callback{
							Timestamp:  2000,
							CallbackID: "cb-1",
							Payload:    "action",
							User:       User{UserID: 1},
						},
					}),
				},
			})
		})

		result, err := c.GetUpdates(context.Background(), GetUpdatesOpts{Limit: 100, Timeout: 30})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Updates) != 2 {
			t.Fatalf("len(Updates) = %d, want 2", len(result.Updates))
		}

		// Parse first update
		var botStarted BotStartedUpdate
		if err := json.Unmarshal(result.Updates[0], &botStarted); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if botStarted.UpdateType != UpdateBotStarted {
			t.Errorf("UpdateType = %q, want %q", botStarted.UpdateType, UpdateBotStarted)
		}
		if botStarted.ChatID != 42 {
			t.Errorf("ChatID = %d, want 42", botStarted.ChatID)
		}

		// Parse second update
		var callback MessageCallbackUpdate
		if err := json.Unmarshal(result.Updates[1], &callback); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if callback.Callback.Payload != "action" {
			t.Errorf("Payload = %q, want %q", callback.Callback.Payload, "action")
		}
	})
}

func mustMarshal(v any) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
