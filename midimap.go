// Map midi events to simulated keypresses
package main

import (
	"fmt"
	"log"

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

	count := portmidi.CountDevices()
	if count == 0 {
		log.Fatal("No devices found, exiting")
	}

	id := portmidi.DefaultInputDeviceID()
	info := portmidi.Info(id)
	fmt.Printf("default input device id: %v\n", id)
	fmt.Printf("default input device info: %v\n", info)

	var bufferSize int64 = 1024
	in, err := portmidi.NewInputStream(id, bufferSize)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	for {
		events, err := in.Read(1024)
		if err != nil {
			// ErrSysExOverflow is returned sporadically when i use the PSR E333 piano keyboard
			// increasing bufferSize does NOT help.
			if err != portmidi.ErrSysExOverflow {
				log.Fatal(err)
			}
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
