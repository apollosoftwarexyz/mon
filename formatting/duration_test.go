package formatting_test

import (
	"testing"
	"time"

	"github.com/apollosoftwarexyz/mon/formatting"
)

func TestDurationUnit_Render(t *testing.T) {
	// Smoke test that the rendering appears similar to TestDuration.
	unit := &formatting.DurationUnit{}
	assertEquals(t, "1.2s", unit.Render(uint64((1*time.Second)+(200*time.Millisecond))))
	assertEquals(t, "24:00:00", unit.Render(uint64(24*time.Hour)))
}

func TestDurationUnit_RenderProgress(t *testing.T) {
	// Smoke test that the rendering appears similar to TestDuration.
	unit := &formatting.DurationUnit{}
	assertEquals(t, "1.2s / 24:00:00", unit.RenderProgress(uint64((1*time.Second)+(200*time.Millisecond)), uint64(24*time.Hour)))
}

func TestDurationUnit_RenderDurationProgress(t *testing.T) {
	// Smoke test that the rendering appears similar to TestDuration.
	unit := &formatting.DurationUnit{}
	assertEquals(t, "1.2s / 24:00:00", unit.RenderDurationProgress((1*time.Second)+(200*time.Millisecond), 24*time.Hour))
	assertEquals(t, "-1.2s / -24:00:00", unit.RenderDurationProgress(-((1*time.Second)+(200*time.Millisecond)), -24*time.Hour))
}

func TestDuration(t *testing.T) {
	assertEquals(t, "0.0s", formatting.Duration(0))
	assertEquals(t, "0.2s", formatting.Duration(200*time.Millisecond))
	assertEquals(t, "1.0s", formatting.Duration(1*time.Second))
	assertEquals(t, "1.2s", formatting.Duration((1*time.Second)+(200*time.Millisecond)))
	assertEquals(t, "10.0s", formatting.Duration(10*time.Second))
	assertEquals(t, "11.3s", formatting.Duration((11*time.Second)+(300*time.Millisecond)))
	assertEquals(t, "59.0s", formatting.Duration(59*time.Second))

	assertEquals(t, "1:00", formatting.Duration(1*time.Minute))
	assertEquals(t, "59:59", formatting.Duration((59*time.Minute)+(59*time.Second)))

	assertEquals(t, "1:00:00", formatting.Duration(1*time.Hour))
	assertEquals(t, "5:00:00", formatting.Duration(5*time.Hour))
	assertEquals(t, "10:00:00", formatting.Duration(10*time.Hour))
	assertEquals(t, "24:00:00", formatting.Duration(24*time.Hour))
	assertEquals(t, "100:00:00", formatting.Duration(100*time.Hour))

	assertEquals(t, "-0.2s", formatting.Duration(-200*time.Millisecond))
	assertEquals(t, "-59.0s", formatting.Duration(-59*time.Second))
	assertEquals(t, "-100:00:00", formatting.Duration(-100*time.Hour))
}
