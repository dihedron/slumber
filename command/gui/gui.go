package gui

import (
	"context"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type GUI struct {
	Window          fyne.Window
	Play            *widget.Button
	Pause           *widget.Button
	Stop            *widget.Button
	Heart           *widget.Icon
	StatusIndicator *canvas.Circle
	Client          *Client
}

func NewGUI(win fyne.Window, client *Client) *GUI {
	g := &GUI{
		Window: win,
		Client: client,
	}

	g.Play = widget.NewButtonWithIcon("Play", theme.MediaPlayIcon(), func() { g.onPlay() })
	g.Pause = widget.NewButtonWithIcon("Pause", theme.MediaPauseIcon(), func() { g.onPause() })
	g.Stop = widget.NewButtonWithIcon("Stop", theme.MediaStopIcon(), func() { g.onStop() })

	g.Heart = widget.NewIcon(theme.InfoIcon())
	g.StatusIndicator = canvas.NewCircle(color.Transparent)
	g.StatusIndicator.Resize(fyne.NewSize(32, 32))

	heartStack := container.NewStack(
		g.StatusIndicator,
		container.NewCenter(g.Heart),
	)

	content := container.NewHBox(
		g.Play,
		g.Pause,
		g.Stop,
		container.NewCenter(heartStack),
	)

	win.SetContent(content)
	return g
}

func (g *GUI) UpdateStatus(status string) {
	switch status {
	case "ACTIVE":
		g.Play.Disable()
		g.Pause.Enable()
		g.Stop.Enable()
		g.StatusIndicator.FillColor = color.RGBA{R: 0, G: 255, B: 0, A: 128} // Green
	case "SHUTOFF":
		g.Play.Enable()
		g.Pause.Disable()
		g.Stop.Disable()
		g.StatusIndicator.FillColor = color.RGBA{R: 255, G: 0, B: 0, A: 128} // Red
	case "PAUSED":
		g.Play.Enable()
		g.Pause.Disable()
		g.Stop.Enable()
		g.StatusIndicator.FillColor = color.RGBA{R: 255, G: 165, B: 0, A: 128} // Orange
	default:
		g.StatusIndicator.FillColor = color.Transparent
	}
	canvas.Refresh(g.StatusIndicator)
}

func (g *GUI) onPlay() {
	ctx := context.Background()
	status, err := g.Client.Status(ctx)
	if err != nil {
		return
	}
	if status == "SHUTOFF" {
		g.Client.Start(ctx)
	} else if status == "PAUSED" {
		g.Client.Unpause(ctx)
	}
}

func (g *GUI) onPause() {
	g.Client.Pause(context.Background())
}

func (g *GUI) onStop() {
	g.Client.Stop(context.Background())
}

func (g *GUI) AnimateHeart() {
	ticker := time.NewTicker(500 * time.Millisecond)
	growing := true
	go func() {
		for range ticker.C {
			if growing {
				g.Heart.Resize(fyne.NewSize(30, 30))
			} else {
				g.Heart.Resize(fyne.NewSize(24, 24))
			}
			growing = !growing
			canvas.Refresh(g.Heart)
		}
	}()
}
