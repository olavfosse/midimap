// Before descending into the packages in this directory, it is advised to read and comprehend the midimap-lang specifications.
package lang

import (
	"bufio"
	"strings"

	"./mapping"
)

// NextMapping parses a mapping as specified in Section 1.2 MAPPINGS of the midimap-lang specification, by parsing lines until a mapping is reached or an io error occurs.
//
// If an io error occurs, NextMapping returns the io error.
// If a parsing error occurs, NextMapping returns the parsing error.
// Otherwise, NextMapping returns m, nil.
func NextMapping(r *bufio.Reader) (m mapping.Mapping, err error) {
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

	return mapping.Parse(line)
}
