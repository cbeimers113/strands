package game

import (
	"fmt"

	"github.com/g3n/engine/gui"
)

// Load the GUI for the game
func LoadGUI() {
	// Set the scene to be managed by the gui manager
	gui.Manager().Set(Scene)

	// Create and add a button to the scene
	btn := gui.NewButton("This is a Button")
	btn.SetPosition(100, 40)
	btn.SetSize(40, 40)
	btn.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("Button pressed.")
	})
	// Scene.Add(btn)
}
