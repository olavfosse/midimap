package mapcm

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fossegrim/midimap/cm/helper"
	"github.com/fossegrim/midimap/driver"
	"github.com/fossegrim/midimap/lang"
	"github.com/fossegrim/midimap/lang/mapping"
	"github.com/fossegrim/midimap/usage"
	"github.com/micmonay/keybd_event"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
)

// MapCM corresponds to the map command modifier. args corresponds to
// the list of arguments listed on the command line after the map
// command modifier.
//
// For documentation about the map command modifier itself, consult
// the (to be written) manual.
// TODO: Write midimap(1) and/or midimap-map(1)
func MapCM(args []string) error {
	portNumber, mapName, err := parseArgs(args)
	if err != nil {
		return err
	}

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

	in, err := helper.GetInByPortNumber(ins, portNumber)
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
			return errors.New("error: insufficient permissions to simulate keypresses")
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
	// TODO: Read chapter 8 in gopl, which hopefully explains this
	exit := make(chan string)

	go rd.ListenTo(in)

	for {
		select {
		case <-exit:
			return nil
		}
	}
}

func parseArgs(args []string) (portNumber uint64, mapName string, err error) {
	if len(args) != 2 {
		err = usage.Usage()
		return
	}
	mapName = args[1]
	portNumber, err = helper.ParsePortNumber(args[0])
	return
}

func mapMIDIMessageToKeyPress(kb keybd_event.KeyBonding, mappings []mapping.Mapping, msg midi.Message) (err error) {
	for _, mapping := range mappings {
		if helper.MatcherMatchesMessage(mapping.Matcher, msg) {
			// NB: We iterate through all mappings regardless of if some earlier mapping matched. This is expected behaviour.
			err = press(kb, mapping.Keycode)
			if err != nil {
				break
			}
		}
	}
	return
}

// getMappingsFromMapName parses a midimap-lang file with a name of mapName and retrieves its mappings.
// If an io-error occurs, the error is returned.
// If the parser fails at parsing some mapping, it describe the problem and moves on to the next mapping.
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

// press simulates pressing k on kb.
func press(kb keybd_event.KeyBonding, k int) (err error) {
	kb.SetKeys(k)
	err = kb.Launching()
	if err != nil {
		return
	}
	kb.Clear()
	return
}
