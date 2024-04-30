package keyboard

import (
	"sync"

	"github.com/g3n/engine/window"
)

// A keyboard is a type that can pipes typing data asynchronously;
// multiple keyboard events send the typed character to the Data chan,
// and another process reads what was typed in order
type Keyboard struct {
	sync.Mutex
	enabled bool
	shift   bool
	ctrl    bool
	data    chan window.Key
	buffer  string
}

// New returns a new keyboard
func New() *Keyboard {
	return &Keyboard{}
}

// listen for key events and decode them
func (k *Keyboard) listen() {
	l := true
	k.data = make(chan window.Key)
	defer close(k.data)

	for l {
		select {
		case ev := <-k.data:
			k.Lock()

			// Check for backspace signal, otherwise try to decode the event
			if ev != window.KeyBackspace {
				k.buffer += decode(ev, k.shift)
			} else if len(k.buffer) > 0 {
				if k.ctrl {
					k.buffer = ""
				} else {
					k.buffer = k.buffer[:len(k.buffer)-1]
				}
			}

			k.Unlock()
		default:
		}

		// Stop listening if the keyboard was turned off
		k.Lock()
		l = k.enabled
		k.Unlock()
	}
}

// Input puts a key event into the data pipe
func (k *Keyboard) Input(keyEv window.Key) {
	if k.GetEnabled() {
		k.data <- keyEv
	}
}

// Read returns the buffer of typed data
func (k *Keyboard) Read() string {
	k.Lock()
	defer k.Unlock()
	return k.buffer
}

// Clear clears the keyboard's buffer
func (k *Keyboard) Clear() {
	k.Lock()
	defer k.Unlock()
	k.buffer = ""
}

// Enable or disable the keyboard
func (k *Keyboard) Enable(enabled bool) {
	k.Lock()
	defer k.Unlock()

	k.enabled = enabled
	if k.enabled {
		go k.listen()
	}
}

// Get the enable status of the keyboard
func (k *Keyboard) GetEnabled() bool {
	k.Lock()
	defer k.Unlock()
	return k.enabled
}

// Shift sets the shift status of the keyboard
func (k *Keyboard) Shift(shift bool) {
	k.Lock()
	defer k.Unlock()
	k.shift = shift
}

// Ctrl sets the ctrl status of the keyboard
func (k *Keyboard) Ctrl(ctrl bool) {
	k.Lock()
	defer k.Unlock()
	k.ctrl = ctrl
}
