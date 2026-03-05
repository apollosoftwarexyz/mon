package formatting_test

import (
	"testing"
	"time"

	"github.com/apollosoftwarexyz/mon/formatting"
	"github.com/stretchr/testify/assert"
)

func TestDurationUnit_Render(t *testing.T) {
	// Smoke test that the rendering appears similar to TestDuration.
	unit := &formatting.DurationUnit{}
	assert.Equal(t, "1.2s", unit.Render(uint64((1*time.Second)+(200*time.Millisecond))))
	assert.Equal(t, "24:00:00", unit.Render(uint64(24*time.Hour)))
}

func TestDurationUnit_RenderProgress(t *testing.T) {
	// Smoke test that the rendering appears similar to TestDuration.
	unit := &formatting.DurationUnit{}
	assert.Equal(t, "1.2s / 24:00:00", unit.RenderProgress(uint64((1*time.Second)+(200*time.Millisecond)), uint64(24*time.Hour)))
}

func TestDurationUnit_RenderDurationProgress(t *testing.T) {
	// Smoke test that the rendering appears similar to TestDuration.
	unit := &formatting.DurationUnit{}
	assert.Equal(t, "1.2s / 24:00:00", unit.RenderDurationProgress((1*time.Second)+(200*time.Millisecond), 24*time.Hour))
	assert.Equal(t, "-1.2s / -24:00:00", unit.RenderDurationProgress(-((1*time.Second)+(200*time.Millisecond)), -24*time.Hour))
}

func TestDuration(t *testing.T) {
	assert.Equal(t, "0.0s", formatting.Duration(0))
	assert.Equal(t, "0.2s", formatting.Duration(200*time.Millisecond))
	assert.Equal(t, "1.0s", formatting.Duration(1*time.Second))
	assert.Equal(t, "1.2s", formatting.Duration((1*time.Second)+(200*time.Millisecond)))
	assert.Equal(t, "10.0s", formatting.Duration(10*time.Second))
	assert.Equal(t, "11.3s", formatting.Duration((11*time.Second)+(300*time.Millisecond)))
	assert.Equal(t, "59.0s", formatting.Duration(59*time.Second))

	assert.Equal(t, "1:00", formatting.Duration(1*time.Minute))
	assert.Equal(t, "59:59", formatting.Duration((59*time.Minute)+(59*time.Second)))

	assert.Equal(t, "1:00:00", formatting.Duration(1*time.Hour))
	assert.Equal(t, "5:00:00", formatting.Duration(5*time.Hour))
	assert.Equal(t, "10:00:00", formatting.Duration(10*time.Hour))
	assert.Equal(t, "24:00:00", formatting.Duration(24*time.Hour))
	assert.Equal(t, "100:00:00", formatting.Duration(100*time.Hour))

	assert.Equal(t, "-0.2s", formatting.Duration(-200*time.Millisecond))
	assert.Equal(t, "-59.0s", formatting.Duration(-59*time.Second))
	assert.Equal(t, "-100:00:00", formatting.Duration(-100*time.Hour))
}
