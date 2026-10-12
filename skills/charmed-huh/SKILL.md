---
name: charmed-huh
description: "Build terminal forms and interactive prompts in Go with Huh v2 (charm.land/huh/v2, charmbracelet): input, text, password, select, multiselect, confirm, note, and file picker fields grouped into pages, with validation, dynamic fields, themes, accessible mode, and a spinner for follow up work. Use when asking a user questions from a CLI, building a wizard, survey, or yes or no question, either standalone with Run or embedded in a Bubble Tea v2 program. Not for shell scripts (charmed-gum) or free form TUIs (charmed-bubbletea)."
license: "MIT"
---

# Huh v2

Build terminal forms and prompts in Go with Huh v2: typed fields, validation,
dynamic options, themes, accessible mode, and forms embedded in Bubble Tea.

Import as `"charm.land/huh/v2"`; the spinner is `"charm.land/huh/v2/spinner"`.
Pinned and verified against v2.0.3. Go blocks headed `// examples/...` are
trimmed from
[complete programs](https://github.com/tpscloud/charmedimsure/tree/main/examples/charmed-huh)
that build against the pinned modules; `// excerpt` blocks are API sketches.

If the project's `go.mod` still requires `github.com/charmbracelet/huh`
(pre-v2), the v2 APIs here do not apply.  Use `go doc` against the
project's version and keep its import paths unless the user asks to
migrate.

## Workflow

1. Declare a variable for each answer.
2. Build fields with a `Title` and `Value(&v)` (or a `Key`), add `Validate`
   where input must be checked, and put them into pages with `NewGroup`.
3. Choose output and mode: `WithOutput(os.Stderr)` when stdout carries
   results; call `WithAccessible` only when `ACCESSIBLE` is set.
4. Standalone, call `form.Run()` and handle `huh.ErrUserAborted` (exit 130)
   and other errors.  Inside Bubble Tea, embed the form instead.
5. Check the bound values after `Run`; accessible mode returns nil without
   running validation on EOF.
6. Try it in a terminal: required fields, back navigation, ctrl+c, dynamic
   options, and once in accessible mode (`TERM=dumb`).  Fix what fails and
   rerun.

## Shape of a form

`Form` holds `Group`s; a group is one page of `Field`s.  The user moves
through fields with enter or tab (shift+tab back), and the form advances to
the next group when every field on the page validates.  Answers land in
variables you bind with `Value(&v)`.

The [standalone example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-huh/standalone/main.go)
collects and prints a local task plan. From a repository checkout, run it in
a terminal with `cd examples && go run ./charmed-huh/standalone`.

```go
// examples/charmed-huh/standalone/main.go
var (
    projects []string
    tasks    []string
    location string
    confirm  bool
)

form := huh.NewForm(
    huh.NewGroup(
        huh.NewNote().
            Title("Local task planner").
            Description("Choose projects and tasks. Use _space_ to toggle, `/` to filter, *ctrl+a* for all."),

        huh.NewMultiSelect[string]().
            Title("Projects").
            Options(
                huh.NewOption("Website", "site").Selected(true),
                huh.NewOption("Documentation", "docs"),
            ).
            Filterable(true).
            Height(8).
            Validate(func(v []string) error {
                if len(v) == 0 {
                    return errors.New("pick at least one project")
                }
                return nil
            }).
            Value(&projects),

        huh.NewMultiSelect[string]().
            Title("Tasks").
            Options(huh.NewOptions("build", "test", "lint", "format")...).
            Value(&tasks),
    ).Title("What"),

    huh.NewGroup(
        huh.NewSelect[string]().
            Title("Working directory").
            Options(
                huh.NewOption("Current directory", "."),
                huh.NewOption("Parent directory", ".."),
            ).
            Value(&location),
    ).Title("Where"),

    huh.NewGroup(
        huh.NewConfirm().Title("Print task plan?").Affirmative("Print").Negative("Cancel").Value(&confirm),
    ),
)

err := form.Run()
switch {
case errors.Is(err, huh.ErrUserAborted):
    fmt.Fprintln(os.Stderr, "cancelled")
    os.Exit(130)                // ctrl+c
case err != nil:
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
}
```

`Run` blocks. The normal TUI renders to stderr; accessible prompts default
to stdout. `RunWithContext(ctx)` adds cancellation to the TUI. A single
field can run on its own: `huh.NewInput().Title("Name").Value(&name).Run()`.

Options are `Option[T]` values: `NewOption(label, value)` for distinct
labels, `NewOptions(values...)` when the label is the value.  `T` is any
comparable type, so selects can yield ints, enums, or struct keys directly.

## Fields and validation

The fields are `NewInput`, `NewText`, `NewSelect[T]`, `NewMultiSelect[T]`,
`NewConfirm`, `NewNote`, and `NewFilePicker`.  All have `Title` and
`Description`; all but `Note` take `Value(&v)`, `Validate`, and `Key`.
Passwords use `NewInput().EchoMode(huh.EchoModePassword)`.

Read [references/fields.md](references/fields.md) for the per field setter
table, default key bindings and `KeyMap`, group layouts, and every form
level option (`WithTimeout`, `WithProgramOptions`, `WithViewHook`, and the
rest).

Built in validators for strings: `ValidateNotEmpty()`,
`ValidateMinLength(n)`, `ValidateMaxLength(n)`, `ValidateLength(min, max)`,
`ValidateOneOf(a, b, c)`.  Validation runs when the field loses focus; a
failing field blocks the page.

## Groups and dynamic fields

`NewGroup(fields...)` with `Title`, `Description`, `WithHide(bool)`, and
`WithHideFunc(func() bool)` to skip a page based on earlier answers
(see the [dynamic example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-huh/dynamic/main.go)):

```go
// examples/charmed-huh/dynamic/main.go
huh.NewGroup(
    huh.NewNote().
        TitleFunc(func() string { return "Details for " + project }, &project).
        Description("This example only prints your choices; it does not run tasks."),
).WithHideFunc(func() bool { return !details })
```

The `*Func` setters take a function and a binding.  huh re-runs the
function whenever the bound value changes, so pass a pointer to whatever
the field depends on (the dynamic example defines `tasksByProject`):

```go
// examples/charmed-huh/dynamic/main.go
huh.NewSelect[string]().
    TitleFunc(func() string { return "Task for " + project }, &project).
    OptionsFunc(func() []huh.Option[string] {
        return huh.NewOptions(tasksByProject[project]...)
    }, &project).
    Value(&task)
```

`PlaceholderFunc`, `SuggestionsFunc`, and `DescriptionFunc` follow the same
pattern.  When a field depends on several values, bind a struct pointer or
a slice.

## Embedding in Bubble Tea

`*huh.Form` has its own component contract: `Update` returns
`huh.Model` (an interface) and `View` returns a `string`.  Assert the
result back to `*huh.Form`, watch `form.State`, and wrap the view in the
outer model's `tea.View`:

```go
// examples/charmed-huh/embedded/main.go
func (m model) Init() tea.Cmd { return m.form.Init() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        // Size the form to what is left inside the frame's border and padding.
        fw, _ := m.styles.frame.GetFrameSize()
        m.form = m.form.WithWidth(max(msg.Width-fw, minFormWidth))
    case tea.KeyPressMsg:
        if msg.String() == "ctrl+c" {
            return m, tea.Quit
        }
    }

    next, cmd := m.form.Update(msg)
    if f, ok := next.(*huh.Form); ok {
        m.form = f
    }

    switch m.form.State {
    case huh.StateCompleted, huh.StateAborted:
        return m, tea.Quit
    }
    return m, cmd
}

func (m model) View() tea.View {
    v := tea.NewView(m.styles.frame.Render(m.form.View()))
    v.AltScreen = true
    return v
}
```

Read [references/embedding.md](references/embedding.md) for reading answers
from the final model, why `tea.Printf` loses output on the alt screen,
`SubmitCmd` and `CancelCmd`, and navigation commands.

## Themes and accessible mode

Theme functions have the signature `func(isDark bool) *huh.Styles`; wrap
them with `huh.ThemeFunc`, as in
`form.WithTheme(huh.ThemeFunc(huh.ThemeDracula))`.  Read
[references/themes.md](references/themes.md) for the built in themes, a
custom theme derived from `ThemeCharm`, and the full `Styles` field list.

`form.WithAccessible(true)` swaps the TUI for line based prompts for screen
readers, and `NewForm` turns it on when `TERM` is `dumb`.  Read
[references/accessible.md](references/accessible.md) for honoring
`ACCESSIBLE` without overriding that default, and for what the accessible
path skips (timeouts, context, hidden groups, validation on EOF).

## Spinner

For work that follows the form, see the
[spinner example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-huh/spinner/main.go).
Its `simulateWork` waits for a timer or cancellation; it performs no file or
network operations:

```go
// examples/charmed-huh/spinner/main.go
err := spinner.New().
    Type(spinner.Dots).
    Title("Simulating local work...").
    Context(ctx).
    ActionWithErr(simulateWork).
    Run()
```

`Action(func())` is the no error form.  Types: `Line`, `Dots`, `MiniDot`,
`Jump`, `Points`, `Pulse`, `Globe`, `Moon`, `Monkey`, `Meter`, `Hamburger`,
`Ellipsis`.  `WithAccessible` and `WithOutput` mirror the form; `WithTheme`
takes a spinner theme (see [references/themes.md](references/themes.md)).
With only `Context` and no action, the spinner runs until the context ends.

## Migrating from v1

Read [references/migration.md](references/migration.md) for the v1 to v2
mapping (themes as `func(isDark bool) *Styles`, removed field level
accessible setters, `EchoMode` for passwords), what is new in v2, and the
steps to check a migrated form.

## Mistakes that cost time

- No `Value(&v)` and no `Key`.  The answer is unreachable after `Run`.
- `huh.NewSelect()` without the type parameter.  It does not compile;
  write `huh.NewSelect[string]()`.
- `TitleFunc(fn, project)` with the value instead of `&project`.  It
  compiles but never refreshes.
- Calling `Run` inside a Bubble Tea `Update`, or on an embedded form.  Two
  programs fight for the terminal.  Embed instead.
- Ignoring `ErrUserAborted`.  Ctrl+c is an error return, not a panic.
  Accessible mode never returns it; check the values there instead.
- Assuming both modes use the same output stream. Set `WithOutput`
  explicitly (`WithOutput(os.Stderr)` keeps stdout for results) when stdout
  must contain only results.
- A `Validate` on multiselect that reads another field's value before that
  field is filled.  Validation runs on blur, in page order.
- Printing a completion message with `tea.Printf` from an alt screen
  program.  It prints nothing; print from the final model after `Run`.
