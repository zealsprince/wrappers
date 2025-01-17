package regex

import (
	"encoding/json"
	"testing"

	"github.com/zealsprince/wrappers"
)

func TestWrapperRegexEmail_Wrap(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		discard     bool
		want        string
		wantDiscard bool
		wantError   bool
	}{
		{
			name:        "Valid email",
			input:       "user@example.com",
			discard:     false,
			want:        "user@example.com",
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Invalid email",
			input:       "user@example",
			discard:     false,
			want:        "", // Value remains unchanged
			wantDiscard: true,
			wantError:   true,
		},
		{
			name:        "Invalid email discard",
			input:       "user@.com",
			discard:     true,
			want:        "",
			wantDiscard: true,
			wantError:   false,
		},
		{
			name:        "Input as WrapperProvider (WrapperRegexEmail)",
			input:       wrappers.NewWithValueUnsafe[*WrapperRegexEmail]("user@example.com"),
			discard:     false,
			want:        "user@example.com",
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Input as WrapperProvider (WrapperRegexPhone)",
			input:       wrappers.NewWithValueUnsafe[*WrapperRegexPhone]("+1234567890"),
			discard:     false,
			want:        "",
			wantDiscard: true,
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the email wrapper
			emailWrapper := *wrappers.New[*WrapperRegexEmail]()

			// Wrap the input
			err := emailWrapper.Wrap(tt.input, tt.discard)
			if (err != nil) != tt.wantError {
				t.Errorf("Wrap() error = %v, wantError %v", err, tt.wantError)
			}

			// Check expected value
			if !emailWrapper.IsDiscarded() && emailWrapper.Get() != tt.want {
				t.Errorf("Get() = %v, want %v", emailWrapper.Get(), tt.want)
			}

			// Check discard status
			if emailWrapper.IsDiscarded() != tt.wantDiscard {
				t.Errorf("IsDiscarded() = %v, want %v", emailWrapper.IsDiscarded(), tt.wantDiscard)
			}
		})
	}
}

func TestWrapperRegexEmail_JSONMarshal(t *testing.T) {
	tests := []struct {
		name         string
		email        string
		wantError    bool
		expectedJSON string
	}{
		{
			name:         "Valid Email Marshalling",
			email:        "john.doe@example.com",
			wantError:    false,
			expectedJSON: `{"email":"john.doe@example.com"}`,
		},
		{
			name:         "Invalid Email Marshalling",
			email:        "dorf@.com",
			wantError:    true,
			expectedJSON: `{"email":null}`,
		},
	}

	// Create a struct for testing
	type Data struct {
		Email *WrapperRegexEmail `json:"email"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the wrappers
			emailWrapper := wrappers.New[*WrapperRegexEmail]()

			// Wrap the email
			err := emailWrapper.Wrap(tt.email, false)
			if err != nil && !tt.wantError {
				t.Fatalf("Unexpected error during wrapping: %v", err)
			}

			data := Data{
				Email: emailWrapper,
			}

			// Marshal to JSON
			jsonData, err := json.Marshal(&data)
			if err != nil {
				t.Fatalf("Failed to marshal Data: %v", err)
			}

			if string(jsonData) != tt.expectedJSON {
				t.Errorf("Marshalled JSON = %v, want %v", string(jsonData), tt.expectedJSON)
			}
		})
	}
}

func TestWrapperRegexEmail_Unmarshal(t *testing.T) {
	tests := []struct {
		name      string
		jsonInput string
		wantError bool
	}{
		{
			name:      "Valid Email",
			jsonInput: `{"email":"valid@email.com"}`,
			wantError: false,
		},
		{
			name:      "Invalid Email Format",
			jsonInput: `{"email":"invalid-email"}`,
			wantError: true,
		},
		{
			name:      "Empty Email Field",
			jsonInput: `{"email":""}`,
			wantError: true,
		},
		{
			name:      "Nil JSON Input",
			jsonInput: ``,
			wantError: true,
		},
	}

	type Data struct {
		Email *WrapperRegexEmail `json:"email"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := Data{} // Initialize the data struct we will unmarshal into

			// Unmarshal from JSON input
			err := json.Unmarshal([]byte(tt.jsonInput), &data)
			if (err != nil) != tt.wantError {
				t.Fatalf("Unmarshal() error = %v, wantError %v", err, tt.wantError)
			} else {
				return // Skip the rest of the test since we can't proceed with an error
			}
		})
	}
}

func TestWrapperRegexEmail_Initialize(t *testing.T) {
	wrapper := &WrapperRegexEmail{}
	wrapper.Initialize()

	if !wrapper.IsInitialized() {
		t.Error("WrapperRegexEmail should be marked as initialized")
	}

	expectedPattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	if wrapper.WrapperRegex.regex.String() != expectedPattern {
		t.Errorf("Expected pattern %s, got %s", expectedPattern, wrapper.WrapperRegex.regex.String())
	}
}
