// examples/charmed-bubbles/viewport/main.go
// A pager: scrollable content sized to the window, mouse wheel enabled.
package main

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	vp    viewport.Model
	ready bool
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := 2
		if !m.ready {
			m.vp = viewport.New(viewport.WithWidth(msg.Width), viewport.WithHeight(msg.Height-headerHeight))
			m.vp.MouseWheelEnabled = true
			m.vp.SoftWrap = true
			m.vp.SetContent(strings.Repeat("Scrollable content line. ", 400))
			m.ready = true
		} else {
			m.vp.SetWidth(msg.Width)
			m.vp.SetHeight(msg.Height - headerHeight)
		}
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.vp, cmd = m.vp.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	if !m.ready {
		return tea.NewView("loading...")
	}
	header := lipgloss.NewStyle().Bold(true).
		Render(fmt.Sprintf("pager  %3.0f%%  (j/k, pgup/pgdn, wheel; q quits)", m.vp.ScrollPercent()*100))
	v := tea.NewView(header + "\n\n" + m.vp.View())
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion // required for MouseWheelEnabled to matter
	return v
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
