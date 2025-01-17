package wrappers

import (
	"encoding/json"
	"testing"
	"time"
)

func TestWrapperTime_Wrap(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		discard     bool
		want        string // RFC3339 formatted string
		wantError   bool
		wantDiscard bool
	}{
		{
			name:        "Wrap time.Time",
			input:       time.Date(2025, 1, 15, 22, 19, 35, 0, time.UTC),
			discard:     false,
			want:        "2025-01-15T22:19:35Z",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap string valid RFC3339",
			input:       "2023-10-10T10:10:10Z",
			discard:     false,
			want:        "2023-10-10T10:10:10Z",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap string invalid RFC3339",
			input:       "10/10/2023",
			discard:     false,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Wrap string invalid RFC3339 discard",
			input:       "10/10/2023",
			discard:     true,
			want:        "",
			wantError:   false,
			wantDiscard: true,
		},
		{
			name:        "Wrap unsupported type",
			input:       true,
			discard:     false,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Wrap unsupported type discard",
			input:       true,
			discard:     true,
			want:        "",
			wantError:   false,
			wantDiscard: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := New[*WrapperTime]()

			err := wrapper.Wrap(tt.input, tt.discard)
			if tt.wantError && err == nil {
				t.Fatalf("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("Did not expect error but got: %v", err)
			}

			if tt.wantDiscard && !wrapper.IsDiscarded() {
				t.Errorf("Expected wrapper to be discarded")
			}

			unwrapped := wrapper.Unwrap()
			if unwrapped != tt.want {
				t.Errorf("Unwrapped value = %v, want %v", unwrapped, tt.want)
			}
		})
	}
}

// TestWrapperTime_JSONMarshal tests JSON marshalling of WrapperTime.
func TestWrapperTime_JSONMarshal(t *testing.T) {
	tests := []struct {
		name         string
		wrapper      *WrapperTime
		expectedJSON string
	}{
		{
			name: "Marshal valid time",
			wrapper: func() *WrapperTime {
				w := New[*WrapperTime]()
				tm := time.Date(2025, time.January, 15, 22, 19, 35, 0, time.UTC)
				w.Wrap(tm, false)
				return w
			}(),
			expectedJSON: `"2025-01-15T22:19:35Z"`,
		},
		{
			name: "Marshal empty time",
			wrapper: func() *WrapperTime {
				w := New[*WrapperTime]()
				w.Wrap(time.Time{}, false)
				return w
			}(),
			expectedJSON: `"0001-01-01T00:00:00Z"`,
		},
		{
			name: "Marshal discarded time",
			wrapper: func() *WrapperTime {
				w := New[*WrapperTime]()
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
				t.Fatalf("Failed to marshal WrapperTime: %v", err)
			}

			if string(jsonData) != tt.expectedJSON {
				t.Errorf("Marshalled JSON = %v, want %v", string(jsonData), tt.expectedJSON)
			}
		})
	}
}

// TestWrapperTime_JSONUnmarshal tests JSON unmarshalling of WrapperTime.
func TestWrapperTime_JSONUnmarshal(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		want        string
		wantDiscard bool
		wantError   bool
	}{
		{
			name:        "Unmarshal valid RFC3339 string",
			jsonInput:   `"2023-10-10T10:10:10Z"`,
			want:        "2023-10-10T10:10:10Z",
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Unmarshal invalid RFC3339 string",
			jsonInput:   `"10/10/2023"`,
			want:        "",
			wantDiscard: true,
			wantError:   true,
		},
		{
			name:        "Unmarshal time.Time object",
			jsonInput:   `"2024-04-20T15:30:00Z"`,
			want:        "2024-04-20T15:30:00Z",
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Unmarshal null",
			jsonInput:   `null`,
			want:        "",
			wantDiscard: true,
			wantError:   false,
		},
		{
			name:        "Unmarshal unsupported type",
			jsonInput:   `true`,
			want:        "",
			wantDiscard: false,
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := New[*WrapperTime]()

			err := json.Unmarshal([]byte(tt.jsonInput), wrapper)
			if tt.wantError && err == nil {
				t.Fatalf("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("Did not expect error but got: %v", err)
			}

			if tt.wantDiscard && !wrapper.IsDiscarded() {
				t.Errorf("Expected wrapper to be discarded")
			}

			unwrapped := wrapper.Unwrap()
			if unwrapped != tt.want {
				t.Errorf("Unwrapped value = %v, want %v", unwrapped, tt.want)
			}
		})
	}
}
