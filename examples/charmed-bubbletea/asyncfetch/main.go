// examples/charmed-bubbletea/asyncfetch/main.go
// I/O belongs in a Cmd. The Cmd returns a Msg; Update reacts to it.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
)

type statusMsg int
type errMsg struct{ err error }

func checkServer(url string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return errMsg{err}
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return errMsg{err}
		}
		defer res.Body.Close()
		return statusMsg(res.StatusCode)
	}
}

type model struct {
	url    string
	status int
	err    error
}

func (m model) Init() tea.Cmd { return checkServer(m.url) }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusMsg:
		m.status = int(msg)
		return m, tea.Quit
	case errMsg:
		m.err = msg.err
		return m, tea.Quit
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Interrupt
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	switch {
	case m.err != nil:
		return tea.NewView("error: " + m.err.Error() + "\n")
	case m.status > 0:
		return tea.NewView(fmt.Sprintf("%s returned %d\n", m.url, m.status))
	default:
		return tea.NewView("checking " + m.url + "...\n")
	}
}

func main() {
	if _, err := tea.NewProgram(model{url: "https://charm.land"}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
