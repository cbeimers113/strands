package state

import (
	"math/rand"
	"strconv"

	"github.com/g3n/engine/core"

	"cbeimers113/strands/internal/chem"
	"cbeimers113/strands/internal/config"
	"cbeimers113/strands/internal/entity"
)

// Represent the game state
type State struct {
	tps        int  // Record the number of world ticks per second
	inMenu     bool // Whether the player is in a menu and everything in the simulation is frozen, including the player
	inSpinMenu bool // Whether we're in a menu where we want the camera to spin
	paused     bool // Whether the simulation physics are paused, but the player can still interact with the simulation
	moving     bool // Whether the player is moving
	fastMove   bool // Whether the player is using fast movement
	showChems  bool // Whether to show the levels of each chemical in the simulation

	Seed       int64                               // The world's random seed value
	Rand       *rand.Rand                          // The simulation's psuedo random number generator
	Clock      *Clock                              // Keep track of in-game time
	LookingAt  entity.Entity                       // What the camera/player is looking at
	Entities   map[int]entity.Entity               // List of entities in the game world
	Quantities map[chem.ElementType]*chem.Quantity // Map of quantities for tracking various substances in the simulation
}

func New(cfg *config.Config, seed int64) *State {
	return &State{
		showChems: true,

		Seed:       seed,
		Rand:       rand.New(rand.NewSource(seed)),
		Clock:      NewClock(cfg, 9, 00, true),
		Entities:   make(map[int]entity.Entity),
		Quantities: make(map[string]*chem.Quantity),
	}
}

// region setters

// Set the number of ticks per second
func (s *State) SetTPS(tps int) {
	s.tps = tps
}

// Set the inMenu state
func (s *State) SetInMenu(inMenu bool) {
	s.inMenu = inMenu
}

// Set the inSpinMenu state
func (s *State) SetInSpinMenu(inSpinMenu bool) {
	s.inSpinMenu = inSpinMenu
	s.SetInMenu(inSpinMenu)
}

// Set the paused state
func (s *State) SetPaused(paused bool) {
	s.paused = paused
}

// Set the movement flag
func (s *State) SetMovement(moving bool) {
	s.moving = moving
}

// Set the fast movement flag
func (s *State) SetFastMovement(fast bool) {
	s.fastMove = fast
}

// Set whether we are showing chemical quantities
func (s *State) SetShowChems(showChems bool) {
	s.showChems = showChems
}

// region getters

// Get the number of ticks per second
func (s State) TPS() int {
	return s.tps
}

// Get the inMenu state
func (s State) InMenu() bool {
	return s.inMenu
}

// Get the inSpinMenu state
func (s State) InSpinMenu() bool {
	return s.inSpinMenu
}

// Get the paused state
func (s State) Paused() bool {
	return s.paused
}

// Get the movement flag
func (s State) Moving() bool {
	return s.moving
}

// Get the fast movement flag
func (s State) FastMovement() bool {
	return s.fastMove
}

// Get whether we are showing chemical quantities
func (s State) ShowChems() bool {
	return s.showChems
}

// Get the entity associated with a node, return nil if there isn't one
func (s State) EntityOf(node *core.Node) entity.Entity {
	if i, err := strconv.Atoi(node.Name()); err == nil {
		return s.Entities[i]
	}

	return nil
}
