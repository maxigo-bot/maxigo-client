package maxigo

import (
	"encoding/json"
	"strings"
	"testing"
)

// helper is a struct used across multiple Optional tests.
type optionalTestStruct struct {
	Name  Optional[string] `json:"name,omitzero"`
	Flag  Optional[bool]   `json:"flag,omitzero"`
	Count Optional[int64]  `json:"count,omitzero"`
}

func TestOptionalSome(t *testing.T) {
	s := Some("hello")
	if !s.Set {
		t.Error("Some: Set = false, want true")
	}
	if s.Value != "hello" {
		t.Errorf("Some: Value = %q, want %q", s.Value, "hello")
	}
}

func TestOptionalIsZero(t *testing.T) {
	var unset Optional[string]
	if !unset.IsZero() {
		t.Error("unset Optional: IsZero() = false, want true")
	}

	set := Some("")
	if set.IsZero() {
		t.Error("set Optional (zero value): IsZero() = true, want false")
	}
}

// TestOptionalStringMarshal verifies three states for Optional[string]:
//   - unset → field omitted
//   - set to "" → field present as ""
//   - set to "hello" → field present as "hello"
func TestOptionalStringMarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    optionalTestStruct
		wantJSON string
	}{
		{
			name:     "string unset — field omitted",
			input:    optionalTestStruct{},
			wantJSON: `{}`,
		},
		{
			name:     "string set to empty — field present as empty string",
			input:    optionalTestStruct{Name: Some("")},
			wantJSON: `{"name":""}`,
		},
		{
			name:     "string set to value — field present with value",
			input:    optionalTestStruct{Name: Some("hello")},
			wantJSON: `{"name":"hello"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("Marshal error: %v", err)
			}
			if string(got) != tt.wantJSON {
				t.Errorf("Marshal =\n  %s\nwant\n  %s", got, tt.wantJSON)
			}
		})
	}
}

// TestOptionalBoolMarshal verifies three states for Optional[bool]:
//   - unset → field omitted
//   - set to false → field present as false
//   - set to true → field present as true
func TestOptionalBoolMarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    optionalTestStruct
		wantJSON string
	}{
		{
			name:     "bool unset — field omitted",
			input:    optionalTestStruct{},
			wantJSON: `{}`,
		},
		{
			name:     "bool set to false — field present as false",
			input:    optionalTestStruct{Flag: Some(false)},
			wantJSON: `{"flag":false}`,
		},
		{
			name:     "bool set to true — field present as true",
			input:    optionalTestStruct{Flag: Some(true)},
			wantJSON: `{"flag":true}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("Marshal error: %v", err)
			}
			if string(got) != tt.wantJSON {
				t.Errorf("Marshal =\n  %s\nwant\n  %s", got, tt.wantJSON)
			}
		})
	}
}

// TestOptionalInt64Marshal verifies three states for Optional[int64]:
//   - unset → field omitted
//   - set to 0 → field present as 0
//   - set to 42 → field present as 42
func TestOptionalInt64Marshal(t *testing.T) {
	tests := []struct {
		name     string
		input    optionalTestStruct
		wantJSON string
	}{
		{
			name:     "int64 unset — field omitted",
			input:    optionalTestStruct{},
			wantJSON: `{}`,
		},
		{
			name:     "int64 set to 0 — field present as 0",
			input:    optionalTestStruct{Count: Some(int64(0))},
			wantJSON: `{"count":0}`,
		},
		{
			name:     "int64 set to 42 — field present as 42",
			input:    optionalTestStruct{Count: Some(int64(42))},
			wantJSON: `{"count":42}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("Marshal error: %v", err)
			}
			if string(got) != tt.wantJSON {
				t.Errorf("Marshal =\n  %s\nwant\n  %s", got, tt.wantJSON)
			}
		})
	}
}

// TestOptionalAllFieldsMarshal verifies that all three fields set together
// produce the correct combined JSON output.
func TestOptionalAllFieldsMarshal(t *testing.T) {
	v := optionalTestStruct{
		Name:  Some("test"),
		Flag:  Some(true),
		Count: Some(int64(99)),
	}

	got, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	want := `{"name":"test","flag":true,"count":99}`
	if string(got) != want {
		t.Errorf("Marshal =\n  %s\nwant\n  %s", got, want)
	}
}

// TestOptionalUnmarshalPresent verifies that a field present in JSON
// results in Set=true with the correct value.
func TestOptionalUnmarshalPresent(t *testing.T) {
	input := `{"name":"world","flag":true,"count":7}`

	var v optionalTestStruct
	if err := json.Unmarshal([]byte(input), &v); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if !v.Name.Set || v.Name.Value != "world" {
		t.Errorf("Name = %+v, want {Value:world Set:true}", v.Name)
	}
	if !v.Flag.Set || v.Flag.Value != true {
		t.Errorf("Flag = %+v, want {Value:true Set:true}", v.Flag)
	}
	if !v.Count.Set || v.Count.Value != 7 {
		t.Errorf("Count = %+v, want {Value:7 Set:true}", v.Count)
	}
}

// TestOptionalUnmarshalAbsent verifies that a field absent from JSON
// results in Set=false.
func TestOptionalUnmarshalAbsent(t *testing.T) {
	input := `{}`

	var v optionalTestStruct
	if err := json.Unmarshal([]byte(input), &v); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if v.Name.Set {
		t.Errorf("Name.Set = true, want false (field was absent)")
	}
	if v.Flag.Set {
		t.Errorf("Flag.Set = true, want false (field was absent)")
	}
	if v.Count.Set {
		t.Errorf("Count.Set = true, want false (field was absent)")
	}
}

// TestOptionalUnmarshalNull verifies that JSON null is treated as unset.
// This prevents accidental overwrites when round-tripping API responses
// through request types (e.g. null description would not become "").
func TestOptionalUnmarshalNull(t *testing.T) {
	input := `{"name":null,"flag":null,"count":null}`

	var v optionalTestStruct
	if err := json.Unmarshal([]byte(input), &v); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	// JSON null should leave Optional as unset (Set=false).
	if v.Name.Set {
		t.Error("Name.Set = true, want false (null should be unset)")
	}
	if v.Flag.Set {
		t.Error("Flag.Set = true, want false (null should be unset)")
	}
	if v.Count.Set {
		t.Error("Count.Set = true, want false (null should be unset)")
	}
}

// TestOptionalUnmarshalZeroValues verifies that zero values in JSON
// are correctly unmarshaled as Set=true with the zero value.
func TestOptionalUnmarshalZeroValues(t *testing.T) {
	input := `{"name":"","flag":false,"count":0}`

	var v optionalTestStruct
	if err := json.Unmarshal([]byte(input), &v); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if !v.Name.Set || v.Name.Value != "" {
		t.Errorf("Name = %+v, want {Value: Set:true}", v.Name)
	}
	if !v.Flag.Set || v.Flag.Value != false {
		t.Errorf("Flag = %+v, want {Value:false Set:true}", v.Flag)
	}
	if !v.Count.Set || v.Count.Value != 0 {
		t.Errorf("Count = %+v, want {Value:0 Set:true}", v.Count)
	}
}

// TestOptionalMarshalUnset verifies that MarshalJSON of an unset Optional
// produces null (used when Optional is marshaled directly, not via omitzero).
func TestOptionalMarshalUnset(t *testing.T) {
	var o Optional[string]
	got, err := json.Marshal(o)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	if string(got) != "null" {
		t.Errorf("Marshal unset = %s, want null", got)
	}
}

// TestOptionalNewMessageBodyRoundTrip verifies a full marshal/unmarshal cycle
// with the real NewMessageBody type used by the API.
func TestOptionalNewMessageBodyRoundTrip(t *testing.T) {
	original := NewMessageBody{
		Text:   Some("Hello, Max!"),
		Notify: Some(false),
		Format: Some(FormatMarkdown),
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	// Verify JSON contains expected fields.
	jsonStr := string(data)
	for _, want := range []string{`"text":"Hello, Max!"`, `"notify":false`, `"format":"markdown"`} {
		if !strings.Contains(jsonStr, want) {
			t.Errorf("JSON %s does not contain %s", jsonStr, want)
		}
	}

	// Verify unset fields (attachments, link) are omitted.
	for _, absent := range []string{`"attachments"`, `"link"`} {
		if strings.Contains(jsonStr, absent) {
			t.Errorf("JSON %s should not contain %s (unset field)", jsonStr, absent)
		}
	}

	// Unmarshal back and verify equality.
	var decoded NewMessageBody
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if !decoded.Text.Set || decoded.Text.Value != "Hello, Max!" {
		t.Errorf("Text = %+v, want {Value:Hello, Max! Set:true}", decoded.Text)
	}
	if !decoded.Notify.Set || decoded.Notify.Value != false {
		t.Errorf("Notify = %+v, want {Value:false Set:true}", decoded.Notify)
	}
	if !decoded.Format.Set || decoded.Format.Value != FormatMarkdown {
		t.Errorf("Format = %+v, want {Value:markdown Set:true}", decoded.Format)
	}
}

// TestOptionalNewMessageBodyTextOnly verifies that a NewMessageBody with only
// text set omits all other optional fields.
func TestOptionalNewMessageBodyTextOnly(t *testing.T) {
	body := NewMessageBody{
		Text: Some("simple text"),
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	want := `{"text":"simple text"}`
	if string(data) != want {
		t.Errorf("Marshal =\n  %s\nwant\n  %s", data, want)
	}
}

// TestOptionalNewMessageBodyEmpty verifies that a completely empty NewMessageBody
// marshals to just {}.
func TestOptionalNewMessageBodyEmpty(t *testing.T) {
	body := NewMessageBody{}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	want := `{}`
	if string(data) != want {
		t.Errorf("Marshal =\n  %s\nwant\n  %s", data, want)
	}
}

// TestOptionalBotPatchMarshal verifies BotPatch marshal with Optional fields.
func TestOptionalBotPatchMarshal(t *testing.T) {
	patch := BotPatch{
		FirstName:   Some("NewBot"),
		Description: Some(""),
	}

	data, err := json.Marshal(patch)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	jsonStr := string(data)

	// first_name should be present.
	if !strings.Contains(jsonStr, `"first_name":"NewBot"`) {
		t.Errorf("JSON %s missing first_name", jsonStr)
	}

	// description should be present even as empty string.
	if !strings.Contains(jsonStr, `"description":""`) {
		t.Errorf("JSON %s missing description (empty string)", jsonStr)
	}

	// name (unset) should be omitted.
	if strings.Contains(jsonStr, `"name"`) {
		t.Errorf("JSON %s should not contain unset name field", jsonStr)
	}
}

// TestOptionalChatPatchMarshal verifies ChatPatch with mixed Optional states.
func TestOptionalChatPatchMarshal(t *testing.T) {
	patch := ChatPatch{
		Title:  Some("New Title"),
		Notify: Some(true),
	}

	data, err := json.Marshal(patch)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	jsonStr := string(data)
	if !strings.Contains(jsonStr, `"title":"New Title"`) {
		t.Errorf("JSON %s missing title", jsonStr)
	}
	if !strings.Contains(jsonStr, `"notify":true`) {
		t.Errorf("JSON %s missing notify", jsonStr)
	}
	// pin (unset) should be omitted.
	if strings.Contains(jsonStr, `"pin"`) {
		t.Errorf("JSON %s should not contain unset pin field", jsonStr)
	}
}

// TestOptionalPinMessageBodyMarshal verifies PinMessageBody with OptBool.
func TestOptionalPinMessageBodyMarshal(t *testing.T) {
	t.Run("notify unset", func(t *testing.T) {
		body := PinMessageBody{MessageID: "msg-123"}
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}
		want := `{"message_id":"msg-123"}`
		if string(data) != want {
			t.Errorf("Marshal = %s, want %s", data, want)
		}
	})

	t.Run("notify false", func(t *testing.T) {
		body := PinMessageBody{MessageID: "msg-123", Notify: Some(false)}
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}
		jsonStr := string(data)
		if !strings.Contains(jsonStr, `"notify":false`) {
			t.Errorf("JSON %s missing notify:false", jsonStr)
		}
	})
}

// TestOptionalCallbackAnswerMarshal verifies CallbackAnswer with OptString.
func TestOptionalCallbackAnswerMarshal(t *testing.T) {
	t.Run("notification unset", func(t *testing.T) {
		answer := CallbackAnswer{}
		data, err := json.Marshal(answer)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}
		if strings.Contains(string(data), `"notification"`) {
			t.Errorf("JSON %s should not contain unset notification", data)
		}
	})

	t.Run("notification set to empty", func(t *testing.T) {
		answer := CallbackAnswer{Notification: Some("")}
		data, err := json.Marshal(answer)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}
		if !strings.Contains(string(data), `"notification":""`) {
			t.Errorf("JSON %s missing notification empty string", data)
		}
	})
}

// TestOptionalButtonMarshal verifies Button with multiple Optional fields.
func TestOptionalButtonMarshal(t *testing.T) {
	btn := Button{
		Type:            "chat",
		Text:            "Create Chat",
		ChatDescription: Some("Welcome!"),
		// start_payload and uuid unset — should be omitted.
	}

	data, err := json.Marshal(btn)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	jsonStr := string(data)
	if !strings.Contains(jsonStr, `"chat_description":"Welcome!"`) {
		t.Errorf("JSON %s missing chat_description", jsonStr)
	}
	if strings.Contains(jsonStr, `"start_payload"`) {
		t.Errorf("JSON %s should not contain unset start_payload", jsonStr)
	}
	if strings.Contains(jsonStr, `"uuid"`) {
		t.Errorf("JSON %s should not contain unset uuid", jsonStr)
	}
}

// TestOptionalTypeAliases verifies that OptString, OptBool, OptInt64 aliases
// work identically to Optional[string], Optional[bool], Optional[int64].
func TestOptionalTypeAliases(t *testing.T) {
	s := Some("test")
	os := OptString(s)
	if !os.Set || os.Value != "test" {
		t.Errorf("OptString = %+v, want {Value:test Set:true}", os)
	}

	b := Some(true)
	ob := OptBool(b)
	if !ob.Set || ob.Value != true {
		t.Errorf("OptBool = %+v, want {Value:true Set:true}", ob)
	}

	i := Some(int64(42))
	oi := OptInt64(i)
	if !oi.Set || oi.Value != 42 {
		t.Errorf("OptInt64 = %+v, want {Value:42 Set:true}", oi)
	}
}
