package main

import (
	"fmt"

	"gitlab.com/gomidi/midi"
)

// portsCommandModifier corresponds to the ports command modifier. args
// corresponds to the list of arguments listed on the command line after the
// ports command modifier.
//
// For documentation about the ports command modifier itself, consult
// midimap(1).
func portsCommandModifier(args []string) error {
	if len(args) != 0 {
		return errUsage
	}

	drv, err := newDriver()
	if err != nil {
		return err
	}
	defer drv.Close()

	ins, err := drv.Ins()
	if err != nil {
		return err
	}

	printIns(ins)
	return nil
}

func printIns(ins []midi.In) {
	for _, in := range ins {
		fmt.Printf("%d\t%s\n", in.Number(), in.String())
	}
}
