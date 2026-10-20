package timex

import (
	"sync"
	"time"
)

// Clock provides the current time.
type Clock interface {
	Now() time.Time
}

var (
	_ Clock = SystemClock{}
	_ Clock = UTCClock{}
	_ Clock = (*ManualClock)(nil)
)

// SystemClock provides the current local time, including its monotonic reading.
type SystemClock struct{}

// Now returns the current local time.
func (SystemClock) Now() time.Time {
	return time.Now()
}

// UTCClock provides the current time in UTC.
// Its results do not contain a monotonic clock reading.
type UTCClock struct{}

// Now returns the current time in UTC.
func (UTCClock) Now() time.Time {
	return time.Now().UTC()
}

// ManualClock is a clock that can be controlled manually.
// Its zero value returns the zero time.
// Methods are safe for concurrent use.
// A ManualClock must not be copied after first use.
type ManualClock struct {
	mu   sync.RWMutex
	curr time.Time
}

// NewManualClock returns a manual clock initialized to curr.
func NewManualClock(curr time.Time) *ManualClock {
	return &ManualClock{curr: curr}
}

// Now returns the clock's current time.
func (c *ManualClock) Now() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.curr
}

// Set replaces the clock's current time.
func (c *ManualClock) Set(t time.Time) {
	c.mu.Lock()
	c.curr = t
	c.mu.Unlock()
}

// Add advances the clock by the given duration.
func (c *ManualClock) Add(d time.Duration) {
	c.mu.Lock()
	c.curr = c.curr.Add(d)
	c.mu.Unlock()
}
