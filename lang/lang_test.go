package lang

import (
	"testing"
)

// Test that parseCOMPARISON parses a valid comparison, where the left operand is data1, correctly.
func TestParseData1Comparison(t *testing.T) {
	wantedLeftOperand := Data1
	wantedOperator := GreaterThanOperator
	var wantedRightOperand int64 = 123
	wantedOk := true

	s := "data1 >123"
	comparison, ok := parseCOMPARISON(s)

	if comparison.LeftOperand != wantedLeftOperand {
		t.Errorf("parseCOMPARISON returns a comparison with an incorrect leftOperand %v, want %v", comparison.LeftOperand, wantedLeftOperand)
	}
	if comparison.Operator != wantedOperator {
		t.Errorf("parseCOMPARISON returns a comparison with an incorrect operator %v, want %v", comparison.Operator, wantedOperator)
	}
	if comparison.RightOperand != wantedRightOperand {
		t.Errorf("parseCOMPARISON returns a comparison with an incorrect rightOperand %d, want %d", comparison.RightOperand, wantedRightOperand)
	}
	if ok != wantedOk {
		t.Errorf("parseCOMPARISON returns incorrect ok %t", ok)
	}
}

// Test that parseCOMPARISON parses a valid comparison, where the left operand is data2, correctly.
func TestParseData2Comparison(t *testing.T) {
	wantedLeftOperand := Data2
	wantedOperator := EqualToOperator
	var wantedRightOperand int64 = 321
	wantedOk := true

	s := "data2 = =  321"
	comparison, ok := parseCOMPARISON(s)

	if comparison.LeftOperand != wantedLeftOperand {
		t.Errorf("parseCOMPARISON returns a comparison with an incorrect leftOperand %v, want %v", comparison.LeftOperand, wantedLeftOperand)
	}
	if comparison.Operator != wantedOperator {
		t.Errorf("parseCOMPARISON returns a comparison with an incorrect operator %v, want %v", comparison.Operator, wantedOperator)
	}
	if comparison.RightOperand != wantedRightOperand {
		t.Errorf("parseCOMPARISON returns a comparison with an incorrect rightOperand %d, want %d", comparison.RightOperand, wantedRightOperand)
	}
	if ok != wantedOk {
		t.Errorf("parseCOMPARISON returns incorrect ok %t", ok)
	}
}

// Test that parseCOMPARISON parses an invalid comparison, lacking the comparison operator, correctly.
func TestParseComparisonSansOperator(t *testing.T) {
	wantedOk := false

	s := "data2456"
	_, ok := parseCOMPARISON(s)

	if ok != wantedOk {
		t.Errorf("parseCOMPARISON returns incorrect ok %t", ok)
	}
}

// Test that parseCOMPARISON parses an invalid comparison, lacking the right operand, correctly.
func TestParseComparisonSansRightOperand(t *testing.T) {
	wantedOk := false

	s := "data1!="
	_, ok := parseCOMPARISON(s)

	if ok != wantedOk {
		t.Errorf("parseCOMPARISON returns incorrect ok %t", ok)
	}
}

// Test that parseCOMPARISON parses an invalid comparison, lacking the left operand, correctly.
func TestParseComparisonSansLeftOperand(t *testing.T) {
	wantedOk := false

	s := "!=54"
	_, ok := parseCOMPARISON(s)

	if ok != wantedOk {
		t.Errorf("parseCOMPARISON returns incorrect ok %t", ok)
	}
}

func areComparisonsEqual(left, right Comparison) bool {
	return left.LeftOperand == right.LeftOperand && left.Operator == right.Operator && left.RightOperand == right.RightOperand
}

// Test that parseMATCHER parses a valid matcher, where the logical operator is a logical and-operator "&&", correctly.
func TestParseLogicalAndMatcher(t *testing.T) {
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
	wantedOk := true

	matcher, ok := parseMATCHER("data1 == 557 && data2 != 365")

	if !areComparisonsEqual(wantedLeftComparison, matcher.LeftComparison) {
		t.Error("parseMATCHER returns a matcher with an invalid LeftComparison")
	}
	if wantedOperator != matcher.Operator {
		t.Errorf("parseMATCHER returns a matcher with an invalid Operator %v, want %v", matcher.Operator, wantedOperator)
	}
	if !areComparisonsEqual(wantedRightComparison, matcher.RightComparison) {
		t.Error("parseMATCHER returns a matcher with an invalid RightComparison")
	}
	if ok != wantedOk {
		t.Errorf("parseMATCHER returns incorrect ok %t", ok)
	}
}

// Test that parseMATCHER parses a valid MATCHER, where the logical operator is a logical or-operator "||", correctly.
func TestParseLogicalOrMatcher(t *testing.T) {
	wantedLeftComparison := Comparison{
		LeftOperand:  Data2,
		Operator:     UnequalToOperator,
		RightOperand: 93,
	}
	wantedOperator := LogicalOrOperator
	wantedRightComparison := Comparison{
		LeftOperand:  Data1,
		Operator:     EqualToOperator,
		RightOperand: 122,
	}
	wantedOk := true

	matcher, ok := parseMATCHER("data2 != 93 || data1 == 122")

	if !areComparisonsEqual(wantedLeftComparison, matcher.LeftComparison) {
		t.Error("parseMATCHER returns a matcher with an invalid LeftComparison")
	}
	if wantedOperator != matcher.Operator {
		t.Errorf("parseMATCHER returns a matcher with an invalid Operator %v, want %v", matcher.Operator, wantedOperator)
	}
	if !areComparisonsEqual(wantedRightComparison, matcher.RightComparison) {
		t.Error("parseMATCHER returns a matcher with an invalid RightComparison")
	}
	if ok != wantedOk {
		t.Errorf("parseMATCHER returns incorrect ok %t", ok)
	}
}

// Test that parseMATCHER parses an invalid matcher, lacking a logical and "&&", correctly.
func TestParseMatcherSansLogicalAnd(t *testing.T) {
	wantedOk := false

	_, ok := parseMATCHER("data1 == 20 data2 != 03")

	if ok != wantedOk {
		t.Errorf("parseMATCHER returns incorrect ok %t", ok)
	}
}

// Test that parseMATCHER parses an invalid matcher, lacking a comparison before logical and "&&", correctly.
func TestParseMatcherSansLeftComparison(t *testing.T) {
	wantedOk := false

	_, ok := parseMATCHER("&&data1<789")

	if ok != wantedOk {
		t.Errorf("parseMATCHER returns incorrect ok %t", ok)
	}
}

// Test that parseMATCHER parses an invalid matcher, lacking a comparison after logical and "&&", correctly.
func TestParseMatcherSansRightComparison(t *testing.T) {
	wantedOk := false

	_, ok := parseMATCHER("data1>53&&")

	if ok != wantedOk {
		t.Errorf("parseMATCHER returns incorrect ok %t", ok)
	}
}

func areMatchersEqual(left, right Matcher) bool {
	areLeftComparisonsEqual := areComparisonsEqual(left.LeftComparison, right.LeftComparison)
	areOperatorsEqual := left.Operator == right.Operator
	areRightComparisonsEqual := areComparisonsEqual(left.RightComparison, right.RightComparison)
	return areLeftComparisonsEqual && areOperatorsEqual && areRightComparisonsEqual
}

func TestParseMapping(t *testing.T) {
	wantedMatcher := Matcher{
		LeftComparison: Comparison{
			LeftOperand:  Data1,
			Operator:     EqualToOperator,
			RightOperand: 1,
		},
		RightComparison: Comparison{
			LeftOperand:  Data2,
			Operator:     UnequalToOperator,
			RightOperand: 2,
		},
	}
	wantedKeyCode := 123
	wantedOk := true

	mapping, ok := parseMAPPING("data1 == 1 && data2 != 2 -> 123")

	if !areMatchersEqual(mapping.Matcher, wantedMatcher) {
		t.Error("parseMAPPING returns mapping with incorrect matcher")
	}
	if mapping.KeyCode != wantedKeyCode {
		t.Error("parseMAPPING returns mapping with incorrect keycode")
	}
	if ok != wantedOk {
		t.Errorf("parseMAPPING returns incorrect ok %t", ok)
	}
}

func TestParseMappingSansSeparator(t *testing.T) {
	wantedOk := false

	_, ok := parseMAPPING("data2 != 3 && data2 >= 2 123")

	if ok != wantedOk {
		t.Errorf("parseMAPPING returns incorrect ok %t", ok)
	}
}

func TestParseMappingSansMatcher(t *testing.T) {
	wantedOk := false

	_, ok := parseMAPPING("-> 123")

	if ok != wantedOk {
		t.Errorf("parseMAPPING returns incorrect ok %t", ok)
	}
}

func TestParseMappingSansKeycode(t *testing.T) {
	wantedOk := false

	_, ok := parseMAPPING("data2 != 3 && data2 ->")

	if ok != wantedOk {
		t.Errorf("parseMAPPING returns incorrect ok %t", ok)
	}
}
