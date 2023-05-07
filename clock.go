package timex

import "time"

// Clock represents interface on time. Useful in tests.
type Clock interface {
	Now() time.Time
}

var RealClock = &systemClock{}

var RealUTCClock = &systemUTCClock{}

var (
	_ Clock = RealClock
	_ Clock = RealUTCClock
	_ Clock = &ManualClock{}
)

type systemClock struct{}

type systemUTCClock struct{}

func (c *systemClock) Now() time.Time    { return time.Now() }
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
