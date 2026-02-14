package maxigo

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestSendMessage(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("method = %q, want POST", r.Method)
			}
			if r.URL.Path != "/messages" {
				t.Errorf("path = %q, want /messages", r.URL.Path)
			}
			if r.URL.Query().Get("chat_id") != "100" {
				t.Errorf("chat_id = %q, want 100", r.URL.Query().Get("chat_id"))
			}

			var body NewMessageBody
			readJSON(t, r, &body)
			if body.Text == nil || *body.Text != "Hello!" {
				t.Errorf("text = %v, want %q", body.Text, "Hello!")
			}

			mid := "msg-123"
			text := "Hello!"
			writeJSON(t, w, sendMessageResult{
				Message: Message{
					Timestamp: 1234567890,
					Body: MessageBody{
						MID:  mid,
						Text: &text,
					},
				},
			})
		})

		msg, err := c.SendMessage(context.Background(), 100, &NewMessageBody{Text: strPtr("Hello!")})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if msg.Body.MID != "msg-123" {
			t.Errorf("MID = %q, want %q", msg.Body.MID, "msg-123")
		}
	})
}

func TestSendMessageToUser(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("user_id") != "42" {
			t.Errorf("user_id = %q, want 42", r.URL.Query().Get("user_id"))
		}
		writeJSON(t, w, sendMessageResult{
			Message: Message{Timestamp: 1},
		})
	})

	_, err := c.SendMessageToUser(context.Background(), 42, &NewMessageBody{Text: strPtr("hi")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSendMessageDisableLinkPreview(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("disable_link_preview") != "true" {
			t.Errorf("disable_link_preview = %q, want true", r.URL.Query().Get("disable_link_preview"))
		}
		writeJSON(t, w, sendMessageResult{
			Message: Message{Timestamp: 1},
		})
	})

	_, err := c.SendMessage(context.Background(), 100, &NewMessageBody{
		Text:               strPtr("https://example.com"),
		DisableLinkPreview: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEditMessage(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %q, want PUT", r.Method)
		}
		if r.URL.Query().Get("message_id") != "mid-1" {
			t.Errorf("message_id = %q, want mid-1", r.URL.Query().Get("message_id"))
		}
		writeJSON(t, w, SimpleQueryResult{Success: true})
	})

	result, err := c.EditMessage(context.Background(), "mid-1", &NewMessageBody{Text: strPtr("edited")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Error("Success should be true")
	}
}

func TestDeleteMessage(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Query().Get("message_id") != "mid-2" {
			t.Errorf("message_id = %q, want mid-2", r.URL.Query().Get("message_id"))
		}
		writeJSON(t, w, SimpleQueryResult{Success: true})
	})

	result, err := c.DeleteMessage(context.Background(), "mid-2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Error("Success should be true")
	}
}

func TestGetMessages(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Query().Get("chat_id") != "200" {
			t.Errorf("chat_id = %q, want 200", r.URL.Query().Get("chat_id"))
		}
		if r.URL.Query().Get("count") != "10" {
			t.Errorf("count = %q, want 10", r.URL.Query().Get("count"))
		}

		text := "Hello"
		writeJSON(t, w, MessageList{
			Messages: []Message{
				{Timestamp: 1, Body: MessageBody{MID: "m1", Text: &text}},
			},
		})
	})

	result, err := c.GetMessages(context.Background(), GetMessagesOpts{ChatID: 200, Count: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Messages) != 1 {
		t.Fatalf("len(Messages) = %d, want 1", len(result.Messages))
	}
	if result.Messages[0].Body.MID != "m1" {
		t.Errorf("MID = %q, want %q", result.Messages[0].Body.MID, "m1")
	}
}

func TestGetMessagesWithIDs(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query().Get("message_ids")
		if ids != "mid-1,mid-2,mid-3" {
			t.Errorf("message_ids = %q, want %q", ids, "mid-1,mid-2,mid-3")
		}
		writeJSON(t, w, MessageList{Messages: []Message{}})
	})

	_, err := c.GetMessages(context.Background(), GetMessagesOpts{MessageIDs: []string{"mid-1", "mid-2", "mid-3"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetMessagesWithTimeRange(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("from") != "1000" {
			t.Errorf("from = %q, want 1000", r.URL.Query().Get("from"))
		}
		if r.URL.Query().Get("to") != "500" {
			t.Errorf("to = %q, want 500", r.URL.Query().Get("to"))
		}
		writeJSON(t, w, MessageList{Messages: []Message{}})
	})

	_, err := c.GetMessages(context.Background(), GetMessagesOpts{ChatID: 100, From: 1000, To: 500})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetMessageByID(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/messages/mid-99" {
			t.Errorf("path = %q, want /messages/mid-99", r.URL.Path)
		}
		writeJSON(t, w, Message{
			Timestamp: 1, Body: MessageBody{MID: "mid-99"},
		})
	})

	msg, err := c.GetMessageByID(context.Background(), "mid-99")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Body.MID != "mid-99" {
		t.Errorf("MID = %q, want %q", msg.Body.MID, "mid-99")
	}
}

func TestAnswerCallback(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("method = %q, want POST", r.Method)
			}
			if r.URL.Path != "/answers" {
				t.Errorf("path = %q, want /answers", r.URL.Path)
			}
			if r.URL.Query().Get("callback_id") != "cb-1" {
				t.Errorf("callback_id = %q, want cb-1", r.URL.Query().Get("callback_id"))
			}

			var answer CallbackAnswer
			readJSON(t, r, &answer)
			if answer.Notification == nil || *answer.Notification != "Done!" {
				t.Errorf("Notification = %v, want Done!", answer.Notification)
			}

			writeJSON(t, w, SimpleQueryResult{Success: true})
		})

		notif := "Done!"
		result, err := c.AnswerCallback(context.Background(), "cb-1", &CallbackAnswer{
			Notification: &notif,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Error("Success should be true")
		}
	})

	t.Run("error", func(t *testing.T) {
		c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
			writeError(t, w, http.StatusMethodNotAllowed, `{"code":"method.not.allowed","message":"method not allowed"}`)
		})

		_, err := c.AnswerCallback(context.Background(), "cb-1", &CallbackAnswer{})
		var e *Error
		if !errors.As(err, &e) {
			t.Fatalf("expected *Error, got %T", err)
		}
		if e.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("StatusCode = %d, want 405", e.StatusCode)
		}
	})
}
