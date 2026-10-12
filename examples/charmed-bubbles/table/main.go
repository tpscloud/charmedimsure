// examples/charmed-bubbles/table/main.go
// A scrollable, selectable table. For a static one use lipgloss/table.
package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	table table.Model
}

func newModel() model {
	columns := []table.Column{
		{Title: "Project", Width: 12},
		{Title: "Status", Width: 10},
	}
	rows := []table.Row{
		{"site", "ready"},
		{"api", "building"},
		{"docs", "ready"},
		{"worker", "queued"},
		{"tests", "running"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true), // focused tables respond to keys
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.Bold(true).BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).BorderForeground(lipgloss.Color("240"))
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(false)
	t.SetStyles(s)

	return model{table: t}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			row := m.table.SelectedRow()
			return m, tea.Printf("selected %s (%s)", row[0], row[1])
		}
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	frame := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240"))
	return tea.NewView(frame.Render(m.table.View()) + "\n" + m.table.HelpView() + "\n")
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
