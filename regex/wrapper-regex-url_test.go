package regex

import (
	"encoding/json"
	"testing"

	"github.com/zealsprince/wrappers"
)

// TestWrapperRegexUrl_Wrap tests the Wrap method of WrapperRegexUrl.
func TestWrapperRegexUrl_Wrap(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		discard     bool
		want        string
		wantError   bool
		wantDiscard bool
	}{
		{
			name:        "Valid URL - HTTP",
			input:       "http://example.com",
			discard:     false,
			want:        "http://example.com",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Valid URL - HTTPS",
			input:       "https://example.com",
			discard:     false,
			want:        "https://example.com",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Valid URL - FTP",
			input:       "ftp://example.com/resource",
			discard:     false,
			want:        "ftp://example.com/resource",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Invalid URL - Missing Scheme",
			input:       "://example.com",
			discard:     false,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Invalid URL - Unsupported Scheme",
			input:       "smtp://example.com",
			discard:     false,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Invalid URL - Spaces",
			input:       "http://example .com",
			discard:     false,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Invalid URL discard",
			input:       "htp://invalid.com",
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
			input: func() *WrapperRegexUrl {
				w := wrappers.New[*WrapperRegexUrl]()
				w.Wrap("https://valid.com", false)
				return w
			}(),
			discard:     false,
			want:        "https://valid.com",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name: "Input as WrapperProvider with invalid value (discarded)",
			input: func() *WrapperRegexUrl {
				w := wrappers.New[*WrapperRegexUrl]()
				w.Wrap("htp://invalid.com", false)
				return w
			}(),
			discard:     false,
			want:        "",
			wantError:   false,
			wantDiscard: true,
		},
		{
			name:        "Unsupported Type without discard",
			input:       7890, // int is unsupported
			discard:     false,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Unsupported Type with discard",
			input:       7890, // int is unsupported
			discard:     true,
			want:        "",
			wantError:   false,
			wantDiscard: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the URL wrapper
			urlWrapper := wrappers.New[*WrapperRegexUrl]()
			urlWrapper.Initialize()

			// Wrap the input
			err := urlWrapper.Wrap(tt.input, tt.discard)
			if (err != nil) != tt.wantError {
				t.Fatalf("Wrap() error = %v, wantError %v", err, tt.wantError)
			}

			// Check if wrapper is discarded as expected
			if urlWrapper.IsDiscarded() != tt.wantDiscard {
				t.Errorf("IsDiscarded() = %v, want %v", urlWrapper.IsDiscarded(), tt.wantDiscard)
			}

			// If not discarded, check the unwrapped value
			if !urlWrapper.IsDiscarded() && urlWrapper.Get() != tt.want {
				t.Errorf("Get() = %v, want %v", urlWrapper.Get(), tt.want)
			}
		})
	}
}

// TestWrapperRegexUrl_JSONMarshal tests JSON marshalling of WrapperRegexUrl.
func TestWrapperRegexUrl_JSONMarshal(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		wantError    bool
		expectedJSON string
	}{
		{
			name:         "Valid URL Marshalling - HTTP",
			url:          "http://example.com",
			wantError:    false,
			expectedJSON: `{"url":"http://example.com"}`,
		},
		{
			name:         "Valid URL Marshalling - HTTPS",
			url:          "https://example.com",
			wantError:    false,
			expectedJSON: `{"url":"https://example.com"}`,
		},
		{
			name:         "Invalid URL Marshalling",
			url:          "htp://invalid.com",
			wantError:    true,
			expectedJSON: `{"url":null}`,
		},
		{
			name:         "Discarded Wrapper",
			url:          "",
			wantError:    false,
			expectedJSON: `{"url":null}`,
		},
	}

	// Create a struct for testing
	type Data struct {
		Url *WrapperRegexUrl `json:"url"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the URL wrapper
			urlWrapper := wrappers.New[*WrapperRegexUrl]()
			urlWrapper.Initialize()

			// Wrap the URL if it's not to be discarded
			if tt.url != "" {
				err := urlWrapper.Wrap(tt.url, false)
				if (err != nil) != tt.wantError {
					t.Fatalf("Wrap() error = %v, wantError %v", err, tt.wantError)
				}
			} else {
				// If URL is empty, we can discard it
				urlWrapper.Discard()
			}

			// Create the data struct
			data := Data{
				Url: urlWrapper,
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

// TestWrapperRegexUrl_Unmarshal tests JSON unmarshalling of WrapperRegexUrl.
func TestWrapperRegexUrl_Unmarshal(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		want        string
		wantError   bool
		wantDiscard bool
	}{
		{
			name:        "Unmarshal valid HTTP URL",
			jsonInput:   `{"url":"http://example.com"}`,
			want:        "http://example.com",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Unmarshal valid HTTPS URL",
			jsonInput:   `{"url":"https://example.com"}`,
			want:        "https://example.com",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Unmarshal invalid URL",
			jsonInput:   `{"url":"htp://invalid.com"}`,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Unmarshal null URL",
			jsonInput:   `{"url":null}`,
			want:        "",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Unmarshal unsupported type",
			jsonInput:   `{"url":12345}`,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
	}

	// Create a struct for testing
	type Data struct {
		Url *WrapperRegexUrl `json:"url"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data Data

			// Initialize the URL wrapper
			data.Url = wrappers.New[*WrapperRegexUrl]()
			data.Url.Initialize()

			// Unmarshal the JSON input
			err := json.Unmarshal([]byte(tt.jsonInput), &data)
			if (err != nil) != tt.wantError {
				t.Fatalf("Unmarshal() error = %v, wantError %v", err, tt.wantError)
			}

			// Check that the URL wrapper is not nil in cases where the payload does not contain the URL
			if data.Url != nil {
				// Check if wrapper is discarded as expected
				if data.Url.IsDiscarded() != tt.wantDiscard {
					t.Errorf("IsDiscarded() = %v, want %v", data.Url.IsDiscarded(), tt.wantDiscard)
				}

				// If not discarded, check the unwrapped value
				if !data.Url.IsDiscarded() && data.Url.Get() != tt.want {
					t.Errorf("Get() = %v, want %v", data.Url.Get(), tt.want)
				}
			}
		})
	}
}
