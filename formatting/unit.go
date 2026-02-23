package formatting

import (
	"fmt"
	"math"
)

// Unit is a way of rendering discrete values.
type Unit interface {
	// Render the value according to the unit.
	//
	// For example, if the unit is bytes and the value is 1,024, this might
	// render the value as "1 KiB".
	Render(value uint64) string

	// RenderProgress according to the unit.
	//
	// This function is specialized to allow for rendering progress in a more
	// optimized way (e.g., "1 / 3 steps" instead of "1 step / 3 steps").
	RenderProgress(current uint64, total uint64) string
}

// StepsUnit renders as discrete numbers of steps.
type StepsUnit struct{}

func (*StepsUnit) Render(value uint64) string {
	if value == 1 {
		return "1 step"
	}

	return fmt.Sprintf("%d steps", value)
}

func (*StepsUnit) RenderProgress(current uint64, total uint64) string {
	return fmt.Sprintf("%d / %d steps", current, total)
}

// BytesUnit renders as discrete numbers of bytes using human-friendly SI
// formatting where appropriate.
type BytesUnit struct{}

func (*BytesUnit) Render(value uint64) string {
	if value == 0 {
		return "0 bytes"
	}

	units := []string{"bytes", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}
	unitIdx := int(math.Log2(float64(value)) / math.Log2(1024))
	if unitIdx >= len(units) {
		unitIdx = len(units) - 1
	}

	if unitIdx == 0 {
		return fmt.Sprintf("%d %s", value, units[unitIdx])
	}

	unitValue := float64(value) / math.Pow(1024, float64(unitIdx))
	return fmt.Sprintf("%0.3f %s", unitValue, units[unitIdx])
}

func (b *BytesUnit) RenderProgress(current uint64, total uint64) string {
	return fmt.Sprintf("%s / %s", b.Render(current), b.Render(total))
}
