package lang

import (
	"testing"
	"errors"
	"fmt"
)

/*
 * parseComparison tests
 */

// Test that parseComparison parses a valid comparison correctly.
func TestParseComparisonValid(t *testing.T) {
	wantedLeftOperand := Data1
	wantedOperator := GreaterThanOperator
	var wantedRightOperand int64 = 123
	
	s := "data1 > 123"
	comparison, err := parseComparison(s)

	if comparison.LeftOperand != wantedLeftOperand {
		t.Errorf("parseComparison returns a comparison with an incorrect leftOperand %v, want %v", comparison.LeftOperand, wantedLeftOperand)
	}
	if comparison.Operator != wantedOperator {
		t.Errorf("parseComparison returns a comparison with an incorrect operator %v, want %v", comparison.Operator, wantedOperator)
	}
	if comparison.RightOperand != wantedRightOperand {
		t.Errorf("parseComparison returns a comparison with an incorrect rightOperand %d, want %d", comparison.RightOperand, wantedRightOperand)
	}
	if err != nil {
		t.Errorf("parseComparison returns non-nil err %q", err)
	}
}

// Test that parseComparison parses an invalid comparison, lacking a left operand, correctly.
func TestParseComparisonInvalidNoLeftOperand(t *testing.T) {
	s := " > 123"

	var wantedErr error = errors.New(fmt.Sprintf("Comparison %q does not have a valid left operand", s))
	
	_, err := parseComparison(s)
	
	if err == nil {
		t.Errorf("parseComparison returns incorrect err nil, want %q", wantedErr)
	} else if err.Error() != wantedErr.Error() {
		t.Errorf("parseComparison returns incorrect err %q, want %q", err, wantedErr)
	}
}

// Test that parseComparison parses an invalid comparison, lacking a operator, correctly.
func TestParseComparisonInvalidNoOperator(t *testing.T) {
	s := "data1 123"

	var wantedErr error = errors.New(fmt.Sprintf("Comparison %q does not have a valid operator", s))
	
	_, err := parseComparison(s)
	
	if err == nil {
		t.Errorf("parseComparison returns incorrect err nil, want %q", wantedErr)
	} else if err.Error() != wantedErr.Error() {
		t.Errorf("parseComparison returns incorrect err %q, want %q", err, wantedErr)
	}
}

// Test that parseComparison parses an invalid comparison, lacking a left operand, correctly.
func TestParseComparisonInvalidNoRightOperand(t *testing.T) {
	s := "data1 >"

	var wantedErr error = errors.New(fmt.Sprintf("Comparison %q does not have a valid right operand", s))
	
	_, err := parseComparison(s)
	
	if err == nil {
		t.Errorf("parseComparison returns incorrect err nil, want %q", wantedErr)
	} else if err.Error() != wantedErr.Error() {
		t.Errorf("parseComparison returns incorrect err %q, want %q", err, wantedErr)
	}
}

/*
 * parseMatcher tests
 */

func areComparisonsEqual(left, right Comparison) bool {
	return left.LeftOperand == right.LeftOperand && left.Operator == right.Operator && left.RightOperand == right.RightOperand
}

// Test that parseMatcher parses a valid matcher correctly.
func TestParseMatcherValid(t *testing.T) {
	wantedLeftComparison := Comparison{
		LeftOperand:  Data1,
		Operator:     EqualToOperator,
		RightOperand: 557,
	}
	wantedOperator := LogicalAndOperator
	wantedRightComparison := Comparison{
		LeftOperand:  Data2,
		Operator:     UnequalToOperator,
		RightOperand: 365,
	}
	
	s := "data1 == 557 && data2 != 365"
	matcher, err := parseMatcher(s)

	if !areComparisonsEqual(matcher.LeftComparison, wantedLeftComparison) {
		t.Errorf("parseMatcher returns a matcher with an incorrect LeftComparison %v, want %v", matcher.LeftComparison, wantedLeftComparison)
	}
	if matcher.Operator != wantedOperator {
		t.Errorf("parseMatcher returns a matcher with an incorrect Operator %d, want %d", matcher.Operator, wantedOperator)
	}
	if !areComparisonsEqual(matcher.RightComparison, wantedRightComparison) {
		t.Errorf("parseMatcher returns a matcher with an incorrect RightComparison %v, want %v", matcher.RightComparison, wantedRightComparison)
	}
	if err != nil {
		t.Errorf("parseMatcher returns non-nil err %q", err)
	}
}

// Test that parseMatcher parses an invalid matcher, with an invalid left comparison, correctly.
func TestParseMatcherInvalidInvalidLeftComparison(t *testing.T) {
	leftComparison := "data1 557"
	wantedErr := errors.New(fmt.Sprintf("Comparison %q does not have a valid operator", leftComparison))
	
	s := leftComparison + " && data2 != 365"
	_, err := parseMatcher(s)

	if err == nil {
		t.Errorf("parseMatcher returns incorrect err nil, want %q", wantedErr)
	} else if err.Error() != wantedErr.Error() {
		t.Errorf("parseMatcher returns incorrect err %q, want %q", err, wantedErr)
	}
}

// Test that parseMatcher parses an invalid matcher, lacking a logical operator, correctly.
func TestParseMatcherInvalidNoLogicalOperator(t *testing.T) {
	s := "data1 557 & & data2 != 365"

	wantedErr := errors.New(fmt.Sprintf("Matcher %q does not have a valid logical operator", s))

	_, err := parseMatcher(s)

	if err == nil {
		t.Errorf("parseMatcher returns incorrect err nil, want %q", wantedErr)
	} else if err.Error() != wantedErr.Error() {
		t.Errorf("parseMatcher returns incorrect err %q, want %q", err, wantedErr)
	}
}

// Test that parseMatcher parses an invalid matcher, with an invalid right comparison, correctly.
func TestParseMatcherInvalidInvalidRightComparison(t *testing.T) {
	rightComparison := "data1 557"
	wantedErr := errors.New(fmt.Sprintf("Comparison %q does not have a valid operator", rightComparison))
	
	s :=  "data2 != 365 && " + rightComparison
	_, err := parseMatcher(s)

	if err == nil {
		t.Errorf("parseMatcher returns incorrect err nil, want %q", wantedErr)
	} else if err.Error() != wantedErr.Error() {
		t.Errorf("parseMatcher returns incorrect err %q, want %q", err, wantedErr)
	}
}
