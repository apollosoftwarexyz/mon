package formatting

import (
	"fmt"
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
	} else if value == 1 {
		return "1 byte"
	}

	// Compute the unit index with ilog2(value)/ilog2(1024)
	// ilog2(1024) = 10 so this has been precomputed.
	units := []string{"bytes", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}
	unitIdx := ilog2(value) / 10

	// Clamp to the maximum defined unit.
	if unitIdx >= uint64(len(units)) {
		unitIdx = uint64(len(units)) - 1
	}

	// If the index is 0, just return the value as-is with the first suffix.
	if unitIdx == 0 {
		return fmt.Sprintf("%d %s", value, units[0])
	}

	// We compute the divisor for the unit which is 1024^unitIdx. That is,
	// unitIdx=1 => we divide by 1024, unitIdx=2 => 1024*1024, etc.,
	unitDivisor := ipow(1024, unitIdx)

	// We perform integer division to get the whole unit value. This is
	// preferable to alternative methods because it preserves precision by doing
	// the division first, then converting to float64.
	unitValue := float64(value / unitDivisor)
	remainder := value % unitDivisor
	decimalRemainder := float64(remainder) / float64(ipow(1024, unitIdx)) // compute the decimal remainder
	decimalRemainderExtra := decimalRemainder * 10000                     // multiply decimal remainder to 3sf + 1 extra precision.
	decimalRemainder = float64(uint64((decimalRemainderExtra + 5) / 10))  // round half-up
	
	return fmt.Sprintf("%0.3f %s", unitValue+(decimalRemainder/1000), units[unitIdx])
}

// ilog2 computes the integer log of n by shifting whilst there are still set
// bits. The result is the counter which is returned as a [uint64].
func ilog2(n uint64) uint64 {
	var result uint64 = 0

	// We don't have a do-while loop in Go, so we just perform the first step
	// immediately.
	n = n >> 1
	for n > 0 {
		n = n >> 1
		result++
	}

	return result
}

// ipow computes the integer power of a by multiplying it together b times.
func ipow(a, b uint64) uint64 {
	var result uint64 = 1
	for b > 0 {
		result *= a
		b--
	}
	return result
}
