// examples/charmed-bubbles/progress/main.go
// Animated progress (SetPercent returns a Cmd; forward FrameMsg) next to a
// static bar drawn with ViewAs.
package main

import (
	"fmt"
	"os"
	"time"

	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type tickMsg time.Time

type model struct {
	animated progress.Model
	static   progress.Model
	percent  float64
}

func newModel() model {
	return model{
		animated: progress.New(progress.WithDefaultBlend(), progress.WithWidth(40)),
		static: progress.New(
			progress.WithColors(lipgloss.Color("#5A56E0"), lipgloss.Color("#EE6FF8")),
			progress.WithScaled(true),
			progress.WithWidth(40),
			progress.WithoutPercentage(),
		),
	}
}

func tick() tea.Cmd {
	return tea.Tick(400*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m model) Init() tea.Cmd { return tick() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	case tickMsg:
		if m.percent >= 1 {
			return m, tea.Quit
		}
		m.percent += 0.1
		// SetPercent starts the spring animation; keep ticking alongside it.
		return m, tea.Batch(m.animated.SetPercent(m.percent), tick())
	case progress.FrameMsg:
		// Animation frames must reach the progress model.
		var cmd tea.Cmd
		m.animated, cmd = m.animated.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() tea.View {
	return tea.NewView(fmt.Sprintf("animated  %s\nstatic    %s  %3.0f%%\n",
		m.animated.View(), m.static.ViewAs(m.percent), m.percent*100))
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
