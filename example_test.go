package timex_test

import (
	"time"

	"github.com/cristalhq/timex"
)

func ExampleMulDur() {
	window := time.Minute

	window = timex.MulDur(window, 1.5)
	if window.Seconds() != 90 {
		panic("want 90 have " + window.String())
	}

	window = time.Minute
	window = time.Duration(float64(window) * 1.5)
	if window.Seconds() != 90 {
		panic("want 90 have " + window.String())
	}

	// Output:
}
