package formatting

import "fmt"

// StepsUnit renders discrete numbers of steps.
//
// This is a generalized unit that can be used as a reasonable default for
// monitoring progress.
type StepsUnit struct{}

func (*StepsUnit) Render(value uint64) string {
	return Steps(value)
}

func (*StepsUnit) RenderProgress(current uint64, total uint64) string {
	return fmt.Sprintf("%d / %d steps", current, total)
}

// Steps formats the given value as a discrete number of steps.
//
// If value is equal to one, the hardcoded string "1 step" is returned instead
// for readability.
func Steps(value uint64) string {
	if value == 1 {
		return "1 step"
	}

	return fmt.Sprintf("%d steps", value)
}
