// examples/charmed-bubbletea/ticker/main.go
// tea.Tick fires once. Return another Tick from Update to keep going.
package main

import (
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
)

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}

type model struct{ remaining int }

func (m model) Init() tea.Cmd { return tick() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.remaining--
		if m.remaining <= 0 {
			return m, tea.Quit
		}
		return m, tick()
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	return tea.NewView(fmt.Sprintf("Quitting in %d seconds. Press q to skip.\n", m.remaining))
}

func main() {
	if _, err := tea.NewProgram(model{remaining: 5}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
