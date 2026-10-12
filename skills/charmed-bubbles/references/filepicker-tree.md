# filepicker and tree in Bubbles v2

Browsers for a directory on disk and for an in memory tree of nodes.  Return to the [main skill](../SKILL.md) for the shared component contract and focus rules.

## filepicker

A directory browser ([examples/charmed-bubbles/filepicker/main.go](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbles/filepicker/main.go)):

```go
// examples/charmed-bubbles/filepicker/main.go
fp := filepicker.New()
fp.AllowedTypes = []string{".go", ".md", ".zip"}
fp.CurrentDirectory, _ = os.Getwd()
fp.SetHeight(15)

func (m model) Init() tea.Cmd { return m.fp.Init() }   // reads the directory

m.fp, cmd = m.fp.Update(msg)
if ok, path := m.fp.DidSelectFile(msg); ok {
    m.selected = path
    return m, tea.Quit
}
if ok, path := m.fp.DidSelectDisabledFile(msg); ok {
    m.err = path + " is not an allowed type"
}
```

Fields: `FileAllowed`, `DirAllowed`, `ShowHidden`, `ShowSize`,
`ShowPermissions`, `AutoHeight`, `Cursor`, `Styles`.  `HighlightedPath()`
is the row under the cursor; `Path` is the selection.

## tree

An interactive tree ([examples/charmed-bubbles/tree/main.go](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbles/tree/main.go)):

```go
// examples/charmed-bubbles/tree/main.go
root := tree.Root("project").Child(
    tree.Root("src").Child("main.go", "config.go", "main_test.go"),
    tree.Root("docs").Child("usage.md").Close(),   // starts collapsed
)
t := tree.New(root, 0, 0)
t.SetStyles(tree.DefaultStyles(isDark))
t.SetOpenCharacter("▾")
t.SetClosedCharacter("▸")

n := m.tree.NodeAtCurrentOffset()   // *tree.Node under the cursor
n.Value(), n.IsOpen(), n.ChildNodes()
```

Default keys: j/k move, enter toggles, l/right opens, h/left closes, g/G,
page keys.  Nodes accept the same `Enumerator`, `Indenter`, and style
setters as `lipgloss/tree`, and `SetNodes` swaps the whole tree.
