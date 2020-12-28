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

type Identifier int

const (
	Part1Identifier Identifier = iota
	Part2Identifier
	NoIdentifer = -1
)

// parseComparison parses a comparison of the form: {part1,part2}{<,<=,==,!=,>=,>}<integer>.
// If s is not of the specified form, ok is false. Otherwise it is true.
func parseComparison(s string) (Identifier, ComparisonOperator, int, bool) {
	var identifier Identifier
	switch {
	case strings.HasPrefix(s, "part1"):
		identifier = Part1Identifier
	case strings.HasPrefix(s, "part2"):
		identifier = Part2Identifier
	default:
		return NoIdentifer, NoOperator, -1, false
	}
	s = s[len("partx"):] // Discard parsed identifier

	var operator ComparisonOperator
	var operatorLength int
	switch {
	case strings.HasPrefix(s, "=="):
		operator = EqualToOperator
		fallthrough
	case strings.HasPrefix(s, "!="):
		operator = UnequalToOperator
		fallthrough
	case strings.HasPrefix(s, "<="):
		operator = LessThanOrEqualToOperator
		fallthrough
	case strings.HasPrefix(s, ">="):
		operator = GreaterThanOrEqualToOperator

		operatorLength = 2
	case strings.HasPrefix(s, ">"):
		operator = LessThanOperator
		fallthrough
	case strings.HasPrefix(s, "<"):
		operator = GreaterThanOperator

		operatorLength = 1
	default:
		return NoIdentifer, NoOperator, -1, false
	}
	s = s[operatorLength:] // Discard parsed operator

	var integer int
	integer, err := strconv.Atoi(s)
	if err != nil {
		return NoIdentifer, NoOperator, -1, false
	}

	return identifier, operator, integer, true
}
