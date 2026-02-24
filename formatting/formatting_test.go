package formatting_test

import "testing"

func assertEquals[T comparable](t *testing.T, expected, actual T) {
	t.Helper()

	if expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}
