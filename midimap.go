// Map midi events to simulated keypresses
package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
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
	case len(args) == 3 && args[1] == "log":
		portNumber, err := strconv.ParseUint(args[2], 10, 0)
		if err != nil {
			return errors.New("error: portnumber must be a valid unsigned integer")
		}
		return log(portNumber)
	default:
		return errors.New("usage:\tmidimap ports\n\tmidimap log portnumber")
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

	printIns(ins)
	return nil
}

// log logs incoming MIDI events from the port with the number portNumber.
func log(portNumber uint64) error {
	drv, err := rtmididrv.New()
	if err != nil {
		return err
	}
	defer drv.Close()

	ins, err := drv.Ins()
	if err != nil {
		return err
	}

	in, ok := getInByPortNumber(ins, portNumber)
	if !ok {
		return errors.New(fmt.Sprintf("error: no MIDI port by number %d", portNumber))
	}
	err = in.Open()
	if err != nil {
		return err
	}

	defer in.Close()

	rd := reader.New(
		reader.NoLogger(),
		reader.Each(func(pos *reader.Position, msg midi.Message) {
			fmt.Printf("got %s\n", msg)
		}),
	)
	err = rd.ListenTo(in)
	if err != nil {
		return err
	}

	// HACK: I couldn't figure out the proper way to read until failure, so I just leave it running for an hour. Should be replaced with waiting until the the device is disconnected or something goes wrong
	d, _ := time.ParseDuration("1h")
	time.Sleep(d)

	return err
}

// getInByPortNumber retrieves the midi.In by number(!= index) portNumber from ins.
// If there is a In with portNumber number in ins, getInByPortNumber returns in, true where in is the in by number number.
// Otherwise, getInByPortNumber return dummyIn, false, where dummyIn is the a In with In zero value.
func getInByPortNumber(ins []midi.In, number uint64) (in midi.In, ok bool) {
	for _, innerIn := range ins {
		if uint64(innerIn.Number()) == number {
			ok = true
			in = innerIn
			return
		}
	}
	return
}

func printIns(ins []midi.In) {
	for _, in := range ins {
		fmt.Printf("%d\t%s\n", in.Number(), in.String())
	}
}
