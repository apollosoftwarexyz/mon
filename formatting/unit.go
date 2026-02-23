package formatting

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
