package state

import (
	"strconv"

	"github.com/g3n/engine/core"

	"cbeimers113/strands/internal/chem"
	"cbeimers113/strands/internal/entity"
)

// Represent the game state
type State struct {
	inMenu       bool // Whether the player is in a menu and everything in the simulation is frozen, including the player
	paused       bool // Whether the simulation physics are paused, but the player can still interact with the simulation
	tps          int  // Record the number of world ticks per second
	menuCooldown int  // Cooldown in ms after which the player can update after closing a menu

	LookingAt  *entity.Entity                     // What the camera/player is looking at
	Entities   map[int]*entity.Entity             // List of entities in the game world
	Quantities map[chem.ElementType]chem.Quantity // Map of quantities of various substances in the simulation
}

func New() *State {
	return &State{
		Entities:   make(map[int]*entity.Entity),
		Quantities: make(map[string]chem.Quantity),
	}
}

// Set the inMenu state
func (s *State) SetInMenu(inMenu bool) {
	s.inMenu = inMenu
	s.SetPaused(inMenu)
}

// Set the paused state
func (s *State) SetPaused(paused bool) {
	s.paused = paused
}

// Set the number of ticks per second
func (s *State) SetTPS(tps int) {
	s.tps = tps
}

// Change the value by the menu cooldown
func (s *State) ChangeMenuCooldown(delta int) {
	s.menuCooldown += delta
	if s.menuCooldown < 0 {
		s.menuCooldown = 0
	}
}

// Get the inMenu state
func (s State) InMenu() bool {
	return s.inMenu
}

// Get the paused state
func (s State) Paused() bool {
	return s.paused
}

// Get the number of ticks per second
func (s State) TPS() int {
	return s.tps
}

// Get the menu cooldown value
func (s State) MenuCooldown() int {
	return s.menuCooldown
}

// Get the entity associated with a node, return nil if there isn't one
func (s State) EntityOf(node *core.Node) *entity.Entity {
	if i, err := strconv.Atoi(node.Name()); err == nil && i < len(s.Entities) {
		return s.Entities[i]
	}

	return nil
}
