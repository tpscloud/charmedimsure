// examples/charmed-log/basic/main.go
// The package level logger and a custom one with options.
package main

import (
	"errors"
	"os"
	"time"

	"charm.land/log/v2"
)

func main() {
	// Package level logger: writes to stderr, info level, timestamps on.
	log.Info("starting", "project", "site")
	log.Debug("hidden at info level")
	log.Warn("disk almost full", "free", "1.2G")
	log.Error("check failed", "err", errors.New("invalid configuration"))
	log.Print("no level prefix, always printed")

	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.Debug("now visible, with caller")

	// A dedicated logger with options.
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
		Prefix:          "tasks",
		Level:           log.InfoLevel,
	})
	logger.Info("checking", "file", "config.json")
	logger.Infof("%d of %d done", 3, 6)

	// With returns a child that carries the pairs on every line.
	projectLog := logger.With("project", "site", "task", "build")
	projectLog.Info("completed")
	projectLog.Warn("output directory already exists")

	// Fatal logs then calls os.Exit(1); keep it for main only.
	if len(os.Args) > 1 && os.Args[1] == "boom" {
		logger.Fatal("giving up", "reason", "requested")
	}
}
