// examples/charmed-log/tea/main.go
// Logging inside a Bubble Tea program: the terminal is taken, so log to a
// file. tea.LogToFileWith accepts the Charm logger.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"github.com/charmbracelet/colorprofile"
)

type model struct{ presses int }

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyPressMsg); ok {
		log.Debug("key", "key", k.String())
		m.presses++
		if k.String() == "q" || k.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	return tea.NewView(fmt.Sprintf("presses: %d (see debug.log; q quits)\n", m.presses))
}

func main() {
	f, err := tea.LogToFileWith("debug.log", "tui", log.Default())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()
	// The default logger detected its color profile for stderr. SetOutput
	// does not re-detect it, so match the profile to the file (plain text).
	log.SetColorProfile(colorprofile.Detect(f, os.Environ()))
	log.SetLevel(log.DebugLevel)
	log.SetReportTimestamp(true)

	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
