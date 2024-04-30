package component

import "github.com/g3n/engine/gui"

type InputField struct {
	*gui.Button
	Text string
}

func NewInputField(width, height float32, text string) *InputField {
	i := &InputField{
		Button: gui.NewButton(text),
		Text:   text,
	}
	i.SetSize(width, height)

	return i
}
