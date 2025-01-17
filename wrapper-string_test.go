package wrappers

import (
	"encoding/json"
	"testing"
)

func TestWrapperString_Wrap(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		discard     bool
		want        string
		wantError   bool
		wantDiscard bool
	}{
		{
			name:        "Wrap string",
			input:       "hello",
			discard:     false,
			want:        "hello",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap integer",
			input:       123,
			discard:     false,
			want:        "123",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap float",
			input:       3.14,
			discard:     false,
			want:        "3.14",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap nested wrapper discarding",
			input:       func() *WrapperString { w := New[*WrapperString](); w.Wrap("nested", true); return w }(),
			discard:     false,
			want:        "nested",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap nested wrapper discarding to true",
			input:       func() *WrapperString { w := New[*WrapperString](); w.Discard(); return w }(),
			discard:     false,
			want:        "",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap unsupported type",
			input:       []int{1, 2, 3},
			discard:     false,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Wrap unsupported type discard",
			input:       []int{1, 2, 3},
			discard:     true,
			want:        "",
			wantError:   false,
			wantDiscard: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := New[*WrapperString]()

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

// TestWrapperString_JSONMarshal tests JSON marshalling of WrapperString.
func TestWrapperString_JSONMarshal(t *testing.T) {
	tests := []struct {
		name         string
		wrapper      *WrapperString
		expectedJSON string
	}{
		{
			name: "Marshal non-empty string",
			wrapper: func() *WrapperString {
				w := New[*WrapperString]()
				w.Wrap("test", false)
				return w
			}(),
			expectedJSON: `"test"`,
		},
		{
			name: "Marshal empty string",
			wrapper: func() *WrapperString {
				w := New[*WrapperString]()
				w.Wrap("", false)
				return w
			}(),
			expectedJSON: `null`,
		},
		{
			name: "Marshal discarded string",
			wrapper: func() *WrapperString {
				w := New[*WrapperString]()
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
				t.Fatalf("Failed to marshal WrapperString: %v", err)
			}

			if string(jsonData) != tt.expectedJSON {
				t.Errorf("Marshalled JSON = %v, want %v", string(jsonData), tt.expectedJSON)
			}
		})
	}
}

// TestWrapperString_JSONUnmarshal tests JSON unmarshalling of WrapperString.
func TestWrapperString_JSONUnmarshal(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		want        string
		wantDiscard bool
		wantError   bool
	}{
		{
			name:        "Unmarshal valid string",
			jsonInput:   `"hello world"`,
			want:        "hello world",
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Unmarshal numeric string",
			jsonInput:   `"12345"`,
			want:        "12345",
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
			name:        "Unmarshal numeric type",
			jsonInput:   `123`,
			want:        "123",
			wantDiscard: false,
			wantError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := New[*WrapperString]()

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
