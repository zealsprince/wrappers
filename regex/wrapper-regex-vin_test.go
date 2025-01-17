package regex

import (
	"encoding/json"
	"testing"

	"github.com/zealsprince/wrappers"
)

// TestWrapperRegexVin_Wrap tests the Wrap method of WrapperRegexVin.
func TestWrapperRegexVin_Wrap(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		discard     bool
		want        string
		wantError   bool
		wantDiscard bool
	}{
		{
			name:        "Valid VIN",
			input:       "1HGCM82633A004352",
			discard:     false,
			want:        "1HGCM82633A004352",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Invalid VIN - too short",
			input:       "1HGCM82633A00435",
			discard:     false,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Invalid VIN - invalid characters",
			input:       "1HGCM82633A00435I", // 'I' is not allowed in VIN
			discard:     false,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Invalid VIN discard",
			input:       "1HGCM82633A00435I",
			discard:     true,
			want:        "",
			wantError:   false,
			wantDiscard: true,
		},
		{
			name:        "Nil Input",
			input:       nil,
			discard:     false,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name: "Input as WrapperProvider with valid value",
			input: func() *WrapperRegexVin {
				w := wrappers.New[*WrapperRegexVin]()
				w.Wrap("1HGCM82633A004352", false)
				return w
			}(),
			discard:     false,
			want:        "1HGCM82633A004352",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name: "Input as WrapperProvider with invalid value (discarded)",
			input: func() *WrapperRegexVin {
				w := wrappers.New[*WrapperRegexVin]()
				w.Wrap("1HGCM82633A00435I", false)
				return w
			}(),
			discard:     false,
			want:        "",
			wantError:   false,
			wantDiscard: true,
		},
		{
			name:        "Unsupported Type without discard",
			input:       true, // bool is unsupported
			discard:     false,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Unsupported Type with discard",
			input:       true, // bool is unsupported
			discard:     true,
			want:        "",
			wantError:   false,
			wantDiscard: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the VIN wrapper
			vinWrapper := wrappers.New[*WrapperRegexVin]()
			vinWrapper.Initialize()

			// Wrap the input
			err := vinWrapper.Wrap(tt.input, tt.discard)
			if (err != nil) != tt.wantError {
				t.Fatalf("Wrap() error = %v, wantError %v", err, tt.wantError)
			}

			// Check if wrapper is discarded as expected
			if vinWrapper.IsDiscarded() != tt.wantDiscard {
				t.Errorf("IsDiscarded() = %v, want %v", vinWrapper.IsDiscarded(), tt.wantDiscard)
			}

			// If not discarded, check the unwrapped value
			if !vinWrapper.IsDiscarded() && vinWrapper.Get() != tt.want {
				t.Errorf("Get() = %v, want %v", vinWrapper.Get(), tt.want)
			}
		})
	}
}

// TestWrapperRegexVin_JSONMarshal tests JSON marshalling of WrapperRegexVin.
func TestWrapperRegexVin_JSONMarshal(t *testing.T) {
	tests := []struct {
		name         string
		vin          string
		wantError    bool
		expectedJSON string
	}{
		{
			name:         "Valid VIN Marshalling",
			vin:          "1HGCM82633A004352",
			wantError:    false,
			expectedJSON: `{"vin":"1HGCM82633A004352"}`,
		},
		{
			name:         "Invalid VIN Marshalling",
			vin:          "1HGCM82633A00435I",
			wantError:    true,
			expectedJSON: `{"vin":null}`,
		},
	}

	// Create a struct for testing
	type Data struct {
		Vin *WrapperRegexVin `json:"vin"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the VIN wrapper
			vinWrapper := wrappers.New[*WrapperRegexVin]()
			vinWrapper.Initialize()

			// Wrap the VIN
			err := vinWrapper.Wrap(tt.vin, false)
			if (err != nil) != tt.wantError {
				t.Fatalf("Wrap() error = %v, wantError %v", err, tt.wantError)
			}

			// Create the data struct
			data := Data{
				Vin: vinWrapper,
			}

			// Marshal to JSON
			jsonData, err := json.Marshal(&data)
			if err != nil {
				t.Fatalf("Failed to marshal Data: %v", err)
			}

			// Compare with expected JSON
			if string(jsonData) != tt.expectedJSON {
				t.Errorf("Marshalled JSON = %v, want %v", string(jsonData), tt.expectedJSON)
			}
		})
	}
}

// TestWrapperRegexVin_Unmarshal tests JSON unmarshalling of WrapperRegexVin.
func TestWrapperRegexVin_Unmarshal(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		want        string
		wantError   bool
		wantDiscard bool
	}{
		{
			name:        "Unmarshal valid VIN",
			jsonInput:   `{"vin":"1HGCM82633A004352"}`,
			want:        "1HGCM82633A004352",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Unmarshal invalid VIN",
			jsonInput:   `{"vin":"1HGCM82633A00435I"}`,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Unmarshal null VIN",
			jsonInput:   `{"vin":null}`,
			want:        "",
			wantError:   false,
			wantDiscard: true,
		},
		{
			name:        "Unmarshal unsupported type VIN",
			jsonInput:   `{"vin":true}`,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
	}

	// Create a struct for testing
	type Data struct {
		Vin *WrapperRegexVin `json:"vin"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data Data

			// Initialize the VIN wrapper
			data.Vin = wrappers.New[*WrapperRegexVin]()
			data.Vin.Initialize()

			// Unmarshal the JSON input
			err := json.Unmarshal([]byte(tt.jsonInput), &data)
			if (err != nil) != tt.wantError {
				t.Fatalf("Unmarshal() error = %v, wantError %v", err, tt.wantError)
			}

			// Make sure that the VIN wrapper is not nil in cases where the input does not contain a VIN
			if data.Vin != nil {
				// Check if wrapper is discarded as expected
				if data.Vin.IsDiscarded() != tt.wantDiscard {
					t.Errorf("IsDiscarded() = %v, want %v", data.Vin.IsDiscarded(), tt.wantDiscard)
				}

				// If not discarded, check the unwrapped value
				if !data.Vin.IsDiscarded() && data.Vin.Get() != tt.want {
					t.Errorf("Get() = %v, want %v", data.Vin.Get(), tt.want)
				}
			}
		})
	}
}
