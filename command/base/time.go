package base

import (
	"encoding/json"
	"time"

	"gopkg.in/yaml.v3"
)

// Duration represents a time duration value, which can be parsed from a string
// using the time.ParseDuration function.
type Duration time.Duration

// String returns the string representation of the Duration value.
func (d *Duration) String() string {
	return time.Duration(*d).String()
}

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

// MarshalJSON marshals the Duration value into a JSON string.
func (d *Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(*d).String())
}

// UnmarshalJSON unmarshals a JSON string into the Duration variable.
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	p, err := time.ParseDuration(v)
	if err == nil {
		*d = Duration(p)
	}
	return err
}

// MarshalText marshals the Duration value into a text string.
func (d *Duration) MarshalText() ([]byte, error) {
	return []byte(time.Duration(*d).String()), nil
}

// UnmarshalText unmarshals a text string into the Duration variable.
func (d *Duration) UnmarshalText(text []byte) error {
	v := string(text)
	p, err := time.ParseDuration(v)
	if err == nil {
		*d = Duration(p)
	}
	return err
}

// MarshalYAML marshals the Duration value into a YAML string.
func (d Duration) MarshalYAML() (any, error) {
	return time.Duration(d).String(), nil
}

// UnmarshalYAML unmarshals a YAML string into the Duration variable.
func (d *Duration) UnmarshalYAML(value *yaml.Node) error {
	var v string
	if err := value.Decode(&v); err != nil {
		return err
	}
	p, err := time.ParseDuration(v)
	if err == nil {
		*d = Duration(p)
	}
	return err
}
