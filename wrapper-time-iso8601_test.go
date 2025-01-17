package wrappers

import (
	"encoding/json"
	"testing"
	"time"
)

// TestWrapperTimeISO8601_Wrap tests the Wrap method of WrapperTimeISO8601.
func TestWrapperTimeISO8601_Wrap(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		discard     bool
		want        string
		wantError   bool
		wantDiscard bool
	}{
		{
			name:        "Wrap valid time.Time",
			input:       time.Date(2025, 1, 16, 7, 8, 30, 677, time.UTC),
			discard:     false,
			want:        "2025-01-16T07:08:30Z",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap valid ISO8601 string",
			input:       "2021-09-01T00:00:00.000Z",
			discard:     false,
			want:        "2021-09-01T00:00:00Z",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap invalid ISO8601 string without discard",
			input:       "invalid_time",
			discard:     false,
			want:        "0001-01-01T00:00:00Z",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Wrap invalid ISO8601 string with discard",
			input:       "invalid_time",
			discard:     true,
			want:        "0001-01-01T00:00:00Z",
			wantError:   false,
			wantDiscard: true,
		},
		{
			name:        "Wrap nil",
			input:       nil,
			discard:     false,
			want:        "0001-01-01T00:00:00Z",
			wantError:   false,
			wantDiscard: true,
		},
		{
			name: "Wrap WrapperTimeISO8601 with valid value",
			input: func() *WrapperTimeISO8601 {
				w := New[*WrapperTimeISO8601]()
				w.Wrap(time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC), false)
				return w
			}(),
			discard:     false,
			want:        "2025-12-31T23:59:59Z",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap WrapperTimeISO8601 with discarded value",
			input:       func() *WrapperTimeISO8601 { w := New[*WrapperTimeISO8601](); w.Discard(); return w }(),
			discard:     false,
			want:        "0001-01-01T00:00:00Z",
			wantError:   false,
			wantDiscard: true,
		},
		{
			name:        "Wrap unsupported type without discard",
			input:       12345, // int is unsupported
			discard:     false,
			want:        "0001-01-01T00:00:00Z",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Wrap unsupported type with discard",
			input:       12345, // int is unsupported
			discard:     true,
			want:        "0001-01-01T00:00:00Z",
			wantError:   false,
			wantDiscard: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := New[*WrapperTimeISO8601]()

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

// TestWrapperTimeISO8601_JSONMarshal tests JSON marshalling of WrapperTimeISO8601.
func TestWrapperTimeISO8601_JSONMarshal(t *testing.T) {
	tests := []struct {
		name         string
		wrapper      *WrapperTimeISO8601
		expectedJSON string
	}{
		{
			name: "Marshal valid time",
			wrapper: func() *WrapperTimeISO8601 {
				w := New[*WrapperTimeISO8601]()
				w.Wrap(time.Date(2025, 1, 16, 7, 8, 30, 677, time.UTC), false)
				return w
			}(),
			expectedJSON: `"2025-01-16T07:08:30Z"`,
		},
		{
			name: "Marshal valid ISO8601 string",
			wrapper: func() *WrapperTimeISO8601 {
				w := New[*WrapperTimeISO8601]()
				w.Wrap("2021-09-01T00:00:00.000Z", false)
				return w
			}(),
			expectedJSON: `"2021-09-01T00:00:00Z"`,
		},
		{
			name: "Marshal discarded wrapper",
			wrapper: func() *WrapperTimeISO8601 {
				w := New[*WrapperTimeISO8601]()
				w.Discard()
				return w
			}(),
			expectedJSON: "null",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.wrapper)
			if err != nil {
				t.Fatalf("Failed to marshal WrapperTimeISO8601: %v", err)
			}

			if string(jsonData) != tt.expectedJSON {
				t.Errorf("Marshalled JSON = %v, want %v", string(jsonData), tt.expectedJSON)
			}
		})
	}
}

// TestWrapperTimeISO8601_JSONUnmarshal tests JSON unmarshalling of WrapperTimeISO8601.
func TestWrapperTimeISO8601_JSONUnmarshal(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		want        string
		wantDiscard bool
		wantError   bool
	}{
		{
			name:        "Unmarshal valid RFC3339 string",
			jsonInput:   `"2025-01-16T07:08:30Z"`,
			want:        "2025-01-16T07:08:30Z",
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Unmarshal valid ISO8601 string with +0000",
			jsonInput:   `"2021-09-01T00:00:00.000+0000"`,
			want:        "2021-09-01T00:00:00Z",
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Unmarshal invalid string",
			jsonInput:   `"invalid_time"`,
			want:        "0001-01-01T00:00:00Z",
			wantDiscard: true,
			wantError:   true,
		},
		{
			name:        "Unmarshal null",
			jsonInput:   `null`,
			want:        "0001-01-01T00:00:00Z",
			wantDiscard: true,
			wantError:   false,
		},
		{
			name:        "Unmarshal unsupported type (number)",
			jsonInput:   `12345`,
			want:        "0001-01-01T00:00:00Z",
			wantDiscard: true,
			wantError:   true,
		},
		{
			name:        "Unmarshal WrapperTimeISO8601 with valid value",
			jsonInput:   `"2025-12-31T23:59:59Z"`,
			want:        "2025-12-31T23:59:59Z",
			wantDiscard: false,
			wantError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := New[*WrapperTimeISO8601]()

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
