package formatting_test

import (
	"testing"

	"github.com/apollosoftwarexyz/mon/formatting"
)

func TestStepsUnit_Render(t *testing.T) {
	unit := &formatting.StepsUnit{}
	assertEquals(t, "0 steps", unit.Render(0))
	assertEquals(t, "1 step", unit.Render(1))
	assertEquals(t, "2 steps", unit.Render(2))
	assertEquals(t, "1024 steps", unit.Render(1024))
}

func TestStepsUnit_RenderProgress(t *testing.T) {
	unit := &formatting.StepsUnit{}
	assertEquals(t, "0 / 0 steps", unit.RenderProgress(0, 0))
	assertEquals(t, "0 / 5 steps", unit.RenderProgress(0, 5))
	assertEquals(t, "1 / 5 steps", unit.RenderProgress(1, 5))
	assertEquals(t, "2 / 5 steps", unit.RenderProgress(2, 5))
	assertEquals(t, "1024 / 5 steps", unit.RenderProgress(1024, 5))

	assertEquals(t, "0 / 1 steps", unit.RenderProgress(0, 1))
	assertEquals(t, "1 / 1 steps", unit.RenderProgress(1, 1))
	assertEquals(t, "1 / 0 steps", unit.RenderProgress(1, 0))
}

func TestSteps(t *testing.T) {
	assertEquals(t, "0 steps", formatting.Steps(0))
	assertEquals(t, "1 step", formatting.Steps(1))
	assertEquals(t, "2 steps", formatting.Steps(2))
	assertEquals(t, "1024 steps", formatting.Steps(1024))
}
