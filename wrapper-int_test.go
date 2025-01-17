package wrappers

import (
	"encoding/json"
	"testing"
)

// Existing Tests for WrapperBool, WrapperFloat, WrapperString, etc.
// ... (Assuming these exist above)

// TestWrapperInt_Wrap tests the Wrap method of WrapperInt.
func TestWrapperInt_Wrap(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		discard     bool
		want        int64
		wantError   bool
		wantDiscard bool
	}{
		{
			name:        "Wrap int",
			input:       42,
			discard:     false,
			want:        42,
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap int16",
			input:       int16(123),
			discard:     false,
			want:        123,
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap int32",
			input:       int32(123456),
			discard:     false,
			want:        123456,
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap int64",
			input:       int64(9876543210),
			discard:     false,
			want:        9876543210,
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap float32",
			input:       float32(3.14),
			discard:     false,
			want:        3, // Truncated to int
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap float64",
			input:       float64(-2.718),
			discard:     false,
			want:        -2, // Truncated to int
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap string with valid int",
			input:       "256",
			discard:     false,
			want:        256,
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap string with invalid int",
			input:       "invalid",
			discard:     false,
			want:        0, // Discarded
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Wrap string with invalid int discard",
			input:       "invalid",
			discard:     true,
			want:        0, // Discarded
			wantError:   false,
			wantDiscard: true,
		},
		{
			name:        "Wrap nil",
			input:       nil,
			discard:     false,
			want:        0, // Discarded
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Wrap unsupported type",
			input:       []int{1, 2, 3},
			discard:     false,
			want:        0, // Discarded
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Wrap unsupported type discard",
			input:       []int{1, 2, 3},
			discard:     true,
			want:        0, // Discarded
			wantError:   false,
			wantDiscard: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := New[*WrapperInt]()

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

			// If no error and 'discard' is true, ensure it's not discarded
			if tt.wantDiscard && !wrapper.IsDiscarded() {
				t.Errorf("Expected wrapper to be discarded")
			}
		})
	}
}

// TestWrapperInt_JSONMarshal tests JSON marshalling of WrapperInt.
func TestWrapperInt_JSONMarshal(t *testing.T) {
	tests := []struct {
		name         string
		wrapper      *WrapperInt
		expectedJSON string
	}{
		{
			name: "Marshal valid int",
			wrapper: func() *WrapperInt {
				w := New[*WrapperInt]()
				w.Wrap(100, false)
				return w
			}(),
			expectedJSON: "100",
		},
		{
			name: "Marshal wrapped string that is a valid int",
			wrapper: func() *WrapperInt {
				w := New[*WrapperInt]()
				w.Wrap("256", false)
				return w
			}(),
			expectedJSON: "256",
		},
		{
			name: "Marshal discarded wrapper",
			wrapper: func() *WrapperInt {
				w := New[*WrapperInt]()
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
				t.Fatalf("Failed to marshal WrapperInt: %v", err)
			}

			if string(jsonData) != tt.expectedJSON {
				t.Errorf("Marshalled JSON = %v, want %v", string(jsonData), tt.expectedJSON)
			}
		})
	}
}

// TestWrapperInt_JSONUnmarshal tests JSON unmarshalling of WrapperInt.
func TestWrapperInt_JSONUnmarshal(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		want        int64
		wantDiscard bool
		wantError   bool
	}{
		{
			name:        "Unmarshal valid JSON int",
			jsonInput:   "12345",
			want:        12345,
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Unmarshal valid JSON float (should truncate)",
			jsonInput:   "78.90",
			want:        78, // Truncated to int
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Unmarshal valid JSON string representing int",
			jsonInput:   `"2048"`,
			want:        2048,
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Unmarshal invalid JSON string",
			jsonInput:   `"not_an_int"`,
			want:        0,
			wantDiscard: true,
			wantError:   true, // Since discard is false during unmarshalling
		},
		{
			name:        "Unmarshal null",
			jsonInput:   "null",
			want:        0,
			wantDiscard: true,
			wantError:   true,
		},
		{
			name:        "Unmarshal invalid JSON type (object)",
			jsonInput:   `{"number": 123}`,
			want:        0,
			wantDiscard: true,
			wantError:   true,
		},
		{
			name:        "Unmarshal invalid JSON type (array)",
			jsonInput:   `[1, 2, 3]`,
			want:        0,
			wantDiscard: true,
			wantError:   true,
		},
		{
			name:        "Unmarshal string invalid int with discard",
			jsonInput:   `"invalid_int"`,
			want:        0,
			wantDiscard: true,
			wantError:   true, // Wrap is called with discard=true, but in Unmarshal it likely calls with discard=false
			// Depending on implementation, adjust
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := New[*WrapperInt]()

			err := json.Unmarshal([]byte(tt.jsonInput), wrapper)
			if tt.wantError {
				if err == nil {
					t.Fatalf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("Did not expect error but got: %v", err)
				}
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
