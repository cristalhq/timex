package timex

import (
	"errors"
	"time"
)

// Zone is an alias for [time.Location].
type Zone = time.Location

// Dur is an alias for [time.Duration].
type Dur = time.Duration

const (
	Week        = 7 * Day
	Day         = 24 * Hour
	HalfDay     = 12 * Hour
	Hour        = time.Hour
	Minute      = time.Minute
	Second      = time.Second
	Millisecond = time.Millisecond
	Microsecond = time.Microsecond
	Nanosecond  = time.Nanosecond
)

// Now returns currect time.
func Now() time.Time {
	return time.Now()
}

// Compare function for two [time.Time] instances.
func Compare(a, b time.Time) int {
	return a.Compare(b)
}

// Parse value using the first matching layout.
// At least one layout is required.
func Parse(value string, layouts ...string) (time.Time, error) {
	return ParseInZone(time.UTC, value, layouts...)
}

// ParseInZone value using the first matching layout.
// At least one layout is required.
func ParseInZone(zone *Zone, value string, layouts ...string) (time.Time, error) {
	if len(layouts) == 0 {
		return time.Time{}, errors.New("timex: no layouts provided")
	}

	var firstErr error
	for _, layout := range layouts {
		t, err := time.ParseInLocation(layout, value, zone)
		if err == nil {
			return t, nil
		}
		if firstErr == nil {
			firstErr = err
		}
	}
	return time.Time{}, firstErr
}

// Max returns the latest nonzero time.
// It returns the zero time if there are no nonzero values.
func Max(values ...time.Time) time.Time {
	var res time.Time
	for _, v := range values {
		if !v.IsZero() && (res.IsZero() || v.After(res)) {
			res = v
		}
	}
	return res
}

// Min returns the earliest nonzero time.
// It returns the zero time if there are no nonzero values.
func Min(values ...time.Time) time.Time {
	var res time.Time
	for _, v := range values {
		if !v.IsZero() && (res.IsZero() || v.Before(res)) {
			res = v
		}
	}
	return res
}

// DropDate replaces the date with January 1, year 0.
// It preserves the clock components and location.
func DropDate(t time.Time) time.Time {
	return time.Date(
		0, time.January, 1,
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond(),
		t.Location(),
	)
}

// DropTime replaces the clock components with midnight.
// It preserves the date and location.
func DropTime(t time.Time) time.Time {
	return time.Date(
		t.Year(), t.Month(), t.Day(),
		0, 0, 0, 0,
		t.Location(),
	)
}

func MustZone(zone string) *Zone {
	loc, err := time.LoadLocation(zone)
	if err != nil {
		panic(err)
	}
	return loc
}
