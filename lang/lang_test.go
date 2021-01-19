package lang

import (
	"fmt"
	"testing"

	"./matcher"
)

// Test that parseKeycode parses a simple valid keycode correctly.
func TestParseKeycode(t *testing.T) {
	var wantedErr error = nil
	wantedKeycode := 1

	s := "1" // ESC
	keycode, err := parseKeycode(s)

	if err != wantedErr {
		t.Errorf("parseKeycode(%q) returns an incorrect error %q, want %v.", s, err, wantedErr)
	}

	if keycode != wantedKeycode {
		t.Errorf("parseKeycode(%q) returns an incorrect keycode %d, want %d.", s, keycode, wantedKeycode)
	}
}

// Test that parseKeycode parses a keycode, with spaces interspersed between the digits, correctly.
func TestParseKeycodeInterspersedSpaces(t *testing.T) {
	s := "1 2 3 4 5 6 7 8 9"
	wantedErr := fmt.Errorf("keycode %q: invalid", s)

	_, err := parseKeycode(s)

	if err == nil {
		t.Errorf("parseKeycode(%q) returns an incorrect error %v, want %q.", s, err, wantedErr)
	} else if err.Error() != wantedErr.Error() {
		t.Errorf("parseKeycode(%q) returns an incorrect error %q, want %q.", s, err, wantedErr)
	}
}

func (m Mapping) equal(n Mapping) bool {
	return m.Matcher.Equal(n.Matcher) && m.Keycode == n.Keycode
}

// Test that parseMapping parses a simple valid mapping correctly.
func TestParseMapping(t *testing.T) {
	var wantedErr error = nil
	wantedMapping := Mapping{
		Matcher: matcher.MatcherWithLogicalOperator{
			LeftMatcher: matcher.MatcherWithoutLogicalOperator{
				LeftOperand:  matcher.Data1,
				Operator:     matcher.EqualToOperator,
				RightOperand: 44,
			},
			Operator: matcher.LogicalAndOperator,
			RightMatcher: matcher.MatcherWithoutLogicalOperator{
				LeftOperand:  matcher.Data2,
				Operator:     matcher.EqualToOperator,
				RightOperand: 64,
			},
		},
		Keycode: 1,
	}

	s := "data1 == 44 && data2 == 64 -> 1"
	mapping, err := parseMapping(s)

	if err != wantedErr {
		t.Errorf("parseMapping(%q) returns an incorrect error %q, want %v.", s, err, wantedErr)
	}

	if !mapping.equal(wantedMapping) {
		t.Errorf("parseMapping(%q) returns an incorrect mapping %v, want %v.", s, mapping, wantedMapping)
	}
}

// Test that parseMapping parses a mapping, lacking a valid separator, correctly.
func TestParseMappingInvalidSeparator(t *testing.T) {
	s := "data1 < 44 && data2 == 64 - > 12" // interspersing spaces in the separator is not allowed.
	wantedErr := fmt.Errorf("mapping %q: no valid separator", s)

	_, err := parseMapping(s)

	if err == nil {
		t.Errorf("parseMapping(%q) returns an incorrect error %v, want %q.", s, err, wantedErr)
	} else if err.Error() != wantedErr.Error() {
		t.Errorf("parseMapping(%q) returns an incorrect error %q, want %q.", s, err, wantedErr)
	}
}
