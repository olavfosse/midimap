package ports

import (
	"fmt"

	"../../driver"
	"../../usage"
	"gitlab.com/gomidi/midi"
)

// Ports corresponds to the ports command modifier. args corresponds
// to the list of arguments listed on the command line after the ports
// command modifier.
//
// For documentation about the ports command modifier itself, consult
// the (to be written) manual.
// TODO: Write midimap(1) and/or midimap-ports(1)
func Ports(args []string) error {
	if len(args) != 0 {
		return usage.Usage()
	}

	drv, err := driver.New()
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
