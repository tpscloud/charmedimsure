// examples/charmed-bubbletea/subscribe/main.go
// Two ways to feed a program from outside: a channel read inside a Cmd
// (re-armed after every message) and Program.Send from any goroutine.
package main

import (
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
)

type eventMsg string

// waitFor returns a Cmd that blocks on the channel. Each delivery is one
// message, so Update must return waitFor again to keep listening.
func waitFor(ch <-chan string) tea.Cmd {
	return func() tea.Msg {
		s, ok := <-ch
		if !ok {
			return tea.Quit()
		}
		return eventMsg(s)
	}
}

type model struct {
	events <-chan string
	log    []string
}

func (m model) Init() tea.Cmd { return waitFor(m.events) }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case eventMsg:
		m.log = append(m.log, string(msg))
		return m, waitFor(m.events)
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	s := "events (q to quit):\n"
	for _, e := range m.log {
		s += "  " + e + "\n"
	}
	return tea.NewView(s)
}

func main() {
	ch := make(chan string)
	p := tea.NewProgram(model{events: ch})

	go func() {
		for i := 1; i <= 3; i++ {
			time.Sleep(700 * time.Millisecond)
			ch <- fmt.Sprintf("channel event %d", i)
		}
		// Program.Send is safe from any goroutine and needs no channel.
		p.Send(eventMsg("sent directly with p.Send"))
		time.Sleep(700 * time.Millisecond)
		close(ch) // waitFor turns a closed channel into tea.Quit
	}()

	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
