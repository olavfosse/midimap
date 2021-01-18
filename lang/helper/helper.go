// The helper package contains random helper functions that are needed by several unrelated lang/ packages.
package helper

import "regexp"

// BeforeAndAfter splits a string into two substrings, one before the leftmost match of a regexp and the other after it.
// If r matches s, BeforeAndAfter returns (left, right, true) where left is the characters prior to the leftmost match and right is the characters to the right of the leftmost match.
// Otherwise, that is if r does not match s, ("", "", false) is returned.
func BeforeAndAfter(r *regexp.Regexp, s string) (string, string, bool) {
	loc := r.FindStringIndex(s)
	if loc == nil {
		return "", "", false
	}

	before, after := s[:loc[0]], s[loc[1]:]
	return before, after, true

}
