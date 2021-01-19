package mapping

import (
	"fmt"
	"testing"

	"../matcher"
)

// Test that Parse parses a simple valid mapping correctly.
func TestParse(t *testing.T) {
	var wantedErr error = nil
	wantedMapping := Mapping{
		Matcher: matcher.MatcherWithLogicalOperator{
			matcher.MatcherWithoutLogicalOperator{matcher.Data1, matcher.EqualToOperator, 44},
			matcher.LogicalAndOperator,
			matcher.MatcherWithoutLogicalOperator{matcher.Data2, matcher.EqualToOperator, 64},
		},
		Keycode: 1,
	}

	s := "data1 == 44 && data2 == 64 -> 1"
	mapping, err := Parse(s)

	if err != wantedErr {
		t.Errorf("Parse(%q) returns an incorrect error %q, want %v.", s, err, wantedErr)
	}

	if !mapping.Equal(wantedMapping) {
		t.Errorf("Parse(%q) returns an incorrect mapping %v, want %v.", s, mapping, wantedMapping)
	}
}

// Test that Parse parses a mapping, lacking a valid separator, correctly.
func TestParseInvalidSeparator(t *testing.T) {
	s := "data1 < 44 && data2 == 64 - > 12" // interspersing spaces in the separator is not allowed.
	wantedErr := fmt.Errorf("mapping %q: no valid separator", s)

	_, err := Parse(s)

	if err == nil {
		t.Errorf("Parse(%q) returns an incorrect error %v, want %q.", s, err, wantedErr)
	} else if err.Error() != wantedErr.Error() {
		t.Errorf("Parse(%q) returns an incorrect error %q, want %q.", s, err, wantedErr)
	}
}
