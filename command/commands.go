package command

import (
	"github.com/dihedron/slumber/command/api"
	"github.com/dihedron/slumber/command/hibernate"
	"github.com/dihedron/slumber/command/monitor"
	"github.com/dihedron/slumber/command/poweroff"
	"github.com/dihedron/slumber/command/version"
)

// Commands is the main container for all the commands of the application.
type Commands struct {
	// Monitor runs the Monitor command
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Monitor monitor.Monitor `command:"monitor" alias:"mon" alias:"m" description:"Monitor for active editor sessions and power off if idle"`
	// API runs the API server
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	API api.API `command:"api" alias:"a" description:"Start the gRPC/REST API server"`
	// Hibernate runs the Hibernate command
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Hibernate hibernate.Hibernate `command:"hibernate" alias:"hib" alias:"h" description:"Hibernate the system" hidden:"true"`
	// PowerOff runs the PowerOff command
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	PowerOff poweroff.PowerOff `command:"poweroff" alias:"powoff" alias:"p" description:"Power off the system" hidden:"true"`
	// Version prints the application version information and exits.
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Version version.Version `command:"version" alias:"ver" alias:"v" description:"Show the command version and exit"`
}
