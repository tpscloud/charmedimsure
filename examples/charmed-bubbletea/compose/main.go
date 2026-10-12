// examples/charmed-bubbletea/compose/main.go
// A parent model owns child bubbles, routes messages, and collects their Cmds.
// Uses key.Binding for bindings and help for the footer.
package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type keyMap struct {
	Submit key.Binding
	Quit   key.Binding
}

// keyMap satisfies help.KeyMap so help.Model can render it.
func (k keyMap) ShortHelp() []key.Binding { return []key.Binding{k.Submit, k.Quit} }
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Submit, k.Quit}}
}

var keys = keyMap{
	Submit: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit")),
	Quit:   key.NewBinding(key.WithKeys("ctrl+c", "esc"), key.WithHelp("esc", "quit")),
}

type model struct {
	input     textinput.Model
	spinner   spinner.Model
	help      help.Model
	submitted string
}

func newModel() model {
	in := textinput.New()
	in.Placeholder = "type something"
	in.SetVirtualCursor(false) // we will place the real terminal cursor ourselves

	sp := spinner.New(
		spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("205"))),
	)
	return model{input: in, spinner: sp, help: help.New()}
}

func (m model) Init() tea.Cmd {
	// Batch runs both: focus the input (starts cursor blink) and start spinning.
	return tea.Batch(m.input.Focus(), m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.input.SetWidth(msg.Width - 4)
		m.help.SetWidth(msg.Width)
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Submit):
			m.submitted = m.input.Value()
			m.input.Reset()
		}
	}

	// Children always get the message. The spinner needs its TickMsg even when
	// the input is focused, and the input ignores messages it does not care about.
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	body := fmt.Sprintf("%s working\n\n%s\n\nlast submitted: %q\n\n%s\n",
		m.spinner.View(), m.input.View(), m.submitted, m.help.View(keys))
	v := tea.NewView(body)
	// The input sits on row 2 of our layout, so offset its cursor by that much.
	if c := m.input.Cursor(); c != nil {
		c.Y += 2
		v.Cursor = c
	}
	return v
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
