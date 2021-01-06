package lang

import (
	"testing"
	"errors"
	"fmt"
)

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
		t.Errorf("parseComparison returns non-nil err")
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
