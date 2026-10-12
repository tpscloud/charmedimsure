# Migrating Bubbles from v1 to v2

Search and replace table for v1 bubbles code.  Return to the [main skill](../SKILL.md) for the shared component contract and focus rules.

| v1                                              | v2                                                   |
|-------------------------------------------------|------------------------------------------------------|
| `github.com/charmbracelet/bubbles/<pkg>`        | `charm.land/bubbles/v2/<pkg>`                        |
| `case tea.KeyMsg:`                              | `case tea.KeyPressMsg:`                              |
| `m.Width = 40`, `m.Height = 20`                 | `m.SetWidth(40)`, `m.SetHeight(20)`; read with `Width()` |
| `textinput.DefaultKeyMap` (var)                 | `textinput.DefaultKeyMap()` (also paginator, textarea) |
| `list.NewModel`, `help.NewModel`, etc           | `New`                                                |
| `list.DefaultStyles()`                          | `list.DefaultStyles(isDark)`                         |
| `list.NewDefaultItemStyles()`                   | `list.NewDefaultItemStyles(isDark)`                  |
| `styles.FilterPrompt`, `styles.FilterCursor`    | `styles.Filter` (a `textinput.Styles`)               |
| `spinner.Tick` (package func)                   | `model.Tick` (method)                                |
| `progress.WithGradient(a, b)`                   | `progress.WithColors(lipgloss.Color(a), lipgloss.Color(b))` |
| `progress.WithDefaultGradient()`                | `progress.WithDefaultBlend()`                        |
| `progress.WithScaledGradient(a, b)`             | `WithColors(...)` plus `WithScaled(true)`            |
| `progress.WithSolidFill(c)`                     | `progress.WithColors(c)`                             |
| `p.FullColor = "#FF0000"`                       | `p.FullColor = lipgloss.Color("#FF0000")`            |
| `stopwatch.NewWithInterval(d)`                  | `stopwatch.New(stopwatch.WithInterval(d))`           |
| `filepicker.DefaultStylesWithRenderer(r)`       | `filepicker.DefaultStyles()`                         |
| `cursor.Blink` field, `BlinkCmd()`              | `IsBlinked`, `Blink()`                               |
| `paginator.UseHLKeys` and friends               | edit `KeyMap` directly                               |
| `runeutil`, `memoization` packages              | internal, not importable                             |

The [v2.2.1 upgrade guide](https://github.com/charmbracelet/bubbles/blob/v2.2.1/UPGRADE_GUIDE_V2.md)
has per component detail.
