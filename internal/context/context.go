package context

import (
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/window"

	"cbeimers113/strands/internal/config"
	"cbeimers113/strands/internal/state"
)

// Collection of shared data and instances
type Context struct {
	Version string             // The version number of this build, tracked in .version at compile time
	App     *app.Application   // The G3N application
	Win     *window.GlfwWindow // The window
	Scene   *core.Node         // The G3N scene containing graphical representation of simulation state
	Cam     *camera.Camera     // The G3N game camera
	Cfg     *config.Config     // The game and sim configuration
	State   *state.State       // The game state
}
