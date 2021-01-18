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
		t.Errorf("Parse(%q) returns incorrect error %q, want %q.", s, err, wantedErr)
		t.Errorf("Parse(%q) returns an incorrect error %q, want %q.", s, err, wantedErr)
	}
}
