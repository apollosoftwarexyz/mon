package formatting

import (
	"fmt"
	"strings"
	"time"
)

// DurationUnit renders a [uint64] number of nanoseconds (the base unit of
// [time.Duration]) as a human-readable duration.
//
// This unit also exposes the [DurationUnit.RenderDurationProgress] method that
// accepts [time.Duration] instead of [uint64]. This allows rendering negative
// progress if required. For the individual unit version, use [Duration].
type DurationUnit struct{}

func (*DurationUnit) Render(value uint64) string {
	return Duration(time.Duration(value))
}

func (*DurationUnit) RenderProgress(current uint64, total uint64) string {
	return fmt.Sprintf("%s / %s", Duration(time.Duration(current)), Duration(time.Duration(total)))
}

func (*DurationUnit) RenderDurationProgress(current time.Duration, total time.Duration) string {
	return fmt.Sprintf("%s / %s", Duration(current), Duration(total))
}

// Duration formats the given value as a period of time, automatically including
// or excluding precision as appropriate.
func Duration(d time.Duration) string {
	var s strings.Builder

	// Handle negative durations by writing the minus sign and then treating the
	// value as absolute.
	if d < 0 {
		s.WriteRune('-')
		d = d.Abs()
	}

	if hours := time.Duration(d.Hours()); hours > 0 {
		s.WriteString(fmt.Sprintf("%d:", hours))
	}

	if minutes := time.Duration(d.Minutes()); minutes > 0 {
		minutesFmt := "%d:"
		if d.Minutes() >= 60 {
			minutesFmt = "%02d:"
		}

		s.WriteString(fmt.Sprintf(minutesFmt, minutes%60))
	}

	if seconds := time.Duration(d.Seconds()); seconds >= 10 {
		s.WriteString(fmt.Sprintf("%02d", seconds%60))
	} else {
		s.WriteString(fmt.Sprintf("%d", seconds%60))
	}

	if d.Seconds() < 60 {
		s.WriteString(fmt.Sprintf(".%01ds", (time.Duration(d.Milliseconds())%1000)/100))
	}

	return s.String()
}
