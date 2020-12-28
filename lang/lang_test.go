package lang

import "testing"

// Test that parseComparison parses a valid comparison, where the left operand is part1, correctly.
func TestParsePart1Comparison(t *testing.T) {
	wantedIdentifier := Part1Identifier
	wantedOperator := GreaterThanOperator
	wantedInteger := 123
	wantedOk := true

	s := "part1>123"
	identifier, operator, integer, ok := parseComparison(s)

	if identifier != wantedIdentifier {
		t.Errorf("parseComparison returns incorrect identifier %v, want %v", identifier, wantedIdentifier)
	}
	if operator != wantedOperator {
		t.Errorf("parseComparison returns incorrect operator %v, want %v", operator, wantedOperator)
	}
	if integer != wantedInteger {
		t.Errorf("parseComparison returns incorrect integer %d, want %d", integer, wantedInteger)
	}
	if ok != wantedOk {
		t.Errorf("parseComparison returns incorrect ok %t, want %t", ok, wantedOk)
	}
}
