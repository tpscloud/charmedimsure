// examples/charmed-bubbles/timer/main.go
// Countdown timer and stopwatch side by side. Both are ID based, so several
// can run in one model without stealing each other's ticks.
package main

import (
	"fmt"
	"os"
	"time"

	"charm.land/bubbles/v2/stopwatch"
	"charm.land/bubbles/v2/timer"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	timer     timer.Model
	stopwatch stopwatch.Model
	expired   bool
}

func newModel() model {
	return model{
		timer:     timer.New(10*time.Second, timer.WithInterval(time.Second)),
		stopwatch: stopwatch.New(stopwatch.WithInterval(100 * time.Millisecond)),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.timer.Init(), m.stopwatch.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "space":
			return m, tea.Batch(m.timer.Toggle(), m.stopwatch.Toggle())
		case "r":
			return m, m.stopwatch.Reset()
		}
	case timer.TimeoutMsg:
		m.expired = true
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.timer, cmd = m.timer.Update(msg)
	cmds = append(cmds, cmd)
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	s := fmt.Sprintf("timer:     %s\nstopwatch: %s\n\nspace pauses both, r resets the stopwatch, q quits\n",
		m.timer.View(), m.stopwatch.View())
	if m.expired {
		s += "timer expired\n"
	}
	return tea.NewView(s)
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
