package component

import (
	"time"

	"github.com/g3n/engine/gui"
)

type InputField struct {
	*gui.Button
	Text string

	last   time.Time
	timer  int64
	freq   int64
	blink  bool
	active bool
}

// Create a new input field with specified width, height, text and cursor blink frequency (in ms)
func NewInputField(width, height float32, text string, freq int64) *InputField {
	button := gui.NewButton(text)

	i := &InputField{
		Button: button,
		Text:   text,
		freq:   freq,
	}
	i.SetSize(width, height)

	return i
}

// Activate enables or disables whether the input field is active / focused
func (i *InputField) Activate(active bool) {
	i.active = active
}

func (i *InputField) Update(text string) {
	var cursor string
	i.Text = text

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

	i.Label.SetText(text + cursor)
}
