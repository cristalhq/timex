package timex

import "time"

// Duration is an alias for time.Duration.
type Duration = time.Duration

// Clock represents interface on time. Useful in tests.
type Clock interface {
	Now() time.Time
}

var (
	RealClock    = &systemClock{}
	RealUTCClock = &systemUTCClock{}

	_ Clock = RealClock
	_ Clock = RealUTCClock
	_ Clock = &ManualClock{}
)

type systemClock struct{}

func (c *systemClock) Now() time.Time { return time.Now() }

type systemUTCClock struct{}

func (c *systemUTCClock) Now() time.Time { return time.Now().UTC() }

// ManualClock represents clock than can be controlled manually.
type ManualClock struct {
	curr time.Time
}

// NewManualClock returns manual Clock with a given time.
func NewManualClock(curr time.Time) *ManualClock {
	return &ManualClock{curr: curr}
}

func (c *ManualClock) Now() time.Time      { return c.curr }
func (c *ManualClock) Set(t time.Time)     { c.curr = t }
func (c *ManualClock) Add(d time.Duration) { c.curr = c.curr.Add(d) }

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
