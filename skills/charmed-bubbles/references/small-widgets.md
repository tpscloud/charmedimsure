# paginator, timer, and stopwatch in Bubbles v2

Small helpers for paging a slice and counting time.  Return to the [main skill](../SKILL.md) for the shared component contract and focus rules.

## paginator

paginator computes bounds; you slice
([examples/charmed-bubbles/paginator/main.go](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbles/paginator/main.go)):

```go
// examples/charmed-bubbles/paginator/main.go
p := paginator.New(paginator.WithPerPage(8))
p.Type = paginator.Dots            // or paginator.Arabic
p.SetTotalPages(len(items))
start, end := m.pager.GetSliceBounds(len(m.items))
for _, it := range m.items[start:end] {
    b.WriteString("  " + it + "\n")
}
```

## timer and stopwatch

timer counts down and emits `timer.TimeoutMsg`; stopwatch counts up.  Both
start from `Init()`, and `Start`, `Stop`, `Toggle` return commands
([examples/charmed-bubbles/timer/main.go](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbles/timer/main.go)):

```go
// examples/charmed-bubbles/timer/main.go
t := timer.New(10*time.Second, timer.WithInterval(time.Second))
sw := stopwatch.New(stopwatch.WithInterval(100 * time.Millisecond))

func (m model) Init() tea.Cmd { return tea.Batch(m.timer.Init(), m.stopwatch.Init()) }

case "space":
    return m, tea.Batch(m.timer.Toggle(), m.stopwatch.Toggle())
case timer.TimeoutMsg:
    m.expired = true
```

`timer.Timedout()`, `timer.Running()`, `stopwatch.Elapsed()`,
`stopwatch.Reset()`.
