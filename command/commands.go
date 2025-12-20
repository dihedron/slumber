package command

import (
	"github.com/dihedron/slumber/command/monitor"
	"github.com/dihedron/slumber/command/version"
)

// Commands is the main container for all the commands of the application.
type Commands struct {
	// Generate runs the Generate command
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Monitor monitor.Monitor `command:"monitor" alias:"mon" alias:"m" description:"Monitor for active editor sessions and power off if idle"`
	// Version prints the application version information and exits.
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Version version.Version `command:"version" alias:"ver" alias:"v" description:"Show the command version and exit"`
}
