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
	_, height := n.app.GetSize()

	// Add labels to screen for notifications
	for i, note := range n.notifications {
		j := float32(len(n.notifications) - i)
		alpha := color.Translucent.A * min(note.timer, LIMIT/2) / (LIMIT / 2)
		bgCol := &math32.Color4{R: 1.0, G: 1.0, B: 1.0, A: alpha}
		obj := gui.NewLabel(note.message)
		y := float32(height) - j*(obj.Height()+10)
		obj.SetPosition(5, y)
		obj.SetPaddings(0, 5, 0, 5)
		obj.SetColor(color.Black)
		obj.SetBgColor4(bgCol)
		scene.Add(obj)
		n.noteObjects = append(n.noteObjects, obj)
	}
}
