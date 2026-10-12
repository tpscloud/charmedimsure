---
name: charmed-bubbles
description: "Use the Bubbles v2 component library (charm.land/bubbles/v2, charmbracelet bubbles) inside a Bubble Tea v2 TUI: list, table, viewport, textinput, textarea, spinner, progress bar, filepicker, tree, paginator, timer, stopwatch, help, and key bindings. Use when adding a ready made widget or TUI component (text input, spinners, scrolling view, file picker) to a terminal app, wiring its focus, size, styles, and messages, or migrating v1 bubbles code. Not for the program loop itself (charmed-bubbletea), raw styling (charmed-lipgloss), or whole forms (charmed-huh)."
license: "MIT"
---

# Bubbles v2

Add ready-made TUI components to a Bubble Tea v2 program: lists, tables,
text inputs, viewports, spinners, progress bars, file pickers, trees, and
timers, each wired through the same Init, Update, and View contract.

Import each component from its own package (`"charm.land/bubbles/v2/list"`,
`.../table`, `.../viewport`, and so on), pinned to v2.2.1. Go blocks headed
`// examples/...` are trimmed from
[complete programs](https://github.com/tpscloud/charmedimsure/tree/main/examples) that build against the
[pinned modules](https://github.com/tpscloud/charmedimsure/blob/main/examples/go.mod).

If the project's `go.mod` still requires `github.com/charmbracelet/bubbles`
(pre-v2), the v2 APIs here do not apply.  Use `go doc` against the
project's version and keep its import paths unless the user asks to
migrate.

## Workflow

1. Pick the component from the [index](#components) and read its reference
   file.  Start from the complete spinner program below.
2. Construct the bubble in your model constructor, apply its styles (most
   take `isDark`), and store it as a value field.
3. Start it from your `Init`: batch its `Init()`, `Focus()`, or `Tick`
   command with your own.
4. In `Update`, size it from `tea.WindowSizeMsg`, send key presses only to
   the focused bubble, forward every other message, store the returned
   bubble back, and batch the returned commands.
5. Place its `View()` string inside your `tea.View`.
6. Run `go vet ./...` and `go build ./...`, then run the program in a
   terminal.  If a bubble is frozen, empty, or ignores keys, check
   [Mistakes that cost time](#mistakes-that-cost-time), fix, and rerun.

## The contract most bubbles follow

An interactive bubble is a struct with `Update(tea.Msg) (Model, tea.Cmd)`
and `View() string`.  Note both differences from your top level model:
`Update` returns the bubble's own concrete type (no type assertion needed),
and `View` returns a plain string that you place inside your `tea.View`.
Not every package fits this shape: `help.Model.View(k KeyMap)` takes the
key map to render, `key` is a plain data package with no model, and
`cursor` is driven by the inputs that embed it.  Do not build an interface
over "all bubbles" that assumes a zero argument `View`.

A complete program with one bubble
([examples/charmed-bubbles/spinner/main.go](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbles/spinner/main.go));
list, table, viewport, and the others slot into the same places:

```go
// examples/charmed-bubbles/spinner/main.go
package main

import (
    "fmt"
    "os"

    "charm.land/bubbles/v2/spinner"
    tea "charm.land/bubbletea/v2"
)

type model struct {
    spinner spinner.Model
}

func (m model) Init() tea.Cmd {
    return m.spinner.Tick // a method value is a tea.Cmd
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if msg, ok := msg.(tea.KeyPressMsg); ok && (msg.String() == "q" || msg.String() == "ctrl+c") {
        return m, tea.Quit
    }
    var cmd tea.Cmd
    m.spinner, cmd = m.spinner.Update(msg) // store it back; bubbles are values
    return m, cmd
}

func (m model) View() tea.View {
    return tea.NewView(m.spinner.View() + " working, q to quit\n")
}

func main() {
    m := model{spinner: spinner.New(spinner.WithSpinner(spinner.Dot))}
    if _, err := tea.NewProgram(m).Run(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

Some bubbles have an `Init() tea.Cmd` that must run (filepicker reads its
directory, timer starts ticking); others start on a method that returns a
`Cmd` (`textinput.Focus()`, `spinner.Tick`).  Batch these from your own
`Init`.

Rules of thumb:

- Forward every message to every bubble unless you have a reason not to.
  Spinners, timers, and progress bars deliver their own tick and frame
  messages through `Update`; if you gate on key presses only, they stall.
- Gate key presses by focus.  Two text inputs both receiving keys will both
  type.
- Sizes are methods in v2: `SetWidth`, `SetHeight`, `SetSize`, `Width()`,
  `Height()`.  Call them from your `tea.WindowSizeMsg` handler.
- Most style constructors take an `isDark bool`: `list.DefaultStyles`,
  `list.NewDefaultItemStyles`, `textinput.DefaultStyles`,
  `textarea.DefaultStyles`, `help.DefaultStyles`, `tree.DefaultStyles`.
  Two do not: `table.DefaultStyles()` and `filepicker.DefaultStyles()` take
  no argument.  Get the bool from `tea.BackgroundColorMsg` (see
  charmed-bubbletea).
- Key maps are functions now: `textinput.DefaultKeyMap()`, not a package
  variable.


## key and help

`key.Binding` pairs the keys that trigger an action with the help text
shown for it.  `key.Matches(msg, binding)` works on a `tea.KeyPressMsg`.
A type with `ShortHelp() []key.Binding` and `FullHelp() [][]key.Binding`
satisfies `help.KeyMap`, and `help.Model.View(keymap)` renders it in one
line (or all groups when `ShowAll` is true).

```go
// examples/charmed-bubbletea/compose/main.go
var keys = keyMap{
    Submit: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit")),
    Quit:   key.NewBinding(key.WithKeys("ctrl+c", "esc"), key.WithHelp("esc", "quit")),
}

case tea.KeyPressMsg:
    switch {
    case key.Matches(msg, keys.Quit):
        return m, tea.Quit
    }

m.help.View(keys)   // "enter submit • esc quit"
```

Disable a binding with `SetEnabled(false)` and it drops out of help too.

## Components

| Component | Reference | Read it for |
|-----------|-----------|-------------|
| textinput, textarea, cursor | [references/text-input.md](references/text-input.md) | focus switching, echo modes, validation, suggestions, the real cursor, cursor modes |
| list | [references/list.md](references/list.md) | items and delegates, filtering, item mutation, default keys |
| table | [references/table.md](references/table.md) | columns, rows, styles, selection, loading delimited text |
| viewport | [references/viewport.md](references/viewport.md) | sizing, wrapping, horizontal scroll, search highlights, gutters |
| spinner, progress | [references/progress-spinner.md](references/progress-spinner.md) | spinner presets, animated and static bars, color blends |
| filepicker, tree | [references/filepicker-tree.md](references/filepicker-tree.md) | allowed types, selection checks, tree nodes, open and close keys |
| paginator, timer, stopwatch | [references/small-widgets.md](references/small-widgets.md) | slice bounds, countdown timeout, start, stop, toggle |

## Migrating from v1

Import paths move to `charm.land/bubbles/v2/<pkg>`, `tea.KeyMsg` cases become
`tea.KeyPressMsg`, width and height fields become `SetWidth`, `SetHeight`,
and `Width()` methods, default key maps become functions, `NewModel`
constructors become `New`, style constructors take `isDark`, `spinner.Tick`
becomes a method, and the progress gradient options become `WithColors`
and `WithDefaultBlend`.  Read [references/migration.md](references/migration.md)
for the full search and replace table and the upstream upgrade guide.  After
migrating, run `go vet ./...` and `go build ./...` and fix what remains.

## Older versions

The skill describes v2.2.1.  The interactive `tree` package was added in
v2.2, so earlier v2 releases do not have it.  The list's `Quit` binding on
`v` (see [references/list.md](references/list.md)) is v2.2.1 behavior;
check the pinned version's `DefaultKeyMap` after a bump.

## Mistakes that cost time

- Not forwarding tick and frame messages: the spinner freezes, the progress
  bar never animates, the timer never fires.
- Forgetting `Init()` on filepicker or timer.  Nothing loads, nothing runs.
- Creating a list or viewport with a zero size and never calling `SetSize`
  from `WindowSizeMsg`.  It renders empty.
- Handling enter yourself while the list is filtering.  Check
  `SettingFilter()` first.
- Intercepting a key the bubble already uses (enter toggles tree nodes,
  space pages the table) and then wondering why the bubble stopped
  responding.
- Assigning the result of `Update` to a new variable and dropping it.
  Bubbles are values; store them back into the model.
