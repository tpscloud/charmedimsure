// examples/charmed-bubbletea/editor/main.go
// Hand the terminal to another interactive program (an editor) and resume.
package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "charm.land/bubbletea/v2"
)

type editorDoneMsg struct{ err error }

func openEditor(path string) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	c := exec.Command(editor, path)
	return tea.ExecProcess(c, func(err error) tea.Msg { return editorDoneMsg{err} })
}

type model struct {
	path string
	err  error
	runs int
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case editorDoneMsg:
		m.err = msg.err
		m.runs++
	case tea.KeyPressMsg:
		switch msg.String() {
		case "e":
			return m, openEditor(m.path)
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	s := fmt.Sprintf("press e to edit %s, q to quit\neditor runs: %d\n", m.path, m.runs)
	if m.err != nil {
		s += "last error: " + m.err.Error() + "\n"
	}
	return tea.NewView(s)
}

func main() {
	f, err := os.CreateTemp("", "tea-editor-*.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	f.Close()
	defer os.Remove(f.Name())

	if _, err := tea.NewProgram(model{path: f.Name()}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
