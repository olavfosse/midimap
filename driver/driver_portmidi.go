// +build portmidi

package driver

import (
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/portmididrv"
)

func New() (midi.Driver, error) {
	return portmididrv.New()
}
