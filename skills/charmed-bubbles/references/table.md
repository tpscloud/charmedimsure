# table in Bubbles v2

A scrolling, selectable grid.  Return to the [main skill](../SKILL.md) for the shared component contract and focus rules.

For a static rendering use `lipgloss/table` instead.  The complete program is
[examples/charmed-bubbles/table/main.go](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbles/table/main.go):

```go
// examples/charmed-bubbles/table/main.go
t := table.New(
    table.WithColumns([]table.Column{{Title: "Project", Width: 12}, {Title: "Status", Width: 10}}),
    table.WithRows([]table.Row{{"site", "ready"}, {"api", "building"}}),
    table.WithFocused(true),
    table.WithHeight(7),
)
s := table.DefaultStyles()
s.Header = s.Header.Bold(true).BorderStyle(lipgloss.NormalBorder()).BorderBottom(true)
s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57"))
t.SetStyles(s)

row := m.table.SelectedRow()       // table.Row is []string
m.table.HelpView()                 // rendered key help
```

`SetRows`, `SetColumns`, `SetCursor`, `MoveUp(n)`, `MoveDown(n)`,
`GotoTop`, `GotoBottom`, `Focus`, `Blur`, `SetWidth`, `SetHeight`.
`FromValues(csv, sep)` loads rows from delimited text.  An unfocused table
ignores keys.
