package maxigo

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestGetChat(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chats/100" {
			t.Errorf("path = %q, want /chats/100", r.URL.Path)
		}
		title := "Test Chat"
		writeJSON(t, w, Chat{
			ChatID: 100,
			Type:   ChatGroup,
			Status: ChatStatusActive,
			Title:  &title,
		})
	})

	chat, err := c.GetChat(context.Background(), 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if chat.ChatID != 100 {
		t.Errorf("ChatID = %d, want 100", chat.ChatID)
	}
	if chat.Type != ChatGroup {
		t.Errorf("Type = %q, want %q", chat.Type, ChatGroup)
	}
}

func TestGetChatByLink(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		// url.PathEscape encodes "/" as "%2F" in the raw path
		escaped := r.URL.EscapedPath()
		if escaped != "/chats/link%2Fwith%2Fslashes" {
			t.Errorf("escaped path = %q, want /chats/link%%2Fwith%%2Fslashes", escaped)
		}
		writeJSON(t, w, Chat{ChatID: 1})
	})

	_, err := c.GetChatByLink(context.Background(), "link/with/slashes")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetChats(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chats" {
			t.Errorf("path = %q, want /chats", r.URL.Path)
		}
		if r.URL.Query().Get("count") != "10" {
			t.Errorf("count = %q, want 10", r.URL.Query().Get("count"))
		}

		marker := int64(999)
		writeJSON(t, w, ChatList{
			Chats:  []Chat{{ChatID: 1}, {ChatID: 2}},
			Marker: &marker,
		})
	})

	result, err := c.GetChats(context.Background(), GetChatsOpts{Count: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Chats) != 2 {
		t.Errorf("len(Chats) = %d, want 2", len(result.Chats))
	}
	if result.Marker == nil || *result.Marker != 999 {
		t.Errorf("Marker = %v, want 999", result.Marker)
	}
}

func TestEditChat(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", r.Method)
		}
		if r.URL.Path != "/chats/100" {
			t.Errorf("path = %q, want /chats/100", r.URL.Path)
		}
		title := "Updated"
		writeJSON(t, w, Chat{ChatID: 100, Title: &title})
	})

	chat, err := c.EditChat(context.Background(), 100, &ChatPatch{Title: Some("Updated")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if chat.Title == nil || *chat.Title != "Updated" {
		t.Errorf("Title = %v, want Updated", chat.Title)
	}
}

func TestGetMembers(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chats/100/members" {
			t.Errorf("path = %q, want /chats/100/members", r.URL.Path)
		}
		writeJSON(t, w, ChatMembersList{
			Members: []ChatMember{
				{UserWithPhoto: UserWithPhoto{User: User{UserID: 1}}},
				{UserWithPhoto: UserWithPhoto{User: User{UserID: 2}}},
			},
		})
	})

	result, err := c.GetMembers(context.Background(), 100, GetMembersOpts{Count: 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Members) != 2 {
		t.Errorf("len(Members) = %d, want 2", len(result.Members))
	}
}

func TestGetMembersWithUserIDs(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query().Get("user_ids")
		if ids != "1,2,3" {
			t.Errorf("user_ids = %q, want %q", ids, "1,2,3")
		}
		writeJSON(t, w, ChatMembersList{Members: []ChatMember{}})
	})

	_, err := c.GetMembers(context.Background(), 100, GetMembersOpts{UserIDs: []int64{1, 2, 3}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetAdmins(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chats/100/members/admins" {
			t.Errorf("path = %q, want /chats/100/members/admins", r.URL.Path)
		}
		writeJSON(t, w, ChatAdminsList{
			Admins: []ChatAdmin{
				{UserID: 1, Permissions: []ChatAdminPermission{PermReadAllMessages}},
			},
		})
	})

	result, err := c.GetAdmins(context.Background(), 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Admins) != 1 {
		t.Errorf("len(Admins) = %d, want 1", len(result.Admins))
	}
}

func TestSetAdmins(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/chats/100/members/admins" {
			t.Errorf("path = %q, want /chats/100/members/admins", r.URL.Path)
		}
		var body ChatAdminsList
		readJSON(t, r, &body)
		if len(body.Admins) != 1 {
			t.Errorf("len(Admins) = %d, want 1", len(body.Admins))
		}
		writeJSON(t, w, SimpleQueryResult{Success: true})
	})

	result, err := c.SetAdmins(context.Background(), 100, &ChatAdminsList{
		Admins: []ChatAdmin{
			{UserID: 42, Permissions: []ChatAdminPermission{PermReadAllMessages, PermWrite}},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Error("Success should be true")
	}
}

func TestRemoveAdmin(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/chats/100/members/admins/42" {
			t.Errorf("path = %q, want /chats/100/members/admins/42", r.URL.Path)
		}
		writeJSON(t, w, SimpleQueryResult{Success: true})
	})

	result, err := c.RemoveAdmin(context.Background(), 100, 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Error("Success should be true")
	}
}

func TestAddMembers(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/chats/100/members" {
			t.Errorf("path = %q, want /chats/100/members", r.URL.Path)
		}
		var body UserIDsList
		readJSON(t, r, &body)
		if len(body.UserIDs) != 2 {
			t.Errorf("len(UserIDs) = %d, want 2", len(body.UserIDs))
		}
		writeJSON(t, w, SimpleQueryResult{Success: true})
	})

	result, err := c.AddMembers(context.Background(), 100, []int64{1, 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Error("Success should be true")
	}
}

func TestRemoveMember(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Query().Get("user_id") != "42" {
			t.Errorf("user_id = %q, want 42", r.URL.Query().Get("user_id"))
		}
		if r.URL.Query().Get("block") != "true" {
			t.Errorf("block = %q, want true", r.URL.Query().Get("block"))
		}
		writeJSON(t, w, SimpleQueryResult{Success: true})
	})

	result, err := c.RemoveMember(context.Background(), 100, 42, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Error("Success should be true")
	}
}

func TestSendAction(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/chats/100/actions" {
			t.Errorf("path = %q, want /chats/100/actions", r.URL.Path)
		}
		var body ActionRequestBody
		readJSON(t, r, &body)
		if body.Action != ActionTypingOn {
			t.Errorf("Action = %q, want %q", body.Action, ActionTypingOn)
		}
		writeJSON(t, w, SimpleQueryResult{Success: true})
	})

	result, err := c.SendAction(context.Background(), 100, ActionTypingOn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Error("Success should be true")
	}
}

func TestGetChatError(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeError(t, w, http.StatusNotFound, `{"code":"not.found","message":"chat not found"}`)
	})

	_, err := c.GetChat(context.Background(), 999)
	var e *Error
	if !errors.As(err, &e) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if e.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d, want 404", e.StatusCode)
	}
	if e.Op != "GetChat" {
		t.Errorf("Op = %q, want GetChat", e.Op)
	}
}

func TestLeaveChat(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/chats/100/members/me" {
			t.Errorf("path = %q, want /chats/100/members/me", r.URL.Path)
		}
		writeJSON(t, w, SimpleQueryResult{Success: true})
	})

	result, err := c.LeaveChat(context.Background(), 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Error("Success should be true")
	}
}

func TestPinMessage(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %q, want PUT", r.Method)
		}
		if r.URL.Path != "/chats/100/pin" {
			t.Errorf("path = %q, want /chats/100/pin", r.URL.Path)
		}
		writeJSON(t, w, SimpleQueryResult{Success: true})
	})

	result, err := c.PinMessage(context.Background(), 100, &PinMessageBody{MessageID: "mid-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Error("Success should be true")
	}
}

func TestDeleteChat(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/chats/100" {
			t.Errorf("path = %q, want /chats/100", r.URL.Path)
		}
		writeJSON(t, w, SimpleQueryResult{Success: true})
	})

	result, err := c.DeleteChat(context.Background(), 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Error("Success should be true")
	}
}

func TestGetMembership(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/chats/100/members/me" {
			t.Errorf("path = %q, want /chats/100/members/me", r.URL.Path)
		}
		writeJSON(t, w, ChatMember{
			UserWithPhoto: UserWithPhoto{User: User{UserID: 12345, FirstName: "Bot"}},
			IsAdmin:       true,
		})
	})

	member, err := c.GetMembership(context.Background(), 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if member.UserID != 12345 {
		t.Errorf("UserID = %d, want 12345", member.UserID)
	}
	if !member.IsAdmin {
		t.Error("IsAdmin should be true")
	}
}

func TestUnpinMessage(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/chats/100/pin" {
			t.Errorf("path = %q, want /chats/100/pin", r.URL.Path)
		}
		writeJSON(t, w, SimpleQueryResult{Success: true})
	})

	result, err := c.UnpinMessage(context.Background(), 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Error("Success should be true")
	}
}

func TestGetPinnedMessage(t *testing.T) {
	c, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/chats/100/pin" {
			t.Errorf("path = %q, want /chats/100/pin", r.URL.Path)
		}
		text := "Pinned!"
		writeJSON(t, w, GetPinnedMessageResult{
			Message: &Message{
				Timestamp: 1000,
				Body:      MessageBody{MID: "mid-pinned", Text: &text},
			},
		})
	})

	result, err := c.GetPinnedMessage(context.Background(), 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Message == nil {
		t.Fatal("Message should not be nil")
	}
	if result.Message.Body.MID != "mid-pinned" {
		t.Errorf("MID = %q, want %q", result.Message.Body.MID, "mid-pinned")
	}
}
