package formatting

import (
	"fmt"
	"math"
)

// BytesUnit renders discrete numbers of bytes using [Bytes].
//
// For clarity, the units are not abbreviated in the progress formatting
// variant.
type BytesUnit struct{}

func (*BytesUnit) Render(value uint64) string {
	return Bytes(value)
}

func (b *BytesUnit) RenderProgress(current uint64, total uint64) string {
	return fmt.Sprintf("%s / %s", b.Render(current), b.Render(total))
}

// Bytes formats the given value as a number of bytes.
//
// Values up to 1,024 bytes are formatted as "<value> bytes". Larger values
// are formatted according to their nearest SI unit.
func Bytes(value uint64) string {
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
