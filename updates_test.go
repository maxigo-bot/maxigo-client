package maxigo

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"
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

func TestGetUpdatesPollDeadline(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("request should not be sent")
	})

	// Deadline shorter than polling timeout (30s) + buffer (5s)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := c.GetUpdates(ctx, GetUpdatesOpts{})
	if !errors.Is(err, ErrPollDeadline) {
		t.Errorf("err = %v, want ErrPollDeadline", err)
	}
}

func TestGetUpdatesAPIError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeError(t, w, http.StatusUnauthorized, `{"code":"verify.token","message":"Invalid access_token"}`)
	})

	_, err := c.GetUpdates(context.Background(), GetUpdatesOpts{})
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
	if e.Op != "GetUpdates" {
		t.Errorf("Op = %q, want GetUpdates", e.Op)
	}
}

func TestNewUpdateTypes(t *testing.T) {
	tests := []struct {
		name       string
		updateType UpdateType
		raw        json.RawMessage
		verify     func(t *testing.T, raw json.RawMessage)
	}{
		{
			name:       "bot_stopped",
			updateType: UpdateBotStopped,
			raw: mustMarshal(BotStoppedUpdate{
				Update: Update{UpdateType: UpdateBotStopped, Timestamp: 1000},
				ChatID: 42,
				User:   User{UserID: 1, FirstName: "Alice"},
			}),
			verify: func(t *testing.T, raw json.RawMessage) {
				var u BotStoppedUpdate
				if err := json.Unmarshal(raw, &u); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if u.UpdateType != UpdateBotStopped {
					t.Errorf("UpdateType = %q, want %q", u.UpdateType, UpdateBotStopped)
				}
				if u.ChatID != 42 {
					t.Errorf("ChatID = %d, want 42", u.ChatID)
				}
				if u.User.UserID != 1 {
					t.Errorf("User.UserID = %d, want 1", u.User.UserID)
				}
			},
		},
		{
			name:       "dialog_muted",
			updateType: UpdateDialogMuted,
			raw: mustMarshal(DialogMutedUpdate{
				Update: Update{UpdateType: UpdateDialogMuted, Timestamp: 2000},
				ChatID: 100,
				User:   User{UserID: 2, FirstName: "Bob"},
			}),
			verify: func(t *testing.T, raw json.RawMessage) {
				var u DialogMutedUpdate
				if err := json.Unmarshal(raw, &u); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if u.UpdateType != UpdateDialogMuted {
					t.Errorf("UpdateType = %q, want %q", u.UpdateType, UpdateDialogMuted)
				}
				if u.ChatID != 100 {
					t.Errorf("ChatID = %d, want 100", u.ChatID)
				}
			},
		},
		{
			name:       "dialog_unmuted",
			updateType: UpdateDialogUnmuted,
			raw: mustMarshal(DialogUnmutedUpdate{
				Update: Update{UpdateType: UpdateDialogUnmuted, Timestamp: 3000},
				ChatID: 100,
				User:   User{UserID: 2, FirstName: "Bob"},
			}),
			verify: func(t *testing.T, raw json.RawMessage) {
				var u DialogUnmutedUpdate
				if err := json.Unmarshal(raw, &u); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if u.UpdateType != UpdateDialogUnmuted {
					t.Errorf("UpdateType = %q, want %q", u.UpdateType, UpdateDialogUnmuted)
				}
			},
		},
		{
			name:       "dialog_cleared",
			updateType: UpdateDialogCleared,
			raw: mustMarshal(DialogClearedUpdate{
				Update: Update{UpdateType: UpdateDialogCleared, Timestamp: 4000},
				ChatID: 200,
				User:   User{UserID: 3, FirstName: "Carol"},
			}),
			verify: func(t *testing.T, raw json.RawMessage) {
				var u DialogClearedUpdate
				if err := json.Unmarshal(raw, &u); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if u.UpdateType != UpdateDialogCleared {
					t.Errorf("UpdateType = %q, want %q", u.UpdateType, UpdateDialogCleared)
				}
				if u.ChatID != 200 {
					t.Errorf("ChatID = %d, want 200", u.ChatID)
				}
			},
		},
		{
			name:       "dialog_removed",
			updateType: UpdateDialogRemoved,
			raw: mustMarshal(DialogRemovedUpdate{
				Update: Update{UpdateType: UpdateDialogRemoved, Timestamp: 5000},
				ChatID: 300,
				User:   User{UserID: 4, FirstName: "Dave"},
			}),
			verify: func(t *testing.T, raw json.RawMessage) {
				var u DialogRemovedUpdate
				if err := json.Unmarshal(raw, &u); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if u.UpdateType != UpdateDialogRemoved {
					t.Errorf("UpdateType = %q, want %q", u.UpdateType, UpdateDialogRemoved)
				}
				if u.User.FirstName != "Dave" {
					t.Errorf("User.FirstName = %q, want %q", u.User.FirstName, "Dave")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify round-trip: marshal → unmarshal produces correct update_type
			var header struct {
				UpdateType UpdateType `json:"update_type"`
			}
			if err := json.Unmarshal(tt.raw, &header); err != nil {
				t.Fatalf("unmarshal header: %v", err)
			}
			if header.UpdateType != tt.updateType {
				t.Errorf("update_type = %q, want %q", header.UpdateType, tt.updateType)
			}
			tt.verify(t, tt.raw)
		})
	}
}

func mustMarshal(v any) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
