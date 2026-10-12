// examples/charmed-bubbletea/mouse/main.go
// Enable the mouse in View, then match the specific message types in Update.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	last   string
	clicks int
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseClickMsg:
		mouse := msg.Mouse()
		if mouse.Button == tea.MouseLeft {
			m.clicks++
		}
		m.last = fmt.Sprintf("click %s at %d,%d", mouse.Button, mouse.X, mouse.Y)
	case tea.MouseWheelMsg:
		m.last = "wheel " + msg.Mouse().Button.String()
	case tea.MouseMotionMsg:
		mouse := msg.Mouse()
		m.last = fmt.Sprintf("motion at %d,%d", mouse.X, mouse.Y)
	case tea.MouseReleaseMsg:
		m.last = "release"
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	v := tea.NewView(fmt.Sprintf("left clicks: %d\nlast event: %s\n\nq to quit\n", m.clicks, m.last))
	v.AltScreen = true
	// CellMotion reports clicks, wheel, and drags. AllMotion also reports
	// hover movement, which is a lot of messages; only ask for it if needed.
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
