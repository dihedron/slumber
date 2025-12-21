package poweroff

import (
	"log/slog"

	"github.com/dihedron/slumber/command/base"
	"github.com/dihedron/slumber/internal/power"
)

type PowerOff struct {
	// Timeout is the idle timeout before action (e.g. 5m, 1h).
	Timeout base.Duration `short:"t" long:"timeout" description:"Idle timeout before action (e.g. 5m, 1h)" default:"10s"`
}

// Execute runs the shutdown command.
func (c *PowerOff) Execute(args []string) error {
	slog.Info("executing poweroff")
	return power.Shutdown()
}
