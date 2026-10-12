# textinput, textarea, and cursor in Bubbles v2

Single line and multi line text editors, and the cursor they embed.  Return to the [main skill](../SKILL.md) for the shared component contract and focus rules.

Single line and multi line editors with the same lifecycle: `Focus()`
returns a `Cmd` (start the cursor blink), `Blur()` stops it, `Value()` reads
the text, `SetValue` and `Reset` write it
([examples/charmed-bubbles/textarea/main.go](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbles/textarea/main.go)):

```go
// examples/charmed-bubbles/textarea/main.go
ti := textinput.New()
ti.Placeholder = "title"
ti.CharLimit = 80
ti.SetStyles(textinput.DefaultStyles(isDark))

ta := textarea.New()
ta.Placeholder = "body"
ta.ShowLineNumbers = true
ta.SetHeight(6)
ta.SetStyles(textarea.DefaultStyles(isDark))

func (m model) Init() tea.Cmd { return m.title.Focus() }

case "tab", "shift+tab":
    m.focus = 1 - m.focus
    if m.focus == 0 {
        m.body.Blur()
        return m, m.title.Focus()
    }
    m.title.Blur()
    return m, m.body.Focus()

// Key presses go to the focused input only; everything else (blink
// ticks, resizes) goes to both.
_, isKey := msg.(tea.KeyPressMsg)
var cmds []tea.Cmd
var cmd tea.Cmd
if !isKey || m.focus == 0 {
    m.title, cmd = m.title.Update(msg)
    cmds = append(cmds, cmd)
}
if !isKey || m.focus == 1 {
    m.body, cmd = m.body.Update(msg)
    cmds = append(cmds, cmd)
}
return m, tea.Batch(cmds...)
```

textinput extras: `EchoMode` (`EchoNormal`, `EchoPassword`, `EchoNone`),
`Validate func(string) error` (sets `Err`), `ShowSuggestions` with
`SetSuggestions` for tab completion, `Prompt`.  textarea extras:
`MaxHeight`, `MaxWidth`, `DynamicHeight`, `EndOfBufferCharacter`,
selection methods, `LineCount`, `Line`, `Column`.

Both draw a simulated cursor by default.  For the terminal's real cursor,
call `SetVirtualCursor(false)`, then in your `View` take `Cursor()`, add
the component's offset on screen, and assign it to `v.Cursor`.  The doc
comment on `Cursor` mentions `tea.NewFrame`, which no longer exists; use
`tea.NewView`.

## cursor

Used inside textinput and textarea; rarely needed directly.  `cursor.New()`
starts in blink mode.  `cursor.Mode` is an enum with three values, not a set
of flags: `cursor.CursorBlink`, `cursor.CursorStatic`, `cursor.CursorHide`.
`SetMode(cursor.CursorStatic)` picks one; values outside the three are
ignored.  `SetMode` returns a `tea.Cmd` (the blink command when switching
to `CursorBlink`, nil otherwise) that you must return from `Update` for the
blink to start.  The package-level `cursor.Blink` is that start command.
The method `m.Blink()` is different: it schedules the next blink and
returns nil unless the mode is `CursorBlink`.
