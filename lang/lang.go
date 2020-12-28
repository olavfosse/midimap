// Functions for parsing midimap-lang
// A prerequisite for understanding this code, is to be familiar with midimap-lang, see midimap-lang.md for a tutorial/specification.
package lang

import (
	"regexp"
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

// parseComparison parses a comparison of the following form.
// {part1,part2}{<,<=,==,!=,>=,>}<integer>
// Spaces may be intersped anywhere without changing the result.
// If s is not of the specified form, it returns comparison, false otherwise it returns comparison, true.
// NB: As of now this function returns a Comparison struct even when it fails to construct it properly. That does feel a bit unclean, but I don't think it justifies using a struct pointer and nil.
func parseComparison(s string) (Comparison, bool) {
	s = strings.ReplaceAll(s, " ", "") // Remove spaces

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

// Currently this requires the matchers to be of the following form:
// {part1,part2}{<,<=,==,!=,>=,>}integer & {part1,part2}{<,<=,==,!=,>=,>}integer
// That only covers a subset of all legal matchers according to the spec.
// I will update this to behave as specified in the spec when I learn how to do polymorphism(presumably with interfaces) in go.
type Matcher struct {
	LeftComparison  Comparison
	RightComparison Comparison
}

// beforeAndAfter splits a string into two substrings, one before the leftmost match of a regexp and the other after it.
// If r matches s, beforeAndAfter returns (left, right, true) where left is the characters prior to the leftmost match and right is the characters to the right of the leftmost match.
// Otherwise, that is if r does not match s, ("", "", false) is returned.
func beforeAndAfter(r *regexp.Regexp, s string) (string, string, bool) {
	loc := r.FindStringIndex(s)
	if loc == nil {
		return "", "", false
	}

	before, after := s[:loc[0]], s[loc[1]:]
	return before, after, true

}

// parseMatcher parses a matcher of the following form.
// comparison&comparison
// Spaces may be intersped anywhere without changing the result.
// If s is of the specified form, it returns matcher, true otherwise it returns matcher, false.
// NB: As of now this function returns a Matcher struct even when it fails to construct it properly. That does feel a bit unclean, but I don't think it justifies using a struct pointer and nil.
func parseMatcher(s string) (Matcher, bool) {
	var matcher Matcher
	// split on &
	r := regexp.MustCompilePOSIX("&")
	left, right, ok := beforeAndAfter(r, s)
	// report if the split failed
	if !ok { // missing logical and "&"
		return matcher, false
	}

	// parse comparisons from before & and after &
	// report if parsing comparisons failed
	leftComparison, ok := parseComparison(left)
	if !ok {
		return matcher, false
	}
	matcher.LeftComparison = leftComparison
	rightComparison, ok := parseComparison(right)
	if !ok {
		return matcher, false
	}
	matcher.RightComparison = rightComparison

	return matcher, true
}
