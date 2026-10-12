// examples/charmed-bubbles/tree/main.go
// An interactive tree (bubbles v2.2 and later): nodes open and close, the
// cursor moves, and the view scrolls. Built on lipgloss/tree for rendering.
package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/tree"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	tree tree.Model
}

func newModel(isDark bool) model {
	root := tree.Root("project").Child(
		tree.Root("src").Child("main.go", "config.go", "main_test.go"),
		tree.Root("docs").Child("usage.md").Close(), // starts collapsed
		tree.Root("assets").Child("logo.svg"),
	)
	t := tree.New(root, 0, 0)
	t.SetStyles(tree.DefaultStyles(isDark))
	t.SetOpenCharacter("▾")
	t.SetClosedCharacter("▸")
	return model{tree: t}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.tree.SetSize(msg.Width, msg.Height-2)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "space": // enter is the tree's own toggle key, so pick another
			n := m.tree.NodeAtCurrentOffset()
			return m, tea.Printf("selected %s", n.Value())
		}
	}
	var cmd tea.Cmd
	m.tree, cmd = m.tree.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	head := lipgloss.NewStyle().Bold(true).Render("j/k move, enter toggles, left/right close/open, space selects, q quits")
	return tea.NewView(head + "\n\n" + m.tree.View())
}

func main() {
	if _, err := tea.NewProgram(newModel(true)).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
