# spinner and progress in Bubbles v2

Animated activity indicators.  Both advance only when their own messages are
forwarded through `Update`.

Return to the [main skill](../SKILL.md) for the shared component contract and
focus rules.

## spinner

```go
// examples/charmed-bubbletea/compose/main.go
sp := spinner.New(
    spinner.WithSpinner(spinner.Dot),
    spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("205"))),
)
func (m model) Init() tea.Cmd {
    return tea.Batch(m.input.Focus(), m.spinner.Tick) // start the spinner
}
m.spinner, cmd = m.spinner.Update(msg) // every message
```

Presets: `Line`, `Dot`, `MiniDot`, `Jump`, `Pulse`, `Points`, `Globe`,
`Moon`, `Monkey`, `Meter`, `Hamburger`, `Ellipsis`.  Spinners carry an ID,
so several in one model ignore each other's ticks.  `Tick` is a method on
the model, not a package function.

## progress

Two modes ([examples/charmed-bubbles/progress/main.go](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbles/progress/main.go)).  Animated: call
`SetPercent(p)` (returns a `Cmd`), and forward `progress.FrameMsg` so the
spring animation advances.  Static: `ViewAs(p)` draws a bar at a value with
no state.

```go
// examples/charmed-bubbles/progress/main.go
animated := progress.New(progress.WithDefaultBlend(), progress.WithWidth(40))
static := progress.New(
    progress.WithColors(lipgloss.Color("#5A56E0"), lipgloss.Color("#EE6FF8")),
    progress.WithScaled(true),
    progress.WithoutPercentage(),
)

case tickMsg:
    return m, tea.Batch(m.animated.SetPercent(m.percent), tick())
case progress.FrameMsg:
    m.animated, cmd = m.animated.Update(msg)
    return m, cmd
```

Colors are `color.Color` values.  `WithColors` with one color is a solid
fill; several make a blend; `WithScaled(true)` stretches the blend across
the filled part only.  `WithColorFunc(func(total, current float64)
color.Color)` colors per cell.  `IncrPercent`, `DecrPercent`,
`SetSpringOptions(frequency, damping)`.
