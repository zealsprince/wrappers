package regex

import (
	"encoding/json"
	"testing"

	"github.com/zealsprince/wrappers"
)

// TestWrapperRegexSepaIban_Wrap tests the Wrap method of WrapperRegexSepaIban.
func TestWrapperRegexSepaIban_Wrap(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		discard     bool
		want        string
		wantError   bool
		wantDiscard bool
	}{
		{
			name:        "Valid SEPA IBAN",
			input:       "DE89370400440532013000",
			discard:     false,
			want:        "DE89370400440532013000",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Invalid SEPA IBAN - invalid characters",
			input:       "DE8937@400440532013000",
			discard:     false,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Invalid SEPA IBAN discard",
			input:       "DE8937@400440532013000",
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
			wantError:   false,
			wantDiscard: true,
		},
		{
			name: "Input as WrapperProvider with valid value",
			input: func() *WrapperRegexSepaIban {
				w := wrappers.New[*WrapperRegexSepaIban]()
				w.Wrap("DE89370400440532013000", false)
				return w
			}(),
			discard:     false,
			want:        "DE89370400440532013000",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name: "Input as WrapperProvider with invalid value",
			input: func() *WrapperRegexSepaIban {
				w := wrappers.New[*WrapperRegexSepaIban]()
				w.Wrap("DE8937@400440532013000", false)
				return w
			}(),
			discard:     false,
			want:        "",
			wantError:   false,
			wantDiscard: true,
		},
		{
			name:        "Unsupported Type without discard",
			input:       67890, // int is unsupported
			discard:     false,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Unsupported Type with discard",
			input:       67890, // int is unsupported
			discard:     true,
			want:        "",
			wantError:   false,
			wantDiscard: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the wrapper
			sepaIbanWrapper := wrappers.New[*WrapperRegexSepaIban]()
			sepaIbanWrapper.Initialize()

			// Wrap the input
			err := sepaIbanWrapper.Wrap(tt.input, tt.discard)
			if (err != nil) != tt.wantError {
				t.Fatalf("Wrap() error = %v, wantError %v", err, tt.wantError)
			}

			// Check if wrapper is discarded as expected
			if sepaIbanWrapper.IsDiscarded() != tt.wantDiscard {
				t.Errorf("IsDiscarded() = %v, want %v", sepaIbanWrapper.IsDiscarded(), tt.wantDiscard)
			}

			// If not discarded, check the unwrapped value
			if !sepaIbanWrapper.IsDiscarded() && sepaIbanWrapper.Get() != tt.want {
				t.Errorf("Get() = %v, want %v", sepaIbanWrapper.Get(), tt.want)
			}
		})
	}
}

// TestWrapperRegexSepaIban_JSONMarshal tests JSON marshalling of WrapperRegexSepaIban.
func TestWrapperRegexSepaIban_JSONMarshal(t *testing.T) {
	tests := []struct {
		name         string
		sepaIban     string
		wantError    bool
		expectedJSON string
	}{
		{
			name:         "Valid SEPA IBAN Marshalling",
			sepaIban:     "DE89370400440532013000",
			wantError:    false,
			expectedJSON: `{"sepa_iban":"DE89370400440532013000"}`,
		},
		{
			name:         "Invalid SEPA IBAN Marshalling",
			sepaIban:     "DE8937@400440532013000",
			wantError:    true,
			expectedJSON: `{"sepa_iban":null}`,
		},
	}

	// Create a struct for testing
	type Data struct {
		SepaIban *WrapperRegexSepaIban `json:"sepa_iban"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the SEPA IBAN wrapper
			sepaIbanWrapper := wrappers.New[*WrapperRegexSepaIban]()
			sepaIbanWrapper.Initialize()

			// Wrap the SEPA IBAN
			err := sepaIbanWrapper.Wrap(tt.sepaIban, false)
			if (err != nil) != tt.wantError {
				t.Fatalf("Wrap() error = %v, wantError %v", err, tt.wantError)
			}

			// Create the data struct
			data := Data{
				SepaIban: sepaIbanWrapper,
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

// TestWrapperRegexSepaIban_Unmarshal tests JSON unmarshalling of WrapperRegexSepaIban.
func TestWrapperRegexSepaIban_Unmarshal(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		want        string
		wantError   bool
		wantDiscard bool
	}{
		{
			name:        "Unmarshal valid SEPA IBAN",
			jsonInput:   `{"sepa_iban":"DE89370400440532013000"}`,
			want:        "DE89370400440532013000",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Unmarshal invalid SEPA IBAN",
			jsonInput:   `{"sepa_iban":"DE8937@400440532013000"}`,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Unmarshal null SEPA IBAN",
			jsonInput:   `{"sepa_iban":null}`,
			want:        "",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Unmarshal invalid type SEPA IBAN",
			jsonInput:   `{"sepa_iban":12345}`,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
	}

	// Create a struct for testing
	type Data struct {
		SepaIban *WrapperRegexSepaIban `json:"sepa_iban"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data Data

			// Initialize the SEPA IBAN wrapper
			data.SepaIban = wrappers.New[*WrapperRegexSepaIban]()
			data.SepaIban.Initialize()

			// Unmarshal the JSON input
			err := json.Unmarshal([]byte(tt.jsonInput), &data)
			if (err != nil) != tt.wantError {
				t.Fatalf("Unmarshal() error = %v, wantError %v", err, tt.wantError)
			}

			// Handle the case where the SEPA IBAN is not in the payload and as such is nil.
			if data.SepaIban != nil {
				// Check if wrapper is discarded as expected
				if data.SepaIban.IsDiscarded() != tt.wantDiscard {
					t.Errorf("IsDiscarded() = %v, want %v", data.SepaIban.IsDiscarded(), tt.wantDiscard)
				}

				// If not discarded, check the unwrapped value
				if !data.SepaIban.IsDiscarded() && data.SepaIban.Get() != tt.want {
					t.Errorf("Get() = %v, want %v", data.SepaIban.Get(), tt.want)
				}
			}
		})
	}
}
