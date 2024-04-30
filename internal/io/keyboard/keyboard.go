package keyboard

// A keyboard is a type that can pipes typing data asynchronously;
// multiple keyboard events send the typed character to the Data chan,
// and another process reads what was typed in order
type Keyboard struct {
	Data chan rune
}

// New returns a new keyboard
func New() *Keyboard {
	return &Keyboard{Data: make(chan rune)}
}

