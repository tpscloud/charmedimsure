// examples/charmed-bubbletea/testing/main_test.go
package main

import (
	"bytes"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

// Unit style: call Update with constructed messages, inspect the model.
func TestCursorMovesAndSelects(t *testing.T) {
	var m tea.Model = model{items: []string{"a", "b", "c"}}

	m, _ = m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	got := m.(model)
	if got.chosen != "c" {
		t.Fatalf("chosen = %q, want c", got.chosen)
	}
	if cmd == nil {
		t.Fatal("enter should return tea.Quit")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("expected the returned command to produce a QuitMsg")
	}
}

// Integration style: run the whole program headless, feeding keystrokes as
// bytes and capturing the rendered output.
func TestProgramHeadless(t *testing.T) {
	in := strings.NewReader("j\r") // down, then enter (carriage return)
	var out bytes.Buffer

	p := tea.NewProgram(model{items: []string{"a", "b", "c"}},
		tea.WithInput(in),
		tea.WithOutput(&out),
		tea.WithWindowSize(40, 10),
		tea.WithoutSignals(),
	)
	final, err := p.Run()
	if err != nil {
		t.Fatal(err)
	}
	if got := final.(model).chosen; got != "b" {
		t.Fatalf("chosen = %q, want b", got)
	}
	if !strings.Contains(out.String(), "> b") {
		t.Fatalf("rendered output never showed the cursor on b:\n%s", out.String())
	}
}

// View is a pure function of the model, so it can be checked without a Program.
func TestViewMarksCursor(t *testing.T) {
	v := model{items: []string{"x", "y"}, cursor: 1}.View()
	if !strings.Contains(v.Content, "> y") {
		t.Fatalf("unexpected view:\n%s", v.Content)
	}
}
