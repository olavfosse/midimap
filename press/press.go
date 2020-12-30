// Simulate pressing various keyboard keys
package press

import (
	"log"

	"github.com/micmonay/keybd_event"
)

var kb keybd_event.KeyBonding

func init() {
	var err error
	kb, err = keybd_event.NewKeyBonding()
	if err != nil {
		log.Fatal(err)
	}
}

func Press(k int) {
	kb.SetKeys(k)
	kb.Launching()
	kb.Clear()
}
