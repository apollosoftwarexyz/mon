package formatting_test

import (
	"testing"

	"github.com/apollosoftwarexyz/mon/formatting"
	"github.com/stretchr/testify/assert"
)

const (
	kibibyte uint64 = 1024
	mebibyte        = 1024 * kibibyte
	gibibyte        = 1024 * mebibyte
	tebibyte        = 1024 * gibibyte
	pebibyte        = 1024 * tebibyte
	exbibyte        = 1024 * pebibyte
)

func TestBytesUnit_Render(t *testing.T) {
	unit := &formatting.BytesUnit{}
	assert.Equal(t, "1.001 GiB", unit.Render(gibibyte+mebibyte))
}

func TestBytesUnit_RenderProgress(t *testing.T) {
	unit := &formatting.BytesUnit{}
	assert.Equal(t, "1.000 GiB / 1.001 GiB", unit.RenderProgress(gibibyte, gibibyte+mebibyte))
}

func TestBytes(t *testing.T) {
	assert.Equal(t, "0 bytes", formatting.Bytes(0))
	assert.Equal(t, "1 byte", formatting.Bytes(1))

	assert.Equal(t, "32 bytes", formatting.Bytes(32))
	assert.Equal(t, "64 bytes", formatting.Bytes(64))
	assert.Equal(t, "702 bytes", formatting.Bytes(702))

	assert.Equal(t, "1.000 KiB", formatting.Bytes(kibibyte))
	assert.Equal(t, "1.001 KiB", formatting.Bytes(kibibyte+1))

	assert.Equal(t, "1023.999 KiB", formatting.Bytes(mebibyte-1))
	assert.Equal(t, "1023.998 KiB", formatting.Bytes(mebibyte-2))

	assert.Equal(t, "1.000 MiB", formatting.Bytes(mebibyte))
	assert.Equal(t, "1.000 MiB", formatting.Bytes(mebibyte+1))
	assert.Equal(t, "1.001 MiB", formatting.Bytes(mebibyte+kibibyte))

	assert.Equal(t, "1023.000 MiB", formatting.Bytes(gibibyte-mebibyte))
	assert.Equal(t, "1023.998 MiB", formatting.Bytes(gibibyte-(2*kibibyte)))
	assert.Equal(t, "1023.999 MiB", formatting.Bytes(gibibyte-kibibyte))
	assert.Equal(t, "1023.999 MiB", formatting.Bytes(gibibyte-1))

	assert.Equal(t, "1.000 GiB", formatting.Bytes(gibibyte))
	assert.Equal(t, "1.000 GiB", formatting.Bytes(gibibyte+kibibyte))
	assert.Equal(t, "1.001 GiB", formatting.Bytes(gibibyte+mebibyte))

	m := 3.123 * float64(mebibyte)
	assert.Equal(t, "3.123 MiB", formatting.Bytes(uint64(m)))

	g := 3.0001 * float64(gibibyte)
	assert.Equal(t, "3.000 GiB", formatting.Bytes(uint64(g)))

	g = 3.001 * float64(gibibyte)
	assert.Equal(t, "3.001 GiB", formatting.Bytes(uint64(g)))

	g = 3.123 * float64(gibibyte)
	assert.Equal(t, "3.123 GiB", formatting.Bytes(uint64(g)))

	g = 3.525 * float64(gibibyte)
	assert.Equal(t, "3.525 GiB", formatting.Bytes(uint64(g)))

	assert.Equal(t, "1.000 TiB", formatting.Bytes(tebibyte))
	assert.Equal(t, "1.000 EiB", formatting.Bytes(exbibyte))

	assert.Equal(t, "7.999 EiB", formatting.Bytes(1<<63-1))
	assert.Equal(t, "8.000 EiB", formatting.Bytes(exbibyte*8))

	assert.Equal(t, "14.000 EiB", formatting.Bytes(exbibyte*14))
	assert.Equal(t, "15.000 EiB", formatting.Bytes(exbibyte*15))
	assert.Equal(t, "15.999 EiB", formatting.Bytes(18446744073709551615))
	assert.Equal(t, "15.999 EiB", formatting.Bytes(1<<64-1))
}

func BenchmarkBytes(b *testing.B) {
	g := 3.525 * float64(gibibyte)
	ug := uint64(g)
	for b.Loop() {
		formatting.Bytes(ug)
	}

}
