// examples/charmed-bubbles/spinner/main.go
// The smallest complete bubble program: one spinner, started from Init,
// fed every message, and stored back into the model.
package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	spinner spinner.Model
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick // a method value is a tea.Cmd
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyPressMsg); ok && (msg.String() == "q" || msg.String() == "ctrl+c") {
		return m, tea.Quit
	}
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg) // store it back; bubbles are values
	return m, cmd
}

func (m model) View() tea.View {
	return tea.NewView(m.spinner.View() + " working, q to quit\n")
}

func main() {
	m := model{spinner: spinner.New(spinner.WithSpinner(spinner.Dot))}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
