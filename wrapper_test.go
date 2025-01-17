package wrappers

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

type ExampleNested struct {
	NestedNumberValue *WrapperInt    `json:"nested_number_value"`
	NestedStringValue *WrapperString `json:"nested_string_value"`
}

type ExampleOmit struct {
	OmitNumberValue *WrapperInt    `json:"omit_number_value,omitempty"`
	OmitStringValue *WrapperString `json:"omit_string_value,omitempty"`
}

// ExampleDiscarder holds Discarder types
type ExampleDiscarder struct {
	DiscarderNumberValue *Discarder[*WrapperInt] `json:"discarder_number_value"`
}

// Updated Example struct to include Discarder
type Example struct {
	NumberValue *WrapperInt    `json:"number_value"`
	StringValue *WrapperString `json:"string_value"`
	BoolValue   *WrapperBool   `json:"bool_value"`
	FloatValue  *WrapperFloat  `json:"float_value"`
	TimeValue   *WrapperTime   `json:"time_value"`

	Nested    ExampleNested    `json:"nested"`
	Omit      *ExampleOmit     `json:"omit,omitempty"`
	Discarder ExampleDiscarder `json:"discarder"`
}

// trimJson trims JSON strings of whitespace and newlines for better comparison while keeping them readable.
func trimJson(jsonStr string) string {
	return strings.ReplaceAll(strings.ReplaceAll(jsonStr, "\n", ""), "\t", "")
}

// compareJSON compares two JSON objects represented as interface{}
func compareJSON(a, b interface{}) bool {
	switch aTyped := a.(type) {
	case map[string]interface{}:
		bTyped, ok := b.(map[string]interface{})
		if !ok {
			return false
		}
		if len(aTyped) != len(bTyped) {
			return false
		}
		for key, aValue := range aTyped {
			bValue, exists := bTyped[key]
			if !exists {
				return false
			}
			if !compareJSON(aValue, bValue) {
				return false
			}
		}
		return true
	case []interface{}:
		bTyped, ok := b.([]interface{})
		if !ok {
			return false
		}
		if len(aTyped) != len(bTyped) {
			return false
		}
		for i := range aTyped {
			if !compareJSON(aTyped[i], bTyped[i]) {
				return false
			}
		}
		return true
	default:
		return a == b
	}
}

// Helper function to compare two Example structs
func compareExamples(t *testing.T, original, unmarshalled *Example) {
	// Compare NumberValue
	if original.NumberValue != nil && unmarshalled.NumberValue != nil {
		if original.NumberValue.IsDiscarded() != unmarshalled.NumberValue.IsDiscarded() {
			t.Errorf("NumberValue discard state mismatch: got %v, want %v",
				unmarshalled.NumberValue.IsDiscarded(), original.NumberValue.IsDiscarded())
		}
		if !original.NumberValue.IsDiscarded() && original.NumberValue.Unwrap() != unmarshalled.NumberValue.Unwrap() {
			t.Errorf("NumberValue mismatch: got %v, want %v",
				unmarshalled.NumberValue.Unwrap(), original.NumberValue.Unwrap())
		}
	}

	// Compare StringValue
	if original.StringValue != nil && unmarshalled.StringValue != nil {
		if original.StringValue.IsDiscarded() != unmarshalled.StringValue.IsDiscarded() {
			t.Errorf("StringValue discard state mismatch: got %v, want %v",
				unmarshalled.StringValue.IsDiscarded(), original.StringValue.IsDiscarded())
		}
		if !original.StringValue.IsDiscarded() && original.StringValue.Unwrap() != unmarshalled.StringValue.Unwrap() {
			t.Errorf("StringValue mismatch: got %v, want %v",
				unmarshalled.StringValue.Unwrap(), original.StringValue.Unwrap())
		}
	}

	// Compare BoolValue
	if original.BoolValue != nil && unmarshalled.BoolValue != nil {
		if original.BoolValue.IsDiscarded() != unmarshalled.BoolValue.IsDiscarded() {
			t.Errorf("BoolValue discard state mismatch: got %v, want %v",
				unmarshalled.BoolValue.IsDiscarded(), original.BoolValue.IsDiscarded())
		}
		if !original.BoolValue.IsDiscarded() && original.BoolValue.Unwrap() != unmarshalled.BoolValue.Unwrap() {
			t.Errorf("BoolValue mismatch: got %v, want %v",
				unmarshalled.BoolValue.Unwrap(), original.BoolValue.Unwrap())
		}
	}

	// Compare FloatValue
	if original.FloatValue != nil && unmarshalled.FloatValue != nil {
		if original.FloatValue.IsDiscarded() != unmarshalled.FloatValue.IsDiscarded() {
			t.Errorf("FloatValue discard state mismatch: got %v, want %v",
				unmarshalled.FloatValue.IsDiscarded(), original.FloatValue.IsDiscarded())
		}
		if !original.FloatValue.IsDiscarded() && original.FloatValue.Unwrap() != unmarshalled.FloatValue.Unwrap() {
			t.Errorf("FloatValue mismatch: got %v, want %v",
				unmarshalled.FloatValue.Unwrap(), original.FloatValue.Unwrap())
		}
	}

	// Compare TimeValue
	if original.TimeValue != nil && unmarshalled.TimeValue != nil {
		if original.TimeValue.IsDiscarded() != unmarshalled.TimeValue.IsDiscarded() {
			t.Errorf("TimeValue discard state mismatch: got %v, want %v",
				unmarshalled.TimeValue.IsDiscarded(), original.TimeValue.IsDiscarded())
		}
		if !original.TimeValue.IsDiscarded() && original.TimeValue.Unwrap() != (unmarshalled.TimeValue.Unwrap()) {
			t.Errorf("TimeValue mismatch: got %v, want %v",
				unmarshalled.TimeValue.Unwrap(), original.TimeValue.Unwrap())
		}
	}

	// Compare Nested.NestedNumberValue
	if original.Nested.NestedNumberValue != nil && unmarshalled.Nested.NestedNumberValue != nil {
		if original.Nested.NestedNumberValue.IsDiscarded() != unmarshalled.Nested.NestedNumberValue.IsDiscarded() {
			t.Errorf("Nested.NestedNumberValue discard state mismatch: got %v, want %v",
				unmarshalled.Nested.NestedNumberValue.IsDiscarded(), original.Nested.NestedNumberValue.IsDiscarded())
		}
		if !original.Nested.NestedNumberValue.IsDiscarded() && original.Nested.NestedNumberValue.Unwrap() != unmarshalled.Nested.NestedNumberValue.Unwrap() {
			t.Errorf("Nested.NestedNumberValue mismatch: got %v, want %v",
				unmarshalled.Nested.NestedNumberValue.Unwrap(), original.Nested.NestedNumberValue.Unwrap())
		}
	}

	// Compare Nested.NestedStringValue
	if original.Nested.NestedStringValue != nil && unmarshalled.Nested.NestedStringValue != nil {
		if original.Nested.NestedStringValue.IsDiscarded() != unmarshalled.Nested.NestedStringValue.IsDiscarded() {
			t.Errorf("Nested.NestedStringValue discard state mismatch: got %v, want %v",
				unmarshalled.Nested.NestedStringValue.IsDiscarded(), original.Nested.NestedStringValue.IsDiscarded())
		}
		if !original.Nested.NestedStringValue.IsDiscarded() && original.Nested.NestedStringValue.Unwrap() != unmarshalled.Nested.NestedStringValue.Unwrap() {
			t.Errorf("Nested.NestedStringValue mismatch: got %v, want %v",
				unmarshalled.Nested.NestedStringValue.Unwrap(), original.Nested.NestedStringValue.Unwrap())
		}
	}

	// Compare Omit Fields
	if original.Omit != nil && unmarshalled.Omit != nil {
		// Compare Omit.OmitNumberValue
		if original.Omit.OmitNumberValue != nil && unmarshalled.Omit.OmitNumberValue != nil {
			if original.Omit.OmitNumberValue.IsDiscarded() != unmarshalled.Omit.OmitNumberValue.IsDiscarded() {
				t.Errorf("Omit.OmitNumberValue discard state mismatch: got %v, want %v",
					unmarshalled.Omit.OmitNumberValue.IsDiscarded(), original.Omit.OmitNumberValue.IsDiscarded())
			}
			if !original.Omit.OmitNumberValue.IsDiscarded() && original.Omit.OmitNumberValue.Unwrap() != unmarshalled.Omit.OmitNumberValue.Unwrap() {
				t.Errorf("Omit.OmitNumberValue mismatch: got %v, want %v",
					unmarshalled.Omit.OmitNumberValue.Unwrap(), original.Omit.OmitNumberValue.Unwrap())
			}
		}

		// Compare Omit.OmitStringValue
		if original.Omit.OmitStringValue != nil && unmarshalled.Omit.OmitStringValue != nil {
			if original.Omit.OmitStringValue.IsDiscarded() != unmarshalled.Omit.OmitStringValue.IsDiscarded() {
				t.Errorf("Omit.OmitStringValue discard state mismatch: got %v, want %v",
					unmarshalled.Omit.OmitStringValue.IsDiscarded(), original.Omit.OmitStringValue.IsDiscarded())
			}
			if !original.Omit.OmitStringValue.IsDiscarded() && original.Omit.OmitStringValue.Unwrap() != unmarshalled.Omit.OmitStringValue.Unwrap() {
				t.Errorf("Omit.OmitStringValue mismatch: got %v, want %v",
					unmarshalled.Omit.OmitStringValue.Unwrap(), original.Omit.OmitStringValue.Unwrap())
			}
		}
	}

	// Compare Discarder Field
	if original.Discarder.DiscarderNumberValue != nil && unmarshalled.Discarder.DiscarderNumberValue != nil {
		discarderOriginal := original.Discarder.DiscarderNumberValue
		discarderUnmarshalled := unmarshalled.Discarder.DiscarderNumberValue

		// Check if the Proxy is nil for both original and unmarshalled Discarder
		if discarderOriginal.Proxy == nil && discarderUnmarshalled.Proxy == nil {
			// Check discard state
			if discarderOriginal.Proxy.IsDiscarded() != discarderUnmarshalled.Proxy.IsDiscarded() {
				t.Errorf("Discarder.DiscarderNumberValue discard state mismatch: got %v, want %v",
					discarderUnmarshalled.Proxy.IsDiscarded(), discarderOriginal.Proxy.IsDiscarded())
			}

			// Check unwrapped value
			if !discarderOriginal.Proxy.IsDiscarded() {
				originalVal := discarderOriginal.Proxy.UnwrapAny()
				unmarshalledVal := discarderUnmarshalled.Proxy.UnwrapAny()
				if originalVal != unmarshalledVal {
					t.Errorf("Discarder.DiscarderNumberValue mismatch: got %v, want %v",
						unmarshalledVal, originalVal)
				}
			}
		}
	}
}

// Helper function to initialize an Example struct with valid data
func initializeValidExample() *Example {
	example := &Example{
		NumberValue: func() *WrapperInt {
			w := New[*WrapperInt]()
			w.Wrap(100, false)
			return w
		}(),
		StringValue: func() *WrapperString {
			w := New[*WrapperString]()
			w.Wrap("Test String", false)
			return w
		}(),
		BoolValue: func() *WrapperBool {
			w := New[*WrapperBool]()
			w.Wrap(true, false)
			return w
		}(),
		FloatValue: func() *WrapperFloat {
			w := New[*WrapperFloat]()
			w.Wrap(123.456, false)
			return w
		}(),
		TimeValue: func() *WrapperTime {
			w := New[*WrapperTime]()
			w.Wrap(time.Date(2025, 1, 16, 6, 34, 8, 685, time.UTC), false)
			return w
		}(),
	}

	// Initialize Nested struct with valid data
	example.Nested.NestedNumberValue = func() *WrapperInt {
		w := New[*WrapperInt]()
		w.Wrap(200, false)
		return w
	}()

	example.Nested.NestedStringValue = func() *WrapperString {
		w := New[*WrapperString]()
		w.Wrap("Nested String", false)
		return w
	}()

	// Initialize Omit struct with valid data
	example.Omit = &ExampleOmit{
		OmitNumberValue: func() *WrapperInt {
			w := New[*WrapperInt]()
			w.Wrap(300, false)
			return w
		}(),
		OmitStringValue: func() *WrapperString {
			w := New[*WrapperString]()
			w.Wrap("Omit String", false)
			return w
		}(),
	}

	// Initialize Discarder with valid data
	discarderWrapper := New[*WrapperInt]()
	discarderWrapper.Wrap(400, false) // Valid value
	example.Discarder.DiscarderNumberValue = NewDiscarder(discarderWrapper)

	return example
}

// Helper function to initialize an Example struct with some invalid and discarded fields
func initializeExampleWithDiscards() *Example {
	example := &Example{
		NumberValue: func() *WrapperInt {
			w := New[*WrapperInt]()
			w.Wrap("invalid_int", true) // Invalid int with discard=true
			return w
		}(),
		StringValue: func() *WrapperString {
			w := New[*WrapperString]()
			w.Wrap("Valid String", false)
			return w
		}(),
		BoolValue: func() *WrapperBool {
			w := New[*WrapperBool]()
			w.Wrap("yes", false) // Valid boolean representation
			return w
		}(),
		FloatValue: func() *WrapperFloat {
			w := New[*WrapperFloat]()
			w.Wrap("invalid_float", true) // Invalid float with discard=true
			return w
		}(),
		TimeValue: func() *WrapperTime {
			w := New[*WrapperTime]()
			w.Wrap("invalid_time_format", true) // Invalid time with discard=true
			return w
		}(),
	}

	// Initialize Nested struct with some invalid data
	example.Nested.NestedNumberValue = func() *WrapperInt {
		w := New[*WrapperInt]()
		w.Wrap(500, false)
		return w
	}()

	example.Nested.NestedStringValue = func() *WrapperString {
		w := New[*WrapperString]()
		w.Wrap("Nested Valid String", false)
		return w
	}()

	// Initialize Omit struct with some invalid data
	example.Omit = &ExampleOmit{
		OmitNumberValue: func() *WrapperInt {
			w := New[*WrapperInt]()
			w.Wrap("invalid_omit_int", true) // Invalid int with discard=true
			return w
		}(),
		// OmitStringValue left as nil to test omitempty
	}

	// Initialize Discarder with invalid data
	discarderWrapper := New[*WrapperInt]()
	discarderWrapper.Wrap("invalid_discarder_int", true) // Invalid int, should be discarded
	example.Discarder.DiscarderNumberValue = NewDiscarder(discarderWrapper)

	return example
}

// TestMarshallingAndUnmarshalling tests marshalling and unmarshalling of the Example struct, including Discarder.
func TestMarshallingAndUnmarshalling(t *testing.T) {
	tests := []struct {
		name           string
		initialExample *Example
		expectedJSON   string
	}{
		{
			name:           "All Fields Valid",
			initialExample: initializeValidExample(),
			expectedJSON: trimJson(`{
				"number_value":100,
				"string_value":"Test String",
				"bool_value":true,
				"float_value":123.456,
				"time_value":"2025-01-16T06:34:08Z",
				"nested":{
					"nested_number_value":200,
					"nested_string_value":"Nested String"
				},
				"omit":{
					"omit_number_value":300,
					"omit_string_value":"Omit String"
				},
				"discarder":{
					"discarder_number_value":400
				}
			}`),
		},
		{
			name:           "Some Fields Invalid and Discarded",
			initialExample: initializeExampleWithDiscards(),
			expectedJSON: trimJson(`{
				"number_value":null,
				"string_value":"Valid String",
				"bool_value":true,
				"float_value":null,
				"time_value":null,
				"nested":{
					"nested_number_value":500,
					"nested_string_value":"Nested Valid String"
				},
				"omit":{
					"omit_number_value":null
				},
				"discarder":{
					"discarder_number_value":null
				}
			}`),
		},
		{
			name: "Omit Struct Completely (All Omit Fields Nil)",
			initialExample: func() *Example {
				example := &Example{
					NumberValue: func() *WrapperInt {
						w := New[*WrapperInt]()
						w.Wrap(50, false)
						return w
					}(),
					StringValue: func() *WrapperString {
						w := New[*WrapperString]()
						w.Wrap("Another Test String", false)
						return w
					}(),
					BoolValue: func() *WrapperBool {
						w := New[*WrapperBool]()
						w.Wrap(false, false)
						return w
					}(),
					FloatValue: func() *WrapperFloat {
						w := New[*WrapperFloat]()
						w.Wrap(789.012, false)
						return w
					}(),
					TimeValue: func() *WrapperTime {
						w := New[*WrapperTime]()
						w.Wrap(time.Date(2025, 1, 16, 6, 34, 8, 685, time.UTC), false)
						return w
					}(),
				}

				// Initialize Nested struct with valid data
				example.Nested.NestedNumberValue = func() *WrapperInt {
					w := New[*WrapperInt]()
					w.Wrap(600, false)
					return w
				}()

				example.Nested.NestedStringValue = func() *WrapperString {
					w := New[*WrapperString]()
					w.Wrap("Nested Complete String", false)
					return w
				}()

				// Omit struct fields left as nil to test omitempty
				// Initialize Discarder with valid data
				discarderWrapper := New[*WrapperInt]()
				discarderWrapper.Wrap(700, false) // Valid value
				example.Discarder.DiscarderNumberValue = NewDiscarder(discarderWrapper)

				return example
			}(),
			expectedJSON: trimJson(`{
				"number_value":50,
				"string_value":"Another Test String",
				"bool_value":false,
				"float_value":789.012,
				"time_value":"2025-01-16T06:34:08Z",
				"nested":{
					"nested_number_value":600,
					"nested_string_value":"Nested Complete String"
				},
				"discarder":{
					"discarder_number_value":700
				}
			}`),
		},
		{
			name: "Discarder Discards Invalid Value",
			initialExample: func() *Example {
				example := initializeValidExample()
				// Overwrite Discarder with invalid data
				example.Discarder.DiscarderNumberValue.Proxy.Wrap("invalid_discarder_int", true)
				return example
			}(),
			expectedJSON: trimJson(`{
				"number_value":100,
				"string_value":"Test String",
				"bool_value":true,
				"float_value":123.456,
				"time_value":"2025-01-16T06:34:08Z",
				"nested":{
					"nested_number_value":200,
					"nested_string_value":"Nested String"
				},
				"omit":{
					"omit_number_value":300,
					"omit_string_value":"Omit String"
				},
				"discarder":{
					"discarder_number_value":null
				}
			}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal the initialExample to JSON
			jsonData, err := json.Marshal(tt.initialExample)
			if err != nil {
				t.Fatalf("Failed to marshal Example struct: %v", err)
			}

			// For better comparison, unmarshal both expectedJSON and actual jsonData into interface{}
			var expected interface{}
			err = json.Unmarshal([]byte(tt.expectedJSON), &expected)
			if err != nil {
				t.Fatalf("Failed to unmarshal expected JSON: %v", err)
			}

			var actual interface{}
			err = json.Unmarshal(jsonData, &actual)
			if err != nil {
				t.Fatalf("Failed to unmarshal actual JSON: %v", err)
			}

			// Compare the two unmarshalled interfaces
			if !compareJSON(expected, actual) {
				t.Errorf("Marshalled JSON mismatch.\nExpected: %v\nActual:   %v", tt.expectedJSON, string(jsonData))
			}

			// Now unmarshal jsonData back into a new Example struct
			var unmarshalled Example
			err = json.Unmarshal(jsonData, &unmarshalled)
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON back into Example struct: %v", err)
			}

			// Compare the initialExample and unmarshalled Example structs
			compareExamples(t, tt.initialExample, &unmarshalled)
		})
	}
}

// TestMarshallingWithOmitFields tests marshalling of Example structs where omit fields are nil or discarded, including Discarder.
func TestMarshallingWithOmitFields(t *testing.T) {
	tests := []struct {
		name           string
		initialExample *Example
		expectedJSON   string
	}{
		{
			name: "Omit Struct with All Omit Fields Set",
			initialExample: func() *Example {
				example := &Example{
					NumberValue: func() *WrapperInt {
						w := New[*WrapperInt]()
						w.Wrap(10, false)
						return w
					}(),
					StringValue: func() *WrapperString {
						w := New[*WrapperString]()
						w.Wrap("Sample String", false)
						return w
					}(),
					BoolValue: func() *WrapperBool {
						w := New[*WrapperBool]()
						w.Wrap(true, false)
						return w
					}(),
					FloatValue: func() *WrapperFloat {
						w := New[*WrapperFloat]()
						w.Wrap(99.99, false)
						return w
					}(),
					TimeValue: func() *WrapperTime {
						w := New[*WrapperTime]()
						w.Wrap(time.Date(2025, 1, 16, 6, 34, 8, 685, time.UTC), false)
						return w
					}(),
				}

				// Initialize Nested struct with valid data
				example.Nested.NestedNumberValue = func() *WrapperInt {
					w := New[*WrapperInt]()
					w.Wrap(20, false)
					return w
				}()

				example.Nested.NestedStringValue = func() *WrapperString {
					w := New[*WrapperString]()
					w.Wrap("Nested Valid String", false)
					return w
				}()

				// Initialize Omit struct with all fields set
				example.Omit = &ExampleOmit{
					OmitNumberValue: func() *WrapperInt {
						w := New[*WrapperInt]()
						w.Wrap(30, false)
						return w
					}(),
					OmitStringValue: func() *WrapperString {
						w := New[*WrapperString]()
						w.Wrap("Omit Valid String", false)
						return w
					}(),
				}

				// Initialize Discarder with valid data
				discarderWrapper := New[*WrapperInt]()
				discarderWrapper.Wrap(40, false) // Valid value
				example.Discarder.DiscarderNumberValue = NewDiscarder(discarderWrapper)

				return example
			}(),
			expectedJSON: trimJson(`{
				"number_value":10,
				"string_value":"Sample String",
				"bool_value":true,
				"float_value":99.99,
				"time_value":"2025-01-16T06:34:08Z",
				"nested":{
					"nested_number_value":20,
					"nested_string_value":"Nested Valid String"
				},
				"omit":{
					"omit_number_value":30,
					"omit_string_value":"Omit Valid String"
				},
				"discarder":{
					"discarder_number_value":40
				}
			}`),
		},
		{
			name: "Omit Struct with Some Omit Fields Discarded",
			initialExample: func() *Example {
				example := &Example{
					NumberValue: func() *WrapperInt {
						w := New[*WrapperInt]()
						w.Wrap(15, false)
						return w
					}(),
					StringValue: func() *WrapperString {
						w := New[*WrapperString]()
						w.Wrap("Another String", false)
						return w
					}(),
					BoolValue: func() *WrapperBool {
						w := New[*WrapperBool]()
						w.Wrap(false, false)
						return w
					}(),
					FloatValue: func() *WrapperFloat {
						w := New[*WrapperFloat]()
						w.Wrap(50.5, false)
						return w
					}(),
					TimeValue: func() *WrapperTime {
						w := New[*WrapperTime]()
						w.Wrap(time.Date(2025, 1, 16, 6, 34, 8, 685, time.UTC), false)
						return w
					}(),
				}

				// Initialize Nested struct with some invalid data
				example.Nested.NestedNumberValue = func() *WrapperInt {
					w := New[*WrapperInt]()
					w.Wrap("invalid_nested_int", true) // Invalid int with discard=true
					return w
				}()

				example.Nested.NestedStringValue = func() *WrapperString {
					w := New[*WrapperString]()
					w.Wrap("Nested Another String", false)
					return w
				}()

				// Initialize Omit struct with some fields discarded
				example.Omit = &ExampleOmit{
					OmitNumberValue: func() *WrapperInt {
						w := New[*WrapperInt]()
						w.Wrap("invalid_omit_int", true) // Invalid int with discard=true
						return w
					}(),
					// OmitStringValue left as nil to test omitempty
				}

				// Initialize Discarder with invalid data
				discarderWrapper := New[*WrapperInt]()
				discarderWrapper.Wrap("invalid_discarder_int", true) // Invalid int, should be discarded
				example.Discarder.DiscarderNumberValue = NewDiscarder(discarderWrapper)

				return example
			}(),
			expectedJSON: `{
				"number_value":15,
				"string_value":"Another String",
				"bool_value":false,
				"float_value":50.5,
				"time_value":"2025-01-16T06:34:08Z",
				"nested":{
					"nested_number_value":null,
					"nested_string_value":"Nested Another String"
				},
				"omit":{
					"omit_number_value":null
				},
				"discarder":{
					"discarder_number_value":null
				}
			}`,
		},
		{
			name: "Omit Struct Fully Omitted (All Omit Fields Discarded)",
			initialExample: func() *Example {
				example := &Example{
					NumberValue: func() *WrapperInt {
						w := New[*WrapperInt]()
						w.Wrap(12, false)
						return w
					}(),
					StringValue: func() *WrapperString {
						w := New[*WrapperString]()
						w.Wrap("Full Omit Test", false)
						return w
					}(),
					BoolValue: func() *WrapperBool {
						w := New[*WrapperBool]()
						w.Wrap(false, false)
						return w
					}(),
					FloatValue: func() *WrapperFloat {
						w := New[*WrapperFloat]()
						w.Wrap(88.88, false)
						return w
					}(),
					TimeValue: func() *WrapperTime {
						w := New[*WrapperTime]()
						w.Wrap(time.Date(2025, 1, 16, 6, 34, 8, 0, time.UTC), false)
						return w
					}(),
				}

				// Initialize Nested struct with valid data
				example.Nested.NestedNumberValue = func() *WrapperInt {
					w := New[*WrapperInt]()
					w.Wrap(24, false)
					return w
				}()

				example.Nested.NestedStringValue = func() *WrapperString {
					w := New[*WrapperString]()
					w.Wrap("Nested Full Omit Test", false)
					return w
				}()

				// Initialize Omit struct with all fields discarded
				example.Omit = &ExampleOmit{
					OmitNumberValue: func() *WrapperInt {
						w := New[*WrapperInt]()
						w.Wrap("invalid_omit_int", true) // Invalid int with discard=true
						return w
					}(),
					OmitStringValue: func() *WrapperString {
						w := New[*WrapperString]()
						w.Wrap(nil, true) // Invalid string with discard=true
						return w
					}(),
				}

				// Initialize Discarder with invalid data
				discarderWrapper := New[*WrapperInt]()
				discarderWrapper.Wrap("invalid_discarder_int", true) // Invalid int, should be discarded
				example.Discarder.DiscarderNumberValue = NewDiscarder(discarderWrapper)

				return example
			}(),
			expectedJSON: trimJson(`{
				"number_value":12,
				"string_value":"Full Omit Test",
				"bool_value":false,
				"float_value":88.88,
				"time_value":"2025-01-16T06:34:08Z",
				"nested":{
					"nested_number_value":24,
					"nested_string_value":"Nested Full Omit Test"
				},
				"omit":{
					"omit_number_value":null,
					"omit_string_value":null
				},
				"discarder":{
					"discarder_number_value":null
				}
			}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal the initialExample to JSON
			jsonData, err := json.Marshal(tt.initialExample)
			if err != nil {
				t.Fatalf("Failed to marshal Example struct: %v", err)
			}

			// Unmarshal both expectedJSON and actual jsonData into interface{}
			var expected interface{}
			err = json.Unmarshal([]byte(tt.expectedJSON), &expected)
			if err != nil {
				t.Fatalf("Failed to unmarshal expected JSON: %v", err)
			}

			var actual interface{}
			err = json.Unmarshal(jsonData, &actual)
			if err != nil {
				t.Fatalf("Failed to unmarshal actual JSON: %v", err)
			}

			// Compare the two unmarshalled interfaces
			if !compareJSON(expected, actual) {
				t.Errorf("Marshalled JSON mismatch.\nExpected: %v\nActual:   %v", tt.expectedJSON, string(jsonData))
			}

			// Now unmarshal jsonData back into a new Example struct
			var unmarshalled Example
			err = json.Unmarshal(jsonData, &unmarshalled)
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON back into Example struct: %v", err)
			}

			// Compare the initialExample and unmarshalled Example structs
			compareExamples(t, tt.initialExample, &unmarshalled)
		})
	}
}

// TestOmitFieldsOmittedWhenNilOrDiscarded tests that fields with omitempty are correctly omitted from JSON when nil or discarded, including Discarder.
func TestOmitFieldsOmittedWhenNilOrDiscarded(t *testing.T) {
	tests := []struct {
		name           string
		initialExample *Example
		expectedJSON   string
	}{
		{
			name: "Omit Struct Completely (All Omit Fields Nil)",
			initialExample: func() *Example {
				example := &Example{
					NumberValue: func() *WrapperInt {
						w := New[*WrapperInt]()
						w.Wrap(5, false)
						return w
					}(),
					StringValue: func() *WrapperString {
						w := New[*WrapperString]()
						w.Wrap("Omit Test", false)
						return w
					}(),
					BoolValue: func() *WrapperBool {
						w := New[*WrapperBool]()
						w.Wrap(false, false)
						return w
					}(),
					FloatValue: func() *WrapperFloat {
						w := New[*WrapperFloat]()
						w.Wrap(55.55, false)
						return w
					}(),
					TimeValue: func() *WrapperTime {
						w := New[*WrapperTime]()
						w.Wrap(time.Date(2025, 1, 16, 6, 34, 8, 0, time.UTC), false)
						return w
					}(),
				}

				// Initialize Nested struct with valid data
				example.Nested.NestedNumberValue = func() *WrapperInt {
					w := New[*WrapperInt]()
					w.Wrap(10, false)
					return w
				}()

				example.Nested.NestedStringValue = func() *WrapperString {
					w := New[*WrapperString]()
					w.Wrap("Nested Omit Test", false)
					return w
				}()

				// Omit struct left with nil fields to test omitempty

				// Initialize Discarder with discarded value
				discarderWrapper := New[*WrapperInt]()
				discarderWrapper.Wrap(nil, true) // Invalid value, should be discarded
				example.Discarder.DiscarderNumberValue = NewDiscarder(discarderWrapper)

				return example
			}(),
			expectedJSON: trimJson(`{
				"number_value":5,
				"string_value":"Omit Test",
				"bool_value":false,
				"float_value":55.55,
				"time_value":"2025-01-16T06:34:08Z",
				"nested":{
					"nested_number_value":10,
					"nested_string_value":"Nested Omit Test"
				},
				"discarder":{
					"discarder_number_value":null
				}
			}`),
		},
		{
			name: "Omit Struct Partially Omitted (One Omit Field Set)",
			initialExample: func() *Example {
				example := &Example{
					NumberValue: func() *WrapperInt {
						w := New[*WrapperInt]()
						w.Wrap(7, false)
						return w
					}(),
					StringValue: func() *WrapperString {
						w := New[*WrapperString]()
						w.Wrap("Partial Omit Test", false)
						return w
					}(),
					BoolValue: func() *WrapperBool {
						w := New[*WrapperBool]()
						w.Wrap(true, false)
						return w
					}(),
					FloatValue: func() *WrapperFloat {
						w := New[*WrapperFloat]()
						w.Wrap(77.77, false)
						return w
					}(),
					TimeValue: func() *WrapperTime {
						w := New[*WrapperTime]()
						w.Wrap(time.Date(2025, 1, 16, 6, 34, 8, 0, time.UTC), false)
						return w
					}(),
				}

				// Initialize Nested struct with valid data
				example.Nested.NestedNumberValue = func() *WrapperInt {
					w := New[*WrapperInt]()
					w.Wrap(14, false)
					return w
				}()

				example.Nested.NestedStringValue = func() *WrapperString {
					w := New[*WrapperString]()
					w.Wrap("Nested Partial Omit Test", false)
					return w
				}()

				// Initialize Omit struct with only OmitNumberValue set
				example.Omit = &ExampleOmit{
					OmitNumberValue: func() *WrapperInt {
						w := New[*WrapperInt]()
						w.Wrap(21, false)
						return w
					}(),
				}

				// Initialize Discarder with discarded value
				discarderWrapper := New[*WrapperInt]()
				discarderWrapper.Wrap("invalid_discarder_int", true) // Invalid int, should be discarded
				example.Discarder.DiscarderNumberValue = NewDiscarder(discarderWrapper)

				return example
			}(),
			expectedJSON: trimJson(`{
				"number_value":7,
				"string_value":"Partial Omit Test",
				"bool_value":true,
				"float_value":77.77,
				"time_value":"2025-01-16T06:34:08Z",
				"nested":{
					"nested_number_value":14,
					"nested_string_value":"Nested Partial Omit Test"
				},
				"omit":{
					"omit_number_value":21
				},
				"discarder":{
					"discarder_number_value":null
				}
			}`),
		},
		{
			name: "Omit Struct Fully Omitted (All Omit Fields Discarded)",
			initialExample: func() *Example {
				example := &Example{
					NumberValue: func() *WrapperInt {
						w := New[*WrapperInt]()
						w.Wrap(12, false)
						return w
					}(),
					StringValue: func() *WrapperString {
						w := New[*WrapperString]()
						w.Wrap("Full Omit Test", false)
						return w
					}(),
					BoolValue: func() *WrapperBool {
						w := New[*WrapperBool]()
						w.Wrap(false, false)
						return w
					}(),
					FloatValue: func() *WrapperFloat {
						w := New[*WrapperFloat]()
						w.Wrap(88.88, false)
						return w
					}(),
					TimeValue: func() *WrapperTime {
						w := New[*WrapperTime]()
						w.Wrap(time.Date(2025, 1, 16, 6, 34, 8, 0, time.UTC), false)
						return w
					}(),
				}

				// Initialize Nested struct with valid data
				example.Nested.NestedNumberValue = func() *WrapperInt {
					w := New[*WrapperInt]()
					w.Wrap(24, false)
					return w
				}()

				example.Nested.NestedStringValue = func() *WrapperString {
					w := New[*WrapperString]()
					w.Wrap("Nested Full Omit Test", false)
					return w
				}()

				// Initialize Omit struct with all fields discarded
				example.Omit = &ExampleOmit{
					OmitNumberValue: func() *WrapperInt {
						w := New[*WrapperInt]()
						w.Wrap("invalid_omit_int", true) // Invalid int with discard=true
						return w
					}(),
					OmitStringValue: func() *WrapperString {
						w := New[*WrapperString]()
						w.Wrap(nil, true) // Invalid string with discard=true
						return w
					}(),
				}

				// Initialize Discarder with invalid data
				discarderWrapper := New[*WrapperInt]()
				discarderWrapper.Wrap("invalid_discarder_int", true) // Invalid int, should be discarded
				example.Discarder.DiscarderNumberValue = NewDiscarder(discarderWrapper)

				return example
			}(),
			expectedJSON: trimJson(`{
				"number_value":12,
				"string_value":"Full Omit Test",
				"bool_value":false,
				"float_value":88.88,
				"time_value":"2025-01-16T06:34:08Z",
				"nested":{
					"nested_number_value":24,
					"nested_string_value":"Nested Full Omit Test"
				},
				"omit":{
					"omit_number_value":null,
					"omit_string_value":null
				},
				"discarder":{
					"discarder_number_value":null
				}
			}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal the initialExample to JSON
			jsonData, err := json.Marshal(tt.initialExample)
			if err != nil {
				t.Fatalf("Failed to marshal Example struct: %v", err)
			}

			// Unmarshal both expectedJSON and actual jsonData into interface{}
			var expected interface{}
			err = json.Unmarshal([]byte(tt.expectedJSON), &expected)
			if err != nil {
				t.Fatalf("Failed to unmarshal expected JSON: %v", err)
			}

			var actual interface{}
			err = json.Unmarshal(jsonData, &actual)
			if err != nil {
				t.Fatalf("Failed to unmarshal actual JSON: %v", err)
			}

			// Compare the two unmarshalled interfaces
			if !compareJSON(expected, actual) {
				t.Errorf("Marshalled JSON mismatch.\nExpected: %v\nActual:   %v", tt.expectedJSON, string(jsonData))
			}

			// Now unmarshal jsonData back into a new Example struct
			var unmarshalled Example
			err = json.Unmarshal(jsonData, &unmarshalled)
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON back into Example struct: %v", err)
			}

			// Compare the initialExample and unmarshalled Example structs
			compareExamples(t, tt.initialExample, &unmarshalled)
		})
	}
}

// TestExample_FullyDiscardedFields tests marshalling and unmarshalling when all fields are discarded, including Discarder.
func TestExample_FullyDiscardedFields(t *testing.T) {
	example := &Example{
		NumberValue: func() *WrapperInt {
			w := New[*WrapperInt]()
			w.Wrap("invalid_int", true) // Discarded
			return w
		}(),
		StringValue: func() *WrapperString {
			w := New[*WrapperString]()
			w.Wrap(12345, false)
			return w
		}(),
		BoolValue: func() *WrapperBool {
			w := New[*WrapperBool]()
			w.Wrap("no", false)
			return w
		}(),
		FloatValue: func() *WrapperFloat {
			w := New[*WrapperFloat]()
			w.Wrap("invalid_float", true) // Discarded
			return w
		}(),
		TimeValue: func() *WrapperTime {
			w := New[*WrapperTime]()
			w.Wrap("invalid_time_format", true) // Discarded
			return w
		}(),
	}

	// Initialize Nested struct with all fields discarded
	example.Nested.NestedNumberValue = func() *WrapperInt {
		w := New[*WrapperInt]()
		w.Wrap("invalid_nested_int", true) // Discarded
		return w
	}()

	example.Nested.NestedStringValue = func() *WrapperString {
		w := New[*WrapperString]()
		w.Wrap(nil, true) // Discarded
		return w
	}()

	// Initialize Omit struct with all fields discarded
	example.Omit = &ExampleOmit{
		OmitNumberValue: func() *WrapperInt {
			w := New[*WrapperInt]()
			w.Wrap("invalid_omit_int", true) // Discarded
			return w
		}(),
		OmitStringValue: func() *WrapperString {
			w := New[*WrapperString]()
			w.Wrap(nil, true) // Discarded
			return w
		}(),
	}

	// Initialize Discarder with invalid data
	discarderWrapper := New[*WrapperInt]()
	discarderWrapper.Wrap("invalid_discarder_int", true) // Invalid int, should be discarded
	example.Discarder.DiscarderNumberValue = NewDiscarder(discarderWrapper)

	// Marshal to JSON
	jsonData, err := json.Marshal(example)
	if err != nil {
		t.Fatalf("Failed to marshal Example struct with fully discarded fields: %v", err)
	}

	// Expected JSON:
	expectedJSON := `{
		"number_value":null,
		"string_value":"12345",
		"bool_value":false,
		"float_value":null,
		"time_value":null,
		"nested":{
			"nested_number_value":null,
			"nested_string_value":null
		},
		"omit":{
			"omit_number_value":null,
			"omit_string_value":null
		},
		"discarder":{
			"discarder_number_value":null
		}
	}`

	// Unmarshal both expectedJSON and actual jsonData into interface{}
	var expected interface{}
	err = json.Unmarshal([]byte(expectedJSON), &expected)
	if err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}

	var actual interface{}
	err = json.Unmarshal(jsonData, &actual)
	if err != nil {
		t.Fatalf("Failed to unmarshal actual JSON: %v", err)
	}

	// Compare the two unmarshalled interfaces
	if !compareJSON(expected, actual) {
		t.Errorf("Marshalled JSON mismatch.\nExpected: %v\nActual:   %v", expectedJSON, string(jsonData))
	}

	// Now unmarshal jsonData back into a new Example struct
	var unmarshalled Example
	err = json.Unmarshal(jsonData, &unmarshalled)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON back into Example struct: %v", err)
	}

	// Compare each field to ensure they're correctly discarded or unmarshalled
	compareExamples(t, example, &unmarshalled)
}
