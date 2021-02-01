package mapping

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fossegrim/midimap/lang/helper"
	"github.com/fossegrim/midimap/lang/keycode"
	"github.com/fossegrim/midimap/lang/matcher"
)

type Mapping struct {
	Matcher matcher.Matcher
	Keycode int
}

func (m Mapping) Equal(n Mapping) bool {
	return m.Matcher.Equal(n.Matcher) && m.Keycode == n.Keycode
}

// Parse parses a mapping as specified in Section 1.2 MAPPINGS of the midimap-lang specification.
//
// If s is a valid mapping as described by the specification, Parse returns mapping, nil.
// Otherwise, Parse returns an error describing why the mapping is invalid.
func Parse(s string) (mapping Mapping, err error) {
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
	mapping.Keycode, err = keycode.Parse(strings.TrimSpace(after))
	return
}
