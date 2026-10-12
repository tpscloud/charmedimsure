// examples/charmed-bubbletea/testing/main.go
// A tiny model used by main_test.go. Tests drive Update directly and also run
// the whole Program headless with WithInput, WithOutput, and WithWindowSize.
package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	items  []string
	cursor int
	chosen string
	width  int
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tea.KeyPressMsg:
		switch msg.String() {
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "enter":
			m.chosen = m.items[m.cursor]
			return m, tea.Quit
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	var b strings.Builder
	for i, it := range m.items {
		prefix := "  "
		if i == m.cursor {
			prefix = "> "
		}
		fmt.Fprintf(&b, "%s%s\n", prefix, it)
	}
	return tea.NewView(b.String())
}

func main() {
	final, err := tea.NewProgram(model{items: []string{"one", "two", "three"}}).Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if m, ok := final.(model); ok && m.chosen != "" {
		fmt.Println("chose", m.chosen)
	}
}
