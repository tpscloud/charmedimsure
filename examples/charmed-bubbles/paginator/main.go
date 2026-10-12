// examples/charmed-bubbles/paginator/main.go
// Page through a slice with the paginator computing the visible bounds.
package main

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/paginator"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	items []string
	pager paginator.Model
}

func newModel() model {
	items := make([]string, 0, 37)
	for i := 1; i <= 37; i++ {
		items = append(items, fmt.Sprintf("item %02d", i))
	}
	p := paginator.New(paginator.WithPerPage(8))
	p.Type = paginator.Dots
	p.SetTotalPages(len(items))
	return model{items: items, pager: p}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyPressMsg); ok && (k.String() == "q" || k.String() == "ctrl+c") {
		return m, tea.Quit
	}
	var cmd tea.Cmd
	m.pager, cmd = m.pager.Update(msg) // left/right, h/l, pgup/pgdn
	return m, cmd
}

func (m model) View() tea.View {
	var b strings.Builder
	start, end := m.pager.GetSliceBounds(len(m.items))
	for _, it := range m.items[start:end] {
		b.WriteString("  " + it + "\n")
	}
	b.WriteString("\n  " + m.pager.View() + "\n  left/right to page, q quits\n")
	return tea.NewView(b.String())
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
