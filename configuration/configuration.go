package configuration

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/dihedron/rawdata"
	"github.com/dihedron/slumberd/pointer"
	"github.com/dihedron/slumberd/timex"
)

type Configuration struct {
	Packages  *string         `json:"packages,omitempty" yaml:"packages,omitempty"`
	Debounce  *timex.Duration `json:"debounce,omitempty" yaml:"debounce,omitempty"`
	Timeout   *timex.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Frequency *timex.Duration `json:"frequency,omitempty" yaml:"frequency,omitempty"`
}

// UnmarshalFlag unmarshals a file path into the Configuration variable.
// This method is used by the go-flags package to handle custom flag types.
// It takes a string value, which is expected to be a file path, and
// populates the Configuration variable from the file, by detecting
// the file type from the file extension and applying the proper unmarshal.
func (c *Configuration) UnmarshalFlag(value string) error {
	slog.Info("unmarshalling configuration from file", "file", value)
	err := rawdata.UnmarshalInto("@"+value, c)
	if err != nil {
		slog.Error("failed to unmarshal configuration from file", "file", value, "error", err)
		return fmt.Errorf("failed to unmarshal configuration from file %s: %w", value, err)
	}
	// fill missing values with defaults
	if c.Packages == nil || *c.Packages == "" {
		slog.Warn("no or invalid packages path specified, using default", "path", c.Packages, "default", "/home/developer/packages.yaml")
		c.Packages = pointer.To("/home/developer/packages.yaml")
	}
	if c.Debounce == nil || *c.Debounce <= 0 {
		slog.Warn("no or invalid debounce specified, using default", "debounce", c.Debounce, "default", timex.Duration(500*time.Millisecond))
		c.Debounce = pointer.To(timex.Duration(500 * time.Millisecond))
	}
	if c.Timeout == nil || *c.Timeout <= 0 {
		slog.Warn("no or invalid timeout specified, using default", "timeout", c.Timeout, "default", timex.Duration(15*time.Minute))
		c.Timeout = pointer.To(timex.Duration(15 * time.Minute))
	}
	if c.Frequency == nil || *c.Frequency <= 0 {
		slog.Warn("no or invalid frequency specified, using default", "frequency", c.Frequency, "default", timex.Duration(1*time.Minute))
		c.Frequency = pointer.To(timex.Duration(time.Minute))
	}

	// check that the packages file exists and is readable
	if _, err := os.Stat(*c.Packages); err != nil {
		if os.IsNotExist(err) {
			slog.Error("packages file does not exist", "file", *c.Packages)
			return fmt.Errorf("packages file %s does not exist", *c.Packages)
		}
		slog.Error("error checking packages file", "file", *c.Packages, "error", err)
		return fmt.Errorf("error checking packages file %s: %w", *c.Packages, err)
	}
	f, err := os.Open(*c.Packages)
	if err != nil {
		slog.Error("packages file is not readable", "file", *c.Packages, "error", err)
		return fmt.Errorf("packages file %s is not readable: %w", *c.Packages, err)
	}
	f.Close()

	// check that activity monitor values are valid
	timeout := time.Duration(*c.Timeout)
	frequency := time.Duration(*c.Frequency)
	if frequency > timeout {
		slog.Warn("frequency is greater than timeout, setting timeout to frequency", "frequency", frequency, "timeout", timeout)
		*c.Timeout = timex.Duration(frequency)
	}

	return nil
}
