package component

import (
	"time"

	"github.com/g3n/engine/gui"
)

type InputField struct {
	*gui.Button

	text   string
	active bool

	// Cursor
	last  time.Time
	timer int64
	freq  int64 //cursor blink frequency (in ms)
	blink bool
}

// Create a new input field with specified width, height and starting text
func NewInputField(width, height float32, text string) *InputField {
	button := gui.NewButton(text)

	i := &InputField{
		Button: button,
		text:   text,

		freq: 500,
	}
	i.SetSize(width, height)

	return i
}

// Activate enables or disables whether the input field is active / focused
func (i *InputField) Activate(active bool) {
	i.active = active
	i.SetEnabled(!active)

	if !i.active && i.blink {
		i.Label.SetText(i.text)
	}
}

func (i *InputField) Update(text string) {
	var cursor string
	i.text = text

	// Keep track of internal timer for blinking
	i.timer += time.Since(i.last).Milliseconds()
	i.last = time.Now()

	// If focused, do the blink
	if i.active {
		if i.timer >= i.freq {
			i.blink = !i.blink
			i.timer = 0
		}

		if i.blink {
			cursor = "|"
		} else {
			cursor = ""
		}

	}

	i.Label.SetText(i.text + cursor)
}
