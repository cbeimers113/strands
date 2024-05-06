package component

import (
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
)

type Popup struct {
	*gui.Panel
	promptLabel  *gui.Label
	submitButton *gui.Button
	cancelButton *gui.Button

	onSubmit func()
	onCancel func()

	text       string // Text for the prompt label
	submitText string // Text for the submit button

	// Popup position and dimensions on screen
	xPos   float32
	yPos   float32
	width  float32
	height float32

	viewType int
	enabled  bool
}

func NewPopup(text, submitText string, xPos, yPos, width, height float32, viewType int, onSubmit func(), onCancel func()) *Popup {
	p := &Popup{
		Panel: gui.NewPanel(width, height),

		onSubmit: onSubmit,
		onCancel: onCancel,

		text:       text,
		submitText: submitText,

		xPos:   xPos,
		yPos:   yPos,
		width:  width,
		height: height,

		viewType: viewType,
	}
	p.Panel.SetPosition(xPos, yPos)
	p.Panel.SetColor(&math32.Color{R: 0.25, G: 0.25, B: 0.25})
	p.Panel.SetBorders(5, 5, 5, 5)
	p.Panel.SetBordersColor(&math32.Color{R: 0.5, G: 0.5, B: 0.5})

	return p
}

func (p *Popup) Open() {
	p.promptLabel = gui.NewLabel(p.text)
	p.promptLabel.SetUserData(p.viewType)
	p.Add(p.promptLabel)

	if p.Width() <= p.promptLabel.Width()+20 {
		p.SetWidth(p.promptLabel.Width() + 20)
	}
	p.promptLabel.SetPosition((p.Width()-p.promptLabel.Width())/2-7.5, 5)

	p.submitButton = gui.NewButton(p.submitText)
	p.submitButton.SetPosition(p.Width()/2-p.submitButton.Width()-12.5, p.Height()-p.submitButton.Height()-15)
	p.submitButton.SetUserData(p.viewType)
	p.submitButton.Subscribe(gui.OnClick, func(name string, ev interface{}) { p.onSubmit() })
	p.Add(p.submitButton)

	p.cancelButton = gui.NewButton("Cancel")
	p.cancelButton.SetPosition(p.Width()/2-7.5, p.Height()-p.cancelButton.Height()-15)
	p.cancelButton.SetUserData(p.viewType)
	p.cancelButton.Subscribe(gui.OnClick, func(name string, ev interface{}) { p.onCancel() })
	p.Add(p.cancelButton)

	p.enabled = true
}

func (p *Popup) SetEnabled(active bool) {
	p.enabled = active
}

func (p Popup) Enabled() bool {
	return p.enabled
}

func (p *Popup) Dispose() {
	p.DisposeChildren(true)
	p.Panel.Dispose()
}
