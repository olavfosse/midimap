// Map midi events to simulated keypresses
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/micmonay/keybd_event"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"

	"./driver"
	"./lang"
	"./lang/mapping"
	"./lang/matcher"
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
	case (len(args) == 3 || len(args) == 4) && args[1] == "log":
		portNumber, err := parsePortNumber(args[2])
		if err != nil {
			return err
		}

		var m matcher.Matcher
		receivedMatcher := false
		if len(args) == 4 {
			m, err = matcher.Parse(args[3])
			if err != nil {
				return err
			}
			receivedMatcher = true
		}
		return dispatchLog(portNumber, receivedMatcher, m)
	case len(args) == 4 && args[1] == "map":
		portNumber, err := parsePortNumber(args[2])
		if err != nil {
			return err
		}
		return dispatchMap(portNumber, args[3])
	default:
		return errors.New("usage:\tmidimap map portnumber mapname\n\tmidimap ports\n\tmidimap log portnumber [matcher]")
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

// dispatchLog logs incoming MIDI events from the port with the number portNumber.
func dispatchLog(portNumber uint64, receivedMatcher bool, matcher matcher.Matcher) error {
	drv, err := driver.New()
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
			if !receivedMatcher || matcherMatchesMessage(matcher, msg) {
				fmt.Println("---===---")
				fmt.Printf("%s\n", msg)
				fmt.Printf("status: %d\n", msg.Raw()[0])
				fmt.Printf("data1: %d\n", msg.Raw()[1])
				fmt.Printf("data2: %d\n", msg.Raw()[2])
			}
		}),
	)

	// I don't understand how and why this snippet works, but it does.
	// I asked for an explanation in https://github.com/vipul-sharma20/midi-macro/issues/1.
	exit := make(chan string)

	go rd.ListenTo(in)

	for {
		select {
		case <-exit:
			return nil
		}
	}
}

func dispatchMap(portNumber uint64, mapName string) error {
	mappings, err := getMappingsFromMapName(mapName)
	if err != nil {
		return err
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

	// I don't understand how and why this snippet works, but it does.
	// I asked for an explanation in https://github.com/vipul-sharma20/midi-macro/issues/1.
	exit := make(chan string)

	go rd.ListenTo(in)

	for {
		select {
		case <-exit:
			return nil
		}
	}
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
func getMappingsFromMapName(mapName string) (mappings []mapping.Mapping, err error) {
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
			fmt.Fprintf(os.Stderr, "lang: %v\n", err)
		}
		mappings = append(mappings, mapping)
	}
	return
}

func mapMIDIMessageToKeyPress(kb keybd_event.KeyBonding, mappings []mapping.Mapping, msg midi.Message) (err error) {
	for _, mapping := range mappings {
		if matcherMatchesMessage(mapping.Matcher, msg) {
			// NB: We iterate through all mappings regardless of if some earlier mapping matched. This is expected behaviour.
			err = press(kb, mapping.Keycode)
			if err != nil {
				break
			}
		}
	}
	return
}

// matcherMatchesMessage reports whether msg matches m.
func matcherMatchesMessage(m matcher.Matcher, msg midi.Message) bool {
	switch m := m.(type) {
	case matcher.MatcherWithoutLogicalOperator:
		var data int64
		// msg.Raw()[0] is status
		// msg.Raw()[1] is data1
		// msg.Raw()[2] is data2
		switch m.LeftOperand {
		case matcher.Data1:
			data = int64(msg.Raw()[1])
		case matcher.Data2:
			data = int64(msg.Raw()[2])
		default:
			panic("unreachable")
		}

		switch m.Operator {
		case matcher.LessThanOperator:
			return data < m.RightOperand
		case matcher.LessThanOrEqualToOperator:
			return data <= m.RightOperand
		case matcher.EqualToOperator:
			return data == m.RightOperand
		case matcher.UnequalToOperator:
			return data != m.RightOperand
		case matcher.GreaterThanOrEqualToOperator:
			return data >= m.RightOperand
		case matcher.GreaterThanOperator:
			return data >= m.RightOperand
		default:
			panic("unreachable")
		}
	case matcher.MatcherWithLogicalOperator:
		switch m.Operator {
		case matcher.LogicalAndOperator:
			return matcherMatchesMessage(m.LeftMatcher, msg) && matcherMatchesMessage(m.RightMatcher, msg)
		case matcher.LogicalOrOperator:
			return matcherMatchesMessage(m.LeftMatcher, msg) || matcherMatchesMessage(m.RightMatcher, msg)
		default:
			panic("unreachable")
		}
	default:
		panic("unreachable")
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
