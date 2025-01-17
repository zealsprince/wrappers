package wrappers

import (
	"encoding/json"
	"testing"
)

// NestedDiscarder is a struct that contains a Discarder wrapping a WrapperString.
type NestedDiscarder struct {
	Value *Discarder[*WrapperString] `json:"value"`
}

func TestDiscarder(t *testing.T) {
	tests := []struct {
		name        string
		wrapper     WrapperProvider
		inputJSON   string
		wantValue   any
		wantDiscard bool
		wantMarshal string
	}{
		{
			name: "WrapperInt - Valid Input",
			wrapper: func() WrapperProvider {
				w := New[*WrapperInt]()
				return w
			}(),
			inputJSON:   "123",
			wantValue:   int64(123),
			wantDiscard: false,
			wantMarshal: "123",
		},
		{
			name: "WrapperInt - Invalid Input with Discard",
			wrapper: func() WrapperProvider {
				w := New[*WrapperInt]()
				return w
			}(),
			inputJSON:   `"abc"`,
			wantValue:   int64(0),
			wantDiscard: true,
			wantMarshal: "null",
		},
		{
			name: "WrapperString - Valid Input",
			wrapper: func() WrapperProvider {
				w := New[*WrapperString]()
				return w
			}(),
			inputJSON:   `"Hello, World!"`,
			wantValue:   "Hello, World!",
			wantDiscard: false,
			wantMarshal: `"Hello, World!"`,
		},
		{
			name: "WrapperString - Invalid Input with Discard",
			wrapper: func() WrapperProvider {
				w := New[*WrapperString]()
				return w
			}(),
			inputJSON:   `123`, // Assuming WrapperString expects a string or convertible types
			wantValue:   "123",
			wantDiscard: false,
			wantMarshal: `"123"`,
		},
		{
			name: "WrapperString - Completely Invalid Input with Discard",
			wrapper: func() WrapperProvider {
				w := New[*WrapperString]()
				return w
			}(),
			inputJSON:   `{"key": "value"}`, // Object is invalid for string
			wantValue:   "",
			wantDiscard: true,
			wantMarshal: "null",
		},
		{
			name: "WrapperInt - Float Input Converted to Int",
			wrapper: func() WrapperProvider {
				w := New[*WrapperInt]()
				return w
			}(),
			inputJSON:   "45.67",
			wantValue:   int64(45),
			wantDiscard: false,
			wantMarshal: "45",
		},
		{
			name: "WrapperInt - String Input Parsed to Int",
			wrapper: func() WrapperProvider {
				w := New[*WrapperInt]()
				return w
			}(),
			inputJSON:   `"789"`,
			wantValue:   int64(789),
			wantDiscard: false,
			wantMarshal: "789",
		},
		{
			name: "WrapperInt - Empty String Input",
			wrapper: func() WrapperProvider {
				w := New[*WrapperInt]()
				return w
			}(),
			inputJSON:   `""`,
			wantValue:   int64(0),
			wantDiscard: true,
			wantMarshal: "null",
		},
		{
			name: "WrapperString - Nil Input",
			wrapper: func() WrapperProvider {
				w := New[*WrapperString]()
				return w
			}(),
			inputJSON:   `null`,
			wantValue:   "",
			wantDiscard: true,
			wantMarshal: "null",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			discarder := NewDiscarder(tt.wrapper)

			// Unmarshal JSON into the Discarder
			err := json.Unmarshal([]byte(tt.inputJSON), discarder)
			if err != nil {
				t.Fatalf("Discarder caught an error which should never happen: %v", err)
			}

			// Check if the value is discarded as expected
			if discarder.Proxy.IsDiscarded() != tt.wantDiscard {
				t.Errorf("Discard status mismatch: got %v, expected %v", discarder.Proxy.IsDiscarded(), tt.wantDiscard)
			}

			// Unwrap and verify the value
			unwrappedValue := discarder.Proxy.UnwrapAny()
			if unwrappedValue != tt.wantValue {
				t.Errorf("Unwrapped value mismatch: got %v, expected %v", unwrappedValue, tt.wantValue)
			}

			// Marshal the Discarder back to JSON
			marshaledJSON, err := json.Marshal(discarder)
			if err != nil {
				t.Fatalf("Marshal returned an unexpected error: %v", err)
			}

			// Compare the marshaled JSON with the expected JSON
			if string(marshaledJSON) != tt.wantMarshal {
				t.Errorf("Marshaled JSON mismatch: got %s, expected %s", string(marshaledJSON), tt.wantMarshal)
			}
		})
	}
}

func TestDiscarderInStruct(t *testing.T) {
	tests := []struct {
		name           string
		inputJSON      string
		wantValue      any
		wantDiscarded  bool
		wantMarshalOut string
	}{
		{
			name:           "Valid String Input",
			inputJSON:      `{"value": "Hello, World!"}`,
			wantValue:      "Hello, World!",
			wantDiscarded:  false,
			wantMarshalOut: `{"value":"Hello, World!"}`,
		},
		{
			name:           "Invalid String Input (Object)",
			inputJSON:      `{"value": {"key": "value"}}`,
			wantValue:      "",
			wantDiscarded:  true,
			wantMarshalOut: `{"value":null}`,
		},
		{
			name:           "Invalid String Input (Number)",
			inputJSON:      `{"value": 123}`,
			wantValue:      "123",
			wantDiscarded:  false,
			wantMarshalOut: `{"value":"123"}`,
		},
		{
			name:           "Null Input",
			inputJSON:      `{"value": null}`,
			wantValue:      "",
			wantDiscarded:  true,
			wantMarshalOut: `{"value":null}`,
		},
		{
			name:           "Empty String Input",
			inputJSON:      `{"value": ""}`,
			wantValue:      "",
			wantDiscarded:  true,
			wantMarshalOut: `{"value":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize a WrapperString
			wrapper := New[*WrapperString]()

			// Initialize the Discarder with the WrapperString
			discarder := NewDiscarder(wrapper)

			// Initialize the NestedDiscarder with the Discarder
			nested := NestedDiscarder{
				Value: discarder,
			}

			// Unmarshal the input JSON into the NestedDiscarder struct
			err := json.Unmarshal([]byte(tt.inputJSON), &nested)
			if err != nil {
				t.Fatalf("We should never receive an error since this is a Discarder: %v", err)
			}

			if nested.Value == nil && tt.inputJSON == `{"value": null}` {
				return
			} else if nested.Value == nil {
				t.Fatalf("Discarder should not be nil")
			}

			// Check if the Discarder is in the expected discarded state
			if nested.Value.Proxy.IsDiscarded() != tt.wantDiscarded {
				t.Errorf("Discarded state mismatch: got %v, want %v", nested.Value.Proxy.IsDiscarded(), tt.wantDiscarded)
			}

			// Unwrap the value and verify it matches the expectation
			unwrappedValue := nested.Value.Proxy.UnwrapAny()
			if unwrappedValue != tt.wantValue {
				t.Errorf("Unwrapped value mismatch: got %v, want %v", unwrappedValue, tt.wantValue)
			}

			// Marshal the NestedDiscarder back to JSON
			marshaledJSON, err := json.Marshal(nested)
			if err != nil {
				t.Fatalf("Failed to marshal NestedDiscarder: %v", err)
			}

			// Compare the marshaled JSON with the expected marshaled output
			if string(marshaledJSON) != tt.wantMarshalOut {
				t.Errorf("Marshaled JSON mismatch:\nGot:      %s\nExpected: %s",
					string(marshaledJSON), tt.wantMarshalOut)
			}
		})
	}
}
