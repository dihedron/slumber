package monitor

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dihedron/slumber/internal/detect"
	"github.com/dihedron/slumber/internal/power"
)

type Monitor struct {
	Timeout string `short:"t" long:"timeout" description:"Idle timeout before action (e.g. 30m, 1h)" default:"30m"`
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Action string `short:"a" long:"action" description:"Action to take: hibernate or shutdown" choice:"hibernate" choice:"shutdown" default:"shutdown"`
}

// Execute runs the monitor command.
func (c *Monitor) Execute(args []string) error {
	slog.Info("starting monitor", "timeout", c.Timeout, "action", c.Action)

	duration, err := time.ParseDuration(c.Timeout)
	if err != nil {
		slog.Error("invalid timeout", "timeout", c.Timeout, "error", err)
		return err
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	lastActive := time.Now()

	// set up signal handling for graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-signals:
			slog.Info("received termination signal, shutting down")
			return nil
		case <-ticker.C:
			editors := detect.IsAnyEditorActive("/proc")
			if len(editors) > 0 {
				slog.Debug("editor sessions active", "editors", editors)
				lastActive = time.Now()
			} else {
				idleTime := time.Since(lastActive)
				slog.Info("no active editor sessions", "idle", idleTime.String())
				if idleTime > duration {
					slog.Warn("idle timeout reached, taking action", "action", c.Action)
					if c.Action == "hibernate" {
						return power.Hibernate()
					} else {
						return power.Shutdown()
					}
				}
			}
		}
	}
}
