package monitor

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dihedron/slumber/command/base"
	"github.com/dihedron/slumber/internal/detect"
)

type Monitor struct {
	// Timeout is the idle timeout before action (e.g. 5m, 1h).
	Timeout base.Duration `short:"t" long:"timeout" description:"Idle timeout before action (e.g. 5m, 1h)" default:"5m"`
	// Frequency is the frequency of checks (e.g. 15s, 1m).
	Frequency base.Duration `short:"f" long:"frequency" description:"Frequency of checks (e.g. 15s, 1m)" default:"15s"`
	// Action is the action to take when the timeout is reached (hibernate or shutdown).
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Action string `short:"a" long:"action" description:"Action to take: hibernate or shutdown" choice:"hibernate" choice:"shutdown" default:"shutdown"`
}

// Execute runs the monitor command.
func (cmd *Monitor) Execute(args []string) error {
	slog.Info("starting monitor", "timeout", cmd.Timeout, "action", cmd.Action)

	timeout := time.Duration(cmd.Timeout)
	frequency := time.Duration(cmd.Frequency)
	if frequency > timeout {
		slog.Warn("frequency is greater than timeout, setting timeout to frequency", "frequency", frequency, "timeout", timeout)
		timeout = frequency
	}

	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	lastActive := time.Now()

	// set up signal handling for graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-signals:
			slog.Info("received termination signal, shutting down")
			fmt.Println("received termination signal, shutting down...")
			return nil
		case <-ticker.C:
			editors := detect.IsAnyEditorActive2("/proc")
			if editors {
				slog.Info("editor sessions active")
				fmt.Println("editor sessions active...")
				lastActive = time.Now()
			} else {
				idleTime := time.Since(lastActive)
				slog.Info("no active editor sessions", "idle", idleTime.String())
				fmt.Println("no active editor sessions...")
				if idleTime > timeout {
					slog.Warn("idle timeout reached, taking action", "action", cmd.Action)
					fmt.Printf("idle timeout reached, taking action: %s", cmd.Action)
					if cmd.Action == "hibernate" {
						slog.Info("hibernating")
						fmt.Println("hibernating...")
						//power.Hibernate()
					} else {
						slog.Info("shutting down")
						fmt.Println("shutting down...")
						//power.Shutdown()
						return nil
					}
				}
			}
		}
	}
}
