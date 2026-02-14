package maxigo

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestGetBot(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("method = %q, want GET", r.Method)
			}
			if r.URL.Path != "/me" {
				t.Errorf("path = %q, want /me", r.URL.Path)
			}
			writeJSON(t, w, BotInfo{
				UserWithPhoto: UserWithPhoto{
					User: User{
						UserID:    12345,
						FirstName: "TestBot",
						IsBot:     true,
					},
				},
				Commands: []BotCommand{
					{Name: "start", Description: strPtr("Start the bot")},
				},
			})
		})

		bot, err := c.GetBot(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if bot.UserID != 12345 {
			t.Errorf("UserID = %d, want 12345", bot.UserID)
		}
		if bot.FirstName != "TestBot" {
			t.Errorf("FirstName = %q, want %q", bot.FirstName, "TestBot")
		}
		if !bot.IsBot {
			t.Error("IsBot should be true")
		}
		if len(bot.Commands) != 1 || bot.Commands[0].Name != "start" {
			t.Errorf("Commands = %v, want [{start}]", bot.Commands)
		}
	})

	t.Run("unauthorized", func(t *testing.T) {
		c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
			writeError(t, w, http.StatusUnauthorized, `{"code":"verify.token","message":"Invalid access_token"}`)
		})

		_, err := c.GetBot(context.Background())
		if err == nil {
			t.Fatal("expected error")
		}
		var e *Error
		if !errors.As(err, &e) {
			t.Fatalf("expected *Error, got %T", err)
		}
		if e.StatusCode != http.StatusUnauthorized {
			t.Errorf("StatusCode = %d, want 401", e.StatusCode)
		}
	})
}

func TestEditBot(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPatch {
				t.Errorf("method = %q, want PATCH", r.Method)
			}
			if r.URL.Path != "/me" {
				t.Errorf("path = %q, want /me", r.URL.Path)
			}

			var patch BotPatch
			readJSON(t, r, &patch)
			if !patch.FirstName.Set || patch.FirstName.Value != "NewName" {
				t.Errorf("FirstName = %v, want NewName", patch.FirstName)
			}

			writeJSON(t, w, BotInfo{
				UserWithPhoto: UserWithPhoto{
					User: User{
						UserID:    12345,
						FirstName: "NewName",
						IsBot:     true,
					},
				},
			})
		})

		bot, err := c.EditBot(context.Background(), &BotPatch{FirstName: Some("NewName")})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if bot.FirstName != "NewName" {
			t.Errorf("FirstName = %q, want %q", bot.FirstName, "NewName")
		}
	})
}
