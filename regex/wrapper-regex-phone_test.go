package regex

import (
	"encoding/json"
	"testing"

	"github.com/zealsprince/wrappers"
)

func TestWrapperRegexPhone_Wrap(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		discard     bool
		want        string
		wantDiscard bool
		wantError   bool
	}{
		{
			name:        "Valid phone number with plus sign",
			input:       "+1234567890",
			discard:     false,
			want:        "+1234567890",
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Valid phone number without plus sign",
			input:       "1234567890",
			discard:     false,
			want:        "1234567890",
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Valid local phone number",
			input:       "01234567890",
			discard:     false,
			want:        "01234567890",
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Invalid phone number",
			input:       "0000000123456789",
			discard:     false,
			want:        "", // Value remains unchanged
			wantDiscard: true,
			wantError:   true,
		},
		{
			name:        "Invalid phone number discard",
			input:       "0000000123456789",
			discard:     true,
			want:        "",
			wantDiscard: true,
			wantError:   false,
		},
		{
			name:        "Nil input",
			input:       nil,
			discard:     false,
			want:        "", // Value remains unchanged
			wantDiscard: true,
			wantError:   true,
		},
		{
			name:        "Input as WrapperProvider (WrapperRegexPhone)",
			input:       wrappers.NewWithValueDiscard[*WrapperRegexPhone]("+1234567890"),
			discard:     false,
			want:        "+1234567890",
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Invalid type input",
			input:       12345,
			discard:     false,
			want:        "", // Value remains unchanged
			wantDiscard: true,
			wantError:   true,
		},
		{
			name:        "Invalid type input with discard",
			input:       12345,
			discard:     true,
			want:        "",
			wantDiscard: true,
			wantError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the phone wrapper
			phoneWrapper := *wrappers.New[*WrapperRegexPhone]()

			// Wrap the input
			err := phoneWrapper.Wrap(tt.input, tt.discard)
			if (err != nil) != tt.wantError {
				t.Errorf("Wrap() error = %v, wantError %v", err, tt.wantError)
			}

			// Check expected value
			if !phoneWrapper.IsDiscarded() && phoneWrapper.Get() != tt.want {
				t.Errorf("Get() = %v, want %v", phoneWrapper.Get(), tt.want)
			}

			// Check discard status
			if phoneWrapper.IsDiscarded() != tt.wantDiscard {
				t.Errorf("IsDiscarded() = %v, want %v", phoneWrapper.IsDiscarded(), tt.wantDiscard)
			}
		})
	}
}

func TestWrapperRegexPhone_JSONMarshal(t *testing.T) {
	tests := []struct {
		name         string
		phone        string
		wantError    bool
		expectedJSON string
	}{
		{
			name:         "Valid phone number marshalling",
			phone:        "+1234567890",
			wantError:    false,
			expectedJSON: `{"phone":"+1234567890"}`,
		},
		{
			name:         "Invalid phone number marshalling",
			phone:        "0000000123456789",
			wantError:    true,
			expectedJSON: `{"phone":""}`,
		},
		{
			name:         "Nil phone number marshalling",
			phone:        "",
			wantError:    true,
			expectedJSON: `{"phone":null}`,
		},
	}

	// Create a struct for testing
	type Data struct {
		Phone *WrapperRegexPhone `json:"phone"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the phone wrapper
			phoneWrapper := wrappers.New[*WrapperRegexPhone]()

			// Wrap the phone number
			err := phoneWrapper.Wrap(tt.phone, false)
			if err != nil && tt.wantError {
				return
			}

			if err != nil && !tt.wantError {
				t.Fatalf("Unexpected error during wrapping: %v", err)
			}

			data := Data{
				Phone: phoneWrapper,
			}

			// Marshal to JSON
			jsonData, err := json.Marshal(data)
			if err != nil {
				t.Fatalf("Unexpected error during marshalling: %v", err)
			}

			// Compare the JSON output
			if string(jsonData) != tt.expectedJSON {
				t.Errorf("Marshalled JSON = %v, want %v", string(jsonData), tt.expectedJSON)
			}
		})
	}
}

func TestWrapperRegexPhone_Unmarshal(t *testing.T) {
	tests := []struct {
		name      string
		jsonInput string
		wantError bool
		want      string
	}{
		{
			name:      "Valid phone number unmarshalling",
			jsonInput: `{"phone":"+1234567890"}`,
			wantError: false,
			want:      "+1234567890",
		},
		{
			name:      "Invalid phone number unmarshalling",
			jsonInput: `{"phone":"0000000123456789"}`,
			wantError: true,
			want:      "",
		},
		{
			name:      "Empty phone field unmarshalling",
			jsonInput: `{"phone":""}`,
			wantError: true,
			want:      "",
		},
		{
			name:      "Nil JSON input unmarshalling",
			jsonInput: `{"phone":null}`,
			wantError: true,
			want:      "",
		},
		{
			name:      "Missing phone field unmarshalling",
			jsonInput: `{}`,
			wantError: false,
			want:      "",
		},
		{
			name:      "Invalid type for phone field unmarshalling",
			jsonInput: `{"phone":12345}`,
			wantError: true,
			want:      "",
		},
	}

	type User struct {
		Phone WrapperRegexPhone `json:"phone"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := User{}

			// Unmarshal from JSON input
			err := json.Unmarshal([]byte(tt.jsonInput), &user)
			if (err != nil) != tt.wantError {
				t.Errorf("Unmarshal() error = %v, wantError %v", err, tt.wantError)
			}

			// Check expected value
			if !user.Phone.IsDiscarded() && user.Phone.Get() != tt.want {
				t.Errorf("Get() = %v, want %v", user.Phone.Get(), tt.want)
			}
		})
	}
}

func TestWrapperRegexPhone_Initialize(t *testing.T) {
	wrapper := &WrapperRegexPhone{}
	wrapper.Initialize()

	if !wrapper.IsInitialized() {
		t.Error("WrapperRegexPhone should be marked as initialized")
	}

	expectedPattern := `^(?:\+?[1-9]\d{1,14}|0\d{1,14})$`
	if wrapper.WrapperRegex.regex.String() != expectedPattern {
		t.Errorf("Expected pattern %s, got %s", expectedPattern, wrapper.WrapperRegex.regex.String())
	}
}
