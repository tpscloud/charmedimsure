// examples/charmed-bubbles/list/main.go
// A filterable list with the default delegate. Items are your own type.
package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// item satisfies list.DefaultItem: Title and Description for the delegate,
// FilterValue for the built in fuzzy filter.
type item struct {
	name, path, status string
}

func (i item) Title() string { return i.name }
func (i item) Description() string {
	return fmt.Sprintf("%s (%s)", i.path, i.status)
}
func (i item) FilterValue() string { return i.name }

type model struct {
	list   list.Model
	chosen string
}

func newModel(isDark bool) model {
	items := []list.Item{
		item{"site", "./site", "ready"},
		item{"api", "./api", "building"},
		item{"docs", "./docs", "ready"},
		item{"worker", "./worker", "queued"},
		item{"cli", "./cli", "ready"},
		item{"tests", "./tests", "running"},
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles = list.NewDefaultItemStyles(isDark)
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("212"))

	l := list.New(items, delegate, 0, 0) // size comes from WindowSizeMsg
	l.Title = "Projects"
	l.Styles = list.DefaultStyles(isDark)
	l.SetStatusBarItemName("project", "projects")
	l.DisableQuitKeybindings() // we handle quitting ourselves
	return model{list: l}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	case tea.KeyPressMsg:
		// While the user is typing a filter, every key belongs to the list.
		if m.list.SettingFilter() {
			break
		}
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if it, ok := m.list.SelectedItem().(item); ok {
				m.chosen = it.name
				return m, tea.Quit
			}
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	v := tea.NewView(m.list.View())
	v.AltScreen = true
	return v
}

func main() {
	final, err := tea.NewProgram(newModel(true)).Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if m, ok := final.(model); ok && m.chosen != "" {
		fmt.Println("chose", m.chosen)
	}
}
