package matcher

import (
	"fmt"
	"testing"
)

func areMatchersEqual(left, right Matcher) bool {
	switch l := left.(type) {
	case MatcherWithoutLogicalOperator:
		r, ok := right.(MatcherWithoutLogicalOperator)
		if !ok {
			return false
		}
		return l.LeftOperand == r.LeftOperand && l.Operator == r.Operator && l.RightOperand == r.RightOperand
	case MatcherWithLogicalOperator:
		r, ok := right.(MatcherWithLogicalOperator)
		if !ok {
			return false
		}
		return areMatchersEqual(l.LeftMatcher, r.LeftMatcher) && l.Operator == r.Operator && areMatchersEqual(l.RightMatcher, r.RightMatcher)
	default:
		return false
	}

}

// Test that Parse parses a simple valid matcher correctly.
func TestParse(t *testing.T) {
	var wantedErr error = nil
	wantedMatcher := MatcherWithLogicalOperator{
		LeftMatcher: MatcherWithoutLogicalOperator{
			LeftOperand:  Data1,
			Operator:     LessThanOperator,
			RightOperand: 2,
		},
		Operator: LogicalAndOperator,
		RightMatcher: MatcherWithoutLogicalOperator{
			LeftOperand:  Data2,
			Operator:     UnequalToOperator,
			RightOperand: 3,
		},
	}

	s := "data1 < 2 && data2 != 3"
	matcher, err := Parse(s)

	if err != wantedErr {
		t.Errorf("Parse(%q) returns an incorrect error %q, want %v.", s, err, wantedErr)
	}

	if !areMatchersEqual(matcher, wantedMatcher) {
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
	wantedMatcher := MatcherWithoutLogicalOperator{
		LeftOperand:  Data1,
		Operator:     EqualToOperator,
		RightOperand: 557,
	}

	s := "data1 == 557"
	matcher, err := Parse(s)

	if err != wantedErr {
		t.Errorf("Parse(%q) returns an incorrect error %q, want %v.", s, err, wantedErr)
	}

	if !areMatchersEqual(matcher, wantedMatcher) {
		t.Errorf("Parse(%q) returns an incorrect matcher %v, want %v.", s, matcher, wantedMatcher)
	}
}

// Test that Parse parses a complex, deeply nested, matcher correctly.
// Important, this test verifies that the logical operators have correct precedence, that is && must have higher precedence than ||.
func TestParseComplex(t *testing.T) {
	var wantedErr error = nil
	wantedMatcher := MatcherWithLogicalOperator{ // data1 == 557 || data1 != 73 && data2 > 20 || data1 < 30 && data2 != 15
		LeftMatcher: MatcherWithLogicalOperator{ // data1 == 557 || data1 != 73
			LeftMatcher: MatcherWithoutLogicalOperator{ // data1 == 557
				LeftOperand:  Data1,
				Operator:     EqualToOperator,
				RightOperand: 557,
			},
			Operator: LogicalOrOperator,
			RightMatcher: MatcherWithoutLogicalOperator{ // data1 != 73
				LeftOperand:  Data1,
				Operator:     UnequalToOperator,
				RightOperand: 73,
			},
		},
		Operator: LogicalAndOperator,
		RightMatcher: MatcherWithLogicalOperator{ // data2 > 20 || data1 < 30 && data2 != 15
			LeftMatcher: MatcherWithLogicalOperator{ // data2 > 20 || data1 < 30
				LeftMatcher: MatcherWithoutLogicalOperator{ // data2 > 20
					LeftOperand:  Data2,
					Operator:     GreaterThanOperator,
					RightOperand: 20,
				},
				Operator: LogicalOrOperator,
				RightMatcher: MatcherWithoutLogicalOperator{ // data1 < 30
					LeftOperand:  Data1,
					Operator:     LessThanOperator,
					RightOperand: 30,
				},
			},
			Operator: LogicalAndOperator,
			RightMatcher: MatcherWithoutLogicalOperator{ // data2 != 15
				LeftOperand:  Data2,
				Operator:     UnequalToOperator,
				RightOperand: 15,
			},
		},
	}

	s := "data1 == 557 || data1 != 73 && data2 > 20 || data1 < 30 && data2 != 15"
	matcher, err := Parse(s)

	if err != wantedErr {
		t.Errorf("Parse(%q) returns an incorrect error %q, want %v.", s, err, wantedErr)
	}

	if !areMatchersEqual(matcher, wantedMatcher) {
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
