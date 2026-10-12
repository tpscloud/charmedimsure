# Migrating from Huh v1

How v1 (`github.com/charmbracelet/huh`) code maps onto v2.0.3
(`charm.land/huh/v2`), and how to check the result.

| v1                                          | v2                                              |
|---------------------------------------------|-------------------------------------------------|
| `github.com/charmbracelet/huh`              | `charm.land/huh/v2`                             |
| `huh.ThemeCharm()` returning `*Theme`       | `huh.ThemeCharm(isDark)` returning `*Styles`    |
| `Theme` struct                              | `Theme` interface; `ThemeFunc(func(bool) *Styles)` |
| `field.WithAccessible(true)`                | removed; use `form.WithAccessible(true)`        |
| `huh/accessibility` package                 | removed                                         |
| `Input.Password(true)`                      | `EchoMode(huh.EchoModePassword)`                |
| `lipgloss.Position` from v1 lipgloss        | from `charm.land/lipgloss/v2`                   |
| outer Bubble Tea model `View() string`      | outer model `View() tea.View`                   |

New in v2: `WithViewHook`, `Width(int)` on `MultiSelect`, exported
`huh.Model`. See the
[v2.0.3 upgrade guide](https://github.com/charmbracelet/huh/blob/v2.0.3/UPGRADE_GUIDE_V2.md),
but use the theme wrapper and outer `tea.View` shown in this skill; some
guide examples do not match the pinned APIs.  The upstream README also still
shows the v1 `View() string` outer model.

## Checking a migration

1. Run your application's tests, fix failures, and rerun.
2. In a terminal, check required-field validation, back navigation,
   cancellation, dynamic options, and completion.
3. Check accessible mode separately ([accessible.md](accessible.md)).
4. Resize an embedded form to verify its width and final output.
