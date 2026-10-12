# Tables

The `charm.land/lipgloss/v2/table` subpackage renders a static, fully rendered
table. For a table the user can scroll and select from, use the Bubbles table
instead (charmed-bubbles).

`table.New()` builds the table
(see the [table example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-lipgloss/table/main.go)):

```go
// examples/charmed-lipgloss/table/main.go
header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).Align(lipgloss.Center)
even := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
odd := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
cell := lipgloss.NewStyle().Padding(0, 1)

t := table.New().
    Border(lipgloss.RoundedBorder()).
    BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("63"))).
    Headers("Project", "Status", "Task").
    Row("site", "ready", "build").
    Row("api", "running", "test").
    Width(44).
    StyleFunc(func(row, col int) lipgloss.Style {
        switch {
        case row == table.HeaderRow:
            return header
        case row%2 == 0:
            return cell.Inherit(even)
        default:
            return cell.Inherit(odd)
        }
    })
lipgloss.Println(t)
```

`StyleFunc` receives `table.HeaderRow` (which is -1) for the header and zero
based indexes for data rows; the version in the upstream doc comment that
tests `row == 0` for the header is stale.

Other knobs: `Rows(...)`, `Data(table.Data)` with `table.NewStringData` or
`table.NewFilter`, `Height(n)` plus `YOffset(n)` for a windowed view,
`Wrap(bool)`, `BorderHeader`, `BorderRow`, `BorderColumn`, and the four outer
sides. `Render()` or `String()` produce the text.
