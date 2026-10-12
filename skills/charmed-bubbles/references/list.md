# list in Bubbles v2

A filterable, paginated list with a status bar, help, and a spinner.  Return to the [main skill](../SKILL.md) for the shared component contract and focus rules.

Items are any type with `FilterValue() string`.  The default delegate also
wants `Title()` and `Description()`.  The list owns filtering (`/`),
pagination, a status bar, help, and a spinner
([examples/charmed-bubbles/list/main.go](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbles/list/main.go)):

```go
// examples/charmed-bubbles/list/main.go
type item struct{ name, path, status string }

func (i item) Title() string       { return i.name }
func (i item) Description() string { return fmt.Sprintf("%s (%s)", i.path, i.status) }
func (i item) FilterValue() string { return i.name }

delegate := list.NewDefaultDelegate()
delegate.Styles = list.NewDefaultItemStyles(isDark)
delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("212"))

l := list.New(items, delegate, 0, 0)
l.Title = "Projects"
l.Styles = list.DefaultStyles(isDark)
l.SetStatusBarItemName("project", "projects")
l.DisableQuitKeybindings()

case tea.WindowSizeMsg:
    m.list.SetSize(msg.Width, msg.Height)
case tea.KeyPressMsg:
    if m.list.SettingFilter() {
        break                      // keys belong to the filter input
    }
    switch msg.String() {
    case "enter":
        if it, ok := m.list.SelectedItem().(item); ok {
            m.chosen = it.name
            return m, tea.Quit
        }
    }
```

Mutate items with `SetItems`, `InsertItem`, `SetItem`, `RemoveItem` (the
first three return a `Cmd`; return it).  `NewStatusMessage(s)` flashes a
message in the status bar.  Turn parts off with `SetShowTitle`,
`SetShowStatusBar`, `SetShowPagination`, `SetShowHelp`,
`SetFilteringEnabled`.  `ShowDescription = false` on the delegate gives one
line items; `SetHeight` and `SetSpacing` on the delegate control item
height.  Write your own `ItemDelegate` (`Render`, `Height`, `Spacing`,
`Update`) for custom rows.

Default keys: up/k, down/j, left/h and right/l for pages, g and G, `/` to
filter, esc to clear it, `?` for full help, ctrl+c to quit (`ForceQuit`).
The `Quit` binding is `v`, labelled "select" in the help, so pressing v
also quits the program and q does nothing.  Call
`DisableQuitKeybindings()` and handle quitting in your own `Update`, or
rebind `KeyMap.Quit`.
