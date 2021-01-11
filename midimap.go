// Map midi events to simulated keypresses
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/micmonay/keybd_event"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/rtmididrv"

	"./lang"
)

func main() {
	err := dispatch(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

// dispatch parses args and dispatches the parsed args to a function corresponding to a command modifier.
// If args matches the usage/synopsis, a command modifier function is called and the return value of the handler function returned.
// Otherwise an error containing the "usage message" is returned.
func dispatch(args []string) error {
	switch {
	case len(args) == 2 && args[1] == "ports":
		return dispatchPorts()
	case len(args) == 3 && args[1] == "log":
		portNumber, err := parsePortNumber(args[2])
		if err != nil {
			return err
		}
		return dispatchLog(portNumber)
	case len(args) == 4 && args[1] == "map":
		portNumber, err := parsePortNumber(args[2])
		if err != nil {
			return err
		}
		return dispatchMap(portNumber, args[3])
	default:
		return errors.New("usage:\tmidimap map portnumber mapname\n\tmidimap ports\n\tmidimap log portnumber")
	}
}

func parsePortNumber(s string) (portNumber uint64, err error) {
	portNumber, err = strconv.ParseUint(s, 10, 0)
	if err != nil {
		err = errors.New("error: portnumber must be a valid unsigned integer")
	}
	return
}

// dispatchPorts outputs a list of the available MIDI devices/ports.
func dispatchPorts() error {
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

// dispatchLog logs incoming MIDI events from the port with the number portNumber.
func dispatchLog(portNumber uint64) error {
	drv, err := rtmididrv.New()
	if err != nil {
		return err
	}
	defer drv.Close()

	ins, err := drv.Ins()
	if err != nil {
		return err
	}

	in, err := getInByPortNumber(ins, portNumber)
	if err != nil {
		return err
	}

	err = in.Open()
	if err != nil {
		return err
	}
	defer in.Close()

	rd := reader.New(
		reader.NoLogger(),
		reader.Each(func(pos *reader.Position, msg midi.Message) {
			fmt.Println("---===---")
			fmt.Printf("%s\n", msg)
			fmt.Printf("status: %d\n", msg.Raw()[0])
			fmt.Printf("data1: %d\n", msg.Raw()[1])
			fmt.Printf("data2: %d\n", msg.Raw()[2])
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

func dispatchMap(portNumber uint64, mapName string) error {
	mappings, err := getMappingsFromMapName(mapName)
	if err != nil {
		return err
	}

	drv, err := rtmididrv.New()
	if err != nil {
		return err
	}
	defer drv.Close()

	ins, err := drv.Ins()
	if err != nil {
		return err
	}

	in, err := getInByPortNumber(ins, portNumber)
	if err != nil {
		return err
	}

	err = in.Open()
	if err != nil {
		return err
	}
	defer in.Close()

	kb, err := keybd_event.NewKeyBonding()

	if err != nil {
		if err.Error() == "permission error for /dev/uinput try cmd : sudo chmod +0666 /dev/uinput" {
			return errors.New("error: insufficient permissions to simulate keypresses\n")
		} else {
			return err
		}
	}

	rd := reader.New(
		reader.NoLogger(),
		reader.Each(func(pos *reader.Position, msg midi.Message) {
			err := mapMIDIMessageToKeyPress(kb, mappings, msg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
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

func press(kb keybd_event.KeyBonding, k int) (err error) {
	kb.SetKeys(k)
	err = kb.Launching()
	if err != nil {
		return
	}
	kb.Clear()
	return
}

// getMappingsFromMapName parses a midimap-lang file with a name of mapName and retrieves its mappings.
// If an io-error occurs, the error is returned.
// If the parser fails at parsing some mapping, it describe the problem and move on to the next mapping.
// No error is returned for parsing errors.
func getMappingsFromMapName(mapName string) (mappings []lang.Mapping, err error) {
	mapFile, err := os.Open(mapName)
	if err != nil {
		return
	}
	r := bufio.NewReader(mapFile)
	for {
		mapping, err := lang.NextMapping(r)

		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		mappings = append(mappings, mapping)
	}
	return
}

func mapMIDIMessageToKeyPress(kb keybd_event.KeyBonding, mappings []lang.Mapping, msg midi.Message) (err error) {
	for _, mapping := range mappings {
		if doesMatcherMatchMessage(mapping.Matcher, msg) {
			err = press(kb, mapping.KeyCode)
			return
		}
	}
	return
}

// doesMatcherMatchMessage returns true if m matches msg, otherwise it returns false.
func doesMatcherMatchMessage(m lang.Matcher, msg midi.Message) bool {
	if m.Operator == lang.LogicalOrOperator {
		return evaluateComparison(m.LeftComparison, msg) || evaluateComparison(m.RightComparison, msg)
	}
	return evaluateComparison(m.LeftComparison, msg) && evaluateComparison(m.RightComparison, msg)
}

// evaluateComparison returns true if the comparison c evaluates to true in the context of event e in midimap-lang. Otherwise it returns false.
func evaluateComparison(c lang.Comparison, msg midi.Message) bool {
	var data int64
	// msg.Raw()[0] is status
	// msg.Raw()[1] is data1
	// msg.Raw()[2] is data2
	if c.LeftOperand == lang.Data1 {
		data = int64(msg.Raw()[1])
	} else if c.LeftOperand == lang.Data2 {
		data = int64(msg.Raw()[2])
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

// getInByPortNumber retrieves the midi.In by number(!= index) portNumber from ins.
// If there is a In with portNumber number in ins, getInByPortNumber returns in, true where in is the in by number number.
// Otherwise, getInByPortNumber return dummyIn, false, where dummyIn is the a In with In zero value.
func getInByPortNumber(ins []midi.In, number uint64) (in midi.In, err error) {
	for _, innerIn := range ins {
		if uint64(innerIn.Number()) == number {
			in = innerIn
			return
		}

	}
	err = fmt.Errorf("error: no MIDI port by number %d", number)
	return
}

func printIns(ins []midi.In) {
	for _, in := range ins {
		fmt.Printf("%d\t%s\n", in.Number(), in.String())
	}
}
