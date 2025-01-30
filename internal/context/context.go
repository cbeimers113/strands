package context

import (
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/window"

	"cbeimers113/strands/internal/config"
	"cbeimers113/strands/internal/io/keyboard"
	"cbeimers113/strands/internal/state"
)

// Collection of shared data and instances
type Context struct {
	Version     string // The version number of this build, tracked in .version at compile time
	SaveFile    string // If set, tells the game to save the state under this filename
	LoadFile    string // If set, tells the game to load the state under this filename
	CurrentSave string // If set, contains the name of the save file we loaded
	RefreshSim  bool   // If set, tells the game to create a new simulation

	App           *app.Application     // The G3N application
	Win           *window.GlfwWindow   // The window
	Scene         *core.Node           // The G3N scene containing graphical representation of simulation state
	Cam           *camera.Camera       // The G3N game camera
	Cfg           *config.Config       // The game and sim configuration
	State         *state.State         // The game state
	Keyboard      *keyboard.Keyboard   // The typing keyboard controller
	Notifications *NotificationManager // Notifications system
}

// Tell the game to save the state under a given filename
func (c *Context) SaveGame(filename string) {
	c.SaveFile = filename
}

// Tell the game to load the state under a given filename
func (c *Context) LoadGame(filename string) {
	c.LoadFile = filename
}

// Tell the game to create a fresh simulation
func (c *Context) NewGame() {
	c.RefreshSim = true
}
