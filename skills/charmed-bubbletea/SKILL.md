---
name: charmed-bubbletea
description: "Build terminal user interfaces (TUI) in Go with Bubble Tea v2 (charm.land/bubbletea/v2, the charmbracelet bubbletea framework) using the Elm architecture Model, Update, View loop, commands for I/O, and declarative tea.View fields for alt screen, mouse, cursor, and title. Use when writing or migrating a Bubble Tea or bubbletea app, handling tea.KeyPressMsg or mouse messages, wiring async work into a TUI, or testing one. Not for ready made widgets (charmed-bubbles), styling (charmed-lipgloss), or forms (charmed-huh)."
license: "MIT"
---

# Bubble Tea v2

Import as `tea "charm.land/bubbletea/v2"`, pinned to v2.1.0. Go blocks headed
`// examples/...` are trimmed from
[complete programs](https://github.com/tpscloud/charmedimsure/tree/main/examples/charmed-bubbletea)
that build against the
[pinned modules](https://github.com/tpscloud/charmedimsure/blob/main/examples/go.mod);
`// excerpt` blocks are API sketches.

If the project's `go.mod` still requires `github.com/charmbracelet/bubbletea`
(pre-v2), the v2 APIs here do not apply.  Use `go doc` against the
project's version and keep its import paths unless the user asks to
migrate.

## Workflow

1. Define a model struct holding all UI state, and its `Init`, `Update`, and
   `View` methods (see [the loop](#the-loop)).
2. Return the first command from `Init` (a fetch, a tick, a focus), or nil.
3. In `Update`, switch on the message type, change state, and return any
   follow up work as a `tea.Cmd`.  Handle `ctrl+c` and `tea.WindowSizeMsg`.
4. Return a `tea.View` from `View`, setting terminal features such as
   `AltScreen` and `MouseMode` as fields.
5. Run with `tea.NewProgram(m).Run()` and check its error.
6. Validate: run `go vet ./...` and `go test ./...` from the application
   module, fix failures, and rerun until both pass.  Then try keys and
   resizing in a real terminal.

## The loop

The v2 delta from the usual Elm loop: `View` returns a `tea.View`, not a
string, and the runtime repaints only what changed.  `Msg` is `any`; `Cmd`
is `func() Msg`, run on its own goroutine.  Keep state changes in `Update`
and blocking work in commands.

Minimal program ([examples/charmed-bubbletea/counter/main.go](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbletea/counter/main.go)):

```go
// examples/charmed-bubbletea/counter/main.go
type model struct{ count int }

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyPressMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
        case "up", "k", "space":
            m.count++
        case "down", "j":
            m.count--
        }
    }
    return m, nil
}

func (m model) View() tea.View {
    return tea.NewView(fmt.Sprintf("Count: %d\n\nup/down or space to change, q to quit\n", m.count))
}

func main() {
    if _, err := tea.NewProgram(model{}).Run(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

`Run` blocks until the program quits and returns the final model, which is
how you read results out after the UI closes (see the
[testing example](references/testing.md)).

## tea.View: declare terminal features, do not toggle them

There are no `WithAltScreen` options and no `EnterAltScreen` commands.  You
set fields on the `tea.View` you return, every frame, and the runtime
reconciles the terminal to match.

```go
// excerpt: every tea.View field in one place; fullscreen and mouse examples set a subset
func (m model) View() tea.View {
    v := tea.NewView(content)          // or: var v tea.View; v.SetContent(content)
    v.AltScreen = true                 // full window mode
    v.WindowTitle = "my app"           // terminal title
    v.MouseMode = tea.MouseModeCellMotion
    v.ReportFocus = true               // get FocusMsg and BlurMsg
    v.Cursor = tea.NewCursor(x, y)     // show the real cursor here; nil hides it
    return v
}
```

Because the fields are re-read on every render, turning something off is just
not setting it.  Conditional alt screen, a cursor that only shows while an
input is focused, a title that tracks state: all of these are plain `if`
statements in `View`.

`tea.WindowSizeMsg` arrives once at start and on every resize; store
`Width` and `Height` and lay out from them in `View`.  Read
[references/screen.md](references/screen.md) for a full screen layout, light
or dark background detection, and the remaining fields (terminal colors, the
native progress bar, keyboard enhancements, `OnMouse`, `ProgramStatus`).

## Keys

Match on `tea.KeyPressMsg`.  `msg.String()` gives a readable name and is the
usual way to switch.  Space is `"space"`, not `" "`.

```go
// excerpt: key names as msg.String() reports them
case tea.KeyPressMsg:
    switch msg.String() {
    case "ctrl+c", "esc":  return m, tea.Quit
    case "enter":          // ...
    case "shift+tab":      // ...
    case "a":              // plain characters are themselves
    }
```

Prefer `key.Binding` from Bubbles for anything beyond a toy: it gives you
one place to define keys and help text, and `key.Matches(msg, binding)`
works directly on a `KeyPressMsg`.  See [references/compose.md](references/compose.md).

Read [references/input.md](references/input.md) for structured matching on
`Code`, `Text`, and `Mod`, key release events, paste messages, and mouse
handling.

## Commands and async work

A `Cmd` is `func() tea.Msg`.  Write a function that does the work and returns
a message, then return that function (not its result) from `Init` or
`Update`.

[examples/charmed-bubbletea/asyncfetch/main.go](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbletea/asyncfetch/main.go):

```go
// examples/charmed-bubbletea/asyncfetch/main.go
type statusMsg int
type errMsg struct{ err error }

func checkServer(url string) tea.Cmd {
    return func() tea.Msg {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
        if err != nil {
            return errMsg{err}
        }
        res, err := http.DefaultClient.Do(req)
        if err != nil {
            return errMsg{err}
        }
        defer res.Body.Close()
        return statusMsg(res.StatusCode)
    }
}

func (m model) Init() tea.Cmd { return checkServer(m.url) }
```

Combinators: `tea.Batch(cmds...)` runs commands concurrently with no
ordering; `tea.Sequence(cmds...)` runs them one after another.  Both drop
nil entries, so collecting child commands into a slice and batching is safe.

Built in commands worth knowing: `tea.Quit`, `tea.Interrupt` (Run returns
`tea.ErrInterrupted`), `tea.Suspend` (ctrl+z behavior, then `ResumeMsg`),
`tea.Println` and `tea.Printf` (print above the UI; no effect in alt
screen), `tea.SetClipboard`, `tea.ReadClipboard`, `tea.RequestWindowSize`,
`tea.RequestBackgroundColor`, `tea.ExecProcess`.  Read
[references/api.md](references/api.md) for every message type, command,
program option, and constant group.

Read [references/async.md](references/async.md) for repeating timers
(`tea.Tick`, `tea.Every`), feeding the program from a channel or
`p.Send`, and handing the terminal to an editor with `tea.ExecProcess`.

## Composing bubbles

The parent stores child components, forwards messages, and batches the commands
they return. Children need every message that concerns them (a spinner needs its
tick even while an input has focus), so forward first and gate only where you must.

```go
// examples/charmed-bubbletea/compose/main.go
var cmd tea.Cmd
m.spinner, cmd = m.spinner.Update(msg)
cmds = append(cmds, cmd)
m.input, cmd = m.input.Update(msg)
cmds = append(cmds, cmd)

return m, tea.Batch(cmds...)
```

Read [references/compose.md](references/compose.md) for the full parent model:
a `key.Binding` key map with help, focus on start, resizing children, and
placing a text input's real cursor in the view.

## Testing and debugging

- Test `Update` with constructed messages and inspect the returned model and command.
- For headless program tests, supply input, captured output, and a fixed window size.
- Log to a file while the renderer owns stdout.

Validate as in workflow step 6.  See [testing and debugging](references/testing.md) for unit tests, headless runs,
color-free output, and file logging.

## Migrating from v1

The shape of the change: the import path moves to `charm.land/bubbletea/v2`,
`View() string` becomes `View() tea.View`, `tea.KeyMsg` becomes
`tea.KeyPressMsg` with `Code`, `Text`, and `Mod` fields, mouse messages
split into concrete types, and every program option or command that toggled
a terminal feature (alt screen, mouse, title, cursor) becomes a `tea.View`
field.  The `tea.View` change is the biggest one.  The search and replace
table, in order, is in [references/api.md](references/api.md#migrating-from-v1); the
[upgrade guide](https://github.com/charmbracelet/bubbletea/blob/v2.1.0/UPGRADE_GUIDE_V2.md)
has the complete list.

After applying the table, validate as in workflow step 6, check remaining
v1-only calls such as `tea.EnterAltScreen`, and exercise timers too.
`tea.KeyMsg` remains valid as an interface in v2; do not remove every occurrence.

Projects pinned to a v2 release before v2.1.0 lack `ProgramStatus` and a
few related calls; [references/api.md](references/api.md#older-v2-releases)
lists them.

## Mistakes that cost time

- Confusing a command with a command factory. `checkServer(m.url)` and `tick()`
  return a `tea.Cmd`; return that command from `Init` or `Update`. `tea.Quit`
  already has the command signature, so return it without calling it. Calling
  a command inside `Update` runs its work synchronously and yields a message
  instead of the required command.
- Forgetting to re-arm: ticks and channel reads deliver one message each.
- Doing I/O in `Update` or `View`.  The UI freezes.
- Ignoring `WindowSizeMsg` and then wondering why layout is wrong at start.
- Printing with `fmt` while running.  Use `tea.Println` (outside alt screen)
  or `tea.LogToFile`.
- Expecting ctrl+c to kill the process.  The terminal is in raw mode; you
  must handle it and return `tea.Quit` or `tea.Interrupt`.
- Matching `" "` for space, or `tea.KeyMsg` as a struct.  Both are v1.
- Sharing the model with a goroutine.  Send a message instead.
- Setting a real cursor from a bubble without offsetting it to the bubble's
  position on screen.
