// Map midi events to simulated keypresses
package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	// Command modifiers
	"github.com/fossegrim/midimap/cm/ports"
	// mapcm's naming scheme is different from the others, which
	// is required since map is a keyword. This is just temporary
	// since map shall be replaced with load+record soon.
	"github.com/fossegrim/midimap/cm/log"
	"github.com/fossegrim/midimap/cm/mapcm"
)

func main() {
	err := mainish()
	if err != nil {
		if err != errUsage {
			fmt.Print("error:")
		}
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

// mainish is like main, except if it encounters an error it returns it instead of presenting it to the user.
func mainish() error {
	if len(os.Args) < 2 {
		return errUsage
	}

	commandModifiers := map[string](func([]string) error){
		"ports": ports.Ports,
		"map":   mapcm.MapCM,
		"log":   log.Log,
	}
	cm, ok := commandModifiers[os.Args[1]]
	if !ok {
		return errUsage
	}

	return cm(os.Args[2:])
}

var errUsage = errors.New(strings.TrimSpace(`
usage:	midimap ports
	midimap map portnumber mapname
	midimap log portnumber [matcher]`))
