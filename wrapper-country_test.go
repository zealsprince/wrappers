// wrapper-country_test.go
package wrappers

import (
	"encoding/json"
	"testing"

	"github.com/biter777/countries"
)

// TestWrapperCountry_Wrap tests the Wrap method of WrapperCountry.
func TestWrapperCountry_Wrap(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		discard     bool
		want        string
		wantError   bool
		wantDiscard bool
	}{
		{
			name:        "Wrap valid CountryCode",
			input:       countries.DE,
			discard:     false,
			want:        "Germany",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap valid Country",
			input:       countries.ByName("Germany"),
			discard:     false,
			want:        "Germany",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap valid Country by string",
			input:       "Canada",
			discard:     false,
			want:        "Canada",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap invalid Country",
			input:       countries.Unknown,
			discard:     false,
			want:        "Unknown",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Wrap invalid Country discard",
			input:       countries.Unknown,
			discard:     true,
			want:        "Unknown",
			wantError:   false,
			wantDiscard: true,
		},
		{
			name:        "Wrap invalid string",
			input:       "invalid_country",
			discard:     false,
			want:        "Unknown",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Wrap nil",
			input:       nil,
			discard:     false,
			want:        "Unknown",
			wantError:   false,
			wantDiscard: true,
		},
		{
			name: "Input as WrapperProvider with valid value",
			input: func() *WrapperCountry {
				w := New[*WrapperCountry]()
				w.Wrap(countries.US, false)
				return w
			}(),
			discard:     false,
			want:        "United States",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name: "Input as WrapperProvider with invalid value (discarded)",
			input: func() *WrapperCountry {
				w := New[*WrapperCountry]()
				w.Wrap(countries.Unknown, false)
				return w
			}(),
			discard:     false,
			want:        "Unknown",
			wantError:   false,
			wantDiscard: true,
		},
		{
			name:        "Unsupported Type without discard",
			input:       12345, // int is unsupported
			discard:     false,
			want:        "Unknown",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Unsupported Type with discard",
			input:       12345, // int is unsupported
			discard:     true,
			want:        "Unknown",
			wantError:   false, // According to Wrap, it should return ErrorType with discard
			wantDiscard: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the country wrapper
			countryWrapper := New[*WrapperCountry]()
			countryWrapper.Initialize()

			// Wrap the input
			err := countryWrapper.Wrap(tt.input, tt.discard)
			if (err != nil) != tt.wantError {
				t.Fatalf("Wrap() error = %v, wantError %v", err, tt.wantError)
			}

			// Check if wrapper is discarded as expected
			if countryWrapper.IsDiscarded() != tt.wantDiscard {
				t.Errorf("IsDiscarded() = %v, want %v", countryWrapper.IsDiscarded(), tt.wantDiscard)
			}

			// If not discarded, check the unwrapped value
			if !countryWrapper.IsDiscarded() && countryWrapper.Unwrap() != tt.want {
				t.Errorf("Get() = %v, want %v", countryWrapper.Unwrap(), tt.want)
			}
		})
	}
}

// TestWrapperCountry_JSONMarshal tests JSON marshalling of WrapperCountry.
func TestWrapperCountry_JSONMarshal(t *testing.T) {
	tests := []struct {
		name         string
		country      string
		wantError    bool
		expectedJSON string
	}{
		{
			name:         "Valid Country Marshalling",
			country:      "DE",
			wantError:    false,
			expectedJSON: `{"country":"Germany"}`,
		},
		{
			name:         "Invalid Country Marshalling",
			country:      "invalid_country",
			wantError:    true,
			expectedJSON: `{"country":null}`,
		},
	}

	// Create a struct for testing
	type Data struct {
		Country *WrapperCountry `json:"country"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the country wrapper
			countryWrapper := New[*WrapperCountry]()
			countryWrapper.Initialize()

			// Wrap the country
			err := countryWrapper.Wrap(tt.country, false)
			if (err != nil) != tt.wantError {
				t.Fatalf("Wrap() error = %v, wantError %v", err, tt.wantError)
			}

			// Create the data struct
			data := Data{
				Country: countryWrapper,
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

// TestWrapperCountry_Unmarshal tests JSON unmarshalling of WrapperCountry.
func TestWrapperCountry_Unmarshal(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		want        string
		wantError   bool
		wantDiscard bool
	}{
		{
			name:        "Unmarshal valid CountryCode",
			jsonInput:   `{"country":"DE"}`,
			want:        "Germany",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Unmarshal invalid CountryCode",
			jsonInput:   `{"country":"invalid_country"}`,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Unmarshal null CountryCode",
			jsonInput:   `{"country":null}`,
			want:        "",
			wantError:   false,
			wantDiscard: true,
		},
		{
			name:        "Unmarshal unsupported type",
			jsonInput:   `{"country":123}`,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
	}

	// Create a struct for testing
	type Data struct {
		Country *WrapperCountry `json:"country"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data Data

			// Initialize the country wrapper
			data.Country = New[*WrapperCountry]()
			data.Country.Initialize()

			// Unmarshal the JSON input
			err := json.Unmarshal([]byte(tt.jsonInput), &data)
			if (err != nil) != tt.wantError {
				t.Fatalf("Unmarshal() error = %v, wantError %v", err, tt.wantError)
			}

			// Make sure that we handle cases where the country is nil
			if data.Country != nil {
				// Check if wrapper is discarded as expected
				if data.Country.IsDiscarded() != tt.wantDiscard {
					t.Errorf("IsDiscarded() = %v, want %v", data.Country.IsDiscarded(), tt.wantDiscard)
				}

				// If not discarded, check the unwrapped value
				if !data.Country.IsDiscarded() && data.Country.Unwrap() != tt.want {
					t.Errorf("Get() = %v, want %v", data.Country.Get(), tt.want)
				}
			}
		})
	}
}
