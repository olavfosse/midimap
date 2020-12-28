// Functions for parsing midimap-lang
// A prerequisite for understanding this code, is to be familiar with midimap-lang, see midimap-lang.md for a tutorial/specification.
package lang

import (
	"strconv"
	"strings"
)

type ComparisonOperator int

const (
	LessThanOperator ComparisonOperator = iota
	LessThanOrEqualToOperator
	EqualToOperator
	UnequalToOperator
	GreaterThanOrEqualToOperator
	GreaterThanOperator
	NoOperator = -1
)

type Part1OrPart2 int

const (
	Part1 Part1OrPart2 = iota
	Part2
	NeitherPart1OrPart2 = -1
)

// parseComparison parses a comparison of the form: {part1,part2}{<,<=,==,!=,>=,>}<integer>.
// If s is not of the specified form, ok is false. Otherwise it is true.
func parseComparison(s string) (Part1OrPart2, ComparisonOperator, int, bool) {
	var leftOperand Part1OrPart2
	switch {
	case strings.HasPrefix(s, "part1"):
		leftOperand = Part1
	case strings.HasPrefix(s, "part2"):
		leftOperand = Part2
	default:
		return NeitherPart1OrPart2, NoOperator, -1, false
	}
	s = s[len("partx"):] // Discard parsed leftOperand

	var operator ComparisonOperator
	var operatorLength int
	switch {
	case strings.HasPrefix(s, "=="):
		operator = EqualToOperator
		operatorLength = 2
	case strings.HasPrefix(s, "!="):
		operator = UnequalToOperator
		operatorLength = 2
	case strings.HasPrefix(s, "<="):
		operator = LessThanOrEqualToOperator
		operatorLength = 2
	case strings.HasPrefix(s, ">="):
		operator = GreaterThanOrEqualToOperator
		operatorLength = 2
	case strings.HasPrefix(s, "<"):
		operator = LessThanOperator
		operatorLength = 1
	case strings.HasPrefix(s, ">"):
		operator = GreaterThanOperator
		operatorLength = 1
	default:
		return NeitherPart1OrPart2, NoOperator, -1, false
	}
	s = s[operatorLength:] // Discard parsed operator

	var integer int
	integer, err := strconv.Atoi(s)
	if err != nil {
		return NeitherPart1OrPart2, NoOperator, -1, false
	}

	return leftOperand, operator, integer, true
}
