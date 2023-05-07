package timex

import "time"

// Duration is an alias for time.Duration.
type Duration = time.Duration

// StopTimer in a proper manner.
func StopTimer(t time.Timer) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
}

// MulDur multiplies duration by a float value.
func MulDur(d time.Duration, scale float64) time.Duration {
	return time.Duration(float64(d) * scale)
}
