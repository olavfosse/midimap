// Map midi events to simulated keypresses
package main

import (
	"fmt"
	"os"

	// Command modifiers
	"./cm/ports"
	// mapcm's naming scheme is different from the others, which
	// is required since map is a keyword. This is just temporary
	// since map shall be replaced with load+record soon.
	"./cm/log"
	"./cm/mapcm"

	"./usage"
)

func main() {
	err := mainish()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

// mainish is like main, except if it encounters an error it returns it instead of presenting it to the user.
func mainish() error {
	if len(os.Args) < 2 {
		return usage.Usage()
	}

	commandModifiers := map[string](func([]string) error){
		"ports": ports.Ports,
		"map":   mapcm.MapCM,
		"log":   log.Log,
	}
	cm, ok := commandModifiers[os.Args[1]]
	if !ok {
		return usage.Usage()
	}

	return cm(os.Args[2:])
}
