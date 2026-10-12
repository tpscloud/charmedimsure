// examples/charmed-log/structured/main.go
// Formatters, styles, and levels from configuration.
package main

import (
	"os"

	"charm.land/lipgloss/v2"
	"charm.land/log/v2"
)

func main() {
	// Machine readable output for files or pipes.
	jsonLogger := log.NewWithOptions(os.Stdout, log.Options{Formatter: log.JSONFormatter, ReportTimestamp: true})
	jsonLogger.Info("json line", "project", "api", "ok", true)

	fmtLogger := log.New(os.Stdout)
	fmtLogger.SetFormatter(log.LogfmtFormatter)
	fmtLogger.Info("logfmt line", "project", "docs", "files", 12)

	// Level from an environment variable or flag.
	lvl, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		lvl = log.InfoLevel
	}
	log.SetLevel(lvl)

	// Restyle the text formatter. Styles are Lip Gloss v2 values.
	styles := log.DefaultStyles()
	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("ERROR").
		Padding(0, 1).
		Background(lipgloss.Color("204")).
		Foreground(lipgloss.Color("0"))
	styles.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	styles.Values["err"] = lipgloss.NewStyle().Bold(true)

	styled := log.New(os.Stderr)
	styled.SetStyles(styles)
	styled.Error("styled error", "err", "checksum mismatch")
}
