# Trees and lists

The `tree` and `list` subpackages of `charm.land/lipgloss/v2` are renderers
with swappable enumerators and styles.

See the [tree and list example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-lipgloss/treelist/main.go):

```go
// examples/charmed-lipgloss/treelist/main.go
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

l := list.New(
    "review the changes",
    "run local checks",
    list.New("format", "test", "build").Enumerator(list.Dash),   // nested list
    "review the output",
).Enumerator(list.Arabic)

checks := list.New("done", "done", "pending").
    Enumerator(func(items list.Items, i int) string {
        if items.At(i).Value() == "pending" {
            return "[ ] "
        }
        return "[x] "
    })
```

List enumerators: `Alphabet`, `Arabic`, `Roman`, `Bullet`, `Dash`,
`Asterisk`. Tree enumerators: `DefaultEnumerator`, `RoundedEnumerator`.

Both take `Indenter` functions and `*StyleFunc` variants that receive the
sibling collection and index, so you can style by position or content.
`Hide(true)` on a node removes it and its children; `Offset(start, end)` trims
siblings.
