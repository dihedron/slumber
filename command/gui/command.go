package gui

import (
	"fmt"
	"log/slog"
	"time"

	"fyne.io/fyne/v2/app"
)

// GUICommand is the command to start the GUI.
type GUICommand struct {
	APIAddress      string `short:"a" long:"api-address" description:"The address of the API server" default:"localhost:9090"`
	Userid          string `short:"u" long:"userid" description:"The user ID to control"`
	RefreshInterval int    `short:"r" long:"refresh" description:"The refresh interval in seconds" default:"120"`
}

// Execute is the main entry point for the gui command.
func (cmd *GUICommand) Execute(args []string) error {
	slog.Info("starting GUI", "api_address", cmd.APIAddress, "userid", cmd.Userid)

	if cmd.Userid == "" {
		return fmt.Errorf("userid is required")
	}

	client, err := NewClient(cmd.APIAddress, cmd.Userid)
	if err != nil {
		return err
	}
	defer client.Close()

	myApp := app.New()
	myWindow := myApp.NewWindow("VM Control - " + cmd.Userid)

	g := NewGUI(myWindow, client)
	g.StartPoller(time.Duration(cmd.RefreshInterval) * time.Second)

	myWindow.ShowAndRun()
	return nil
}
