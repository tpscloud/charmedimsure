// examples/charmed-bubbles/textarea/main.go
// Two inputs with focus switching. Only the focused component gets key
// presses; both get everything else (blink ticks, resizes).
package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	title textinput.Model
	body  textarea.Model
	focus int // 0 = title, 1 = body
}

func newModel(isDark bool) model {
	ti := textinput.New()
	ti.Placeholder = "title"
	ti.CharLimit = 80
	ti.SetStyles(textinput.DefaultStyles(isDark))

	ta := textarea.New()
	ta.Placeholder = "body (tab switches fields, ctrl+s submits)"
	ta.ShowLineNumbers = true
	ta.SetHeight(6)
	ta.SetStyles(textarea.DefaultStyles(isDark))

	return model{title: ti, body: ta}
}

func (m model) Init() tea.Cmd { return m.title.Focus() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.title.SetWidth(msg.Width - 4)
		m.body.SetWidth(msg.Width - 4)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "ctrl+s":
			return m, tea.Sequence(
				tea.Printf("title=%q body=%d chars", m.title.Value(), len(m.body.Value())),
				tea.Quit,
			)
		case "tab", "shift+tab":
			m.focus = 1 - m.focus
			if m.focus == 0 {
				m.body.Blur()
				return m, m.title.Focus()
			}
			m.title.Blur()
			return m, m.body.Focus()
		}
	}

	// Key presses go to the focused input only; everything else goes to both.
	_, isKey := msg.(tea.KeyPressMsg)
	var cmds []tea.Cmd
	var cmd tea.Cmd
	if !isKey || m.focus == 0 {
		m.title, cmd = m.title.Update(msg)
		cmds = append(cmds, cmd)
	}
	if !isKey || m.focus == 1 {
		m.body, cmd = m.body.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	label := lipgloss.NewStyle().Faint(true)
	return tea.NewView(label.Render("Title") + "\n" + m.title.View() + "\n\n" +
		label.Render("Body") + "\n" + m.body.View() + "\n")
}

func main() {
	if _, err := tea.NewProgram(newModel(true)).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
