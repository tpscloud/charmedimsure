// examples/charmed-huh/embedded/main.go
// A huh form as a child of a Bubble Tea v2 model. huh's Update returns its
// own Model interface and its View returns a string, so wrap accordingly.
// The result is printed after Run returns: tea.Printf prints nothing while
// the alternate screen is active.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

// minFormWidth keeps the form usable when the terminal is very narrow.
const minFormWidth = 20

type model struct {
	form   *huh.Form
	styles styles
}

type styles struct {
	frame lipgloss.Style
}

func newModel() model {
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("project").
				Title("Choose your project").
				Options(huh.NewOptions("site", "api", "docs")...),
			huh.NewSelect[int]().
				Key("priority").
				Title("Choose your priority").
				Options(huh.NewOptions(1, 2, 3)...),
		),
	).WithShowHelp(false)

	return model{
		form:   f,
		styles: styles{frame: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2)},
	}
}

func (m model) Init() tea.Cmd { return m.form.Init() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Size the form to what is left inside the frame's border and padding.
		fw, _ := m.styles.frame.GetFrameSize()
		m.form = m.form.WithWidth(max(msg.Width-fw, minFormWidth))
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	// Form.Update returns huh.Model (an interface); assert back to *huh.Form.
	next, cmd := m.form.Update(msg)
	if f, ok := next.(*huh.Form); ok {
		m.form = f
	}

	switch m.form.State {
	case huh.StateCompleted, huh.StateAborted:
		return m, tea.Quit
	}
	return m, cmd
}

func (m model) View() tea.View {
	// huh renders to a string; the outer program owns the tea.View.
	v := tea.NewView(m.styles.frame.Render(m.form.View()))
	v.AltScreen = true
	return v
}

func main() {
	final, err := tea.NewProgram(newModel()).Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// Read the answers from the final model once the terminal is restored.
	if m, ok := final.(model); ok && m.form.State == huh.StateCompleted {
		fmt.Printf("You picked %s at priority %d\n", m.form.GetString("project"), m.form.GetInt("priority"))
	}
}
