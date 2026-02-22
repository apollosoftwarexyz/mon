package animations

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type Keyframes interface {
	// Length of the animation, in frames.
	Length() int

	// Frame returns the animation frame at the current index.
	//
	// If the frames are not an even multiple of the [Keyframes.Width], the last
	// frame is padded with blanks accordingly.
	//
	// If the frame index, i, is greater than the length of the animation, an
	// empty frame of the animation's frame width will be returned instead.
	Frame(i int) string
}

type keyframes struct {
	frames      []string // frames is the list of individual keyframes.
	widestFrame int      // widestFrame is the width, in characters of the widest frame.
}

func (k *keyframes) Length() int { return len(k.frames) }
func (k *keyframes) Frame(i int) string {
	if i >= len(k.frames) {
		return strings.Repeat(" ", k.widestFrame)
	}

	return fmt.Sprintf("%-"+strconv.Itoa(k.widestFrame)+"s", k.frames[i])
}

// fromStrings converts the given slice of strings into individual frames.
func fromStrings(s []string) Keyframes {

	widestFrame := 0

	for _, frame := range s {
		l := len(frame)
		if l > widestFrame {
			widestFrame = l
		}
	}

	return &keyframes{
		frames:      s,
		widestFrame: widestFrame,
	}

}

// fromUtf8String converts the given string of UTF-8 runes into keyframes where
// each individual UTF-8 rune is a keyframe.
func fromUtf8String(s string) Keyframes {

	frames := make([]string, 0, utf8.RuneCountInString(s))

	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		frames = append(frames, string(r))
		s = s[size:]
	}

	return &keyframes{
		frames:      frames,
		widestFrame: 1,
	}

}

// A is an animation made up of [Keyframes] - a string of keyframes where each
// n characters (where n is the width) is a frame to be displayed for the given
// interval.
type A struct {
	Keyframes
	Interval time.Duration // Interval between frames.
}

// Duration of the animation to play one complete cycle, as a [time.Duration].
func (a *A) Duration() time.Duration {
	return a.Interval * time.Duration(a.Length())
}

func (a *A) Frame(t time.Duration) string {
	// Wrap the duration to the time of the animation - such that after the
	// animation is completed, it loops.
	wrappedDuration := t % a.Duration()

	// Then compute the frame by dividing the elapsed duration by the interval
	// each frame is displayed for.
	return a.Keyframes.Frame(int(wrappedDuration / a.Interval))
}
