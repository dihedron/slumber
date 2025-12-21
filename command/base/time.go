package base

import "time"

// Duration represents a time duration value, which can be parsed from a string
// using the time.ParseDuration function.
type Duration time.Duration

// UnmarshalFlag unmarshals a string value into the Duration variable.
// This method is used by the go-flags package to handle custom flag types.
// It takes a string value, which is expected to be in a format that can be
// parsed by the time.ParseDuration function (e.g., "30m", "1h"), and
// populates the Duration variable accordingly.
func (d *Duration) UnmarshalFlag(value string) error {
	p, err := time.ParseDuration(value)
	if err == nil {
		*d = Duration(p)
	}
	return err
}
