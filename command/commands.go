package command

import (
	"github.com/dihedron/slumberd/command/daemon"
	"github.com/dihedron/slumberd/command/poweroff"
	"github.com/dihedron/slumberd/command/version"
)

// Commands is the main container for all the commands of the application.
type Commands struct {
	// Daemon runs the daemon command
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Daemon daemon.Daemon `command:"daemon" alias:"d" description:"Monitor for active editor sessions and power off if idle"`
	// PowerOff runs the PowerOff command
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	PowerOff poweroff.PowerOff `command:"poweroff" alias:"powoff" alias:"p" description:"Power off the system" hidden:"true"`
	// Version prints the application version information and exits.
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Version version.Version `command:"version" alias:"ver" alias:"v" description:"Show the command version and exit"`
}
