package maxigo_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/maxigo-bot/maxigo-client"
)

func ExampleNew() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(maxigo.BotInfo{
			UserWithPhoto: maxigo.UserWithPhoto{
				User: maxigo.User{UserID: 1, FirstName: "TestBot", IsBot: true},
			},
		})
	}))
	defer srv.Close()

	client, err := maxigo.New("test-token", maxigo.WithBaseURL(srv.URL))
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	bot, err := client.GetBot(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(bot.FirstName)
	// Output: TestBot
}

func ExampleClient_SendMessage() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		resp := struct {
			Message maxigo.Message `json:"message"`
		}{
			Message: maxigo.Message{
				Body: maxigo.MessageBody{MID: "mid-1", Text: strPtr("Hello!")},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client, _ := maxigo.New("test-token", maxigo.WithBaseURL(srv.URL))
	msg, err := client.SendMessage(context.Background(), 12345, &maxigo.NewMessageBody{
		Text: maxigo.Some("Hello!"),
	})
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(*msg.Body.Text)
	// Output: Hello!
}

func ExampleClient_AnswerCallback() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(maxigo.SimpleQueryResult{Success: true})
	}))
	defer srv.Close()

	client, _ := maxigo.New("test-token", maxigo.WithBaseURL(srv.URL))
	result, err := client.AnswerCallback(context.Background(), "cb-123", &maxigo.CallbackAnswer{
		Notification: maxigo.Some("Done!"),
	})
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(result.Success)
	// Output: true
}

func ExampleClient_UploadPhoto() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/uploads":
			// GetUploadURL — return the same server
			_ = json.NewEncoder(w).Encode(maxigo.UploadEndpoint{
				URL: "http://" + r.Host + "/do-upload",
			})
		default:
			// Actual upload
			_ = json.NewEncoder(w).Encode(maxigo.PhotoTokens{
				Photos: map[string]maxigo.PhotoToken{"default": {Token: "photo-tok"}},
			})
		}
	}))
	defer srv.Close()

	client, _ := maxigo.New("test-token", maxigo.WithBaseURL(srv.URL))
	tokens, err := client.UploadPhoto(context.Background(), "photo.jpg", strings.NewReader("image data"))
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(tokens.Photos["default"].Token)
	// Output: photo-tok
}

func ExampleNewCallbackButton() {
	btn := maxigo.NewCallbackButton("OK", "confirm")
	fmt.Printf("type=%s text=%s payload=%s\n", btn.Type, btn.Text, btn.Payload)
	// Output: type=callback text=OK payload=confirm
}

func ExampleNewLinkButton() {
	btn := maxigo.NewLinkButton("Open", "https://example.com")
	fmt.Printf("type=%s url=%s\n", btn.Type, btn.URL)
	// Output: type=link url=https://example.com
}

func ExampleSome() {
	opt := maxigo.Some("hello")
	fmt.Printf("set=%t value=%s\n", opt.Set, opt.Value)

	var empty maxigo.OptString
	fmt.Printf("set=%t value=%q\n", empty.Set, empty.Value)
	// Output:
	// set=true value=hello
	// set=false value=""
}

func ExampleNewInlineKeyboardAttachment() {
	kb := maxigo.NewInlineKeyboardAttachment([][]maxigo.Button{
		{
			maxigo.NewCallbackButton("Yes", "yes"),
			maxigo.NewCallbackButton("No", "no"),
		},
	})
	fmt.Println(kb.Type)
	// Output: inline_keyboard
}

func strPtr(s string) *string { return &s }
