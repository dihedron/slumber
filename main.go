package main

import (
	"os"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	CPUProfile string         `long:"cpu-profile" description:"write cpu profile to file"`
	MemProfile string         `long:"mem-profile" description:"write memory profile to file"`
	Debug      bool           `long:"debug" description:"enable debug logging"`
	Monitor    MonitorCommand `command:"monitor" description:"Monitor for active editor sessions and power off if idle"`
}

type MonitorCommand struct {
	Timeout string `short:"t" long:"timeout" description:"Idle timeout before action (e.g. 30m, 1h)" default:"30m"`
	Action  string `short:"a" long:"action" description:"Action to take: hibernate or shutdown" choice:"hibernate" choice:"shutdown" default:"shutdown"`
}

var options Options
var parser = flags.NewParser(&options, flags.Default)

func main() {
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	cpuProfile = options.CPUProfile
	memProfile = options.MemProfile

	startProfiling()
	defer stopProfiling()

	// The monitor command will be executed by flags.Parse() if it implements the flags.Commander interface.
}
