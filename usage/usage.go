package usage

import (
	"errors"
	"strings"
)

func Usage() error {
	return errors.New(strings.TrimSpace(`
usage:	midimap ports
	midimap map portnumber mapname
	midimap log portnumber [matcher]`))
}
