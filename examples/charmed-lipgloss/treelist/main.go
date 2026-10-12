// examples/charmed-lipgloss/treelist/main.go
// Trees and lists are static renderers with pluggable enumerators.
package main

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/list"
	"charm.land/lipgloss/v2/tree"
)

func main() {
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	item := lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	t := tree.Root("project").
		Child(
			tree.Root("src").Child("main.go", "config.go"),
			tree.Root("docs").Child("guide.md"),
			"assets",
		).
		Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(dim).
		ItemStyle(item)
	lipgloss.Println(t)

	l := list.New(
		"review the changes",
		"run local checks",
		list.New("format", "test", "build").Enumerator(list.Dash),
		"review the output",
	).Enumerator(list.Arabic).EnumeratorStyle(dim)
	lipgloss.Println(l)

	// A custom enumerator gets the sibling collection and the index.
	checks := list.New("done", "done", "pending").
		Enumerator(func(items list.Items, i int) string {
			if items.At(i).Value() == "pending" {
				return "[ ] "
			}
			return "[x] "
		})
	lipgloss.Println(checks)
}
