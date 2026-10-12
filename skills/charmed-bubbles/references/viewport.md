# viewport in Bubbles v2

A scrollable window over a block of text.  Return to the [main skill](../SKILL.md) for the shared component contract and focus rules.

Create it on the first `tea.WindowSizeMsg` and resize it on later ones
([examples/charmed-bubbles/viewport/main.go](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbles/viewport/main.go)):

```go
// examples/charmed-bubbles/viewport/main.go
case tea.WindowSizeMsg:
    if !m.ready {
        m.vp = viewport.New(viewport.WithWidth(msg.Width), viewport.WithHeight(msg.Height-headerHeight))
        m.vp.MouseWheelEnabled = true
        m.vp.SoftWrap = true
        m.vp.SetContent(content)
        m.ready = true
    } else {
        m.vp.SetWidth(msg.Width)
        m.vp.SetHeight(msg.Height - headerHeight)
    }

// View needs v.MouseMode set for the wheel to reach the viewport.
```

`ScrollPercent()`, `AtTop()`, `AtBottom()`, `GotoTop()`, `GotoBottom()`,
`SetYOffset(n)`, `ScrollDown(n)`.  `SoftWrap` wraps long lines; with it
off the viewport scrolls horizontally (`ScrollLeft`, `ScrollRight`).
`SetHighlights([][]int)` plus `HighlightNext` and `HighlightPrevious`
implement search hits.  `LeftGutterFunc` draws line numbers or markers.
`FillHeight` pads short content to the full height.
