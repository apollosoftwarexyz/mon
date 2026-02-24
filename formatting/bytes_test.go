package formatting_test

import (
	"testing"

	"github.com/apollosoftwarexyz/mon/formatting"
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
	assertEquals(t, "1.001 GiB", unit.Render(gibibyte+mebibyte))
}

func TestBytesUnit_RenderProgress(t *testing.T) {
	unit := &formatting.BytesUnit{}
	assertEquals(t, "1.000 GiB / 1.001 GiB", unit.RenderProgress(gibibyte, gibibyte+mebibyte))
}

func TestBytes(t *testing.T) {
	assertEquals(t, "0 bytes", formatting.Bytes(0))
	assertEquals(t, "1 byte", formatting.Bytes(1))

	assertEquals(t, "32 bytes", formatting.Bytes(32))
	assertEquals(t, "64 bytes", formatting.Bytes(64))
	assertEquals(t, "702 bytes", formatting.Bytes(702))

	assertEquals(t, "1.000 KiB", formatting.Bytes(kibibyte))
	assertEquals(t, "1.001 KiB", formatting.Bytes(kibibyte+1))

	assertEquals(t, "1.000 MiB", formatting.Bytes(mebibyte))
	assertEquals(t, "1.000 MiB", formatting.Bytes(mebibyte+1))
	assertEquals(t, "1.001 MiB", formatting.Bytes(mebibyte+kibibyte))

	assertEquals(t, "1.000 GiB", formatting.Bytes(gibibyte))
	assertEquals(t, "1.000 GiB", formatting.Bytes(gibibyte+kibibyte))
	assertEquals(t, "1.001 GiB", formatting.Bytes(gibibyte+mebibyte))

	m := 3.123 * float64(mebibyte)
	assertEquals(t, "3.123 MiB", formatting.Bytes(uint64(m)))

	g := 3.0001 * float64(gibibyte)
	assertEquals(t, "3.000 GiB", formatting.Bytes(uint64(g)))

	g = 3.001 * float64(gibibyte)
	assertEquals(t, "3.001 GiB", formatting.Bytes(uint64(g)))

	g = 3.123 * float64(gibibyte)
	assertEquals(t, "3.123 GiB", formatting.Bytes(uint64(g)))

	g = 3.525 * float64(gibibyte)
	assertEquals(t, "3.525 GiB", formatting.Bytes(uint64(g)))

	assertEquals(t, "1.000 TiB", formatting.Bytes(tebibyte))
	assertEquals(t, "1.000 EiB", formatting.Bytes(exbibyte))

	assertEquals(t, "14.000 EiB", formatting.Bytes(exbibyte*14))
}

func BenchmarkBytes(b *testing.B) {
	g := 3.525 * float64(gibibyte)
	ug := uint64(g)
	for b.Loop() {
		formatting.Bytes(ug)
	}

}
