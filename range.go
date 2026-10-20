package timex

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Range represents a half-open interval [Start, Finish).
// Equal endpoints form a valid empty range.
// The zero value is a valid empty range.
type Range struct {
	Start  time.Time
	Finish time.Time
}

// NewRange returns an error if finish is before start.
func NewRange(start, finish time.Time) (Range, error) {
	r := Range{Start: start, Finish: finish}
	if !r.IsValid() {
		return Range{}, fmt.Errorf("timex: finish before start")
	}
	return r, nil
}

func ParseRange(s string) (Range, error) {
	return ParseRangeInZone(s, time.UTC)
}

func ParseRangeInZone(s string, zone *time.Location) (Range, error) {
	if zone == nil {
		return Range{}, fmt.Errorf("timex: nil zone")
	}

	s = strings.TrimSpace(s)

	offset := 0
	if hasRangeDate(s) {
		offset = 10
	}

	idx := strings.IndexByte(s[offset:], '-')
	if idx < 0 {
		return Range{}, fmt.Errorf("timex: missing range separator in %q", s)
	}
	idx += offset

	startText := strings.TrimSpace(s[:idx])
	finishText := strings.TrimSpace(s[idx+1:])

	start, err := parseRangeEndpointInZone(startText, zone)
	if err != nil {
		return Range{}, fmt.Errorf("timex: invalid range start: %w", err)
	}

	// An omitted end date inherits the start date.
	if !hasRangeDate(finishText) && hasRangeDate(startText) {
		finishText = startText[:10] + " " + finishText
	}

	finish, err := parseRangeEndpointInZone(finishText, zone)
	if err != nil {
		return Range{}, fmt.Errorf("timex: invalid range end: %w", err)
	}

	return NewRange(start, finish)
}

func parseRangeEndpointInZone(s string, zone *time.Location) (time.Time, error) {
	date := "0001-01-01"

	if hasRangeDate(s) {
		date = s[:10]
		s = strings.TrimSpace(s[11:])
	}

	var clockLayout string
	switch strings.Count(s, ":") {
	case 0:
		clockLayout = "15"
	case 1:
		clockLayout = "15:04"
	case 2:
		clockLayout = "15:04:05"
	default:
		return time.Time{}, fmt.Errorf("timex: invalid clock %q", s)
	}

	return time.ParseInLocation(
		"2006-01-02 "+clockLayout,
		date+" "+s,
		zone,
	)
}

// IsValid reports whether Finish is at or after Start.
func (r Range) IsValid() bool {
	return !r.Finish.Before(r.Start)
}

// IsEmpty reports whether the endpoints represent the same instant.
func (r Range) IsEmpty() bool {
	return r.Start.Equal(r.Finish)
}

// Has reports whether t belongs to [Start, Finish).
// Invalid and empty ranges contain no instants.
func (r Range) Has(t time.Time) bool {
	return r.IsValid() && !t.Before(r.Start) && t.Before(r.Finish)
}

// Overlaps reports whether the ranges share a non-empty interval.
// Adjacent, empty, and invalid ranges do not overlap.
func (r Range) Overlaps(other Range) bool {
	if !r.IsValid() || !other.IsValid() || r.IsEmpty() || other.IsEmpty() {
		return false
	}
	return r.Start.Before(other.Finish) && other.Start.Before(r.Finish)
}

// Intersect returns the non-empty intersection of the ranges.
// If there is no intersection, it returns Range{} and false.
func (r Range) Intersect(other Range) (Range, bool) {
	if !r.Overlaps(other) {
		return Range{}, false
	}

	start := r.Start
	if other.Start.After(start) {
		start = other.Start
	}

	finish := r.Finish
	if other.Finish.Before(finish) {
		finish = other.Finish
	}

	return Range{Start: start, Finish: finish}, true
}

// Span returns Finish.Sub(Start).
// Invalid ranges produce a negative duration.
// Results outside Span's range saturate, following time.Time.Sub.
func (r Range) Span() Span {
	return Span(r.Finish.Sub(r.Start))
}

// String formats the range using its endpoints' local dates and clocks.
//
// The date is omitted when both endpoints have date 0001-01-01.
// A shared date is printed once; different dates are printed separately.
// Seconds are included for both endpoints if either has seconds or
// nanoseconds. Fractional seconds are printed only when non-zero.
//
// Timezone information is omitted. Use endpoints in the same location
// when an unambiguous local display is required.
func (r Range) String() string {
	start, finish := r.Start, r.Finish

	sy, sm, sd := start.Date()
	ey, em, ed := finish.Date()

	sameDay := sy == ey && sm == em && sd == ed
	noDate := sameDay && sy == 1 && sm == 1 && sd == 1

	clock := "15:04"
	if start.Second() != 0 || finish.Second() != 0 ||
		start.Nanosecond() != 0 || finish.Nanosecond() != 0 {
		clock = "15:04:05.999999999"
	}

	startLayout, finishLayout, sep := clock, clock, "-"

	switch {
	case noDate:
		// use defaults
	case sameDay:
		startLayout = "2006-01-02 " + clock
	default:
		startLayout = "2006-01-02 " + clock
		finishLayout = startLayout
		sep = " - "
	}

	return start.Format(startLayout) + sep + finish.Format(finishLayout)
}

// MarshalJSON encodes the range as a JSON string in UTC.
func (r Range) MarshalJSON() ([]byte, error) {
	text, err := r.MarshalText()
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(text))
}

func (r *Range) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	return r.UnmarshalText([]byte(s))
}

// MarshalText encodes the range with both endpoints in UTC.
func (r Range) MarshalText() ([]byte, error) {
	if !r.IsValid() {
		return nil, fmt.Errorf("timex: finish before start")
	}

	utc := Range{
		Start:  r.Start.UTC(),
		Finish: r.Finish.UTC(),
	}
	return []byte(utc.String()), nil
}

func (r *Range) UnmarshalText(data []byte) error {
	parsed, err := ParseRange(string(data))
	if err != nil {
		return err
	}
	*r = parsed
	return nil
}

func hasRangeDate(s string) bool {
	return len(s) >= 11 && s[4] == '-' && s[7] == '-' && s[10] == ' '
}
