package timex

import (
	"encoding/json"
	"testing"
	"time"
)

func TestParseRange(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		start, finish string
	}{
		{
			name:   "zero",
			input:  "00:00-00:00",
			start:  "0001-01-01T00:00:00Z",
			finish: "0001-01-01T00:00:00Z",
		},
		{
			name:   "clock only",
			input:  "18:30-22:00",
			start:  "0001-01-01T18:30:00Z",
			finish: "0001-01-01T22:00:00Z",
		},
		{
			name:   "whitespace",
			input:  " \t18:30  -  22:00\n",
			start:  "0001-01-01T18:30:00Z",
			finish: "0001-01-01T22:00:00Z",
		},
		{
			name:   "seconds",
			input:  "18:30:15-22:00:05",
			start:  "0001-01-01T18:30:15Z",
			finish: "0001-01-01T22:00:05Z",
		},
		{
			name:   "mixed clock precision",
			input:  "18:30-22:00:05.123",
			start:  "0001-01-01T18:30:00Z",
			finish: "0001-01-01T22:00:05.123Z",
		},
		{
			name:   "nanoseconds",
			input:  "18:30:00.000000001-22:00:00.999999999",
			start:  "0001-01-01T18:30:00.000000001Z",
			finish: "0001-01-01T22:00:00.999999999Z",
		},
		{
			name:   "shared date",
			input:  "2026-10-17 18:30-22:00",
			start:  "2026-10-17T18:30:00Z",
			finish: "2026-10-17T22:00:00Z",
		},
		{
			name:   "shared date with seconds",
			input:  "2026-10-17 18:30:15.1234-22:00:00",
			start:  "2026-10-17T18:30:15.1234Z",
			finish: "2026-10-17T22:00:00Z",
		},
		{
			name:   "explicit dates",
			input:  "2026-10-17 18:30 - 2026-10-18 02:00",
			start:  "2026-10-17T18:30:00Z",
			finish: "2026-10-18T02:00:00Z",
		},
		{
			name:   "explicit dates without separator spaces",
			input:  "2026-10-17 18:30-2026-10-18 02:00",
			start:  "2026-10-17T18:30:00Z",
			finish: "2026-10-18T02:00:00Z",
		},
		{
			name:   "year boundary",
			input:  "2026-12-31 23:59:59.999999999 - 2027-01-01 00:00",
			start:  "2026-12-31T23:59:59.999999999Z",
			finish: "2027-01-01T00:00:00Z",
		},
		{
			name:   "leap day",
			input:  "2028-02-29 09:00-10:00",
			start:  "2028-02-29T09:00:00Z",
			finish: "2028-02-29T10:00:00Z",
		},
		{
			name:   "hours only",
			input:  "18-20",
			start:  "0001-01-01T18:00:00Z",
			finish: "0001-01-01T20:00:00Z",
		},
		{
			name:   "minutes to hours",
			input:  "18:30-20",
			start:  "0001-01-01T18:30:00Z",
			finish: "0001-01-01T20:00:00Z",
		},
		{
			name:   "hours to seconds",
			input:  "18-20:30:40",
			start:  "0001-01-01T18:00:00Z",
			finish: "0001-01-01T20:30:40Z",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRange(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			start := rangeTestTime(t, tt.start)
			finish := rangeTestTime(t, tt.finish)

			if !got.Start.Equal(start) || !got.Finish.Equal(finish) {
				t.Fatalf(
					"ParseRange(%q) = [%s, %s), want [%s, %s)",
					tt.input, got.Start, got.Finish, start, finish,
				)
			}
			if got.Start.Location() != time.UTC ||
				got.Finish.Location() != time.UTC {
				t.Error("parsed endpoints must use UTC")
			}

			// Formatting and parsing must preserve the parsed value.
			again, err := ParseRange(got.String())
			if err != nil {
				t.Fatalf("parse formatted range: %v", err)
			}
			if !again.Start.Equal(got.Start) ||
				!again.Finish.Equal(got.Finish) {
				t.Fatal("String/ParseRange round trip changed endpoints")
			}
		})
	}
}

func TestParseRangeInZone(t *testing.T) {
	warsaw, err := time.LoadLocation("Europe/Warsaw")
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name       string
		input      string
		zone       *time.Location
		start, end string
	}{
		{
			name:  "UTC",
			input: "2026-10-17 18:30-22",
			zone:  time.UTC,
			start: "2026-10-17T18:30:00Z",
			end:   "2026-10-17T22:00:00Z",
		},
		{
			name:  "Warsaw summer offset",
			input: "2026-10-17 18:30-22",
			zone:  warsaw,
			start: "2026-10-17T18:30:00+02:00",
			end:   "2026-10-17T22:00:00+02:00",
		},
		{
			name:  "Warsaw winter offset",
			input: "2026-12-17 18-20:30:40",
			zone:  warsaw,
			start: "2026-12-17T18:00:00+01:00",
			end:   "2026-12-17T20:30:40+01:00",
		},
		{
			name:  "fixed zone with fractional seconds",
			input: "2026-10-17 18:30:00.123-20",
			zone:  time.FixedZone("UTC+0530", 5*3600+30*60),
			start: "2026-10-17T18:30:00.123+05:30",
			end:   "2026-10-17T20:00:00+05:30",
		},
		{
			name:  "clock only",
			input: "18-20:30:40",
			zone:  time.FixedZone("UTC-0400", -4*3600),
			start: "0001-01-01T18:00:00-04:00",
			end:   "0001-01-01T20:30:40-04:00",
		},
		{
			name:  "explicit dates",
			input: "2026-10-17 18:30 - 2026-10-18 02:00",
			zone:  warsaw,
			start: "2026-10-17T18:30:00+02:00",
			end:   "2026-10-18T02:00:00+02:00",
		},
		{
			name:  "spring DST transition",
			input: "2026-03-29 01:30-03:30",
			zone:  warsaw,
			start: "2026-03-29T01:30:00+01:00",
			end:   "2026-03-29T03:30:00+02:00",
		},
		{
			name:  "autumn DST transition",
			input: "2026-10-25 01:30-03:30",
			zone:  warsaw,
			start: "2026-10-25T01:30:00+02:00",
			end:   "2026-10-25T03:30:00+01:00",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRangeInZone(tt.input, tt.zone)
			if err != nil {
				t.Fatal(err)
			}

			start := rangeTestTime(t, tt.start)
			end := rangeTestTime(t, tt.end)

			if !got.Start.Equal(start) || !got.Finish.Equal(end) {
				t.Fatalf(
					"got [%s, %s), want [%s, %s)",
					got.Start, got.Finish, start, end,
				)
			}
			if got.Start.Location() != tt.zone ||
				got.Finish.Location() != tt.zone {
				t.Fatal("endpoints did not retain the supplied zone")
			}

			// Validates actual elapsed time across DST transitions.
			if got.Span() != Span(end.Sub(start)) {
				t.Fatalf("unexpected duration: %v", got.Span())
			}

			again, err := ParseRangeInZone(got.String(), tt.zone)
			if err != nil {
				t.Fatal(err)
			}
			if !again.Start.Equal(got.Start) ||
				!again.Finish.Equal(got.Finish) {
				t.Fatal("String round trip changed endpoints")
			}
		})
	}
}

func TestParseRangeInZoneInvalid(t *testing.T) {
	zone := time.FixedZone("UTC+0200", 2*3600)

	testCases := []struct {
		name  string
		input string
		zone  *time.Location
	}{
		{"nil zone", "18-20", nil},
		{"empty", "", zone},
		{"missing separator", "18:30", zone},
		{"invalid start", "bad-20", zone},
		{"invalid end", "18-bad", zone},
		{"invalid start clock", "18:00:00:00-20", zone},
		{"invalid end clock", "18-20:00:00:00", zone},
		{"invalid date", "2026-02-29 18-20", zone},
		{"invalid end date", "2026-10-17 18 - 2026-02-29 20", zone},
		{"reversed", "20-18", zone},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRangeInZone(tt.input, tt.zone)
			if err == nil {
				t.Fatal("expected an error")
			}
			if got != (Range{}) {
				t.Fatalf("got %v, want zero range", got)
			}
		})
	}
}

func TestParseRangeDefaultsToUTC(t *testing.T) {
	const input = "2026-10-17 18:30-22"

	got, err := ParseRange(input)
	if err != nil {
		t.Fatal(err)
	}
	want, err := ParseRangeInZone(input, time.UTC)
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Fatalf("ParseRange() = %v, want %v", got, want)
	}
	if got.Start.Location() != time.UTC || got.Finish.Location() != time.UTC {
		t.Fatal("ParseRange must use UTC")
	}
}

func TestParseRangeInvalid(t *testing.T) {
	testCases := []string{
		"",
		"   ",
		"-",
		"18:30",
		"18:30-",
		"-22:00",
		"18:30/22:00",
		"18:30–22:00",
		"foo-bar",
		"24:00-25:00",
		"18:60-22:00",
		"18:30:60-22:00",
		"18:30-22:00 trailing",
		"18:30-22:00-23:00",
		"22:00-18:30",
		"18:30-02:00",
		"2026-02-29 09:00-10:00",
		"2026-04-31 09:00-10:00",
		"2026-13-01 09:00-10:00",
		"2026-10-17 22:00-18:30",
		"2026-10-18 02:00 - 2026-10-17 18:30",
	}

	for _, input := range testCases {
		t.Run(input, func(t *testing.T) {
			if got, err := ParseRange(input); err == nil {
				t.Fatalf("ParseRange(%q) = %v, want error", input, got)
			}
		})
	}
}

func TestRangeString(t *testing.T) {
	testCases := []struct {
		name          string
		start, finish string
		want          string
	}{
		{
			name:   "zero",
			start:  "0001-01-01T00:00:00Z",
			finish: "0001-01-01T00:00:00Z",
			want:   "00:00-00:00",
		},
		{
			name:   "clock only",
			start:  "0001-01-01T18:30:00Z",
			finish: "0001-01-01T22:00:00Z",
			want:   "18:30-22:00",
		},
		{
			name:   "start seconds",
			start:  "0001-01-01T18:30:15Z",
			finish: "0001-01-01T22:00:00Z",
			want:   "18:30:15-22:00:00",
		},
		{
			name:   "end seconds",
			start:  "0001-01-01T18:30:00Z",
			finish: "0001-01-01T22:00:05Z",
			want:   "18:30:00-22:00:05",
		},
		{
			name:   "same date",
			start:  "2026-10-17T18:30:00Z",
			finish: "2026-10-17T22:00:00Z",
			want:   "2026-10-17 18:30-22:00",
		},
		{
			name:   "fractional seconds",
			start:  "2026-10-17T18:30:15.123400Z",
			finish: "2026-10-17T22:00:00Z",
			want:   "2026-10-17 18:30:15.1234-22:00:00",
		},
		{
			name:   "nanoseconds without whole seconds",
			start:  "2026-10-17T18:30:00Z",
			finish: "2026-10-17T22:00:00.000000001Z",
			want:   "2026-10-17 18:30:00-22:00:00.000000001",
		},
		{
			name:   "different dates",
			start:  "2026-10-17T18:30:00Z",
			finish: "2026-10-18T02:00:00Z",
			want:   "2026-10-17 18:30 - 2026-10-18 02:00",
		},
		{
			name:   "different dates with seconds",
			start:  "2026-10-17T18:30:15Z",
			finish: "2026-10-18T02:00:00Z",
			want:   "2026-10-17 18:30:15 - 2026-10-18 02:00:00",
		},
		{
			name:   "different months with same day number",
			start:  "2026-10-17T18:30:00Z",
			finish: "2026-11-17T22:00:00Z",
			want:   "2026-10-17 18:30 - 2026-11-17 22:00",
		},
		{
			name:   "local clock",
			start:  "2026-10-17T18:30:00+02:00",
			finish: "2026-10-17T22:00:00+02:00",
			want:   "2026-10-17 18:30-22:00",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := Range{
				Start:  rangeTestTime(t, tt.start),
				Finish: rangeTestTime(t, tt.finish),
			}
			if got := r.String(); got != tt.want {
				t.Fatalf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRangeValidity(t *testing.T) {
	start := rangeTestTime(t, "2026-10-17T18:00:00Z")

	testCases := []struct {
		name  string
		r     Range
		valid bool
		empty bool
	}{
		{"zero", Range{}, true, true},
		{"nonempty", Range{start, start.Add(time.Hour)}, true, false},
		{"empty", Range{start, start}, true, true},
		{"reversed", Range{start.Add(time.Hour), start}, false, false},
		{"same instant in another zone", Range{start, start.In(time.FixedZone("X", 7200))}, true, true},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.IsValid(); got != tt.valid {
				t.Errorf("Valid() = %v, want %v", got, tt.valid)
			}
			if got := tt.r.IsEmpty(); got != tt.empty {
				t.Errorf("Empty() = %v, want %v", got, tt.empty)
			}

			got, err := NewRange(tt.r.Start, tt.r.Finish)
			if (err == nil) != tt.valid {
				t.Fatalf("NewRange() error = %v, valid = %v", err, tt.valid)
			}
			if err == nil &&
				(!got.Start.Equal(tt.r.Start) || !got.Finish.Equal(tt.r.Finish)) {
				t.Fatal("NewRange() changed endpoints")
			}
		})
	}
}

func TestRangeHas(t *testing.T) {
	start := rangeTestTime(t, "2026-10-17T18:00:00Z")
	r := Range{start, start.Add(time.Hour)}

	testCases := []struct {
		name string
		at   time.Time
		want bool
	}{
		{"before", start.Add(-time.Nanosecond), false},
		{"start", start, true},
		{"inside", start.Add(30 * time.Minute), true},
		{"before end", r.Finish.Add(-time.Nanosecond), true},
		{"end", r.Finish, false},
		{"after", r.Finish.Add(time.Nanosecond), false},
		{"another zone", start.In(time.FixedZone("X", 7200)), true},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Has(tt.at); got != tt.want {
				t.Fatalf("Has() = %v, want %v", got, tt.want)
			}
		})
	}

	if (Range{start, start}).Has(start) {
		t.Error("empty range Has its endpoint")
	}
	if (Range{r.Finish, r.Start}).Has(start) {
		t.Error("invalid range Has an instant")
	}
}

func TestRangeIntersection(t *testing.T) {
	base := rangeTestTime(t, "2026-10-17T00:00:00Z")
	at := func(hour int) time.Time {
		return base.Add(time.Duration(hour) * time.Hour)
	}
	r := Range{at(10), at(14)}

	testCases := []struct {
		name          string
		other         Range
		start, finish int
		ok            bool
	}{
		{"identical", Range{at(10), at(14)}, 10, 14, true},
		{"inside", Range{at(11), at(13)}, 11, 13, true},
		{"surrounding", Range{at(9), at(15)}, 10, 14, true},
		{"left overlap", Range{at(9), at(12)}, 10, 12, true},
		{"right overlap", Range{at(12), at(16)}, 12, 14, true},
		{"left adjacent", Range{at(8), at(10)}, 0, 0, false},
		{"right adjacent", Range{at(14), at(16)}, 0, 0, false},
		{"disjoint", Range{at(16), at(18)}, 0, 0, false},
		{"empty inside", Range{at(12), at(12)}, 0, 0, false},
		{"invalid", Range{at(13), at(11)}, 0, 0, false},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Both operations should be symmetric.
			for _, pair := range [][2]Range{{r, tt.other}, {tt.other, r}} {
				a, b := pair[0], pair[1]

				if got := a.Overlaps(b); got != tt.ok {
					t.Errorf("Overlaps() = %v, want %v", got, tt.ok)
				}

				got, ok := a.Intersect(b)
				if ok != tt.ok {
					t.Fatalf("Intersect() ok = %v, want %v", ok, tt.ok)
				}
				if ok {
					if !got.Start.Equal(at(tt.start)) ||
						!got.Finish.Equal(at(tt.finish)) {
						t.Errorf("Intersect() = %v, unexpected endpoints", got)
					}
				} else if got != (Range{}) {
					t.Errorf("Intersect() = %v, want zero range", got)
				}
			}
		})
	}
}

func TestRangeDuration(t *testing.T) {
	start := rangeTestTime(t, "2026-10-17T18:00:00Z")

	for _, want := range []time.Duration{
		0,
		time.Nanosecond,
		90 * time.Minute,
		-time.Hour,
	} {
		r := Range{start, start.Add(want)}
		if got := r.Span(); got != Span(want) {
			t.Errorf("Duration() = %v, want %v", got, want)
		}
	}
}

func TestRangeMarshal(t *testing.T) {
	r := Range{
		Start:  rangeTestTime(t, "2026-10-17T18:30:00Z"),
		Finish: rangeTestTime(t, "2026-10-17T22:00:00Z"),
	}

	text, err := r.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(text), "2026-10-17 18:30-22:00"; got != want {
		t.Errorf("MarshalText() = %q, want %q", got, want)
	}

	// Check both value and pointer marshaling.
	for _, value := range []any{r, &r} {
		data, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := string(data), `"2026-10-17 18:30-22:00"`; got != want {
			t.Errorf("json.Marshal() = %s, want %s", got, want)
		}
	}
}

func TestRangeUnmarshalText(t *testing.T) {
	const valid = "2026-10-17 18:30-22:00"

	want := Range{
		Start:  rangeTestTime(t, "2026-10-17T18:30:00Z"),
		Finish: rangeTestTime(t, "2026-10-17T22:00:00Z"),
	}

	var got Range
	if err := got.UnmarshalText([]byte(valid)); err != nil {
		t.Fatal(err)
	}
	if !got.Start.Equal(want.Start) || !got.Finish.Equal(want.Finish) {
		t.Fatalf("UnmarshalText() = %v, want %v", got, want)
	}

	for _, input := range []string{
		"",
		"invalid",
		"18:30:00:00-22:00",
		"18:30-22:00:00:00",
		"22-18",
	} {
		t.Run(input, func(t *testing.T) {
			got := want
			if err := got.UnmarshalText([]byte(input)); err == nil {
				t.Fatal("expected an error")
			}
			if got != want {
				t.Fatal("failed UnmarshalText changed the receiver")
			}
		})
	}
}

func TestRangeUnmarshalJSON(t *testing.T) {
	want := Range{
		Start:  rangeTestTime(t, "2026-10-17T18:30:00Z"),
		Finish: rangeTestTime(t, "2026-10-17T22:00:00Z"),
	}

	var got Range
	if err := json.Unmarshal(
		[]byte(`"2026-10-17 18:30-22:00"`), &got,
	); err != nil {
		t.Fatal(err)
	}
	if !got.Start.Equal(want.Start) || !got.Finish.Equal(want.Finish) {
		t.Fatalf("UnmarshalJSON() = %v, want %v", got, want)
	}

	for _, input := range []string{
		``,
		`"unterminated`,
		`123`,
		`true`,
		`{}`,
		`[]`,
		`null`,
		`"invalid"`,
		`"22-18"`,
		`"18:30-22:00:00:00"`,
	} {
		t.Run(input, func(t *testing.T) {
			got := want

			// Call directly so malformed JSON reaches this method.
			if err := got.UnmarshalJSON([]byte(input)); err == nil {
				t.Fatal("expected an error")
			}
			if got != want {
				t.Fatal("failed UnmarshalJSON changed the receiver")
			}
		})
	}
}

func FuzzParseRange(f *testing.F) {
	for _, seed := range []string{
		"",
		"-",
		"00:00-00:00",
		"18:30-22:00",
		" 18:30 - 22:00 ",
		"18:30:15.123456789-22:00:00",
		"2026-10-17 18:30-22:00",
		"2026-10-17 18:30 - 2026-10-18 02:00",
		"2028-02-29 09:00-10:00",
		"22:00-18:30",
		"24:00-25:00",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		r, err := ParseRange(input)
		if err != nil {
			return
		}

		if !r.IsValid() {
			t.Fatalf("ParseRange(%q) returned an invalid range", input)
		}

		check := func(name string, got Range) {
			t.Helper()
			if !got.Start.Equal(r.Start) || !got.Finish.Equal(r.Finish) {
				t.Fatalf(
					"%s changed endpoints: [%s, %s) -> [%s, %s)",
					name, r.Start, r.Finish, got.Start, got.Finish,
				)
			}
		}

		parsed, err := ParseRange(r.String())
		if err != nil {
			t.Fatalf("parse String() %q: %v", r.String(), err)
		}
		check("String round trip", parsed)

		text, err := r.MarshalText()
		if err != nil {
			t.Fatal(err)
		}
		var fromText Range
		if err := fromText.UnmarshalText(text); err != nil {
			t.Fatalf("unmarshal text %q: %v", text, err)
		}
		check("text round trip", fromText)

		data, err := json.Marshal(r)
		if err != nil {
			t.Fatalf("marshal JSON: %v", err)
		}

		var fromJSON Range
		if err := json.Unmarshal(data, &fromJSON); err != nil {
			t.Fatalf("unmarshal JSON %s: %v", data, err)
		}
		check("JSON round trip", fromJSON)
	})
}

func rangeTestTime(t *testing.T, s string) time.Time {
	t.Helper()

	v, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t.Fatal(err)
	}
	return v
}
