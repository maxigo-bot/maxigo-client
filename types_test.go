package maxigo

import "testing"

func TestNewPhotoAttachment(t *testing.T) {
	url := "https://example.com/photo.jpg"
	a := NewPhotoAttachment(PhotoAttachmentRequestPayload{URL: &url})
	if a.Type != "image" {
		t.Errorf("Type = %q, want %q", a.Type, "image")
	}
	p, ok := a.Payload.(PhotoAttachmentRequestPayload)
	if !ok {
		t.Fatalf("Payload type = %T, want PhotoAttachmentRequestPayload", a.Payload)
	}
	if p.URL == nil || *p.URL != url {
		t.Errorf("Payload.URL = %v, want %q", p.URL, url)
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
	name := "John"
	a := NewContactAttachment(ContactAttachmentRequestPayload{Name: &name})
	if a.Type != "contact" {
		t.Errorf("Type = %q, want %q", a.Type, "contact")
	}
	p, ok := a.Payload.(ContactAttachmentRequestPayload)
	if !ok {
		t.Fatalf("Payload type = %T, want ContactAttachmentRequestPayload", a.Payload)
	}
	if p.Name == nil || *p.Name != "John" {
		t.Errorf("Payload.Name = %v, want %q", p.Name, "John")
	}
}

func TestNewShareAttachment(t *testing.T) {
	url := "https://example.com"
	a := NewShareAttachment(ShareAttachmentPayload{URL: &url})
	if a.Type != "share" {
		t.Errorf("Type = %q, want %q", a.Type, "share")
	}
	p, ok := a.Payload.(ShareAttachmentPayload)
	if !ok {
		t.Fatalf("Payload type = %T, want ShareAttachmentPayload", a.Payload)
	}
	if p.URL == nil || *p.URL != url {
		t.Errorf("Payload.URL = %v, want %q", p.URL, url)
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
