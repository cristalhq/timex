package timex

import (
	"encoding/json"
	"testing"
)

func TestTimestampJSONString(t *testing.T) {
	jsn := `{
		"at":"1685725500000",
		"ts":1685725500000
	}`

	var obj struct {
		At Timestamp `json:"at"`
		Ts Timestamp `json:"ts"`
	}
	if err := json.Unmarshal([]byte(jsn), &obj); err != nil {
		t.Error(err)
	}

	if obj.At != TimestampFromMilli(1685725500000) {
		t.Fatalf("have %v, want %v", obj.At, 1685725500000)
	}
	if obj.Ts != TimestampFromMilli(1685725500000) {
		t.Fatalf("have %v, want %v", obj.Ts, 1685725500000)
	}
}
