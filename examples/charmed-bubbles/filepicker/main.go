// examples/charmed-bubbles/filepicker/main.go
// Browse directories and pick a file. Init reads the starting directory, so
// the parent must return it.
package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/filepicker"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	fp       filepicker.Model
	selected string
	err      string
}

func newModel() model {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".go", ".md", ".zip"}
	fp.CurrentDirectory, _ = os.Getwd()
	fp.ShowHidden = false
	fp.ShowSize = true
	fp.SetHeight(15)
	return model{fp: fp}
}

func (m model) Init() tea.Cmd { return m.fp.Init() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.fp, cmd = m.fp.Update(msg)

	// Ask the picker whether this message completed a selection.
	if ok, path := m.fp.DidSelectFile(msg); ok {
		m.selected = path
		return m, tea.Quit
	}
	if ok, path := m.fp.DidSelectDisabledFile(msg); ok {
		m.err = path + " is not an allowed type"
	}
	return m, cmd
}

func (m model) View() tea.View {
	s := "Pick a file (q quits)\n"
	if m.err != "" {
		s += m.err + "\n"
	}
	return tea.NewView(s + "\n" + m.fp.View())
}

func main() {
	final, err := tea.NewProgram(newModel()).Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if m, ok := final.(model); ok && m.selected != "" {
		fmt.Println("selected", m.selected)
	}
}
