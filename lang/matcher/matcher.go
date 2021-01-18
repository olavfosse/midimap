package matcher

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"../helper"
)

// Parse parses a matcher as specified in Section 1.2.1 MATCHERS of the midimap-lang specification.
//
// If s is a valid matcher as described by the specification, Parse returns matcher, nil.
// Otherwise, Parse returns an error describing why the matcher is invalid.
// s may not contain any leading or trailing space characters.
func Parse(s string) (Matcher, error) {
	left, right, operator := beforeAndAfterLogicalOperator(s)
	if operator == NoLogicalOperator {
		return parseMatcherWithoutLogicalOperator(s)
	}

	var matcher MatcherWithLogicalOperator
	matcher.Operator = operator
	var err error
	matcher.LeftMatcher, err = Parse(strings.TrimSpace(left))
	if err != nil {
		return matcher, err
	}
	matcher.RightMatcher, err = Parse(strings.TrimSpace(right))
	return matcher, err
}

// parseMatcherWithoutLogicalOperator parses a matcher without a logical operator.
//
// If s is a valid matcher, parseMatcherWithoutLogicalOperator returns matcher, nil.
// Otherwise, parseMatcherWithoutLogicalOperator returns an error describing why the matcher without
// s may not contain any leading or trailing space characters or any logical operators, that is && or ||.
func parseMatcherWithoutLogicalOperator(s string) (m MatcherWithoutLogicalOperator, err error) {
	unParsed := s // the characters of s which are yet to be parsed

	switch {
	case strings.HasPrefix(unParsed, "data1"):
		m.LeftOperand = Data1
	case strings.HasPrefix(unParsed, "data2"):
		m.LeftOperand = Data2
	default:
		err = fmt.Errorf("matcher %q: no valid left operand", s)
		return
	}
	unParsed = unParsed[len("datax"):] // Discard parsed leftOperand

	skipToNonSpaceCharacter(&unParsed)
	var operatorLength int
	switch {
	case strings.HasPrefix(unParsed, "=="):
		m.Operator = EqualToOperator
		operatorLength = 2
	case strings.HasPrefix(unParsed, "!="):
		m.Operator = UnequalToOperator
		operatorLength = 2
	case strings.HasPrefix(unParsed, "<="):
		m.Operator = LessThanOrEqualToOperator
		operatorLength = 2
	case strings.HasPrefix(unParsed, ">="):
		m.Operator = GreaterThanOrEqualToOperator
		operatorLength = 2
	case strings.HasPrefix(unParsed, "<"):
		m.Operator = LessThanOperator
		operatorLength = 1
	case strings.HasPrefix(unParsed, ">"):
		m.Operator = GreaterThanOperator
		operatorLength = 1
	default:
		err = fmt.Errorf("matcher %q: no valid comparison operator", s)
		return
	}
	unParsed = unParsed[operatorLength:] // Discard parsed operator

	skipToNonSpaceCharacter(&unParsed)
	n, err := strconv.ParseInt(unParsed, 10, 64)
	if err != nil {
		err = fmt.Errorf("matcher %q: no valid right operand", s)
		return
	}
	m.RightOperand = int64(n)

	return
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

// beforeAndAfterLogicalOperator splits a string into two substrings, one before a logical operator and the other after it.
// If s contains "&&" beforeAndAfterLogicalOperator returns before, after, LogicalAndOperator, where before is the substring of s which appear before "&&" and after is the substring of s which appear after "&&".
// Otherwise if s contains "||" beforeAndAfterLogicalOperator returns before, after, LogicalOrOperator, where before is the substring of s which appear before "||" and after is the substring of s which appear after it.
// If s contains neither "||" nor "&&", beforeAndAfterLogicalOperator returns "", "", NoLogicalOperator
func beforeAndAfterLogicalOperator(s string) (string, string, LogicalOperator) {
	logicalAndRegexp := regexp.MustCompilePOSIX("&&")
	logicalOrRegexp := regexp.MustCompilePOSIX(`\|\|`)

	before, after, ok := helper.BeforeAndAfter(logicalAndRegexp, s)
	if ok {
		return before, after, LogicalAndOperator
	}

	before, after, ok = helper.BeforeAndAfter(logicalOrRegexp, s)
	if ok {
		return before, after, LogicalOrOperator
	}

	return "", "", NoLogicalOperator
}

// Matcher is a discriminated union of MatcherWithoutLogicalOperator and MatcherWithLogicalOperator.
//
// This method of representing a syntax tree is based on the following article.
// https://eli.thegreenplace.net/2018/go-and-algebraic-data-types/
// Please read it if you encounter difficulties understanding something.
type Matcher interface {
	isMatcher()
}

// MatcherWithoutLogicalOperator represents a simple matcher with no logical operator, such as:
// data1 == 1
// data2 < 4
// data1 != 37
type MatcherWithoutLogicalOperator struct {
	LeftOperand  Data1OrData2
	Operator     ComparisonOperator
	RightOperand int64
}

func (_ MatcherWithoutLogicalOperator) isMatcher() {}

type Data1OrData2 int

const (
	Data1 Data1OrData2 = iota
	Data2
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

// MatcherWithLogicalOperator represents a matcher with at least one logical operator, such as:
// data1 == 1 || data1 > 10
// (data1 != 1 && data1 != 51) || data2 > 10
// It can represent a matcher of abitrary nesting.
type MatcherWithLogicalOperator struct {
	LeftMatcher  Matcher
	Operator     LogicalOperator
	RightMatcher Matcher
}

func (_ MatcherWithLogicalOperator) isMatcher() {}

type LogicalOperator int

const (
	LogicalAndOperator LogicalOperator = iota
	LogicalOrOperator
	NoLogicalOperator LogicalOperator = -1
)
