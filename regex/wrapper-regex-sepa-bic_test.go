package regex

import (
	"encoding/json"
	"testing"

	"github.com/zealsprince/wrappers"
)

// TestWrapperRegexSepaBic_Wrap tests the Wrap method of WrapperRegexSepaBic.
func TestWrapperRegexSepaBic_Wrap(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		discard     bool
		want        string
		wantError   bool
		wantDiscard bool
	}{
		{
			name:        "Valid SEPA BIC",
			input:       "DEUTDEFF",
			discard:     false,
			want:        "DEUTDEFF",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Invalid SEPA BIC - too short",
			input:       "DEUTDE",
			discard:     false,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Invalid SEPA BIC - invalid characters",
			input:       "DEUTD3FF",
			discard:     false,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Invalid SEPA BIC discard",
			input:       "DEUTD3FF",
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
			input: func() *WrapperRegexSepaBic {
				w := wrappers.New[*WrapperRegexSepaBic]()
				w.Wrap("DEUTDEFF", false)
				return w
			}(),
			discard:     false,
			want:        "DEUTDEFF",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name: "Input as WrapperProvider with invalid value (discard)",
			input: func() *WrapperRegexSepaBic {
				w := wrappers.New[*WrapperRegexSepaBic]()
				w.Wrap("DEUTD3FFFFFFFFFF", false)
				return w
			}(),
			discard:     false,
			want:        "",
			wantError:   false,
			wantDiscard: true,
		},
		{
			name:        "Unsupported Type without discard",
			input:       12345, // int is unsupported
			discard:     false,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Unsupported Type with discard",
			input:       12345, // int is unsupported
			discard:     true,
			want:        "",
			wantError:   false,
			wantDiscard: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the wrapper
			sepaBicWrapper := wrappers.New[*WrapperRegexSepaBic]()
			sepaBicWrapper.Initialize()

			// Wrap the input
			err := sepaBicWrapper.Wrap(tt.input, tt.discard)
			if (err != nil) != tt.wantError {
				t.Fatalf("Wrap() error = %v, wantError %v", err, tt.wantError)
			}

			// Check if wrapper is discarded as expected
			if sepaBicWrapper.IsDiscarded() != tt.wantDiscard {
				t.Errorf("IsDiscarded() = %v, want %v", sepaBicWrapper.IsDiscarded(), tt.wantDiscard)
			}

			// If not discarded, check the unwrapped value
			if !sepaBicWrapper.IsDiscarded() && sepaBicWrapper.Get() != tt.want {
				t.Errorf("Get() = %v, want %v", sepaBicWrapper.Get(), tt.want)
			}
		})
	}
}

// TestWrapperRegexSepaBic_JSONMarshal tests JSON marshalling of WrapperRegexSepaBic.
func TestWrapperRegexSepaBic_JSONMarshal(t *testing.T) {
	tests := []struct {
		name         string
		sepaBic      string
		wantError    bool
		expectedJSON string
	}{
		{
			name:         "Valid SEPA BIC Marshalling",
			sepaBic:      "DEUTDEFF",
			wantError:    false,
			expectedJSON: `{"sepa_bic":"DEUTDEFF"}`,
		},
		{
			name:         "Invalid SEPA BIC Marshalling",
			sepaBic:      "DEUTD3FF",
			wantError:    true,
			expectedJSON: `{"sepa_bic":null}`,
		},
	}

	// Create a struct for testing
	type Data struct {
		SepaBic *WrapperRegexSepaBic `json:"sepa_bic"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the SEPA BIC wrapper
			sepaBicWrapper := wrappers.New[*WrapperRegexSepaBic]()
			sepaBicWrapper.Initialize()

			// Wrap the SEPA BIC
			err := sepaBicWrapper.Wrap(tt.sepaBic, false)
			if (err != nil) != tt.wantError {
				t.Fatalf("Wrap() error = %v, wantError %v", err, tt.wantError)
			}

			// Create the data struct
			data := Data{
				SepaBic: sepaBicWrapper,
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

// TestWrapperRegexSepaBic_Unmarshal tests JSON unmarshalling of WrapperRegexSepaBic.
func TestWrapperRegexSepaBic_Unmarshal(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		want        string
		wantError   bool
		wantDiscard bool
	}{
		{
			name:        "Unmarshal valid SEPA BIC",
			jsonInput:   `{"sepa_bic":"DEUTDEFF"}`,
			want:        "DEUTDEFF",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Unmarshal invalid SEPA BIC",
			jsonInput:   `{"sepa_bic":"DEUTD3FF"}`,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Unmarshal null SEPA BIC",
			jsonInput:   `{"sepa_bic":null}`,
			want:        "",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Unmarshal invalid type SEPA BIC",
			jsonInput:   `{"sepa_bic":12345}`,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
	}

	// Create a struct for testing
	type Data struct {
		SepaBic *WrapperRegexSepaBic `json:"sepa_bic"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data Data

			// Initialize the SEPA BIC wrapper
			data.SepaBic = wrappers.New[*WrapperRegexSepaBic]()
			data.SepaBic.Initialize()

			// Unmarshal the JSON input
			err := json.Unmarshal([]byte(tt.jsonInput), &data)
			if (err != nil) != tt.wantError {
				t.Fatalf("Unmarshal() error = %v, wantError %v", err, tt.wantError)
			}

			// Validate that the incoming data actually includes the SEPA BIC value.
			if data.SepaBic != nil {
				// Check if wrapper is discarded as expected
				if data.SepaBic.IsDiscarded() != tt.wantDiscard {
					t.Errorf("IsDiscarded() = %v, want %v", data.SepaBic.IsDiscarded(), tt.wantDiscard)
				}

				// If not discarded, check the unwrapped value
				if !data.SepaBic.IsDiscarded() && data.SepaBic.Get() != tt.want {
					t.Errorf("Get() = %v, want %v", data.SepaBic.Get(), tt.want)
				}
			}
		})
	}
}
