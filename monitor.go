package main

import (
	"log/slog"
	"time"

	"github.com/dihedron/slumber/internal/detect"
	"github.com/dihedron/slumber/internal/power"
)

func (c *MonitorCommand) Execute(args []string) error {
	slog.Info("starting monitor", "timeout", c.Timeout, "action", c.Action)

	duration, err := time.ParseDuration(c.Timeout)
	if err != nil {
		slog.Error("invalid timeout", "timeout", c.Timeout, "error", err)
		return err
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	lastActive := time.Now()

	for range ticker.C {
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
	return nil
}
