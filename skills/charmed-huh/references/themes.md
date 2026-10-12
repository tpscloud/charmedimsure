# Themes

How Huh v2.0.3 themes are built, the full `Styles` field list, and the
spinner's separate theme type.

Theme functions have the signature `func(isDark bool) *huh.Styles`. Huh asks
the terminal for its background and calls the theme with the answer. Wrap
custom and built in functions with `huh.ThemeFunc`, for example
`form.WithTheme(huh.ThemeFunc(huh.ThemeDracula))`. Built in functions include
`ThemeCharm`, `ThemeDracula`, `ThemeCatppuccin`, `ThemeBase16`, and
`ThemeBase`. See the complete
[theme example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-huh/theme/main.go):

```go
// examples/charmed-huh/theme/main.go
func myTheme(isDark bool) *huh.Styles {
    s := huh.ThemeCharm(isDark)
    ld := lipgloss.LightDark(isDark)
    accent := ld(lipgloss.Color("#5A2D82"), lipgloss.Color("#C4A7F5"))
    s.Focused.Title = s.Focused.Title.Foreground(accent).Bold(true)
    s.Focused.SelectedOption = s.Focused.SelectedOption.Foreground(accent)
    s.Focused.FocusedButton = s.Focused.FocusedButton.Background(accent)
    return s
}

form.WithTheme(huh.ThemeFunc(myTheme))
```

## Styles fields

`Styles` has `Form`, `Group`, `FieldSeparator`, `Focused` and `Blurred`, and
`Help`.  `Focused` and `Blurred` are `FieldStyles`: `Base`, `Title`,
`Description`, `ErrorIndicator`, `ErrorMessage`, `SelectSelector`, `Option`,
`MultiSelectSelector`, `SelectedOption`, `SelectedPrefix`,
`UnselectedOption`, `UnselectedPrefix`, `TextInput`, `FocusedButton`,
`BlurredButton`, `Card`, `NoteTitle`, `Next`, `Directory`, `File`.

Themes can be set per group or per field too.

## Spinner themes

The spinner's `WithTheme` takes a `spinner.Theme` (for example
`spinner.ThemeFunc(spinner.ThemeDefault)`), not a form theme.
