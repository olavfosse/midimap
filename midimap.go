package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func main() {
	err := mainish()
	if err != nil {
		if err != errUsage {
			fmt.Print("error: ")
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
	switch os.Args[1] {
	case "ports":
		return portsCommandModifier(os.Args[2:])
	case "map":
		return mapCommandModifier(os.Args[2:])
	case "log":
		return logCommandModifier(os.Args[2:])
	default:
		return errUsage
	}
}

var errUsage = errors.New(strings.TrimSpace(`
usage:	midimap ports
	midimap map portnumber mapname
	midimap log portnumber [matcher]`))

// 	midimap map portnumber mapname
//	midimap log portnumber [matcher]`))
