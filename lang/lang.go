// Functions for parsing midimap-lang code.
//
// Before reading this code it is advised that you read through and comprehend the midimap-lang specification.
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
)

type Data1OrData2 int

const (
	Data1 Data1OrData2 = iota
	Data2
)

type Comparison struct {
	LeftOperand  Data1OrData2
	Operator     ComparisonOperator
	RightOperand int64
}

// skipToNonSpaceCharacter moves the start of *s to the next non-space character.
// If there is a non-space character in *s, *s is modified such that the first character in *s is that non-space character, and true is returned.
// Otherwise *s is not modified and false is returned.
func skipToNonSpaceCharacter(s *string) (ok bool) {
	for i, c := range *s {
		if c != ' ' {
			*s = (*s)[i:]
			return true
		}
	}
	return false
}
	

// parseComparison parses a comparison as specified in Section 1.2.1.1 COMPARISONS of the midimap-lang specification.
// If s is a valid comparison as described by the specification, parseComparison returns comparison, true.
// Otherwise parseComparison returns comparison, false.
func parseComparison(s string) (Comparison, error) {
	unParsed := s // the characters of s which are yet to be parsed

	skipToNonSpaceCharacter(&unParsed)
	var comparison Comparison
	switch {
	case strings.HasPrefix(unParsed, "data1"):
		comparison.LeftOperand = Data1
	case strings.HasPrefix(unParsed, "data2"):
		comparison.LeftOperand = Data2
	default:
		return comparison, fmt.Errorf("Comparison %q does not have a valid left operand", s)
	}
	unParsed = unParsed[len("datax"):] // Discard parsed leftOperand

	skipToNonSpaceCharacter(&unParsed)
	var operatorLength int
	switch {
	case strings.HasPrefix(unParsed, "=="):
		comparison.Operator = EqualToOperator
		operatorLength = 2
	case strings.HasPrefix(unParsed, "!="):
		comparison.Operator = UnequalToOperator
		operatorLength = 2
	case strings.HasPrefix(unParsed, "<="):
		comparison.Operator = LessThanOrEqualToOperator
		operatorLength = 2
	case strings.HasPrefix(unParsed, ">="):
		comparison.Operator = GreaterThanOrEqualToOperator
		operatorLength = 2
	case strings.HasPrefix(unParsed, "<"):
		comparison.Operator = LessThanOperator
		operatorLength = 1
	case strings.HasPrefix(unParsed, ">"):
		comparison.Operator = GreaterThanOperator
		operatorLength = 1
	default:
		return comparison, fmt.Errorf("Comparison %q does not have a valid operator", s)
	}
	unParsed = unParsed[operatorLength:] // Discard parsed operator

	skipToNonSpaceCharacter(&unParsed)
	n, err := strconv.ParseInt(unParsed, 10, 64)
	if err != nil {
		return comparison, fmt.Errorf("Comparison %q does not have a valid right operand", s)
	}
	comparison.RightOperand = int64(n)
	
	return comparison, nil
}

type LogicalOperator int

const (
	LogicalAndOperator LogicalOperator = iota
	LogicalOrOperator
	NoLogicalOperator LogicalOperator = -1
)

type Matcher struct {
	LeftComparison  Comparison
	Operator        LogicalOperator
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

// beforeAndAfterLogicalOperator splits a string into two substrings, one before a logical operator and the other after it.
// If s contains "&&" beforeAndAfterLogicalOperator returns before, after, LogicalAndOperator, where before is the substring of s which appear before "&&" and after is the substring of s which appear after "&&".
// Otherwise if s contains "||" beforeAndAfterLogicalOperator returns before, after, LogicalOrOperator, where before is the substring of s which appear before "||" and after is the substring of s which appear after it.
// If s contains neither "||" nor "&&", beforeAndAfterLogicalOperator returns "", "", NoLogicalOperator
func beforeAndAfterLogicalOperator(s string) (string, string, LogicalOperator) {
	logicalAndRegexp := regexp.MustCompilePOSIX("&&")
	logicalOrRegexp := regexp.MustCompilePOSIX(`\|\|`)

	before, after, ok := beforeAndAfter(logicalAndRegexp, s)
	if ok {
		return before, after, LogicalAndOperator
	}

	before, after, ok = beforeAndAfter(logicalOrRegexp, s)
	if ok {
		return before, after, LogicalOrOperator
	}

	return "", "", NoLogicalOperator
}

// parseMATCHER parses a MATCHER as specified in Section 1.2.1 MATCHERS of the midimap-lang specification.
// If s is a valid MATCHER as described by the specification, parseMATCHER returns matcher, true.
// Otherwise parseMATCHER returns matcher, false.
func parseMATCHER(s string) (Matcher, bool) {
	var matcher Matcher

	var left, right string
	left, right, matcher.Operator = beforeAndAfterLogicalOperator(s)
	if matcher.Operator == NoLogicalOperator {
		return matcher, false
	}

	var err error
	matcher.LeftComparison, err = parseComparison(left)
	if err != nil {
		return matcher, false
	}
	matcher.RightComparison, err = parseComparison(right)
	if err != nil {
		return matcher, false
	}

	return matcher, true
}

// parseKEYCODE parses a KEYCODE as specified in Section 1.2.2 KEYCODES of the midimap-lang specification.
// If s is a valid KEYCODE as described by the specification, parseKEYCODE returns keyCode, true.
// Otherwise parseKEYCODE returns keyCode, false.
func parseKEYCODE(s string) (int, bool) {
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

// parseMAPPING parses a MAPPING as specified in Section 1.2 MAPPINGS of the midimap-lang specification.
// If s is a valid MAPPING as described by the specification, parseMAPPING returns mapping, true.
// Otherwise parseMAPPING returns mapping, false.
func parseMAPPING(s string) (Mapping, bool) {
	var mapping Mapping
	r := regexp.MustCompilePOSIX("- *>")
	before, after, ok := beforeAndAfter(r, s)
	if !ok {
		return mapping, false
	}
	mapping.Matcher, ok = parseMATCHER(before)
	if !ok {
		return mapping, false
	}
	mapping.KeyCode, ok = parseKEYCODE(after)
	if !ok {
		return mapping, false
	}
	return mapping, true
}

// NextMAPPING attemps to parse the next MAPPING, as specified in Section 1.2 MAPPINGS of the midimap-lang specification, from r by parsing lines until a MAPPING is reached or an io error occurs.
// If an io error occured NextMAPPING returns mapping, ioError.
// If an invalid mapping is reached NextMAPPING returns mapping, err, where err is an error describing how the mapping is invalid.
// Otherwise NextMAPPING returns mapping, nil.
func NextMAPPING(r *bufio.Reader) (Mapping, error) {
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
	mapping, ok := parseMAPPING(line)
	var err error = nil
	if !ok {
		err = errors.New(fmt.Sprintf("invalid mapping %q", line))
	}

	return mapping, err
}
