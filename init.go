package main

import (
	"log/slog"
	"os"
	"runtime"
	"runtime/pprof"
)

var (
	cpuProfile string
	memProfile string
	cpuFile    *os.File
	memFile    *os.File
)

func init() {
	// Initialize logging
	level := slog.LevelInfo
	if os.Getenv("SLUMBER_DEBUG") != "" {
		level = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)
}

func startProfiling() {
	if cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			slog.Error("could not create CPU profile", "error", err)
			return
		}
		cpuFile = f
		if err := pprof.StartCPUProfile(f); err != nil {
			slog.Error("could not start CPU profile", "error", err)
			return
		}
		slog.Info("CPU profiling started", "file", cpuProfile)
	}
}

func stopProfiling() {
	if cpuFile != nil {
		pprof.StopCPUProfile()
		cpuFile.Close()
		slog.Info("CPU profiling stopped")
	}
	if memProfile != "" {
		f, err := os.Create(memProfile)
		if err != nil {
			slog.Error("could not create memory profile", "error", err)
			return
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			slog.Error("could not write memory profile", "error", err)
		}
		slog.Info("memory profile written", "file", memProfile)
	}
}
