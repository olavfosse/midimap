package lang

import "testing"

// Test that parseComparison parses a valid comparison, where the left operand is part1, correctly.
func TestParsePart1Comparison(t *testing.T) {
	wantedLeftOperand := Part1
	wantedOperator := GreaterThanOperator
	wantedRightOperand := 123
	wantedOk := true

	s := "part1>123"
	leftOperand, operator, rightOperand, ok := parseComparison(s)

	if leftOperand != wantedLeftOperand {
		t.Errorf("parseComparison returns incorrect leftOperand %v, want %v", leftOperand, wantedLeftOperand)
	}
	if operator != wantedOperator {
		t.Errorf("parseComparison returns incorrect operator %v, want %v", operator, wantedOperator)
	}
	if rightOperand != wantedRightOperand {
		t.Errorf("parseComparison returns incorrect rightOperand %d, want %d", rightOperand, wantedRightOperand)
	}
	if ok != wantedOk {
		t.Errorf("parseComparison returns incorrect ok %t, want %t", ok, wantedOk)
	}
}

// Test that parseComparison parses a valid comparison, where the left operand is part2, correctly.
func TestParsePart2Comparison(t *testing.T) {
	wantedLeftOperand := Part2
	wantedOperator := EqualToOperator
	wantedRightOperand := 321
	wantedOk := true

	s := "part2==321"
	leftOperand, operator, rightOperand, ok := parseComparison(s)

	if leftOperand != wantedLeftOperand {
		t.Errorf("parseComparison returns incorrect leftOperand %v, want %v", leftOperand, wantedLeftOperand)
	}
	if operator != wantedOperator {
		t.Errorf("parseComparison returns incorrect operator %v, want %v", operator, wantedOperator)
	}
	if rightOperand != wantedRightOperand {
		t.Errorf("parseComparison returns incorrect rightOperand %d, want %d", rightOperand, wantedRightOperand)
	}
	if ok != wantedOk {
		t.Errorf("parseComparison returns incorrect ok %t, want %t", ok, wantedOk)
	}
}

// Test that parseComparison parses an invalid comparison, lacking the comparison operator, correctly.
func TestParseComparisonSansOperator(t *testing.T) {
	wantedOk := false

	s := "part2456"
	_, _, _, ok := parseComparison(s)

	if ok != wantedOk {
		t.Errorf("parseComparison returns incorrect ok %t, want %t", ok, wantedOk)
	}
}

// Test that parseComparison parses an invalid comparison, lacking the right operand, correctly.
func TestParseComparisonSansRightOperand(t *testing.T) {
	wantedOk := false

	s := "part1!="
	_, _, _, ok := parseComparison(s)

	if ok != wantedOk {
		t.Errorf("parseComparison returns incorrect ok %t, want %t", ok, wantedOk)
	}
}

// Test that parseComparison parses an invalid comparison, lacking the left operand, correctly.
func TestParseComparisonSansLeftOperand(t *testing.T) {
	wantedOk := false

	s := "!=54"
	_, _, _, ok := parseComparison(s)

	if ok != wantedOk {
		t.Errorf("parseComparison returns incorrect ok %t, want %t", ok, wantedOk)
	}
}
