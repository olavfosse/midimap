// Map midi events to simulated keypresses
package main

import (
	"fmt"
	"os"

	"./press"
	"github.com/micmonay/keybd_event"
	"github.com/rakyll/portmidi"
)

// the corresponding keyCode of a midi event E with data1 field D, is data1KeyCode[D]
var data1KeyCode = make(map[int64]int)

func init() {
	// left upper tom drum skin
	data1KeyCode[48] = keybd_event.VK_E

	// sharp drum skin
	data1KeyCode[38] = keybd_event.VK_F

	// floor tom drum skin
	data1KeyCode[43] = keybd_event.VK_J

	// right upper tom drum skin
	data1KeyCode[45] = keybd_event.VK_I
}

func main() {
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

	for {
		events, err := in.Read(1024)
		if err != nil && err != portmidi.ErrSysExOverflow {
			// ErrSysExOverflow is returned sporadically when i use the PSR E333 piano keyboard
			// increasing bufferSize does NOT help.
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		for _, event := range events {
			mapMidiEventToKeyPress(event)
		}
	}
}

func mapMidiEventToKeyPress(e portmidi.Event) {
	if e.Data2 == 64 { // ignore "finished hitting" event
		return
	}
	keyCode, ok := data1KeyCode[e.Data1]
	if ok {
		press.Press(keyCode)
	}
}
