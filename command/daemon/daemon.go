package daemon

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dihedron/slumberd/configuration"
	"github.com/dihedron/slumberd/internal/detect"
)

// Daemon is the daemon command.
type Daemon struct {
	// Configuration is the configuration for the daemon.
	Configuration *configuration.Configuration `short:"c" long:"configuration" description:"Configuration file" default:"/home/ubuntu/packages.yaml"`
}

// Execute runs the daemon command.
func (cmd *Daemon) Execute(args []string) error {
	slog.Info("starting daemon", "timeout", cmd.Configuration.Timeout.String(), "frequency", cmd.Configuration.Frequency.String(), "packages", *cmd.Configuration.Packages)

	timeout := time.Duration(*cmd.Configuration.Timeout)
	frequency := time.Duration(*cmd.Configuration.Frequency)

	// set up ticker to run every frequency and check for active editors
	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	// set up signal handling for graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	lastActive := time.Now()

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
				fmt.Printf("no active editor sessions... idle: %s\n", idleTime.String())
				if idleTime > timeout {
					slog.Warn("idle timeout reached, shutting down...")
					fmt.Println("shutting down...")
					//power.Shutdown()
					return nil
				}
			}
		}
	}
}
