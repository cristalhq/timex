package timezone

import (
	"time"
	_ "time/tzdata"
)

// Must loads a named location and panics if loading fails.
// Use it for known names during initialization.
// For runtime input, use time.LoadLocation.
func Must(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}
