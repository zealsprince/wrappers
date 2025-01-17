package wrappers

import (
	"encoding/json"
	"testing"
)

func TestWrapperBool_Wrap(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		discard   bool
		want      bool
		wantError bool
	}{
		{
			name:      "Wrap boolean true",
			input:     true,
			discard:   false,
			want:      true,
			wantError: false,
		},
		{
			name:      "Wrap boolean false",
			input:     false,
			discard:   false,
			want:      false,
			wantError: false,
		},
		{
			name:      "Wrap integer non-zero",
			input:     1,
			discard:   false,
			want:      true,
			wantError: false,
		},
		{
			name:      "Wrap integer zero",
			input:     0,
			discard:   false,
			want:      false,
			wantError: false,
		},
		{
			name:      "Wrap float non-zero",
			input:     3.14,
			discard:   false,
			want:      true,
			wantError: false,
		},
		{
			name:      "Wrap float zero",
			input:     0.0,
			discard:   false,
			want:      false,
			wantError: false,
		},
		{
			name:      "Wrap string 'true'",
			input:     "true",
			discard:   false,
			want:      true,
			wantError: false,
		},
		{
			name:      "Wrap string 'false'",
			input:     "false",
			discard:   false,
			want:      false,
			wantError: false,
		},
		{
			name:      "Wrap string 'yes'",
			input:     "yes",
			discard:   false,
			want:      true,
			wantError: false,
		},
		{
			name:      "Wrap string 'no'",
			input:     "no",
			discard:   false,
			want:      false,
			wantError: false,
		},
		{
			name:      "Wrap string '1'",
			input:     "1",
			discard:   false,
			want:      true,
			wantError: false,
		},
		{
			name:      "Wrap string '0'",
			input:     "0",
			discard:   false,
			want:      false,
			wantError: false,
		},
		{
			name:      "Wrap invalid string",
			input:     "maybe",
			discard:   false,
			want:      false, // Value is not set; default to false
			wantError: true,
		},
		{
			name:      "Wrap invalid string discard",
			input:     "maybe",
			discard:   true,
			want:      false, // Discarded means should return default
			wantError: false,
		},
		{
			name:      "Wrap unsupported type",
			input:     []int{1, 2, 3},
			discard:   false,
			want:      false,
			wantError: true,
		},
		{
			name:      "Wrap unsupported type discard",
			input:     []int{1, 2, 3},
			discard:   true,
			want:      false,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := New[*WrapperBool]()

			err := wrapper.Wrap(tt.input, tt.discard)
			if tt.wantError && err == nil {
				t.Fatalf("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("Did not expect error but got: %v", err)
			}

			unwrapped := wrapper.Unwrap()
			if unwrapped != tt.want {
				t.Errorf("Unwrapped value = %v, want %v", unwrapped, tt.want)
			}
		})
	}
}

// TestWrapperBool_JSONMarshal tests JSON marshalling of WrapperBool.
func TestWrapperBool_JSONMarshal(t *testing.T) {
	tests := []struct {
		name         string
		wrapper      *WrapperBool
		expectedJSON string
	}{
		{
			name: "Marshal true",
			wrapper: func() *WrapperBool {
				w := New[*WrapperBool]()
				w.Wrap(true, false)
				return w
			}(),
			expectedJSON: "true",
		},
		{
			name: "Marshal false",
			wrapper: func() *WrapperBool {
				w := New[*WrapperBool]()
				w.Wrap(false, false)
				return w
			}(),
			expectedJSON: "false",
		},
		{
			name: "Marshal discarded",
			wrapper: func() *WrapperBool {
				w := New[*WrapperBool]()
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
				t.Fatalf("Failed to marshal WrapperBool: %v", err)
			}

			if string(jsonData) != tt.expectedJSON {
				t.Errorf("Marshalled JSON = %v, want %v", string(jsonData), tt.expectedJSON)
			}
		})
	}
}

// TestWrapperBool_JSONUnmarshal tests JSON unmarshalling of WrapperBool.
func TestWrapperBool_JSONUnmarshal(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		want        bool
		wantDiscard bool
		wantError   bool
	}{
		{
			name:        "Unmarshal true",
			jsonInput:   "true",
			want:        true,
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Unmarshal false",
			jsonInput:   "false",
			want:        false,
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Unmarshal invalid string",
			jsonInput:   `"maybe"`,
			want:        false,
			wantDiscard: true,
			wantError:   true,
		},
		{
			name:        "Unmarshal null",
			jsonInput:   "null",
			want:        false,
			wantDiscard: true,
			wantError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := New[*WrapperBool]()

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
