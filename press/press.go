// Simulate pressing various keyboard keys
package press

import (
	"log"
	"runtime"
	"time"

	"github.com/micmonay/keybd_event"
)

var kb keybd_event.KeyBonding

func init() {
	var err error
	kb, err = keybd_event.NewKeyBonding()
	if err != nil {
		log.Fatal(err)
	}

	// For linux, it is important to wait 2 seconds, see https://github.com/micmonay/keybd_event/issues/25 for why
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}
}

func pressKey(k int) {
	kb.SetKeys(k)
	kb.Launching()
	kb.Clear()
}

func E() { pressKey(keybd_event.VK_E) }
func F() { pressKey(keybd_event.VK_F) }
func J() { pressKey(keybd_event.VK_J) }
func I() { pressKey(keybd_event.VK_I) }
