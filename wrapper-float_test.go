package wrappers

import (
	"encoding/json"
	"math"
	"testing"
)

func TestWrapperFloat_Wrap(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		discard     bool
		want        float64
		wantError   bool
		wantDiscard bool
	}{
		{
			name:        "Wrap float64",
			input:       3.1415,
			discard:     false,
			want:        3.1415,
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap float32",
			input:       float32(2.718),
			discard:     false,
			want:        2.718,
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap int",
			input:       42,
			discard:     false,
			want:        42.0,
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap string valid float",
			input:       "123.456",
			discard:     false,
			want:        123.456,
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap string invalid float",
			input:       "invalid",
			discard:     false,
			want:        0.0,
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Wrap string invalid float discard",
			input:       "invalid",
			discard:     true,
			want:        0.0,
			wantError:   false,
			wantDiscard: true,
		},
		{
			name:        "Wrap unsupported type",
			input:       true,
			discard:     false,
			want:        0.0,
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Wrap unsupported type discard",
			input:       true,
			discard:     true,
			want:        0.0,
			wantError:   false,
			wantDiscard: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := New[*WrapperFloat]()

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

			// Due to floating point precision, we need to use an epsilon for comparison.
			// To explain: if the input is a float32, the unwrapped value will be a float64.
			// In the case of float32, it has a precision of about 7 decimal digits, while float64 has a precision of
			// about 15-16 decimal digits. When you convert from float32 to float64, you may see slight differences
			// due to how these values are stored in binary format.

			// The value float32(2.718) is actually stored as 2.7179999351501465 in binary representation,
			// which is the closest representable value in float32.
			const epsilon = 1e-6
			unwrapped := wrapper.Unwrap()
			if math.Abs(unwrapped-tt.want) > epsilon {
				t.Errorf("Unwrapped value = %v, want %v", unwrapped, tt.want)
			}
		})
	}
}

// TestWrapperFloat_JSONMarshal tests JSON marshalling of WrapperFloat.
func TestWrapperFloat_JSONMarshal(t *testing.T) {
	tests := []struct {
		name         string
		wrapper      *WrapperFloat
		expectedJSON string
	}{
		{
			name: "Marshal positive float",
			wrapper: func() *WrapperFloat {
				w := New[*WrapperFloat]()
				w.Wrap(3.14, false)
				return w
			}(),
			expectedJSON: "3.14",
		},
		{
			name: "Marshal negative float",
			wrapper: func() *WrapperFloat {
				w := New[*WrapperFloat]()
				w.Wrap(-2.71, false)
				return w
			}(),
			expectedJSON: "-2.71",
		},
		{
			name: "Marshal discarded float",
			wrapper: func() *WrapperFloat {
				w := New[*WrapperFloat]()
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
				t.Fatalf("Failed to marshal WrapperFloat: %v", err)
			}

			if string(jsonData) != tt.expectedJSON {
				t.Errorf("Marshalled JSON = %v, want %v", string(jsonData), tt.expectedJSON)
			}
		})
	}
}

// TestWrapperFloat_JSONUnmarshal tests JSON unmarshalling of WrapperFloat.
func TestWrapperFloat_JSONUnmarshal(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		want        float64
		wantDiscard bool
		wantError   bool
	}{
		{
			name:        "Unmarshal positive float",
			jsonInput:   "3.14",
			want:        3.14,
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Unmarshal negative float",
			jsonInput:   "-2.71",
			want:        -2.71,
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Unmarshal valid string float",
			jsonInput:   `"123.456"`,
			want:        123.456,
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Unmarshal invalid string float",
			jsonInput:   `"invalid"`,
			want:        0.0,
			wantDiscard: true,
			wantError:   true,
		},
		{
			name:        "Unmarshal null",
			jsonInput:   "null",
			want:        0.0,
			wantDiscard: true,
			wantError:   false,
		},
		{
			name:        "Unmarshal unsupported type",
			jsonInput:   `"true"`,
			want:        0.0,
			wantDiscard: false,
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := New[*WrapperFloat]()

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
