package maxigo

import (
	"encoding/json"
	"testing"
)

func TestButtonConstructors(t *testing.T) {
	t.Run("NewCallbackButton", func(t *testing.T) {
		b := NewCallbackButton("OK", "ok_payload")
		if b.Type != "callback" {
			t.Errorf("Type = %q, want %q", b.Type, "callback")
		}
		if b.Text != "OK" {
			t.Errorf("Text = %q, want %q", b.Text, "OK")
		}
		if b.Payload != "ok_payload" {
			t.Errorf("Payload = %q, want %q", b.Payload, "ok_payload")
		}
		if b.Intent != "" {
			t.Errorf("Intent = %q, want empty", b.Intent)
		}
	})

	t.Run("NewCallbackButtonWithIntent", func(t *testing.T) {
		b := NewCallbackButtonWithIntent("Yes", "yes", IntentPositive)
		if b.Type != "callback" {
			t.Errorf("Type = %q, want %q", b.Type, "callback")
		}
		if b.Payload != "yes" {
			t.Errorf("Payload = %q, want %q", b.Payload, "yes")
		}
		if b.Intent != IntentPositive {
			t.Errorf("Intent = %q, want %q", b.Intent, IntentPositive)
		}
	})

	t.Run("NewLinkButton", func(t *testing.T) {
		b := NewLinkButton("Open", "https://example.com")
		if b.Type != "link" {
			t.Errorf("Type = %q, want %q", b.Type, "link")
		}
		if b.Text != "Open" {
			t.Errorf("Text = %q, want %q", b.Text, "Open")
		}
		if b.URL != "https://example.com" {
			t.Errorf("URL = %q, want %q", b.URL, "https://example.com")
		}
	})

	t.Run("NewRequestContactButton", func(t *testing.T) {
		b := NewRequestContactButton("Share contact")
		if b.Type != "request_contact" {
			t.Errorf("Type = %q, want %q", b.Type, "request_contact")
		}
		if b.Text != "Share contact" {
			t.Errorf("Text = %q, want %q", b.Text, "Share contact")
		}
	})

	t.Run("NewRequestGeoLocationButton", func(t *testing.T) {
		b := NewRequestGeoLocationButton("Send location", true)
		if b.Type != "request_geo_location" {
			t.Errorf("Type = %q, want %q", b.Type, "request_geo_location")
		}
		if b.Text != "Send location" {
			t.Errorf("Text = %q, want %q", b.Text, "Send location")
		}
		if !b.Quick {
			t.Error("Quick = false, want true")
		}
	})

	t.Run("NewChatButton", func(t *testing.T) {
		b := NewChatButton("Create chat", "My Chat")
		if b.Type != "chat" {
			t.Errorf("Type = %q, want %q", b.Type, "chat")
		}
		if b.Text != "Create chat" {
			t.Errorf("Text = %q, want %q", b.Text, "Create chat")
		}
		if b.ChatTitle != "My Chat" {
			t.Errorf("ChatTitle = %q, want %q", b.ChatTitle, "My Chat")
		}
	})

	t.Run("NewMessageButton", func(t *testing.T) {
		b := NewMessageButton("Send")
		if b.Type != "message" {
			t.Errorf("Type = %q, want %q", b.Type, "message")
		}
		if b.Text != "Send" {
			t.Errorf("Text = %q, want %q", b.Text, "Send")
		}
	})
}

func TestNewPhotoAttachment(t *testing.T) {
	a := NewPhotoAttachment(PhotoAttachmentRequestPayload{URL: Some("https://example.com/photo.jpg")})
	if a.Type != "image" {
		t.Errorf("Type = %q, want %q", a.Type, "image")
	}
	p, ok := a.Payload.(PhotoAttachmentRequestPayload)
	if !ok {
		t.Fatalf("Payload type = %T, want PhotoAttachmentRequestPayload", a.Payload)
	}
	if !p.URL.Set || p.URL.Value != "https://example.com/photo.jpg" {
		t.Errorf("Payload.URL = %v, want %q", p.URL, "https://example.com/photo.jpg")
	}
}

func TestNewVideoAttachment(t *testing.T) {
	a := NewVideoAttachment(UploadedInfo{Token: "video-tok"})
	if a.Type != "video" {
		t.Errorf("Type = %q, want %q", a.Type, "video")
	}
	p, ok := a.Payload.(UploadedInfo)
	if !ok {
		t.Fatalf("Payload type = %T, want UploadedInfo", a.Payload)
	}
	if p.Token != "video-tok" {
		t.Errorf("Payload.Token = %q, want %q", p.Token, "video-tok")
	}
}

func TestNewAudioAttachment(t *testing.T) {
	a := NewAudioAttachment(UploadedInfo{Token: "audio-tok"})
	if a.Type != "audio" {
		t.Errorf("Type = %q, want %q", a.Type, "audio")
	}
	p, ok := a.Payload.(UploadedInfo)
	if !ok {
		t.Fatalf("Payload type = %T, want UploadedInfo", a.Payload)
	}
	if p.Token != "audio-tok" {
		t.Errorf("Payload.Token = %q, want %q", p.Token, "audio-tok")
	}
}

func TestNewFileAttachment(t *testing.T) {
	a := NewFileAttachment(UploadedInfo{Token: "file-tok"})
	if a.Type != "file" {
		t.Errorf("Type = %q, want %q", a.Type, "file")
	}
	p, ok := a.Payload.(UploadedInfo)
	if !ok {
		t.Fatalf("Payload type = %T, want UploadedInfo", a.Payload)
	}
	if p.Token != "file-tok" {
		t.Errorf("Payload.Token = %q, want %q", p.Token, "file-tok")
	}
}

func TestNewStickerAttachment(t *testing.T) {
	a := NewStickerAttachment(StickerAttachmentRequestPayload{Code: "smile"})
	if a.Type != "sticker" {
		t.Errorf("Type = %q, want %q", a.Type, "sticker")
	}
	p, ok := a.Payload.(StickerAttachmentRequestPayload)
	if !ok {
		t.Fatalf("Payload type = %T, want StickerAttachmentRequestPayload", a.Payload)
	}
	if p.Code != "smile" {
		t.Errorf("Payload.Code = %q, want %q", p.Code, "smile")
	}
}

func TestNewContactAttachment(t *testing.T) {
	a := NewContactAttachment(ContactAttachmentRequestPayload{Name: Some("John")})
	if a.Type != "contact" {
		t.Errorf("Type = %q, want %q", a.Type, "contact")
	}
	p, ok := a.Payload.(ContactAttachmentRequestPayload)
	if !ok {
		t.Fatalf("Payload type = %T, want ContactAttachmentRequestPayload", a.Payload)
	}
	if !p.Name.Set || p.Name.Value != "John" {
		t.Errorf("Payload.Name = %v, want %q", p.Name, "John")
	}
}

func TestNewShareAttachment(t *testing.T) {
	a := NewShareAttachment(ShareAttachmentPayload{URL: Some("https://example.com")})
	if a.Type != "share" {
		t.Errorf("Type = %q, want %q", a.Type, "share")
	}
	p, ok := a.Payload.(ShareAttachmentPayload)
	if !ok {
		t.Fatalf("Payload type = %T, want ShareAttachmentPayload", a.Payload)
	}
	if !p.URL.Set || p.URL.Value != "https://example.com" {
		t.Errorf("Payload.URL = %v, want %q", p.URL, "https://example.com")
	}
}

func TestNewInlineKeyboardAttachment(t *testing.T) {
	buttons := [][]Button{
		{{Type: "callback", Text: "OK", Payload: "ok"}},
	}
	a := NewInlineKeyboardAttachment(buttons)
	if a.Type != "inline_keyboard" {
		t.Errorf("Type = %q, want %q", a.Type, "inline_keyboard")
	}
	p, ok := a.Payload.(Keyboard)
	if !ok {
		t.Fatalf("Payload type = %T, want Keyboard", a.Payload)
	}
	if len(p.Buttons) != 1 || len(p.Buttons[0]) != 1 {
		t.Fatalf("Buttons shape = %v, want 1x1", p.Buttons)
	}
	if p.Buttons[0][0].Text != "OK" {
		t.Errorf("Button.Text = %q, want %q", p.Buttons[0][0].Text, "OK")
	}
}

func TestNewLocationAttachment(t *testing.T) {
	a := NewLocationAttachment(55.7558, 37.6173)
	if a.Type != "location" {
		t.Errorf("Type = %q, want %q", a.Type, "location")
	}
	if a.Latitude != 55.7558 {
		t.Errorf("Latitude = %f, want %f", a.Latitude, 55.7558)
	}
	if a.Longitude != 37.6173 {
		t.Errorf("Longitude = %f, want %f", a.Longitude, 37.6173)
	}
}

func TestParseAttachments(t *testing.T) {
	tests := []struct {
		name   string
		input  any
		check  func(t *testing.T, att Attachment)
	}{
		{
			name: "image",
			input: PhotoAttachment{
				AttachmentType: AttachmentType{Type: "image"},
				Payload:        PhotoAttachmentPayload{PhotoID: 42, Token: "tok", URL: "https://img.example.com/1.jpg"},
			},
			check: func(t *testing.T, att Attachment) {
				a, ok := att.(*PhotoAttachment)
				if !ok {
					t.Fatalf("type = %T, want *PhotoAttachment", att)
				}
				if a.Payload.PhotoID != 42 {
					t.Errorf("PhotoID = %d, want 42", a.Payload.PhotoID)
				}
				if a.Payload.URL != "https://img.example.com/1.jpg" {
					t.Errorf("URL = %q, want %q", a.Payload.URL, "https://img.example.com/1.jpg")
				}
			},
		},
		{
			name: "video",
			input: VideoAttachment{
				AttachmentType: AttachmentType{Type: "video"},
				Payload:        MediaAttachmentPayload{URL: "https://video.example.com/v.mp4", Token: "vtok"},
			},
			check: func(t *testing.T, att Attachment) {
				a, ok := att.(*VideoAttachment)
				if !ok {
					t.Fatalf("type = %T, want *VideoAttachment", att)
				}
				if a.Payload.URL != "https://video.example.com/v.mp4" {
					t.Errorf("URL = %q", a.Payload.URL)
				}
			},
		},
		{
			name: "audio",
			input: AudioAttachment{
				AttachmentType: AttachmentType{Type: "audio"},
				Payload:        MediaAttachmentPayload{URL: "https://audio.example.com/a.ogg", Token: "atok"},
			},
			check: func(t *testing.T, att Attachment) {
				a, ok := att.(*AudioAttachment)
				if !ok {
					t.Fatalf("type = %T, want *AudioAttachment", att)
				}
				if a.Payload.Token != "atok" {
					t.Errorf("Token = %q, want %q", a.Payload.Token, "atok")
				}
			},
		},
		{
			name: "file",
			input: FileAttachment{
				AttachmentType: AttachmentType{Type: "file"},
				Payload:        FileAttachmentPayload{URL: "https://files.example.com/f.pdf", Token: "ftok"},
				Filename:       "report.pdf",
				Size:           1024,
			},
			check: func(t *testing.T, att Attachment) {
				a, ok := att.(*FileAttachment)
				if !ok {
					t.Fatalf("type = %T, want *FileAttachment", att)
				}
				if a.Filename != "report.pdf" {
					t.Errorf("Filename = %q, want %q", a.Filename, "report.pdf")
				}
				if a.Size != 1024 {
					t.Errorf("Size = %d, want 1024", a.Size)
				}
			},
		},
		{
			name: "sticker",
			input: StickerAttachment{
				AttachmentType: AttachmentType{Type: "sticker"},
				Payload:        StickerAttachmentPayload{URL: "https://stk.example.com/s.png", Code: "smile"},
				Width:          128,
				Height:         128,
			},
			check: func(t *testing.T, att Attachment) {
				a, ok := att.(*StickerAttachment)
				if !ok {
					t.Fatalf("type = %T, want *StickerAttachment", att)
				}
				if a.Payload.Code != "smile" {
					t.Errorf("Code = %q, want %q", a.Payload.Code, "smile")
				}
				if a.Width != 128 || a.Height != 128 {
					t.Errorf("Size = %dx%d, want 128x128", a.Width, a.Height)
				}
			},
		},
		{
			name: "contact",
			input: ContactAttachment{
				AttachmentType: AttachmentType{Type: "contact"},
				Payload: ContactAttachmentPayload{
					MaxInfo: &User{UserID: 99, FirstName: "John"},
				},
			},
			check: func(t *testing.T, att Attachment) {
				a, ok := att.(*ContactAttachment)
				if !ok {
					t.Fatalf("type = %T, want *ContactAttachment", att)
				}
				if a.Payload.MaxInfo == nil {
					t.Fatal("MaxInfo is nil")
				}
				if a.Payload.MaxInfo.FirstName != "John" {
					t.Errorf("FirstName = %q, want %q", a.Payload.MaxInfo.FirstName, "John")
				}
			},
		},
		{
			name: "share",
			input: ShareAttachment{
				AttachmentType: AttachmentType{Type: "share"},
				Payload:        ShareAttachmentPayload{URL: Some("https://example.com")},
			},
			check: func(t *testing.T, att Attachment) {
				a, ok := att.(*ShareAttachment)
				if !ok {
					t.Fatalf("type = %T, want *ShareAttachment", att)
				}
				if !a.Payload.URL.Set || a.Payload.URL.Value != "https://example.com" {
					t.Errorf("URL = %v, want %q", a.Payload.URL, "https://example.com")
				}
			},
		},
		{
			name: "location",
			input: LocationAttachment{
				AttachmentType: AttachmentType{Type: "location"},
				Latitude:       55.7558,
				Longitude:      37.6173,
			},
			check: func(t *testing.T, att Attachment) {
				a, ok := att.(*LocationAttachment)
				if !ok {
					t.Fatalf("type = %T, want *LocationAttachment", att)
				}
				if a.Latitude != 55.7558 {
					t.Errorf("Latitude = %f, want 55.7558", a.Latitude)
				}
				if a.Longitude != 37.6173 {
					t.Errorf("Longitude = %f, want 37.6173", a.Longitude)
				}
			},
		},
		{
			name: "data",
			input: DataAttachment{
				AttachmentType: AttachmentType{Type: "data"},
				Data:           `{"key":"value"}`,
			},
			check: func(t *testing.T, att Attachment) {
				a, ok := att.(*DataAttachment)
				if !ok {
					t.Fatalf("type = %T, want *DataAttachment", att)
				}
				if a.Data != `{"key":"value"}` {
					t.Errorf("Data = %q", a.Data)
				}
			},
		},
		{
			name: "inline_keyboard",
			input: InlineKeyboardAttachment{
				AttachmentType: AttachmentType{Type: "inline_keyboard"},
				Payload: Keyboard{
					Buttons: [][]Button{
						{NewCallbackButton("OK", "ok")},
					},
				},
			},
			check: func(t *testing.T, att Attachment) {
				a, ok := att.(*InlineKeyboardAttachment)
				if !ok {
					t.Fatalf("type = %T, want *InlineKeyboardAttachment", att)
				}
				if len(a.Payload.Buttons) != 1 || len(a.Payload.Buttons[0]) != 1 {
					t.Fatalf("Buttons shape = %v, want 1x1", a.Payload.Buttons)
				}
				if a.Payload.Buttons[0][0].Text != "OK" {
					t.Errorf("Button.Text = %q, want %q", a.Payload.Buttons[0][0].Text, "OK")
				}
			},
		},
		{
			name: "reply_keyboard",
			input: ReplyKeyboardAttachment{
				AttachmentType: AttachmentType{Type: "reply_keyboard"},
				Buttons: [][]ReplyButton{
					{{Text: "Yes"}, {Text: "No"}},
				},
			},
			check: func(t *testing.T, att Attachment) {
				a, ok := att.(*ReplyKeyboardAttachment)
				if !ok {
					t.Fatalf("type = %T, want *ReplyKeyboardAttachment", att)
				}
				if len(a.Buttons) != 1 || len(a.Buttons[0]) != 2 {
					t.Fatalf("Buttons shape = %v, want 1x2", a.Buttons)
				}
				if a.Buttons[0][0].Text != "Yes" {
					t.Errorf("Button[0][0].Text = %q, want %q", a.Buttons[0][0].Text, "Yes")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mb := MessageBody{
				Attachments: []json.RawMessage{mustMarshal(tt.input)},
			}

			attachments, err := mb.ParseAttachments()
			if err != nil {
				t.Fatalf("ParseAttachments() error: %v", err)
			}
			if len(attachments) != 1 {
				t.Fatalf("len = %d, want 1", len(attachments))
			}

			if attachments[0].GetType() != tt.name {
				t.Errorf("GetType() = %q, want %q", attachments[0].GetType(), tt.name)
			}

			tt.check(t, attachments[0])
		})
	}

	t.Run("unknown type is skipped", func(t *testing.T) {
		mb := MessageBody{
			Attachments: []json.RawMessage{
				json.RawMessage(`{"type":"future_type","data":"something"}`),
			},
		}
		attachments, err := mb.ParseAttachments()
		if err != nil {
			t.Fatalf("ParseAttachments() error: %v", err)
		}
		if len(attachments) != 0 {
			t.Errorf("len = %d, want 0", len(attachments))
		}
	})

	t.Run("empty attachments returns nil", func(t *testing.T) {
		mb := MessageBody{}
		attachments, err := mb.ParseAttachments()
		if err != nil {
			t.Fatalf("ParseAttachments() error: %v", err)
		}
		if attachments != nil {
			t.Errorf("got %v, want nil", attachments)
		}
	})

	t.Run("mixed known and unknown types", func(t *testing.T) {
		mb := MessageBody{
			Attachments: []json.RawMessage{
				mustMarshal(LocationAttachment{
					AttachmentType: AttachmentType{Type: "location"},
					Latitude:       1.0,
					Longitude:      2.0,
				}),
				json.RawMessage(`{"type":"unknown_v2"}`),
				mustMarshal(DataAttachment{
					AttachmentType: AttachmentType{Type: "data"},
					Data:           "test",
				}),
			},
		}
		attachments, err := mb.ParseAttachments()
		if err != nil {
			t.Fatalf("ParseAttachments() error: %v", err)
		}
		if len(attachments) != 2 {
			t.Fatalf("len = %d, want 2", len(attachments))
		}
		if _, ok := attachments[0].(*LocationAttachment); !ok {
			t.Errorf("[0] type = %T, want *LocationAttachment", attachments[0])
		}
		if _, ok := attachments[1].(*DataAttachment); !ok {
			t.Errorf("[1] type = %T, want *DataAttachment", attachments[1])
		}
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		mb := MessageBody{
			Attachments: []json.RawMessage{
				json.RawMessage(`{invalid`),
			},
		}
		_, err := mb.ParseAttachments()
		if err == nil {
			t.Fatal("expected error for invalid JSON")
		}
	})
}
