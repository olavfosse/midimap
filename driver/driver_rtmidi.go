// +build rtmidi

package driver

import (
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/rtmididrv"
)

func New() (midi.Driver, error) {
	return rtmididrv.New()
}
