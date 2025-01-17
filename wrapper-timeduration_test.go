package wrappers

import (
	"encoding/json"
	"testing"
	"time"
)

// TestWrapperTimeDuration_Wrap tests the Wrap method of WrapperTimeDuration.
func TestWrapperTimeDuration_Wrap(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		discard     bool
		want        string
		wantError   bool
		wantDiscard bool
	}{
		{
			name:        "Wrap valid time.Duration",
			input:       time.Hour + 30*time.Minute, // 1h30m0s
			discard:     false,
			want:        "1h30m0s",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap valid duration string",
			input:       "2h45m",
			discard:     false,
			want:        "2h45m0s",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap invalid duration string without discard",
			input:       "invalid_duration",
			discard:     false,
			want:        "0s",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Wrap invalid duration string with discard",
			input:       "invalid_duration",
			discard:     true,
			want:        "0s",
			wantError:   false,
			wantDiscard: true,
		},
		{
			name:        "Wrap nil",
			input:       nil,
			discard:     false,
			want:        "0s",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Wrap WrapperTimeDuration with valid value",
			input:       func() *WrapperTimeDuration { w := New[*WrapperTimeDuration](); w.Wrap(time.Minute*15, false); return w }(),
			discard:     false,
			want:        "15m0s",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap WrapperTimeDuration with discarded value",
			input:       func() *WrapperTimeDuration { w := New[*WrapperTimeDuration](); w.Discard(); return w }(),
			discard:     false,
			want:        "0s",
			wantError:   false,
			wantDiscard: true,
		},
		{
			name:        "Wrap unsupported type without discard",
			input:       "not_a_duration", // string handled
			discard:     false,
			want:        "0s",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Wrap unsupported type with discard",
			input:       []int{1, 2, 3}, // unsupported
			discard:     true,
			want:        "0s",
			wantError:   false,
			wantDiscard: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := New[*WrapperTimeDuration]()

			err := wrapper.Wrap(tt.input, tt.discard)
			if tt.wantError && err == nil {
				t.Fatalf("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("Did not expect error but got: %v", err)
			}

			if tt.wantDiscard != wrapper.IsDiscarded() {
				t.Errorf("Discarded state = %v, want %v", wrapper.IsDiscarded(), tt.wantDiscard)
			}

			unwrapped := wrapper.Unwrap()
			if unwrapped != tt.want {
				t.Errorf("Unwrapped value = %v, want %v", unwrapped, tt.want)
			}
		})
	}
}

// TestWrapperTimeDuration_JSONMarshal tests JSON marshalling of WrapperTimeDuration.
func TestWrapperTimeDuration_JSONMarshal(t *testing.T) {
	tests := []struct {
		name         string
		wrapper      *WrapperTimeDuration
		expectedJSON string
	}{
		{
			name: "Marshal valid duration",
			wrapper: func() *WrapperTimeDuration {
				w := New[*WrapperTimeDuration]()
				w.Wrap(time.Minute*45, false) // 45m0s
				return w
			}(),
			expectedJSON: `"45m0s"`,
		},
		{
			name: "Marshal valid duration string",
			wrapper: func() *WrapperTimeDuration {
				w := New[*WrapperTimeDuration]()
				w.Wrap("30m", false)
				return w
			}(),
			expectedJSON: `"30m0s"`,
		},
		{
			name: "Marshal discarded wrapper",
			wrapper: func() *WrapperTimeDuration {
				w := New[*WrapperTimeDuration]()
				w.Discard()
				return w
			}(),
			expectedJSON: "null",
		},
		{
			name: "Marshal zero duration",
			wrapper: func() *WrapperTimeDuration {
				w := New[*WrapperTimeDuration]()
				w.Wrap(0, false)
				return w
			}(),
			expectedJSON: `"0s"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.wrapper)
			if err != nil {
				t.Fatalf("Failed to marshal WrapperTimeDuration: %v", err)
			}

			if string(jsonData) != tt.expectedJSON {
				t.Errorf("Marshalled JSON = %v, want %v", string(jsonData), tt.expectedJSON)
			}
		})
	}
}

// TestWrapperTimeDuration_JSONUnmarshal tests JSON unmarshalling of WrapperTimeDuration.
func TestWrapperTimeDuration_JSONUnmarshal(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		want        string
		wantDiscard bool
		wantError   bool
	}{
		{
			name:        "Unmarshal valid duration string",
			jsonInput:   `"2h30m"`,
			want:        "2h30m0s",
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Unmarshal valid duration number (seconds)",
			jsonInput:   `3600`, // 3.6 microseconds
			want:        "3.6µs",
			wantDiscard: false,
			wantError:   false, // Unsupported type during unmarshalling
		},
		{
			name:        "Unmarshal invalid duration string without discard",
			jsonInput:   `"invalid_duration"`,
			want:        "0s",
			wantDiscard: true,
			wantError:   true,
		},
		{
			name:        "Unmarshal invalid duration string with discard",
			jsonInput:   `"invalid_duration"`,
			want:        "0s",
			wantDiscard: true,
			wantError:   true,
		},
		{
			name:        "Unmarshal null",
			jsonInput:   `null`,
			want:        "0s",
			wantDiscard: true,
			wantError:   true,
		},
		{
			name:        "Unmarshal unsupported type (object)",
			jsonInput:   `{"duration": "1h30m"}`,
			want:        "0s",
			wantDiscard: true,
			wantError:   true,
		},
		{
			name:        "Unmarshal WrapperTimeDuration with valid time.Duration",
			jsonInput:   `"1h45m"`,
			want:        "1h45m0s",
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Unmarshal WrapperTimeDuration with zero duration",
			jsonInput:   `"0s"`,
			want:        "0s",
			wantDiscard: false,
			wantError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := New[*WrapperTimeDuration]()

			err := json.Unmarshal([]byte(tt.jsonInput), wrapper)
			if tt.wantError && err == nil {
				t.Fatalf("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("Did not expect error but got: %v", err)
			}

			if tt.wantDiscard != wrapper.IsDiscarded() {
				t.Errorf("Discarded state = %v, want %v", wrapper.IsDiscarded(), tt.wantDiscard)
			}

			unwrapped := wrapper.Unwrap()
			if unwrapped != tt.want {
				t.Errorf("Unwrapped value = %v, want %v", unwrapped, tt.want)
			}
		})
	}
}
