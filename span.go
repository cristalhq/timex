package timex

import (
	"encoding/json"
	"strings"
	"time"
)

// Span represents a duration in nanoseconds.
// Similar to [time.Duration] but with improved marshaling.
type Span time.Duration

// ParseSpan parses a duration using time.ParseDuration.
func ParseSpan(s string) (Span, error) {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, err
	}
	return Span(d), nil
}

// String returns the duration formatted like [time.Duration.String]
// Omits zero minutes and zero seconds.
// For example, 1h0m0s becomes 1h, and 1m0s becomes 1m.
func (s Span) String() string {
	ss := s.Dur().String()
	ss = strings.ReplaceAll(ss, "m0s", "m")
	ss = strings.ReplaceAll(ss, "h0m", "h")
	return ss
}

// Dur returns a time.Duration.
func (s Span) Dur() time.Duration {
	return time.Duration(s)
}

// Mul returns s multiplied by x.
//
// The calculation uses float64 and may lose precision for large durations.
// Both x and the product must be finite, and the product must be in
// [-2^63, 2^63). These requirements are not checked.
func (s Span) Mul(x float64) Span {
	return Span(float64(s) * x)
}

func (s Span) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Span) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		return err
	}
	return s.UnmarshalText([]byte(text))
}

func (s Span) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *Span) UnmarshalText(text []byte) error {
	v, err := ParseSpan(string(text))
	if err != nil {
		return err
	}
	*s = v
	return nil
}

// Spans is a slice of [Span].
type Spans []Span

func (ss Spans) String() string {
	return strings.Join(ss.Strings(), ",")
}

func (ss Spans) Strings() []string {
	res := make([]string, len(ss))
	for i := range ss {
		res[i] = ss[i].String()
	}
	return res
}

// Durs returns each span as a [time.Duration].
func (ss Spans) Durs() []time.Duration {
	res := make([]time.Duration, len(ss))
	for i := range ss {
		res[i] = ss[i].Dur()
	}
	return res
}
