// Functions for parsing midimap-lang
// A prerequisite for understanding this code, is to be familiar with midimap-lang, see midimap-lang.md for a tutorial/specification.
package lang

import (
	"bufio"
	"errors"
	"fmt"
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
)

type Comparison struct {
	LeftOperand  Part1OrPart2
	Operator     ComparisonOperator
	RightOperand int64
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

	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return comparison, false
	}
	comparison.RightOperand = int64(n)

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
	matcher.LeftComparison, ok = parseComparison(left)
	if !ok {
		return matcher, false
	}
	matcher.RightComparison, ok = parseComparison(right)
	if !ok {
		return matcher, false
	}

	return matcher, true
}

// parseKeyCode parses a KeyCode of the following form.
// integer
// Spaces may be be intersped anywhere without changing the result.
// If s is of the specified form, it returns keyCode, true, otherwise it returns keyCode, false.
func parseKeyCode(s string) (int, bool) {
	s = strings.ReplaceAll(s, " ", "") // Remove spaces
	keyCode, err := strconv.Atoi(s)
	if err != nil {
		return keyCode, false
	}
	return keyCode, true
}

type Mapping struct {
	Matcher Matcher
	KeyCode int
}

// parseMapping parses a mapping of the following form.
// matcher->keycode
// Spaces may be intersped anywhere without changing the result.
// If s is of the specified form, it returns mapping, true otherwise it returns mapping, false.
// NB: As of now this function returns a Mapping struct even when it fails to construct it properly. That does feel a bit unclean, but I don't think it justifies using a struct pointer and nil.
func parseMapping(s string) (Mapping, bool) {
	var mapping Mapping
	r := regexp.MustCompilePOSIX("- *>")
	before, after, ok := beforeAndAfter(r, s)
	if !ok {
		return mapping, false
	}
	mapping.Matcher, ok = parseMatcher(before)
	if !ok {
		return mapping, false
	}
	mapping.KeyCode, ok = parseKeyCode(after)
	if !ok {
		return mapping, false
	}
	return mapping, true
}

// NextMapping parses a midimap-lang mapping from r.
// NextMapping reads lines until it encounters a non-comment line or EOF.
// If EOF is reached NextMapping returns (nil, io.EOF).
// If line is not a valid mapping return (nil, err), where err is an error describing where the error occured.
func NextMapping(r *bufio.Reader) (Mapping, error) {
	var mapping Mapping
	var line string
	for {
		s, err := r.ReadString('\n')
		if err != nil {
			return mapping, err
		}
		line = s[:len(s)-1]
		// skip comments
		if !strings.HasPrefix(line, "#") {
			break
		}
	}

	// I wonder if there is an idiom to return all the return values of a called function.
	// It sounds a bit sugarish, so probably not.
	mapping, ok := parseMapping(line)
	var err error = nil
	if !ok {
		err = errors.New(fmt.Sprintf("invalid mapping %q", line))
	}

	return mapping, err
}
