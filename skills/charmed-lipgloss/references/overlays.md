# Overlays with layers

Layers and the compositor draw one block on top of another, so a modal dialog
can sit over a full screen view without slicing strings by hand.

The [overlay example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-lipgloss/overlay/main.go)
is a complete Bubble Tea program. The fragments below span its `compositor`,
`View`, and mouse handler: `m.width` and `m.height` come from
`tea.WindowSizeMsg`; `bg` and `box` are rendered strings,
`w, h := lipgloss.Size(box)`, and `mouse := msg.Mouse()`. From a repository
checkout, run it in a terminal:

```sh
cd examples
go run ./charmed-lipgloss/overlay
```

```go
// examples/charmed-lipgloss/overlay/main.go
base := lipgloss.NewLayer(bg).ID("background")

dialog := lipgloss.NewLayer(box).
    ID("dialog").
    X((m.width - w) / 2).
    Y((m.height - h) / 2).
    Z(1)                              // higher Z draws on top

comp := lipgloss.NewCompositor(base, dialog)
v := tea.NewView(comp.Render())

// Mouse hit testing against the same layers:
hit := comp.Hit(mouse.X, mouse.Y)
if hit.ID() == "dialog" {
    m.showDialog = false              // clicked inside the dialog
}
```

Layers nest (`NewLayer(content, children...)`), and `Compositor.GetLayer(id)`
finds one by ID. `Canvas` is the lower level cell buffer if you need to draw
arbitrary `uv.Drawable`s.

To check an overlay, resize the window and confirm layer placement and mouse
hit testing.
