package matcher

import (
	"fmt"
	"testing"
)

// Test that Parse parses a simple valid matcher correctly.
func TestParse(t *testing.T) {
	var wantedErr error = nil
	wantedMatcher := MatcherWithLogicalOperator{
		MatcherWithoutLogicalOperator{Data1, LessThanOperator, 2},
		LogicalAndOperator,
		MatcherWithoutLogicalOperator{Data2, UnequalToOperator, 3},
	}

	s := "data1 < 2 && data2 != 3"
	matcher, err := Parse(s)

	if err != wantedErr {
		t.Errorf("Parse(%q) returns an incorrect error %q, want %v.", s, err, wantedErr)
	}

	if !matcher.Equal(wantedMatcher) {
		t.Errorf("Parse(%q) returns an incorrect matcher %v, want %v.", s, matcher, wantedMatcher)
	}
}

// Test that Parse parses a matcher, with an invalid left matcher, correctly.
func TestParseInvalidLeftMatcher(t *testing.T) {
	leftMatcher := "data1 557"
	wantedErr := fmt.Errorf("matcher %q: no valid comparison operator", leftMatcher)

	s := leftMatcher + " && data2 != 365"
	_, err := Parse(s)

	if err == nil {
		t.Errorf("Parse(%q) returns an incorrect error %v, want %q.", s, err, wantedErr)
	} else if err.Error() != wantedErr.Error() {
		t.Errorf("Parse(%q) returns an incorrect error %q, want %q.", s, err, wantedErr)
	}
}

// Test that Parse parses a matcher, without a logical operator, correctly.
func TestParseWithoutLogicalOperator(t *testing.T) {
	var wantedErr error = nil
	wantedMatcher := MatcherWithoutLogicalOperator{Data1, EqualToOperator, 557}

	s := "data1 == 557"
	matcher, err := Parse(s)

	if err != wantedErr {
		t.Errorf("Parse(%q) returns an incorrect error %q, want %v.", s, err, wantedErr)
	}

	if !matcher.Equal(wantedMatcher) {
		t.Errorf("Parse(%q) returns an incorrect matcher %v, want %v.", s, matcher, wantedMatcher)
	}
}

// Test that Parse parses a complex, deeply nested, matcher correctly.
// Important, this test verifies that the logical operators have correct precedence, that is && must have higher precedence than ||.
func TestParseComplex(t *testing.T) {
	var wantedErr error = nil
	wantedMatcher := MatcherWithLogicalOperator{
		MatcherWithLogicalOperator{
			MatcherWithoutLogicalOperator{Data1, EqualToOperator, 557},
			LogicalOrOperator,
			MatcherWithoutLogicalOperator{Data1, UnequalToOperator, 73},
		},
		LogicalAndOperator,
		MatcherWithLogicalOperator{
			MatcherWithLogicalOperator{
				MatcherWithoutLogicalOperator{Data2, GreaterThanOperator, 20},
				LogicalOrOperator,
				MatcherWithoutLogicalOperator{Data1, LessThanOperator, 30},
			},
			LogicalAndOperator,
			MatcherWithoutLogicalOperator{Data2, UnequalToOperator, 15},
		},
	}

	s := "data1 == 557 || data1 != 73 && data2 > 20 || data1 < 30 && data2 != 15"
	matcher, err := Parse(s)

	if err != wantedErr {
		t.Errorf("Parse(%q) returns an incorrect error %q, want %v.", s, err, wantedErr)
	}

	if !matcher.Equal(wantedMatcher) {
		t.Errorf("Parse(%q) returns an incorrect matcher %v, want %v.", s, matcher, wantedMatcher)
	}
}

// Test that Parse parses a matcher, without a logical operator, lacking a left operand, correctly.
func TestParseNoLeftOperand(t *testing.T) {
	s := " > 123"

	var wantedErr error = fmt.Errorf("matcher %q: no valid left operand", s)

	_, err := Parse(s)

	if err == nil {
		t.Errorf("Parse(%q) returns an incorrect error %v, want %q.", s, err, wantedErr)
	} else if err.Error() != wantedErr.Error() {
		t.Errorf("Parse(%q) returns an incorrect error %q, want %q.", s, err, wantedErr)
	}
}

// Test that Parse parses a matcher, without a logical operator, lacking a right operand, correctly.
func TestParseNo(t *testing.T) {
	s := "data1 >"

	var wantedErr error = fmt.Errorf("matcher %q: no valid right operand", s)

	_, err := Parse(s)

	if err == nil {
		t.Errorf("Parse(%q) returns incorrect error %v, want %q.", s, err, wantedErr)
	} else if err.Error() != wantedErr.Error() {
		t.Errorf("Parse(%q) returns incorrect error %q, want %q.", s, err, wantedErr)
	}
}
