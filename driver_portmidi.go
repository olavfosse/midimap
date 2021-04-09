// +build portmidi

package main

import (
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/portmididrv"
)

func newDriver() (midi.Driver, error) {
	return portmididrv.New()
}
