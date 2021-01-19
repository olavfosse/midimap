package keycode

import (
	"fmt"
	"testing"
)

// Test that Parse parses a simple valid keycode correctly.
func TestParse(t *testing.T) {
	var wantedErr error = nil
	wantedKeycode := 1

	s := "1" // ESC
	keycode, err := Parse(s)

	if err != wantedErr {
		t.Errorf("Parse(%q) returns an incorrect error %q, want %v.", s, err, wantedErr)
	}

	if keycode != wantedKeycode {
		t.Errorf("Parse(%q) returns an incorrect keycode %d, want %d.", s, keycode, wantedKeycode)
	}
}

// Test that Parse parses a keycode, with spaces interspersed between the digits, correctly.
func TestParseInterspersedSpaces(t *testing.T) {
	s := "1 2 3 4 5 6 7 8 9"
	wantedErr := fmt.Errorf("keycode %q: invalid", s)

	_, err := Parse(s)

	if err == nil {
		t.Errorf("Parse(%q) returns an incorrect error %v, want %q.", s, err, wantedErr)
	} else if err.Error() != wantedErr.Error() {
		t.Errorf("Parse(%q) returns an incorrect error %q, want %q.", s, err, wantedErr)
	}
}
