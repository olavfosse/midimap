// Functions for parsing midimap-lang code.
//
// Before reading this code it is advised that you read through and comprehend the midimap-lang specification.
package lang

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"./helper"
	"./matcher"
)

// parseKeycode parses a keycode as specified in Section 1.2.2 KEYCODES of the midimap-lang specification.
//
// If s is a valid keycode as described by the specification, parseKeycode returns keycode, nil.
// Otherwise, parseKeycode returns an error describing why the keycode is invalid.
// s may not contain any leading or trailing spaces.
func parseKeycode(s string) (int, error) {
	keycode, err := strconv.Atoi(s)
	if err != nil {
		return keycode, fmt.Errorf("keycode %q: invalid", s)
	}
	return keycode, nil
}

type Mapping struct {
	Matcher matcher.Matcher
	Keycode int
}

// parseMapping parses a mapping as specified in Section 1.2 MAPPINGS of the midimap-lang specification.
//
// If s is a valid mapping as described by the specification, parseMapping returns mapping, nil.
// Otherwise, parseMapping returns an error describing why the mapping is invalid.
func parseMapping(s string) (mapping Mapping, err error) {
	r := regexp.MustCompilePOSIX("->")
	before, after, ok := helper.BeforeAndAfter(r, s)
	if !ok {
		err = fmt.Errorf("mapping %q: no valid separator", s)
		return
	}
	mapping.Matcher, err = matcher.Parse(strings.TrimSpace(before))
	if err != nil {
		return
	}
	mapping.Keycode, err = parseKeycode(strings.TrimSpace(after))
	return
}

// NextMapping parses a mapping as specified in Section 1.2 MAPPINGS of the midimap-lang specification, by parsing lines until a mapping is reached or an io error occurs.
//
// If an io error occurs, NextMapping returns the io error.
// If a parsing error occurs, NextMapping returns the parsing error.
// Otherwise, NextMapping returns mapping, nil.
func NextMapping(r *bufio.Reader) (mapping Mapping, err error) {
	var line string
	for {
		var s string
		s, err = r.ReadString('\n')
		if err != nil {
			return
		}
		line = s[:len(s)-1]
		// skip comments
		if !strings.HasPrefix(line, "#") {
			break
		}
	}

	return parseMapping(line)
}
