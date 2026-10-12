// examples/charmed-log/slog/main.go
// The logger is a slog.Handler and can back the standard library log too.
package main

import (
	"log/slog"
	"os"

	"charm.land/log/v2"
)

func main() {
	handler := log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		Level:           log.DebugLevel,
	})

	// Hand it to slog. Code that receives sl logs through Charm; package level
	// slog.Info still goes to slog's own default handler. To redirect the
	// whole process (a decision for main, not a library), call
	// slog.SetDefault(sl), which also captures the standard library log package.
	sl := slog.New(handler)
	sl.Info("from slog", "project", "worker")
	sl.With("group", "checks").Debug("child logger")
	slog.Info("package level slog, not routed through Charm")

	// Or the standard library logger, for code that takes *log.Logger.
	std := handler.StandardLog(log.StandardLogOptions{ForceLevel: log.WarnLevel})
	std.Println("from the standard logger, forced to warn")
}
