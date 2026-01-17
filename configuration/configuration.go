package configuration

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/dihedron/rawdata"
	"github.com/dihedron/slumberd/command/base"
	"github.com/dihedron/slumberd/pointer"
)

type Configuration struct {
	Timeout   *base.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Frequency *base.Duration `json:"frequency,omitempty" yaml:"frequency,omitempty"`
	Packages  *string        `json:"packages,omitempty" yaml:"packages,omitempty"`
	Action    *string        `json:"action,omitempty" yaml:"action,omitempty"`
}

// UnmarshalFlag unmarshals a file path into the Configuration variable.
// This method is used by the go-flags package to handle custom flag types.
// It takes a string value, which is expected to be a file path, and
// populates the Configuration variable from the file, by detecting
// the file type from the file extension and applying the proper unmarshal.
func (c *Configuration) UnmarshalFlag(value string) error {
	err := rawdata.UnmarshalInto("@"+value, c)
	if err != nil {
		slog.Error("failed to unmarshal configuration from file", "file", value, "error", err)
		return fmt.Errorf("failed to unmarshal configuration from file %s: %w", value, err)
	}
	// fill missing values with defaults
	if c.Timeout == nil || *c.Timeout <= 0 {
		slog.Warn("no or invalid timeut specified, using default", "timeout", c.Timeout, "default", base.Duration(15*time.Minute))
		c.Timeout = pointer.To(base.Duration(15 * time.Minute))
	}
	if c.Frequency == nil || *c.Frequency <= 0 {
		slog.Warn("no or invalid frequency specified, using default", "frequency", c.Frequency, "default", base.Duration(1*time.Minute))
		c.Frequency = pointer.To(base.Duration(time.Minute))
	}
	if c.Packages == nil || *c.Packages == "" {
		slog.Warn("no or invalid packages specified, using default", "packages", c.Packages, "default", "/home/ubuntu/packages.yaml")
		c.Packages = pointer.To("/home/ubuntu/packages.yaml")
	}
	if c.Action == nil || *c.Action == "" {
		slog.Warn("no or invalid action specified, using default", "action", c.Action, "default", "shutdown")
		c.Action = pointer.To("shutdown")
	}
	timeout := time.Duration(*c.Timeout)
	frequency := time.Duration(*c.Frequency)
	if frequency > timeout {
		slog.Warn("frequency is greater than timeout, setting timeout to frequency", "frequency", frequency, "timeout", timeout)
		*c.Timeout = base.Duration(frequency)
	}

	return nil
}
