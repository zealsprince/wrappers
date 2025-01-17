// wrapper-enum_test.go
package enum

import (
	"encoding/json"
	"testing"

	"github.com/zealsprince/wrappers"
)

func TestWrapperEnumCardinalDirections_Wrap(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		discard     bool
		want        string
		wantDiscard bool
		wantError   bool
	}{
		{
			name:        "Wrap valid enum value (DirectionNorth)",
			input:       DirectionNorth,
			discard:     false,
			want:        "north",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap valid enum value (DirectionEast)",
			input:       DirectionEast,
			discard:     false,
			want:        "east",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap valid string 'south'",
			input:       "south",
			discard:     false,
			want:        "south",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap valid string 'west'",
			input:       "west",
			discard:     false,
			want:        "west",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap invalid string 'northeast'",
			input:       "northeast",
			discard:     false,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Wrap invalid string 'northeast' discard",
			input:       "northeast",
			discard:     true,
			want:        "",
			wantError:   false,
			wantDiscard: true,
		},
		{
			name:        "Wrap invalid type (int)",
			input:       123,
			discard:     false,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Wrap invalid type (int) discard",
			input:       123,
			discard:     true,
			want:        "",
			wantError:   false,
			wantDiscard: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := &WrapperEnumCardinalDirections{}
			wrapper.Initialize() // Ensure validValues are set.

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

			// Additional check: If discard is true and error occurred, wrapper should be discarded
			if tt.wantError && tt.discard {
				if !wrapper.IsDiscarded() {
					t.Errorf("Expected wrapper to be discarded, but it was not")
				}
			}

			// If no error and discard is true, ensure that the discarded flag was raised
			if !tt.wantError && tt.discard && tt.wantDiscard {
				if !wrapper.IsDiscarded() {
					t.Errorf("Expected wrapper to be discarded")
				}
			}
		})
	}
}

// TestWrapperEnumCardinalDirections_JSONMarshal tests JSON marshalling of WrapperEnumCardinalDirections.
func TestWrapperEnumCardinalDirections_JSONMarshal(t *testing.T) {
	tests := []struct {
		name         string
		wrapper      *WrapperEnumCardinalDirections
		expectedJSON string
	}{
		{
			name: "Marshal valid enum value (north)",
			wrapper: func() *WrapperEnumCardinalDirections {
				w := &WrapperEnumCardinalDirections{}
				w.Initialize()
				w.Wrap(DirectionNorth, false)
				return w
			}(),
			expectedJSON: `"north"`,
		},
		{
			name: "Marshal valid enum value (east)",
			wrapper: func() *WrapperEnumCardinalDirections {
				w := &WrapperEnumCardinalDirections{}
				w.Initialize()
				w.Wrap(DirectionEast, false)
				return w
			}(),
			expectedJSON: `"east"`,
		},
		{
			name: "Marshal valid enum value (south)",
			wrapper: func() *WrapperEnumCardinalDirections {
				w := &WrapperEnumCardinalDirections{}
				w.Initialize()
				w.Wrap(DirectionSouth, false)
				return w
			}(),
			expectedJSON: `"south"`,
		},
		{
			name: "Marshal valid enum value (west)",
			wrapper: func() *WrapperEnumCardinalDirections {
				w := &WrapperEnumCardinalDirections{}
				w.Initialize()
				w.Wrap(DirectionWest, false)
				return w
			}(),
			expectedJSON: `"west"`,
		},
		{
			name: "Marshal discarded value",
			wrapper: func() *WrapperEnumCardinalDirections {
				w := &WrapperEnumCardinalDirections{}
				w.Initialize()
				w.Discard()
				return w
			}(),
			expectedJSON: `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.wrapper)
			if err != nil {
				t.Fatalf("Failed to marshal WrapperEnumCardinalDirections: %v", err)
			}

			if string(jsonData) != tt.expectedJSON {
				t.Errorf("Marshalled JSON = %v, want %v", string(jsonData), tt.expectedJSON)
			}
		})
	}
}

// TestWrapperEnumCardinalDirections_JSONUnmarshal tests JSON unmarshalling of WrapperEnumCardinalDirections.
func TestWrapperEnumCardinalDirections_JSONUnmarshal(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		want        string
		wantDiscard bool
		wantError   bool
	}{
		{
			name:        "Unmarshal valid enum 'north'",
			jsonInput:   `"north"`,
			want:        "north",
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Unmarshal valid enum 'east'",
			jsonInput:   `"east"`,
			want:        "east",
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Unmarshal valid enum 'south'",
			jsonInput:   `"south"`,
			want:        "south",
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Unmarshal valid enum 'west'",
			jsonInput:   `"west"`,
			want:        "west",
			wantDiscard: false,
			wantError:   false,
		},
		{
			name:        "Unmarshal invalid enum 'northeast'",
			jsonInput:   `"northeast"`,
			want:        "",
			wantDiscard: true,
			wantError:   true,
		},
		{
			name:        "Unmarshal null",
			jsonInput:   `null`,
			want:        "",
			wantDiscard: true,
			wantError:   true,
		},
		{
			name:        "Unmarshal invalid type (number)",
			jsonInput:   `123`,
			want:        "",
			wantDiscard: true,
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := &WrapperEnumCardinalDirections{}
			// Initialization occurs in UnmarshalJSON if not initialized

			err := json.Unmarshal([]byte(tt.jsonInput), wrapper)

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

			if tt.wantDiscard != wrapper.IsDiscarded() {
				t.Errorf("Discarded state = %v, want %v", wrapper.IsDiscarded(), tt.wantDiscard)
			}
		})
	}
}

// Additional Tests for WrapperEnum – Ensure invalid T types are handled correctly.
// This test is optional and showcases using WrapperEnum with another enum type.
// Define a custom enum type
type Season string

const (
	SeasonSpring Season = "spring"
	SeasonSummer Season = "summer"
	SeasonAutumn Season = "autumn"
	SeasonWinter Season = "winter"
)

type WrapperEnumSeasons struct {
	WrapperEnum[Season]
}

// Initialize with valid seasons
func (wrapper *WrapperEnumSeasons) Initialize() {
	wrapper.name = wrappers.Name("WrapperEnumSeasons")
	wrapper.SetValidValues(
		[]Season{
			SeasonSpring,
			SeasonSummer,
			SeasonAutumn,
			SeasonWinter,
		},
	)

	wrapper.WrapperBase.Initialize()
}

func TestWrapperEnum_CustomEnum(t *testing.T) {

	tests := []struct {
		name        string
		input       any
		discard     bool
		want        string
		wantError   bool
		wantDiscard bool
	}{
		{
			name:        "Wrap valid enum value (SeasonSpring)",
			input:       SeasonSpring,
			discard:     false,
			want:        "spring",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap valid string 'summer'",
			input:       "summer",
			discard:     false,
			want:        "summer",
			wantError:   false,
			wantDiscard: false,
		},
		{
			name:        "Wrap invalid string 'monsoon'",
			input:       "monsoon",
			discard:     false,
			want:        "",
			wantError:   true,
			wantDiscard: true,
		},
		{
			name:        "Wrap invalid string 'monsoon' discard",
			input:       "monsoon",
			discard:     true,
			want:        "",
			wantError:   false,
			wantDiscard: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := &WrapperEnumSeasons{}
			wrapper.Initialize()

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

			if tt.wantError && tt.discard {
				if !wrapper.IsDiscarded() {
					t.Errorf("Expected wrapper to be discarded, but it was not")
				}
			}

			if !tt.wantError && tt.discard && tt.wantDiscard {
				if !wrapper.IsDiscarded() {
					t.Errorf("Expected wrapper to be discarded")
				}
			}
		})
	}
}
