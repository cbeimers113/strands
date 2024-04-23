package state_test

import (
	"testing"

	"github.com/g3n/engine/core"
	"github.com/stretchr/testify/assert"

	"cbeimers113/strands/internal/entity"
	"cbeimers113/strands/internal/state"
)

func Test_State(t *testing.T) {
	s := state.State{Entities: make(map[int]*entity.Entity)}
	assert.False(t, s.InMenu())
	assert.False(t, s.Paused())

	s.SetInMenu(true)
	s.SetPaused(true)
	assert.True(t, s.InMenu())
	assert.True(t, s.Paused())

	s.SetInMenu(false)
	s.SetPaused(false)
	assert.False(t, s.InMenu())
	assert.False(t, s.Paused())

	node := &core.Node{}
	node.SetName("0")
	e := entity.New(entity.NewLeaf(), entity.Plant, s.Entities)
	assert.Equal(t, e, s.EntityOf(node))

	node.SetName("1")
	assert.Nil(t, s.EntityOf(node))
}
