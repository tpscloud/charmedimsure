# Timers, outside events, and external programs in Bubble Tea v2

Patterns for work that arrives over time or from outside the program, all
built on commands.  Return to the [main skill](../SKILL.md) for the loop and
`tea.Batch`.

## Timers

`tea.Tick` and `tea.Every` each send one message.  To keep ticking, return
another tick from the handler ([examples/charmed-bubbletea/ticker/main.go](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbletea/ticker/main.go)):

```go
// examples/charmed-bubbletea/ticker/main.go
type tickMsg time.Time

func tick() tea.Cmd {
    return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m model) Init() tea.Cmd { return tick() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tickMsg:
        m.remaining--
        if m.remaining <= 0 {
            return m, tea.Quit
        }
        return m, tick()
    case tea.KeyPressMsg:
        if msg.String() == "ctrl+c" || msg.String() == "q" {
            return m, tea.Quit
        }
    }
    return m, nil
}
```

`Tick` measures from when it was invoked.  `Every` aligns to the system
clock, so a one second `Every` fires on the second boundary.

## Feeding the program from outside

Two patterns ([examples/charmed-bubbletea/subscribe/main.go](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbletea/subscribe/main.go)).  A channel read
inside a `Cmd`, re-armed after every message:

```go
// examples/charmed-bubbletea/subscribe/main.go
func waitFor(ch <-chan string) tea.Cmd {
    return func() tea.Msg {
        s, ok := <-ch
        if !ok {
            return tea.Quit()
        }
        return eventMsg(s)
    }
}

func (m model) Init() tea.Cmd { return waitFor(m.events) }

case eventMsg:
    m.log = append(m.log, string(msg))
    return m, waitFor(m.events)   // listen again
```

Or `p.Send(msg)` from any goroutine.  `Send` is safe to call concurrently
and blocks until the program has started.

## Running another program

`tea.ExecProcess` releases the terminal, runs an `*exec.Cmd`, and resumes
([examples/charmed-bubbletea/editor/main.go](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbletea/editor/main.go)):

```go
// examples/charmed-bubbletea/editor/main.go
type editorDoneMsg struct{ err error }

func openEditor(path string) tea.Cmd {
    editor := os.Getenv("EDITOR")
    if editor == "" {
        editor = "vi"
    }
    c := exec.Command(editor, path)
    return tea.ExecProcess(c, func(err error) tea.Msg { return editorDoneMsg{err} })
}
```

Keep the fallback: with `EDITOR` unset, `exec.Command("", path)` fails at
launch.  An `EDITOR` value with arguments (`code --wait`) needs splitting
into program and arguments (`strings.Fields` covers the common case); do
not paper over it by running the value through a shell.
