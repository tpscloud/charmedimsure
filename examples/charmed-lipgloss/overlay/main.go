// examples/charmed-lipgloss/overlay/main.go
// Layers and the Compositor draw one block on top of another, which is how
// you get a modal dialog over a full screen view in Bubble Tea v2. Hit
// testing tells you which layer a mouse click landed on.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	width, height int
	showDialog    bool
	last          string
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "space", "enter":
			m.showDialog = !m.showDialog
		}
	case tea.MouseClickMsg:
		mouse := msg.Mouse()
		hit := m.compositor().Hit(mouse.X, mouse.Y)
		m.last = fmt.Sprintf("clicked layer %q", hit.ID())
		if hit.ID() == "dialog" {
			m.showDialog = false
		}
	}
	return m, nil
}

func (m model) compositor() *lipgloss.Compositor {
	bg := lipgloss.NewStyle().Width(m.width).Height(m.height).
		Render(fmt.Sprintf("background layer\n\nspace toggles the dialog, q quits\n%s", m.last))
	base := lipgloss.NewLayer(bg).ID("background")

	if !m.showDialog {
		return lipgloss.NewCompositor(base)
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("212")).
		Padding(1, 3).
		Render("a dialog on its own layer\n\nclick it to close")
	w, h := lipgloss.Size(box)
	dialog := lipgloss.NewLayer(box).
		ID("dialog").
		X((m.width - w) / 2).
		Y((m.height - h) / 2).
		Z(1)

	return lipgloss.NewCompositor(base, dialog)
}

func (m model) View() tea.View {
	v := tea.NewView(m.compositor().Render())
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
