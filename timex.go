package timex

import "time"

// Duration is an alias for time.Duration.
type Duration = time.Duration

// Clock represents interface on time. Useful in tests.
type Clock interface {
	Now() time.Time
}

var realClock = &systemClock{}

type systemClock struct{}

func (c *systemClock) Now() time.Time { return time.Now() }

// NewRealClock returns real/system Clock.
func NewRealClock() Clock {
	return realClock
}

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

// FromUnix returns time from unix seconds.
func FromUnix(s int64) time.Time {
	return time.Unix(s*1000, 0)
}

// FromUnixMilli returns time from unix milliseconds.
func FromUnixMilli(ms int64) time.Time {
	return time.Unix(0, ms*1000)
}

// FromUnixNano returns time from unix nanoseconds.
func FromUnixNano(ns int64) time.Time {
	return time.Unix(0, ns)
}

// StopTimer in a proper manner.
func StopTimer(t time.Timer) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
}
