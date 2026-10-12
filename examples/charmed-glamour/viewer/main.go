// examples/charmed-glamour/viewer/main.go
// A markdown pager: glamour renders to a string, the viewport scrolls it,
// and a resize or a background color answer re-renders at the current
// width and theme. WindowSizeMsg and BackgroundColorMsg arrive in either
// order, so both handlers re-render.
package main

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
)

// minWrap keeps WithWordWrap positive on tiny terminals; 0 disables wrapping.
const minWrap = 10

type model struct {
	md     string
	vp     viewport.Model
	width  int
	ready  bool
	isDark bool
}

func (m model) Init() tea.Cmd { return tea.RequestBackgroundColor }

func (m model) render() (string, error) {
	style := "light"
	if m.isDark {
		style = "dark"
	}
	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle(style),
		glamour.WithWordWrap(max(m.width-2, minWrap)),
	)
	if err != nil {
		return "", err
	}
	return r.Render(m.md)
}

// rerender is called from both message handlers below.
func (m *model) rerender() {
	out, err := m.render()
	if err != nil {
		out = "render error: " + err.Error()
	}
	m.vp.SetContent(out)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		m.isDark = msg.IsDark()
		if m.ready { // the size may have arrived first
			m.rerender()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		height := max(msg.Height-1, 1)
		if !m.ready {
			m.vp = viewport.New(viewport.WithWidth(msg.Width), viewport.WithHeight(height))
			m.vp.MouseWheelEnabled = true
			m.ready = true
		} else {
			m.vp.SetWidth(msg.Width)
			m.vp.SetHeight(height)
		}
		m.rerender()
	case tea.KeyPressMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.vp, cmd = m.vp.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	if !m.ready {
		return tea.NewView("loading...")
	}
	v := tea.NewView(m.vp.View() + fmt.Sprintf("\n %3.0f%%  q quits", m.vp.ScrollPercent()*100))
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

func main() {
	md := "# No file given\n\nPass a markdown file as the first argument.\n"
	if len(os.Args) > 1 {
		b, err := os.ReadFile(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		md = string(b)
	}
	md = strings.ReplaceAll(md, "\r\n", "\n")
	if _, err := tea.NewProgram(model{md: md}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
