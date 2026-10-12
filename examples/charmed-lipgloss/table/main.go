// examples/charmed-lipgloss/table/main.go
// Static tables. For a scrollable, selectable table use charmed-bubbles.
package main

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

func main() {
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
		Row("docs", "pending", "check links").
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

	// Rows can also come from any type that implements table.Data; the
	// StringData helper covers the common case of [][]string.
	data := table.NewStringData([]string{"a", "1"}, []string{"b", "2"})
	lipgloss.Println(table.New().Headers("Key", "Value").Data(data).Border(lipgloss.ASCIIBorder()))
}
