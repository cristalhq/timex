package timex

import (
	"errors"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name    string
		value   string
		layouts []string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "no layouts",
			value:   "2026-10-20",
			wantErr: true,
		},
		{
			name:    "first layout matches",
			value:   "2026-10-20",
			layouts: []string{"2006-01-02", time.RFC3339},
			want:    time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			name:    "second layout matches",
			value:   "2026-10-20",
			layouts: []string{time.RFC3339, "2006-01-02"},
			want:    time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "third layout matches",
			value: "2026-10-20",
			layouts: []string{
				time.RFC3339,
				"02/01/2006",
				"2006-01-02",
			},
			want: time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			name:    "all layouts fail",
			value:   "invalid",
			layouts: []string{time.RFC3339, "2006-01-02", "15:04"},
			wantErr: true,
		},
		{
			name:    "empty value fails",
			layouts: []string{"2006-01-02"},
			wantErr: true,
		},
		{
			name:    "empty layout matches empty value",
			layouts: []string{""},
			want:    time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:    "first matching layout wins",
			value:   "01/02/2026",
			layouts: []string{"02/01/2006", "01/02/2006"},
			want:    time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:    "offset and nanoseconds",
			value:   "2026-10-20T12:34:56.123456789+02:00",
			layouts: []string{time.RFC3339Nano},
			want: time.Date(
				2026, 10, 20, 12, 34, 56, 123456789,
				time.FixedZone("", 2*60*60),
			),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.value, tt.layouts...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !got.Equal(tt.want) {
				t.Fatalf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseReturnsFirstError(t *testing.T) {
	const value = "invalid"
	layouts := []string{time.RFC3339, "2006-01-02", "15:04"}

	_, err := Parse(value, layouts...)

	var got *time.ParseError
	if !errors.As(err, &got) {
		t.Fatalf("error = %v, want *time.ParseError", err)
	}

	_, firstErr := time.Parse(layouts[0], value)
	want := firstErr.(*time.ParseError)

	if *got != *want {
		t.Fatalf("error = %+v, want %+v", got, want)
	}
}

func TestMinMax(t *testing.T) {
	zero := time.Time{}
	early := time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)
	middle := time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)
	late := middle.Add(time.Nanosecond)
	equivalent := middle.In(time.FixedZone("UTC+02", 2*60*60))

	testCases := []struct {
		name   string
		values []time.Time
		min    time.Time
		max    time.Time
	}{
		{name: "empty"},
		{name: "all zero", values: []time.Time{zero, zero}},
		{
			name:   "single",
			values: []time.Time{middle},
			min:    middle,
			max:    middle,
		},
		{
			name:   "ascending",
			values: []time.Time{early, middle, late},
			min:    early,
			max:    late,
		},
		{
			name:   "descending",
			values: []time.Time{late, middle, early},
			min:    early,
			max:    late,
		},
		{
			name:   "ignore zero",
			values: []time.Time{zero, middle, zero, early, late, zero},
			min:    early,
			max:    late,
		},
		{
			name:   "before zero time",
			values: []time.Time{zero, early},
			min:    early,
			max:    early,
		},
		{
			name:   "equal instants preserve first value",
			values: []time.Time{equivalent, middle},
			min:    equivalent,
			max:    equivalent,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Exact comparison also verifies which equal value is retained.
			if got := Min(tt.values...); got != tt.min {
				t.Errorf("Min() = %v, want %v", got, tt.min)
			}
			if got := Max(tt.values...); got != tt.max {
				t.Errorf("Max() = %v, want %v", got, tt.max)
			}
		})
	}
}

func TestDropDateAndTime(t *testing.T) {
	loc := time.FixedZone("UTC+02", 2*60*60)

	testCases := []struct {
		name string
		in   time.Time
	}{
		{name: "zero"},
		{
			name: "UTC with nanoseconds",
			in:   time.Date(2026, 10, 20, 12, 34, 56, 123456789, time.UTC),
		},
		{
			name: "fixed offset",
			in:   time.Date(2024, 2, 29, 23, 59, 58, 987654321, loc),
		},
		{
			name: "midnight",
			in:   time.Date(2026, 10, 20, 0, 0, 0, 0, loc),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			dateDropped := DropDate(tt.in)
			y, m, d := dateDropped.Date()
			if y != 0 || m != time.January || d != 1 {
				t.Errorf("DropDate() date = %d-%02d-%02d", y, m, d)
			}
			if dateDropped.Hour() != tt.in.Hour() ||
				dateDropped.Minute() != tt.in.Minute() ||
				dateDropped.Second() != tt.in.Second() ||
				dateDropped.Nanosecond() != tt.in.Nanosecond() {
				t.Errorf("DropDate() changed clock: %v", dateDropped)
			}
			if dateDropped.Location() != tt.in.Location() {
				t.Error("DropDate() changed location")
			}
			if got := DropDate(dateDropped); got != dateDropped {
				t.Error("DropDate() is not idempotent")
			}

			timeDropped := DropTime(tt.in)
			gotY, gotM, gotD := timeDropped.Date()
			wantY, wantM, wantD := tt.in.Date()
			if gotY != wantY || gotM != wantM || gotD != wantD {
				t.Errorf("DropTime() changed date: %v", timeDropped)
			}
			if timeDropped.Hour() != 0 ||
				timeDropped.Minute() != 0 ||
				timeDropped.Second() != 0 ||
				timeDropped.Nanosecond() != 0 {
				t.Errorf("DropTime() clock is not midnight: %v", timeDropped)
			}
			if timeDropped.Location() != tt.in.Location() {
				t.Error("DropTime() changed location")
			}
			if got := DropTime(timeDropped); got != timeDropped {
				t.Error("DropTime() is not idempotent")
			}
		})
	}
}
