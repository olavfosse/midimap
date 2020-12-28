package lang

import (
	"testing"
)

// Test that parseComparison parses a valid comparison, where the left operand is part1, correctly.
func TestParsePart1Comparison(t *testing.T) {
	wantedLeftOperand := Part1
	wantedOperator := GreaterThanOperator
	wantedRightOperand := 123
	wantedOk := true

	s := "part1 >123"
	comparison, ok := parseComparison(s)

	if comparison.LeftOperand != wantedLeftOperand {
		t.Errorf("parseComparison returns a comparison with an incorrect leftOperand %v, want %v", comparison.LeftOperand, wantedLeftOperand)
	}
	if comparison.Operator != wantedOperator {
		t.Errorf("parseComparison returns a comparison with an incorrect operator %v, want %v", comparison.Operator, wantedOperator)
	}
	if comparison.RightOperand != wantedRightOperand {
		t.Errorf("parseComparison returns a comparison with an incorrect rightOperand %d, want %d", comparison.RightOperand, wantedRightOperand)
	}
	if ok != wantedOk {
		t.Errorf("parseComparison returns incorrect ok %t", ok)
	}
}

// Test that parseComparison parses a valid comparison, where the left operand is part2, correctly.
func TestParsePart2Comparison(t *testing.T) {
	wantedLeftOperand := Part2
	wantedOperator := EqualToOperator
	wantedRightOperand := 321
	wantedOk := true

	s := "part2 = =  321"
	comparison, ok := parseComparison(s)

	if comparison.LeftOperand != wantedLeftOperand {
		t.Errorf("parseComparison returns a comparison with an incorrect leftOperand %v, want %v", comparison.LeftOperand, wantedLeftOperand)
	}
	if comparison.Operator != wantedOperator {
		t.Errorf("parseComparison returns a comparison with an incorrect operator %v, want %v", comparison.Operator, wantedOperator)
	}
	if comparison.RightOperand != wantedRightOperand {
		t.Errorf("parseComparison returns a comparison with an incorrect rightOperand %d, want %d", comparison.RightOperand, wantedRightOperand)
	}
	if ok != wantedOk {
		t.Errorf("parseComparison returns incorrect ok %t", ok)
	}
}

// Test that parseComparison parses an invalid comparison, lacking the comparison operator, correctly.
func TestParseComparisonSansOperator(t *testing.T) {
	wantedOk := false

	s := "part2456"
	_, ok := parseComparison(s)

	if ok != wantedOk {
		t.Errorf("parseComparison returns incorrect ok %t", ok)
	}
}

// Test that parseComparison parses an invalid comparison, lacking the right operand, correctly.
func TestParseComparisonSansRightOperand(t *testing.T) {
	wantedOk := false

	s := "part1!="
	_, ok := parseComparison(s)

	if ok != wantedOk {
		t.Errorf("parseComparison returns incorrect ok %t", ok)
	}
}

// Test that parseComparison parses an invalid comparison, lacking the left operand, correctly.
func TestParseComparisonSansLeftOperand(t *testing.T) {
	wantedOk := false

	s := "!=54"
	_, ok := parseComparison(s)

	if ok != wantedOk {
		t.Errorf("parseComparison returns incorrect ok %t", ok)
	}
}

func areComparisonsEqual(left, right Comparison) bool {
	return left.LeftOperand == right.LeftOperand && left.Operator == right.Operator && left.RightOperand == right.RightOperand
}

// Test that parseMatcher parses a valid matcher correctly.
func TestParseMatcher(t *testing.T) {
	wantedLeftComparison := Comparison{
		LeftOperand:  Part1,
		Operator:     EqualToOperator,
		RightOperand: 557,
	}
	wantedRightComparison := Comparison{
		LeftOperand:  Part2,
		Operator:     UnequalToOperator,
		RightOperand: 365,
	}
	wantedOk := true

	matcher, ok := parseMatcher("part1 == 557 & part2 != 365")

	if !areComparisonsEqual(wantedLeftComparison, matcher.LeftComparison) {
		t.Error("parseMatcher returns a matcher with an invalid LeftComparison")
	}
	if !areComparisonsEqual(wantedRightComparison, matcher.RightComparison) {
		t.Error("parseMatcher returns a matcher with an invalid RightComparison")
	}
	if ok != wantedOk {
		t.Errorf("parseMatcher returns incorrect ok %t", ok)
	}
}

// Test that parseMatcher parses an invalid matcher, lacking a logical and "&", correctly.
func TestParseMatcherSansLogicalAnd(t *testing.T) {
	wantedOk := false

	_, ok := parseMatcher("part1 == 20 part2 != 03")

	if ok != wantedOk {
		t.Errorf("parseMatcher returns incorrect ok %t", ok)
	}
}

// Test that parseMatcher parses an invalid matcher, lacking a comparison before logical and "&", correctly.
func TestParseMatcherSansLeftComparison(t *testing.T) {
	wantedOk := false

	_, ok := parseMatcher("&part1<789")

	if ok != wantedOk {
		t.Errorf("parseMatcher returns incorrect ok %t", ok)
	}
}

// Test that parseMatcher parses an invalid matcher, lacking a comparison after logical and "&", correctly.
func TestParseMatcherSansRightComparison(t *testing.T) {
	wantedOk := false

	_, ok := parseMatcher("part1>53&")

	if ok != wantedOk {
		t.Errorf("parseMatcher returns incorrect ok %t", ok)
	}
}

func areMatchersEqual(left, right Matcher) bool {
	return areComparisonsEqual(left.LeftComparison, right.LeftComparison) && areComparisonsEqual(left.RightComparison, right.RightComparison)
}

func TestParseMapping(t *testing.T) {
	wantedMatcher := Matcher{
		LeftComparison: Comparison{
			LeftOperand:  Part1,
			Operator:     EqualToOperator,
			RightOperand: 1,
		},
		RightComparison: Comparison{
			LeftOperand:  Part2,
			Operator:     UnequalToOperator,
			RightOperand: 2,
		},
	}
	wantedKeyCode := 123
	wantedOk := true

	mapping, ok := parseMapping("part1 == 1 & part2 != 2 -> 123")

	if !areMatchersEqual(mapping.Matcher, wantedMatcher) {
		t.Error("parseMapping returns mapping with incorrect matcher")
	}
	if mapping.KeyCode != wantedKeyCode {
		t.Error("parseMapping returns mapping with incorrect keycode")
	}
	if ok != wantedOk {
		t.Errorf("parseMapping returns incorrect ok %t", ok)
	}
}
