package context

import (
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"

	"cbeimers113/strands/internal/gui/color"
)

// The time limit in ms before a notification disappears
const LIMIT = 7000

type notification struct {
	message string
	timer   float32
}

type NotificationManager struct {
	getScene      func() *core.Node
	app           *app.Application
	notifications []notification
	noteObjects   []*gui.Label
}

func NewNotificationManager(getScene func() *core.Node, app *app.Application) *NotificationManager {
	return &NotificationManager{
		getScene:      getScene,
		app:           app,
		notifications: make([]notification, 0),
		noteObjects:   make([]*gui.Label, 0),
	}
}

func (n *NotificationManager) Push(message string) {
	n.notifications = append(n.notifications, notification{
		message: message,
		timer:   LIMIT,
	})
}

func (n *NotificationManager) Update(delta float32) {
	notes := make([]notification, 0)
	for _, note := range n.notifications {
		note.timer -= delta

		if note.timer <= 0 {
			continue
		}

		notes = append(notes, note)
	}

	n.notifications = notes
}

func (n *NotificationManager) Render() {
	scene := n.getScene()

	// Remove existing notification objects
	for _, obj := range n.noteObjects {
		obj.DisposeChildren(true)
		obj.Dispose()
		scene.Remove(obj)
	}

	n.noteObjects = make([]*gui.Label, 0)
	width, _ := n.app.GetSize()

	// Add labels to screen for notifications
	for i, note := range n.notifications {
		alpha := color.Opaque.A * min(note.timer, LIMIT/2) / (LIMIT / 2)
		fgCol := &math32.Color4{R: 0.0, G: 0.0, B: 0.0, A: alpha}
		bgCol := &math32.Color4{R: 1.0, G: 1.0, B: 1.0, A: alpha}
		obj := gui.NewLabel(note.message)
		y := float32(i)*obj.Height() + 10*float32(i+1)
		obj.SetPosition(float32(width)-obj.Width()-15, y)
		obj.SetPaddings(0, 5, 0, 5)
		obj.SetColor4(fgCol)
		obj.SetBgColor4(bgCol)
		scene.Add(obj)
		n.noteObjects = append(n.noteObjects, obj)
	}
}
