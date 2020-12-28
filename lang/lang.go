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

type Comparison struct {
	LeftOperand  Part1OrPart2
	Operator     ComparisonOperator
	RightOperand int
}

// parseComparison parses a comparison of the form: {part1,part2}{<,<=,==,!=,>=,>}<integer>.
// If s is not of the specified form, it returns comparison, false otherwise it returns comparison, true.
// NB: As of now this function returns a Comparison struct even when it fails to construct it properly. That does feel a bit unclean, but I don't think it justifies using a struct pointer and nil.
func parseComparison(s string) (Comparison, bool) {
	var comparison Comparison
	switch {
	case strings.HasPrefix(s, "part1"):
		comparison.LeftOperand = Part1
	case strings.HasPrefix(s, "part2"):
		comparison.LeftOperand = Part2
	default:
		return comparison, false
	}
	s = s[len("partx"):] // Discard parsed leftOperand

	var operatorLength int
	switch {
	case strings.HasPrefix(s, "=="):
		comparison.Operator = EqualToOperator
		operatorLength = 2
	case strings.HasPrefix(s, "!="):
		comparison.Operator = UnequalToOperator
		operatorLength = 2
	case strings.HasPrefix(s, "<="):
		comparison.Operator = LessThanOrEqualToOperator
		operatorLength = 2
	case strings.HasPrefix(s, ">="):
		comparison.Operator = GreaterThanOrEqualToOperator
		operatorLength = 2
	case strings.HasPrefix(s, "<"):
		comparison.Operator = LessThanOperator
		operatorLength = 1
	case strings.HasPrefix(s, ">"):
		comparison.Operator = GreaterThanOperator
		operatorLength = 1
	default:
		return comparison, false
	}
	s = s[operatorLength:] // Discard parsed operator

	// I initially thought the two lines following this paragraph could be written more simply as the following.
	// comparison.RightOperand, err := strconv.Atoi(s)
	// Unfortunately, that caused an "expected identifier on left side of :=" error.
	// I am not quite sure why exactly that happened, but it seems like struct field are not "identifiers"
	var err error
	comparison.RightOperand, err = strconv.Atoi(s)
	if err != nil {
		return comparison, false
	}

	return comparison, true
}
