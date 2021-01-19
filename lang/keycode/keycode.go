package keycode

import (
	"fmt"
	"strconv"
)

// Parse parses a keycode as specified in Section 1.2.2 KEYCODES of the midimap-lang specification.
//
// If s is a valid keycode as described by the specification, Parse returns keycode, nil.
// Otherwise, Parse returns an error describing why the keycode is invalid.
// s may not contain any leading or trailing spaces.
func Parse(s string) (int, error) {
	keycode, err := strconv.Atoi(s)
	if err != nil {
		return keycode, fmt.Errorf("keycode %q: invalid", s)
	}
	return keycode, nil
}
