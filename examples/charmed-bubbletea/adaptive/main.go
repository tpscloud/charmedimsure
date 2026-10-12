// examples/charmed-bubbletea/adaptive/main.go
// Ask the terminal for its background color, then build light or dark styles.
// Bubble Tea downsamples colors to the terminal's profile on output, so
// Lip Gloss styles here can use true color freely. The UI renders from the
// first frame with a default guess; some terminals never answer.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type styles struct {
	title lipgloss.Style
	body  lipgloss.Style
}

func newStyles(isDark bool) styles {
	ld := lipgloss.LightDark(isDark)
	return styles{
		title: lipgloss.NewStyle().Bold(true).Foreground(ld(lipgloss.Color("#5A2D82"), lipgloss.Color("#C4A7F5"))),
		body:  lipgloss.NewStyle().Foreground(ld(lipgloss.Color("#333333"), lipgloss.Color("#DDDDDD"))),
	}
}

type model struct {
	styles   styles
	isDark   bool
	detected bool
}

func newModel() model {
	// Start with a guess so the first frame is already styled; the reply
	// (if one comes) replaces it.
	return model{styles: newStyles(true), isDark: true}
}

func (m model) Init() tea.Cmd { return tea.RequestBackgroundColor }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		m.isDark, m.detected = msg.IsDark(), true
		m.styles = newStyles(m.isDark)
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	source := "assumed until the terminal answers"
	if m.detected {
		source = "reported by the terminal"
	}
	return tea.NewView(m.styles.title.Render("adaptive colors") + "\n" +
		m.styles.body.Render(fmt.Sprintf("dark background: %v (%s; q to quit)", m.isDark, source)) + "\n")
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
