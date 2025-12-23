package gui

import (
	"context"
	"log/slog"
	"time"
)

func (g *GUI) StartPoller(refreshInterval time.Duration) {
	// 1. Animate heart
	g.AnimateHeart()

	// 2. Poll status every refreshInterval
	go func() {
		ticker := time.NewTicker(refreshInterval)
		defer ticker.Stop()

		// Initial check
		status, err := g.Client.Status(context.Background())
		if err != nil {
			slog.Error("failed to get initial status", "error", err)
		} else {
			g.UpdateStatus(status)
		}

		for range ticker.C {
			status, err := g.Client.Status(context.Background())
			if err != nil {
				slog.Error("failed to poll status", "error", err)
				continue
			}
			g.UpdateStatus(status)
		}
	}()
}
