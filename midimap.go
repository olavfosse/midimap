// Map midi events to simulated keypresses
package main

import (
	"errors"
	"fmt"
	"os"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/rtmididrv"
)

func main() {
	err := dispatch(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

// dispatch parses args and decides what action the program should perform.
// If args matches the synopsis, a handler function is called and the return value of the handler function returned.
// Otherwise an error containing the "usage message" is returned.
func dispatch(args []string) error {
	// Parse arguments
	switch {
	case len(args) == 2 && args[1] == "ports":
		return ports()
	// case len(args) == 2 && args[1] == "log":
	// 	return log()
	default:
		return errors.New("usage:\tmidimap ports")
	}
}

// ports outputs a list of the available MIDI devices/ports.
func ports() error {
	drv, err := rtmididrv.New()
	if err != nil {
		return err
	}
	defer drv.Close()

	ins, err := drv.Ins()
	if err != nil {
		return err
	}

	printInPorts(ins)
	return nil
}

func printInPorts(ports []midi.In) {
	for _, port := range ports {
		fmt.Printf("%d\t%s\n", port.Number(), port.String())
	}
}
