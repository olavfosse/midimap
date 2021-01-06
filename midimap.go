// Map midi events to simulated keypresses
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"./lang"
	"./press"
	"github.com/rakyll/portmidi"
)

func main() {
	// Parse arguments
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: midimap map")
		os.Exit(1)
	}
	mapFileName := os.Args[1]

	// Parse midi map
	mapFile, err := os.Open(mapFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	r := bufio.NewReader(mapFile)
	var mappings []lang.Mapping
	for {
		mapping, err := lang.NextMapping(r)
		if err == io.EOF {
			break
		}
		if err != nil {
			// For some reason "ALSA lib seq.c:4176:(snd_seq_event_input_feed) poll: Interrupted system call" is spammed to stderr when running the program. I am yet to figure out the nature of this bug and how to remove it.
			// Printing these actual errors to stdout is a temporary workaround until this gets fixed.
			// The "fake alarm" ALSA errors can be discarded by redirecting stderr to /dev/null when running the command.
			fmt.Fprintf(os.Stdout, "%v\n", err)
		}
		mappings = append(mappings, mapping)
	}

	// Initialize midi device
	portmidi.Initialize()
	defer portmidi.Terminate()

	var id portmidi.DeviceID

	count := portmidi.CountDevices()
	switch count {
	case 0:
		fmt.Fprintln(os.Stderr, "midimap: no devices found, exiting")
		os.Exit(1)
	case 1:
		id = portmidi.DefaultInputDeviceID()
	default:
		// list alternatives
		alternatives, sep := "", ""
		maxID := portmidi.DeviceID(count - 1)

		for i := portmidi.DeviceID(0); i <= maxID; i++ {
			fmt.Printf("%d: %s\n", i, portmidi.Info(i))
			alternatives += fmt.Sprintf("%s%d", sep, i)
			sep = ","
		}

		// prompt for a choice
		for {
			fmt.Printf("{%s}: ", alternatives)

			_, err := fmt.Scanf("%d", &id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			if id < 0 || id > maxID {
				fmt.Printf("Please answer {%s}\n", alternatives)
			} else {
				break
			}
		}
	}

	info := portmidi.Info(id)

	fmt.Printf("input device id: %v\n", id)
	fmt.Printf("input device info: %v\n", info)

	var bufferSize int64 = 1024
	in, err := portmidi.NewInputStream(id, bufferSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer in.Close()

	// Process MIDI events
	for {
		events, err := in.Read(1024)
		if err != nil && err != portmidi.ErrSysExOverflow {
			// ErrSysExOverflow is returned sporadically when i use the PSR E333 piano keyboard
			// increasing bufferSize does NOT help.
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		for _, event := range events {
			mapMidiEventToKeyPress(mappings, event)
		}
	}
}

// evaluateComparison returns true if the comparison c evaluates to true in the context of event e in midimap-lang. Otherwise it returns false.
func evaluateComparison(c lang.Comparison, e portmidi.Event) bool {
	var data int64
	if c.LeftOperand == lang.Data1 {
		data = e.Data1
	} else if c.LeftOperand == lang.Data2 {
		data = e.Data2
	} else {
		fmt.Fprintf(os.Stderr, "This should never happen, you have found a bug")
		os.Exit(1)
	}

	switch c.Operator {
	case lang.LessThanOperator:
		return data < c.RightOperand
	case lang.LessThanOrEqualToOperator:
		return data <= c.RightOperand
	case lang.EqualToOperator:
		return data == c.RightOperand
	case lang.UnequalToOperator:
		return data != c.RightOperand
	case lang.GreaterThanOrEqualToOperator:
		return data >= c.RightOperand
	case lang.GreaterThanOperator:
		return data >= c.RightOperand
	default:
		fmt.Fprintf(os.Stderr, "This should never happen, you have found a bug")
		os.Exit(1)
		return false // Although, this line cannot be run it is required for the code to compile.
	}
}

// doesMatcherMatchEvent returns true if m matches e, otherwise it returns false.
func doesMatcherMatchEvent(m lang.Matcher, e portmidi.Event) bool {
	if m.Operator == lang.LogicalOrOperator {
		return evaluateComparison(m.LeftComparison, e) || evaluateComparison(m.RightComparison, e)
	}
	return evaluateComparison(m.LeftComparison, e) && evaluateComparison(m.RightComparison, e)
}

func mapMidiEventToKeyPress(mappings []lang.Mapping, e portmidi.Event) {
	for _, m := range mappings {
		if doesMatcherMatchEvent(m.Matcher, e) {
			press.Press(m.KeyCode)
		}
	}
}
