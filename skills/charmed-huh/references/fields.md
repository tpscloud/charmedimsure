# Fields, keys, groups, and form options

Setter reference for every Huh v2.0.3 field type, the default key bindings,
group layouts, and form level options.

## Common setters

Fields have `Title`, `Description`, `Key` (name it for `form.Get*`), plus
`TitleFunc` and `DescriptionFunc` for dynamic text.  Two exceptions: `Note`
has no `Key`, and `FilePicker` has no `TitleFunc` or `DescriptionFunc`.
Input fields also have `Value(&v)` and `Validate(func(T) error)`; `Note` has
no answer and does not expose those setters. `Run()` runs a field alone.

## Field types

| Field | Value type | Notable setters |
|-------|-----------|-----------------|
| `NewInput()` | `string` | `Placeholder`, `Prompt`, `CharLimit`, `EchoMode(huh.EchoModePassword)`, `Suggestions` (ctrl+e completes), `Inline` |
| `NewText()` | `string` | `Lines`, `CharLimit`, `ShowLineNumbers`, `Editor("vim")` and `EditorExtension("md")` (ctrl+e opens it), `ExternalEditor(false)` |
| `NewSelect[T]()` | `T` | `Options`, `Height`, `Inline` (horizontal), `Filtering(true)` starts filtered |
| `NewMultiSelect[T]()` | `[]T` | `Options`, `Limit`, `Height`, `Width`, `Filterable`, `Filtering` |
| `NewConfirm()` | `bool` | `Affirmative`, `Negative`, `Inline`, `WithButtonAlignment(lipgloss.Left)` |
| `NewNote()` | none | `Height`, `Next(true)` adds a button, `NextLabel` |
| `NewFilePicker()` | `string` | `CurrentDirectory`, `AllowedTypes`, `FileAllowed`, `DirAllowed`, `ShowHidden`, `ShowSize`, `ShowPermissions`, `Height`, `Picking(true)` |

Note descriptions render light markdown: `_italic_`, `*bold*`, and backtick
code.

## Default keys

Enter or tab moves to the next field, shift+tab back.  Select and
multiselect move with up/down, j/k, ctrl+p/ctrl+n; `/` filters, g and G
jump, ctrl+u and ctrl+d move half a page.  Multiselect toggles with space or
x and selects all or none with ctrl+a.  Confirm toggles with left/right or
h/l and accepts y or n.  Ctrl+c aborts.

All of this lives in `huh.KeyMap`; change it with `NewDefaultKeyMap()`,
edit, then `WithKeyMap`.

## Group layouts

`form.WithLayout(huh.LayoutDefault)` shows one group at a time;
`huh.LayoutStack` shows them all; `huh.LayoutColumns(n)` and
`huh.LayoutGrid(rows, cols)` tile them.

## Form level options

`WithShowHelp`, `WithShowErrors`, `WithWidth`, `WithHeight`,
`WithTimeout(d)` (returns `huh.ErrTimeout`), `WithOutput(os.Stdout)`,
`WithInput`, `WithTheme`, `WithKeyMap`, `WithAccessible(bool)`,
`WithProgramOptions(tea.ProgramOption...)`, and
`WithViewHook(func(tea.View) tea.View)` to set alt screen or the mouse mode
on the underlying program.

`WithProgramOptions` replaces earlier program options, including the default
stderr output, so pass `tea.WithOutput` in the same call if you need it.
