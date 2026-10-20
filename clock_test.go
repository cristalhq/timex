package timex

import (
	"sync"
	"testing"
	"time"
)

func TestSystemClock(t *testing.T) {
	before := time.Now()
	got := (SystemClock{}).Now()
	after := time.Now()

	if got.Before(before) || got.After(after) {
		t.Fatalf("Now() = %v; want between %v and %v", got, before, after)
	}
	if got.Location() != time.Local {
		t.Fatalf("location = %v; want Local", got.Location())
	}
}

func TestUTCClock(t *testing.T) {
	before := time.Now()
	got := (UTCClock{}).Now()
	after := time.Now()

	if got.Before(before) || got.After(after) {
		t.Fatalf("Now() = %v; want between %v and %v", got, before, after)
	}
	if got.Location() != time.UTC {
		t.Fatalf("location = %v; want UTC", got.Location())
	}
	if got != got.Round(0) {
		t.Fatal("Now() contains a monotonic clock reading")
	}
}

func TestManualClockZeroValue(t *testing.T) {
	var clock ManualClock

	if got := clock.Now(); got != (time.Time{}) {
		t.Fatalf("Now() = %v; want zero time", got)
	}

	want := time.Time{}.Add(time.Hour)
	clock.Add(time.Hour)

	if got := clock.Now(); got != want {
		t.Fatalf("Now() = %v; want %v", got, want)
	}
}

func TestNewManualClock(t *testing.T) {
	want := time.Now()
	clock := NewManualClock(want)

	if got := clock.Now(); got != want {
		t.Fatalf("Now() = %v; want exact initial time %v", got, want)
	}
}

func TestManualClockSet(t *testing.T) {
	zone := time.FixedZone("test", 2*60*60)
	values := []time.Time{
		time.Now(),
		time.Date(2026, time.October, 20, 17, 0, 0, 123, zone),
		{},
	}

	var clock ManualClock
	for _, want := range values {
		clock.Set(want)

		if got := clock.Now(); got != want {
			t.Fatalf("Now() = %v; want exact set time %v", got, want)
		}
	}
}

func TestManualClockAdd(t *testing.T) {
	start := time.Date(2026, time.October, 20, 23, 30, 0, 0, time.UTC)

	tests := []struct {
		name string
		d    time.Duration
	}{
		{"positive", 2 * time.Hour},
		{"negative", -48 * time.Hour},
		{"zero", 0},
		{"nanosecond", time.Nanosecond},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clock := NewManualClock(start)
			clock.Add(tt.d)

			want := start.Add(tt.d)
			if got := clock.Now(); got != want {
				t.Fatalf("Now() = %v; want %v", got, want)
			}
		})
	}
}

func TestManualClockConcurrentAdd(t *testing.T) {
	const workers, iterations = 16, 1_000

	start := time.Date(2026, time.October, 20, 0, 0, 0, 0, time.UTC)
	clock := NewManualClock(start)

	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			for j := 0; j < iterations; j++ {
				clock.Add(time.Nanosecond)
				_ = clock.Now()
			}
		})
	}
	wg.Wait()

	want := start.Add(time.Duration(workers*iterations) * time.Nanosecond)
	if got := clock.Now(); got != want {
		t.Fatalf("Now() = %v; want %v", got, want)
	}
}

func TestManualClockConcurrentSet(t *testing.T) {
	const workers, iterations = 16, 1_000

	initial := time.Date(2026, time.October, 20, 0, 0, 0, 0, time.UTC)
	replacement := initial.Add(time.Hour)
	clock := NewManualClock(initial)

	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			for range iterations {
				clock.Set(replacement)
				if got := clock.Now(); got != replacement {
					t.Errorf("Now() = %v; want %v", got, replacement)
					return
				}
			}
		})
	}
	wg.Wait()
}
