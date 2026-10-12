# Keys, paste, and mouse in Bubble Tea v2

Structured key matching, key release, paste, and mouse messages.  Return to
the [main skill](../SKILL.md) for `msg.String()` matching and `key.Binding`.

## Structured key matching

For structured matching use the fields of `tea.Key` (a `KeyPressMsg` is a
`Key`): `Code` is a rune, either a printable character or a constant such
as `tea.KeyEnter`, `tea.KeyTab`, `tea.KeyEscape`, `tea.KeyUp`, `tea.KeyF1`;
`Text` holds the printable text and is empty for special keys; `Mod` is a
bitmask you test with `msg.Mod.Contains(tea.ModCtrl)` or compare directly.

```go
// excerpt: structured matching on the tea.Key fields
if msg.Code == 'c' && msg.Mod == tea.ModCtrl {
    return m, tea.Quit // ctrl+c
}
```

`tea.KeyMsg` is an interface covering both `KeyPressMsg` and `KeyReleaseMsg`.
Releases only arrive if you set `v.KeyboardEnhancements.ReportEventTypes =
true` and the terminal supports it.

The full list of key and modifier constants is in [api.md](api.md#key-constants).

## Paste

Pasted text arrives as `tea.PasteMsg` with a `Content` field, bracketed by
`PasteStartMsg` and `PasteEndMsg`.

## Mouse

Set `v.MouseMode` in `View`, then match the concrete message types
([examples/charmed-bubbletea/mouse/main.go](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbletea/mouse/main.go)):

```go
// examples/charmed-bubbletea/mouse/main.go
case tea.MouseClickMsg:
    mouse := msg.Mouse()
    if mouse.Button == tea.MouseLeft {
        m.clicks++
    }
    m.last = fmt.Sprintf("click %s at %d,%d", mouse.Button, mouse.X, mouse.Y)
case tea.MouseWheelMsg:
    m.last = "wheel " + msg.Mouse().Button.String()
case tea.MouseMotionMsg:
    // drags with CellMotion; all movement with AllMotion
case tea.MouseReleaseMsg:
```

`tea.MouseMsg` is an interface with a `Mouse()` method that returns the
`tea.Mouse` struct (`X`, `Y`, `Button`, `Mod`).  Coordinates are zero based
from the top left.  Buttons are `tea.MouseLeft`, `tea.MouseRight`,
`tea.MouseMiddle`, `tea.MouseWheelUp`, `tea.MouseWheelDown`, and so on.
`MouseModeCellMotion` gives clicks, wheel, and drags; `MouseModeAllMotion`
adds hover movement.
