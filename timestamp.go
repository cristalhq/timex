package timex

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

type Timestamp struct {
	time.Time
}

// TimestampFromString returns Timestamp fomr the string representing Unix time.
func TimestampFromString(s string) (Timestamp, error) {
	msec, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return Timestamp{}, err
	}
	return TimestampFromMilli(msec), nil
}

// TimestampFromMilli returns Timestamp from the Unix time in milliseconds.
func TimestampFromMilli(msec int64) Timestamp {
	return AsTimestamp(time.UnixMilli(msec))
}

// AsTimestamp returns time.Time wrapped as Timestamp.
func AsTimestamp(t time.Time) Timestamp {
	return Timestamp{Time: t.UTC()}
}

// MarshalJSON implements json.Marshaler interface.
func (ts Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(ts.UnixMilli())
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (ts *Timestamp) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return errEmptyInput
	}

	var s string
	if b[0] == '"' && b[len(b)-1] == '"' {
		s = string(b[1 : len(b)-1])
	} else {
		s = string(b)
	}

	msec, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	*ts = AsTimestamp(time.UnixMilli(msec))
	return nil
}

var errEmptyInput = errors.New("empty timestamp")
