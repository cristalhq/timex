package timex

import (
	"encoding"
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

var (
	_ encoding.TextMarshaler   = Span(0)
	_ encoding.TextUnmarshaler = (*Span)(nil)
)

func TestParseSpan(t *testing.T) {
	testCases := []struct {
		input string
		want  Span
	}{
		{"0", 0},
		{"0s", 0},
		{"1ns", Span(time.Nanosecond)},
		{"1us", Span(time.Microsecond)},
		{"1µs", Span(time.Microsecond)},
		{"1μs", Span(time.Microsecond)},
		{"1.5ms", Span(1500 * time.Microsecond)},
		{"1m", Span(time.Minute)},
		{"1h30m", Span(90 * time.Minute)},
		{"+2s", Span(2 * time.Second)},
		{"-1.5s", Span(-1500 * time.Millisecond)},
		{"2562047h47m16.854775807s", Span(1<<63 - 1)},
		{"-2562047h47m16.854775808s", Span(-1 << 63)},
	}

	for _, tt := range testCases {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseSpan(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("got %d ns, want %d ns", got, tt.want)
			}
		})
	}
}

func TestParseSpanInvalid(t *testing.T) {
	testCases := []string{
		"",
		" ",
		"1",
		"1d",
		"invalid",
		"1s trailing",
		" 1s",
		"1s ",
		"2562047h47m16.854775808s",
		"-2562047h47m16.854775809s",
	}

	for _, input := range testCases {
		t.Run(input, func(t *testing.T) {
			got, err := ParseSpan(input)
			if err == nil {
				t.Fatal("expected an error")
			}
			if got != 0 {
				t.Fatalf("got %d ns, want zero", got)
			}
		})
	}
}

func TestSpanString(t *testing.T) {
	testCases := []struct {
		value time.Duration
		want  string
	}{
		{0, "0s"},
		{time.Nanosecond, "1ns"},
		{time.Microsecond, "1µs"},
		{1500 * time.Millisecond, "1.5s"},
		{time.Minute, "1m"},
		{time.Minute + time.Second, "1m1s"},
		{time.Minute + time.Millisecond, "1m0.001s"},
		{time.Hour, "1h"},
		{time.Hour + time.Second, "1h1s"},
		{time.Hour + time.Millisecond, "1h0.001s"},
		{time.Hour + time.Minute, "1h1m"},
		{time.Hour + time.Minute + time.Second, "1h1m1s"},
		{24 * time.Hour, "24h"},
		{-time.Hour, "-1h"},
		{-time.Hour - time.Second, "-1h1s"},
		{-time.Minute, "-1m"},
		{time.Duration(1<<63 - 1), "2562047h47m16.854775807s"},
		{time.Duration(-1 << 63), "-2562047h47m16.854775808s"},
	}

	for _, tt := range testCases {
		t.Run(tt.want, func(t *testing.T) {
			s := Span(tt.value)
			if got := s.String(); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
			if got := s.Dur(); got != tt.value {
				t.Fatalf("Dur() = %d, want %d", got, tt.value)
			}
		})
	}
}

func TestSpanMul(t *testing.T) {
	testCases := []struct {
		name string
		s    Span
		x    float64
		want Span
	}{
		{"zero span", 0, 2, 0},
		{"zero multiplier", 10, 0, 0},
		{"identity", 10, 1, 10},
		{"double", 10, 2, 20},
		{"fraction", 10, 0.5, 5},
		{"truncate positive", 3, 0.5, 1},
		{"truncate negative", -3, 0.5, -1},
		{"negative multiplier", 3, -0.5, -1},
		{"both negative", -3, -0.5, 1},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Mul(tt.x); got != tt.want {
				t.Fatalf("got %d ns, want %d ns", got, tt.want)
			}
		})
	}
}

func TestSpanMarshalJSON(t *testing.T) {
	testCases := []struct {
		s    Span
		want string
	}{
		{0, `"0s"`},
		{Span(time.Minute), `"1m"`},
		{Span(time.Hour + time.Minute), `"1h1m"`},
		{Span(-time.Second), `"-1s"`},
		{Span(time.Microsecond), `"1µs"`},
	}

	for _, tt := range testCases {
		t.Run(tt.want, func(t *testing.T) {
			got, err := tt.s.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tt.want {
				t.Fatalf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestSpanUnmarshalJSON(t *testing.T) {
	tests := []struct {
		input string
		want  Span
	}{
		{`"0s"`, 0},
		{`"1m"`, Span(time.Minute)},
		{`"1h30m"`, Span(90 * time.Minute)},
		{`"-1.5s"`, Span(-1500 * time.Millisecond)},
		{`"1\u00b5s"`, Span(time.Microsecond)},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			s := Span(time.Hour)
			if err := s.UnmarshalJSON([]byte(tt.input)); err != nil {
				t.Fatal(err)
			}
			if s != tt.want {
				t.Fatalf("got %d ns, want %d ns", s, tt.want)
			}
		})
	}
}

func TestSpanUnmarshalJSONInvalid(t *testing.T) {
	testCases := []string{
		``,
		`null`,
		`123`,
		`true`,
		`{}`,
		`[]`,
		`""`,
		`"invalid"`,
		`"1d"`,
		`"1s`,
		`"1s" trailing`,
	}

	for _, input := range testCases {
		t.Run(input, func(t *testing.T) {
			want := Span(time.Hour)
			s := want

			if err := s.UnmarshalJSON([]byte(input)); err == nil {
				t.Fatal("expected an error")
			}
			if s != want {
				t.Fatalf("receiver changed: got %d ns, want %d ns", s, want)
			}
		})
	}
}

func TestSpanMarshalText(t *testing.T) {
	s := Span(time.Hour + time.Minute)

	text, err := s.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	if got := string(text); got != "1h1m" {
		t.Fatalf("got %q, want %q", got, "1h1m")
	}

	var decoded Span
	if err := decoded.UnmarshalText(text); err != nil {
		t.Fatal(err)
	}
	if decoded != s {
		t.Fatalf("got %d ns, want %d ns", decoded, s)
	}
}

func TestSpanUnmarshalText(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		s := Span(time.Hour)
		if err := s.UnmarshalText([]byte("1m30s")); err != nil {
			t.Fatal(err)
		}
		if want := Span(90 * time.Second); s != want {
			t.Fatalf("got %d ns, want %d ns", s, want)
		}
	})

	for _, input := range []string{"", "invalid", "1d"} {
		t.Run("invalid/"+input, func(t *testing.T) {
			want := Span(time.Hour)
			s := want

			if err := s.UnmarshalText([]byte(input)); err == nil {
				t.Fatal("expected an error")
			}
			if s != want {
				t.Fatalf("receiver changed: got %d ns, want %d ns", s, want)
			}
		})
	}
}

func TestSpanJSON(t *testing.T) {
	type config struct {
		Timeout Span `json:"timeout"`
	}

	want := config{Timeout: Span(90 * time.Minute)}
	data, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(data); got != `{"timeout":"1h30m"}` {
		t.Fatalf("unexpected JSON: %s", got)
	}

	var got config
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	for _, input := range []string{
		`{"timeout":"invalid"}`,
		`{"timeout":123}`,
	} {
		if err := json.Unmarshal([]byte(input), &got); err == nil {
			t.Fatalf("expected an error for %s", input)
		}
	}
}

func TestSpans(t *testing.T) {
	testCases := []struct {
		name    string
		spans   Spans
		text    string
		strings []string
		durs    []time.Duration
	}{
		{"nil", nil, "", []string{}, []time.Duration{}},
		{"empty", Spans{}, "", []string{}, []time.Duration{}},
		{
			"values",
			Spans{0, Span(time.Minute), Span(-time.Hour)},
			"0s,1m,-1h",
			[]string{"0s", "1m", "-1h"},
			[]time.Duration{0, time.Minute, -time.Hour},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.spans.String(); got != tt.text {
				t.Fatalf("String() = %q, want %q", got, tt.text)
			}
			if got := tt.spans.Strings(); !reflect.DeepEqual(got, tt.strings) {
				t.Fatalf("Strings() = %#v, want %#v", got, tt.strings)
			}
			if got := tt.spans.Durs(); !reflect.DeepEqual(got, tt.durs) {
				t.Fatalf("Durs() = %#v, want %#v", got, tt.durs)
			}
		})
	}
}

func TestSpansConversionsIndependent(t *testing.T) {
	ss := Spans{Span(time.Minute)}

	strings := ss.Strings()
	durs := ss.Durs()

	strings[0] = "changed"
	durs[0] = time.Hour

	if ss[0] != Span(time.Minute) {
		t.Fatal("modifying conversion results changed the source")
	}

	ss[0] = Span(time.Second)
	if strings[0] != "changed" || durs[0] != time.Hour {
		t.Fatal("modifying the source changed conversion results")
	}
}

func FuzzSpanRoundTrip(f *testing.F) {
	for _, n := range []int64{
		0, 1, -1,
		int64(time.Minute),
		int64(time.Hour + time.Millisecond),
		1<<63 - 1,
		-1 << 63,
	} {
		f.Add(n)
	}

	f.Fuzz(func(t *testing.T, n int64) {
		want := Span(n)

		got, err := ParseSpan(want.String())
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("got %d ns, want %d ns", got, want)
		}
	})
}
