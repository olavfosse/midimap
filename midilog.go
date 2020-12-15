// Log midi events of default device to stdout
package main

import (
	"fmt"
	"log"

	"github.com/rakyll/portmidi"
)

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
			fmt.Printf("data1: %d\n", event.Data1) // which key / drum
			fmt.Printf("data2: %d\n", event.Data2) // how hard
			fmt.Println("---===---")
		}
	}
}
