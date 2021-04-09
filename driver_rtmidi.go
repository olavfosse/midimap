// +build rtmidi

package main

import (
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/rtmididrv"
)

func newDriver() (midi.Driver, error) {
	return rtmididrv.New()
}
