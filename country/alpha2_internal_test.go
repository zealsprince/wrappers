package country

import "testing"

// isAlpha2 guards against pseudo-entries in the countries package. The character
// check is not reachable through Parse, since ByName only ever returns a real
// code or a sentinel, so it gets exercised here directly rather than deleted.
func TestIsAlpha2(t *testing.T) {
	tests := []struct {
		code string
		want bool
	}{
		{"DE", true},
		{"US", true},
		{"de", false},
		{"D1", false},
		{"1E", false},
		{"D", false},
		{"", false},
		{"DEU", false},
		{"None", false},
		{"International", false},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			if got := isAlpha2(tt.code); got != tt.want {
				t.Errorf("isAlpha2(%q) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}
