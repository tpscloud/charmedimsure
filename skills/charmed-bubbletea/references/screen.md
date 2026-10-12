# Screen, size, and colors in Bubble Tea v2

The less common `tea.View` fields, full window layouts, and light or dark
background detection.  Return to the [main skill](../SKILL.md) for the core
`tea.View` fields; [api.md](api.md#view-fields) has the field table.

## Other tea.View fields

`BackgroundColor`, `ForegroundColor` (terminal colors, nil
resets), `ProgressBar` (native terminal progress, where supported),
`DisableBracketedPasteMode`, `KeyboardEnhancements` (ask for key release
events), `OnMouse` (intercept mouse messages against the last render), and
`ProgramStatus` (an OSC 7501 report of idle, working,
blocked, done, or failed for terminals that support it; it is not cleared
on exit, so report done or error before quitting).

## Window size and full screen

`tea.WindowSizeMsg` arrives once at start and again on every resize.  Store
`Width` and `Height` and use them in `View`
([examples/charmed-bubbletea/fullscreen/main.go](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbletea/fullscreen/main.go)):

```go
// examples/charmed-bubbletea/fullscreen/main.go
case tea.WindowSizeMsg:
    m.width, m.height = msg.Width, msg.Height

func (m model) View() tea.View {
    box := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 3).
        Render(fmt.Sprintf("terminal is %d x %d\npress q to quit", m.width, m.height))
    v := tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box))
    v.AltScreen = true
    v.WindowTitle = "fullscreen example"
    return v
}
```

## Light and dark backgrounds

Bubble Tea downsamples colors to the terminal's profile on output, so styles
may use true color.  What it cannot guess is whether the background is light
or dark; ask for it ([examples/charmed-bubbletea/adaptive/main.go](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbletea/adaptive/main.go)):

```go
// examples/charmed-bubbletea/adaptive/main.go
func newModel() model {
    // Start with a guess so the first frame is already styled; the reply
    // (if one comes) replaces it.
    return model{styles: newStyles(true), isDark: true}
}

func (m model) Init() tea.Cmd { return tea.RequestBackgroundColor }

case tea.BackgroundColorMsg:
    m.isDark, m.detected = msg.IsDark(), true
    m.styles = newStyles(m.isDark)

func newStyles(isDark bool) styles {
    ld := lipgloss.LightDark(isDark)
    return styles{
        title: lipgloss.NewStyle().Bold(true).
            Foreground(ld(lipgloss.Color("#5A2D82"), lipgloss.Color("#C4A7F5"))),
    }
}
```

Render the real UI from the first frame with a default guess; some
terminals never reply, and a "detecting..." placeholder would then stay up
forever.  Most Bubbles style constructors take the same bool (see
charmed-bubbles for the two that do not).
