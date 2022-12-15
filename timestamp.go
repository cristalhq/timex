package timex

import (
	"encoding/json"
	"strconv"
	"time"
)

type Timestamp struct {
	time.Time
}

// AsTimestamp returns time.Time wrapped as Timestamp.
func AsTimestamp(t time.Time) Timestamp {
	return Timestamp{Time: t}
}

// MarshalJSON implements json.Marshaler interface.
func (ts Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(ts.UnixMilli())
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (ts *Timestamp) UnmarshalJSON(b []byte) error {
	ms, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}

	*ts = AsTimestamp(time.UnixMilli(ms))
	return nil
}
