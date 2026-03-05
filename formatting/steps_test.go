package formatting_test

import (
	"testing"

	"github.com/apollosoftwarexyz/mon/formatting"
	"github.com/stretchr/testify/assert"
)

func TestStepsUnit_Render(t *testing.T) {
	unit := &formatting.StepsUnit{}
	assert.Equal(t, "0 steps", unit.Render(0))
	assert.Equal(t, "1 step", unit.Render(1))
	assert.Equal(t, "2 steps", unit.Render(2))
	assert.Equal(t, "1024 steps", unit.Render(1024))
}

func TestStepsUnit_RenderProgress(t *testing.T) {
	unit := &formatting.StepsUnit{}
	assert.Equal(t, "0 / 0 steps", unit.RenderProgress(0, 0))
	assert.Equal(t, "0 / 5 steps", unit.RenderProgress(0, 5))
	assert.Equal(t, "1 / 5 steps", unit.RenderProgress(1, 5))
	assert.Equal(t, "2 / 5 steps", unit.RenderProgress(2, 5))
	assert.Equal(t, "1024 / 5 steps", unit.RenderProgress(1024, 5))

	assert.Equal(t, "0 / 1 steps", unit.RenderProgress(0, 1))
	assert.Equal(t, "1 / 1 steps", unit.RenderProgress(1, 1))
	assert.Equal(t, "1 / 0 steps", unit.RenderProgress(1, 0))
}

func TestSteps(t *testing.T) {
	assert.Equal(t, "0 steps", formatting.Steps(0))
	assert.Equal(t, "1 step", formatting.Steps(1))
	assert.Equal(t, "2 steps", formatting.Steps(2))
	assert.Equal(t, "1024 steps", formatting.Steps(1024))
}
