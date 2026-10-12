# Composing bubbles in Bubble Tea v2

A parent model that holds Bubbles components, routes keys through
`key.Binding`, renders help, and places a real cursor.  Return to the
[main skill](../SKILL.md) for the loop; charmed-bubbles covers each
component.

[examples/charmed-bubbletea/compose/main.go](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbletea/compose/main.go):

```go
// examples/charmed-bubbletea/compose/main.go
type keyMap struct {
    Submit key.Binding
    Quit   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding { return []key.Binding{k.Submit, k.Quit} }
func (k keyMap) FullHelp() [][]key.Binding { return [][]key.Binding{{k.Submit, k.Quit}} }

var keys = keyMap{
    Submit: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit")),
    Quit:   key.NewBinding(key.WithKeys("ctrl+c", "esc"), key.WithHelp("esc", "quit")),
}

func (m model) Init() tea.Cmd {
    return tea.Batch(m.input.Focus(), m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.input.SetWidth(msg.Width - 4)
        m.help.SetWidth(msg.Width)
    case tea.KeyPressMsg:
        switch {
        case key.Matches(msg, keys.Quit):
            return m, tea.Quit
        case key.Matches(msg, keys.Submit):
            m.submitted = m.input.Value()
            m.input.Reset()
        }
    }

    var cmd tea.Cmd
    m.spinner, cmd = m.spinner.Update(msg)
    cmds = append(cmds, cmd)
    m.input, cmd = m.input.Update(msg)
    cmds = append(cmds, cmd)

    return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
    body := fmt.Sprintf("%s working\n\n%s\n\nlast submitted: %q\n\n%s\n",
        m.spinner.View(), m.input.View(), m.submitted, m.help.View(keys))
    v := tea.NewView(body)
    if c := m.input.Cursor(); c != nil {
        c.Y += 2                  // the input is on row 2 of this layout
        v.Cursor = c
    }
    return v
}
```

Bubbles' `Update` methods return their own concrete type, so
`m.input, cmd = m.input.Update(msg)` needs no type assertion.
A text input can either draw a fake cursor inside its own
string (the default) or hand you a real `*tea.Cursor` via `Cursor()` once you
call `SetVirtualCursor(false)`; you then offset it by where the input sits in
your layout and put it on the view.
