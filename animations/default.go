package animations

import "time"

// Default animation.
func Default() *A {
	return &A{
		Keyframes: fromUtf8String("‧⁚⁝‖⁙‖⁝⁚‧"),
		Interval:  time.Millisecond * 100,
	}
}

// Ellipsis animation.
func Ellipsis() *A {
	return &A{
		Keyframes: fromStrings([]string{"", ".", "..", "..."}),
		Interval:  time.Millisecond * 500,
	}
}
