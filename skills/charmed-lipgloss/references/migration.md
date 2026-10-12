# Migrating from Lip Gloss v1

How v1 (`github.com/charmbracelet/lipgloss`) code maps onto v2.0.6
(`charm.land/lipgloss/v2`), and how to check the result.

## What changed

v2 removed the renderer. `Render` always emits full color, and downsampling
to what the terminal supports happens at output: automatically inside Bubble
Tea v2, and through `lipgloss.Print*` standalone.

There is no `AdaptiveColor` in the root package; use `LightDark`. The
`compat` subpackage keeps v1's `AdaptiveColor`, `CompleteColor`, and
`CompleteAdaptiveColor` types for quick migrations. They read stdin and
stdout through package globals and are not thread safe, so prefer `LightDark`
in new code.

New in v2: `UnderlineStyle` with `UnderlineColor`, and layers with a
compositor for overlays ([overlays.md](overlays.md)).

| v1                                            | v2                                                    |
|-----------------------------------------------|-------------------------------------------------------|
| `github.com/charmbracelet/lipgloss`           | `charm.land/lipgloss/v2` (subpackages likewise)       |
| `lipgloss.Color` as a string type             | `lipgloss.Color(s)` function returning `color.Color`  |
| `lipgloss.TerminalColor` in signatures        | `image/color.Color`                                   |
| `lipgloss.AdaptiveColor{Light, Dark}`         | `lipgloss.LightDark(isDark)(light, dark)` or `compat.AdaptiveColor` |
| `lipgloss.CompleteColor{...}`                 | `lipgloss.Complete(profile)(ansi, ansi256, truecolor)` or `compat.CompleteColor` |
| `lipgloss.HasDarkBackground()`                | `lipgloss.HasDarkBackground(os.Stdin, os.Stdout)`     |
| `DefaultRenderer()`, `NewRenderer`, `SetColorProfile` | removed; print with `lipgloss.Println`        |
| `fmt.Println(style.Render(s))`                | `lipgloss.Println(style.Render(s))`                   |
| `WithWhitespaceForeground`, `WithWhitespaceBackground` | `WithWhitespaceStyle(style)`                 |
| `lipgloss.ANSIColor(n)`                       | still valid (an alias of `ansi.IndexedColor`); `lipgloss.Color("n")` also works |

The [v2.0.6 upgrade guide](https://github.com/charmbracelet/lipgloss/blob/v2.0.6/UPGRADE_GUIDE_V2.md)
has the full list.

## Checking a migration

1. Run your application's tests, fix failures, and rerun.
2. Check rendered widths with `lipgloss.Size`.
3. Inspect colors, wrapping, and borders in a terminal.
4. For overlays, resize the window and check layer placement and mouse hit
   testing.

A v1 `lipgloss.Style` passed to a v2 library (or the reverse) fails to
compile; the types are in different packages and the error names both paths.
