package lang

import (
	"fmt"
	"testing"
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

	var wantedErr error = fmt.Errorf("Comparison %q does not have a valid left operand", s)

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

	var wantedErr error = fmt.Errorf("Comparison %q does not have a valid operator", s)

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

	var wantedErr error = fmt.Errorf("Comparison %q does not have a valid right operand", s)

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
	wantedErr := fmt.Errorf("Comparison %q does not have a valid operator", leftComparison)

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

	wantedErr := fmt.Errorf("Matcher %q does not have a valid logical operator", s)

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
	wantedErr := fmt.Errorf("Comparison %q does not have a valid operator", rightComparison)

	s := "data2 != 365 && " + rightComparison
	_, err := parseMatcher(s)

	if err == nil {
		t.Errorf("parseMatcher returns incorrect err nil, want %q", wantedErr)
	} else if err.Error() != wantedErr.Error() {
		t.Errorf("parseMatcher returns incorrect err %q, want %q", err, wantedErr)
	}
}

/*
 * parseKeyCode tests
 */

// Test that parseKeyCode parses a valid key code correctly.
func TestParseKeyCodeValid(t *testing.T) {
	wantedKeyCode := 1

	s := "1" // ESC

	keyCode, err := parseKeyCode(s)

	if keyCode != wantedKeyCode {
		t.Errorf("parseKeyCode returns incorrect key code %d, want %d", keyCode, wantedKeyCode)
	}
	if err != nil {
		t.Errorf("parseKeyCode returns non-nil err %q", err)
	}
}

// Test that parseKeyCode parses an invalid keyCode, with spaces interspersed between the digits, correctly.
func TestParseKeyCodeInvalidInterspersedSpaces(t *testing.T) {
	s := "1 2 3 4 5 6 7 8 9"

	wantedErr := fmt.Errorf("Key code %q is invalid", s)

	_, err := parseKeyCode(s)

	if err == nil {
		t.Errorf("parseKeyCode returns incorrect err nil, want %q", wantedErr)
	} else if err.Error() != wantedErr.Error() {
		t.Errorf("parseKeyCode returns incorrect err %q, want %q", err, wantedErr)
	}
}

/*
 * parseMapping tests
 */

func areMatchersEqual(left, right Matcher) bool {
	areLeftComparisonsEqual := areComparisonsEqual(left.LeftComparison, right.LeftComparison)
	areOperatorsEqual := left.Operator == right.Operator
	areRightComparisonsEqual := areComparisonsEqual(left.RightComparison, right.RightComparison)
	return areLeftComparisonsEqual && areOperatorsEqual && areRightComparisonsEqual
}

// Test that parseMapping parses a valid mapping correctly.
func TestParseMappingValid(t *testing.T) {
	s := "data1 == 44 && data2 == 64 -> 1"

	wantedMatcher := Matcher{
		LeftComparison: Comparison{
			LeftOperand:  Data1,
			Operator:     EqualToOperator,
			RightOperand: 44,
		},
		Operator: LogicalAndOperator,
		RightComparison: Comparison{
			LeftOperand:  Data2,
			Operator:     EqualToOperator,
			RightOperand: 64,
		},
	}
	wantedKeyCode := 1

	mapping, err := parseMapping(s)

	if !areMatchersEqual(mapping.Matcher, wantedMatcher) {
		t.Errorf("parseMapping returns mapping with incorrect matcher %v, want %v", mapping.Matcher, wantedMatcher)
	}
	if mapping.KeyCode != wantedKeyCode {
		t.Errorf("parseMapping returns mapping with incorrect key code %d, want %d", mapping.KeyCode, wantedKeyCode)
	}
	if err != nil {
		t.Errorf("parseMapping returns non-nil error %q", err)
	}
}

// Test that parseMapping parses an invalid mapping, with an invalid matcher, correctly.
func TestParseMappingInvalidInvalidMatcher(t *testing.T) {
	matcher := "data1 < 44 data2 == 64"
	s := matcher + " -> 1"

	wantedErr := fmt.Errorf("Matcher %q does not have a valid logical operator", matcher)

	_, err := parseMapping(s)

	if err == nil {
		t.Errorf("parseMapping returns err nil, want %q", wantedErr)
	} else if err.Error() != wantedErr.Error() {
		t.Errorf("parseMapping returns incorrect err %q, want %q", err, wantedErr)
	}
}

// Test that parseMapping parses an invalid mapping, with an invalid separator, correctly.
func TestParseMappingInvalidInvalidSeparator(t *testing.T) {
	s := "data1 < 44 && data2 == 64 - > 12" // interspersing spaces in the separator is not allowed.

	wantedErr := fmt.Errorf("Mapping %q does not have a valid separator", s)

	_, err := parseMapping(s)

	if err == nil {
		t.Errorf("parseMapping returns err nil, want %q", wantedErr)
	} else if err.Error() != wantedErr.Error() {
		t.Errorf("parseMapping returns incorrect err %q, want %q", err, wantedErr)
	}
}

// Test that parseMapping parses an invalid mapping, with an invalid keycode, correctly.
func TestParseMappingInvalidInvalidKeycode(t *testing.T) {
	keyCode := "1 2"
	s := "data1 < 44 && data2 == 64 -> " + keyCode

	wantedErr := fmt.Errorf("Key code %q is invalid", keyCode)

	_, err := parseMapping(s)

	if err == nil {
		t.Errorf("parseMapping returns err nil, want %q", wantedErr)
	} else if err.Error() != wantedErr.Error() {
		t.Errorf("parseMapping returns incorrect err %q, want %q", err, wantedErr)
	}
}
