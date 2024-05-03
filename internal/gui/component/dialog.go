package component

import (
	"github.com/g3n/engine/gui"
)

type Dialog struct {
	*gui.Panel
	inputField   *InputField
	submitButton *gui.Button
	cancelButton *gui.Button

	onClick  func()
	onSubmit func()
	onCancel func()

	text     string
	enabled  bool
	xPos     float32
	yPos     float32
	width    float32
	height   float32
	viewType int
}

// Create a new dialog field with the given input box dimensions and starting text
func NewDialog(text string, xPos, yPos, width, height float32, viewType int, onClick func(), onSubmit func(), onCancel func()) *Dialog {
	d := &Dialog{
		Panel: gui.NewPanel(width, height),

		onClick:  onClick,
		onSubmit: onSubmit,
		onCancel: onCancel,

		text:     text,
		xPos:     xPos,
		yPos:     yPos,
		width:    width,
		height:   height,
		viewType: viewType,
	}
	d.Panel.SetPosition(xPos, yPos)

	return d
}

func (d *Dialog) Open() {
	d.inputField = NewInputField(d.width, d.height, d.text)
	d.inputField.SetWidth(d.width)
	d.inputField.SetHeight(d.height)
	d.inputField.SetUserData(d.viewType)
	d.inputField.Button.Subscribe(gui.OnClick, func(name string, ev interface{}) { d.onClick() })
	d.Add(d.inputField)

	d.submitButton = gui.NewButton("Save")
	d.submitButton.SetPosition(5+d.inputField.Position().X+d.inputField.Width(), 0)
	d.submitButton.SetUserData(d.viewType)
	d.submitButton.Subscribe(gui.OnClick, func(name string, ev interface{}) { d.onSubmit() })
	d.Add(d.submitButton)

	d.cancelButton = gui.NewButton("Cancel")
	d.cancelButton.SetPosition(5+d.submitButton.Position().X+d.submitButton.Width(), 0)
	d.cancelButton.SetUserData(d.viewType)
	d.cancelButton.Subscribe(gui.OnClick, func(name string, ev interface{}) { d.onCancel() })
	d.Add(d.cancelButton)

	d.enabled = true
	d.SetWidth(d.GetPanel().Width() + d.submitButton.GetPanel().Width() + d.cancelButton.GetPanel().Width() + 10)
}

func (d *Dialog) Update(text string) {
	if d.inputField != nil {
		d.inputField.Update(text)
	}
}

func (d *Dialog) SetEnabled(active bool) {
	d.enabled = active

	if d.inputField != nil {
		d.inputField.Activate(active)
	}
}

func (d Dialog) Enabled() bool {
	return d.enabled
}
