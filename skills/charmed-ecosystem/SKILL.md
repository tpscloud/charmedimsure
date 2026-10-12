---
name: charmed-ecosystem
description: "Map of the Charm (charmbracelet) terminal stack and how the pieces fit: which library to reach for (Bubble Tea, Bubbles, Lip Gloss, Huh, Glamour, Log, Wish, Harmonica) or which CLI (gum, vhs, glow, freeze, mods, crush), the v2 charm.land import paths and pinned versions, the dependency order, how the libraries hand data to each other, and a consolidated v1 to v2 migration checklist. Use when starting a TUI or terminal project in Go, choosing a Go TUI stack, deciding between Charm libraries, or planning a charmbracelet v1 to v2 upgrade. For any single library's API use its own charmed-* skill."
license: "MIT"
---

# The Charm ecosystem

Choose the right Charm library or CLI for a terminal project, get the v2
`charm.land` import paths and pinned versions, and see how the pieces fit
together before reaching for a library skill.

Charm publishes Go libraries for terminal programs and a set of CLIs built on
them. The core v2 UI libraries import from `charm.land`; supporting libraries
can still use GitHub module paths. This skill routes to the per-library
charmed-* skills and their
[examples](https://github.com/tpscloud/charmedimsure/tree/main/examples),
which build against the pinned modules.

## Workflow

1. Match the task to a row in "Which library" below and note the skills it
   names.
2. Check the project's `go.mod`. `github.com/charmbracelet/<name>` paths for
   the UI libraries mean the pre-v2 APIs; keep them unless the user asks to
   migrate, and then follow "Migrating from v1".
3. Add the pinned versions with the `go get` command below.
4. Open the per-library skill for API detail and start from the example it
   links.
5. Run `go test ./...` in the application module. On a type mismatch between
   Charm libraries, run `go list -m all` and look for a v1 and v2 mix of the
   libraries that share types.

## Import paths and versions

| Library | Import | Pinned here | Skill |
|---------|--------|-------------|-------|
| [Bubble Tea](https://github.com/charmbracelet/bubbletea) | `charm.land/bubbletea/v2` (alias `tea`) | v2.1.0 | [charmed-bubbletea](https://github.com/tpscloud/charmedimsure/blob/main/skills/charmed-bubbletea/SKILL.md) |
| [Bubbles](https://github.com/charmbracelet/bubbles) | `charm.land/bubbles/v2/<component>` | v2.2.1 | [charmed-bubbles](https://github.com/tpscloud/charmedimsure/blob/main/skills/charmed-bubbles/SKILL.md) |
| [Lip Gloss](https://github.com/charmbracelet/lipgloss) | `charm.land/lipgloss/v2` (+ `/table`, `/tree`, `/list`, `/compat`) | v2.0.6 | [charmed-lipgloss](https://github.com/tpscloud/charmedimsure/blob/main/skills/charmed-lipgloss/SKILL.md) |
| [Huh](https://github.com/charmbracelet/huh) | `charm.land/huh/v2` (+ `/spinner`) | v2.0.3 | [charmed-huh](https://github.com/tpscloud/charmedimsure/blob/main/skills/charmed-huh/SKILL.md) |
| [Glamour](https://github.com/charmbracelet/glamour) | `charm.land/glamour/v2` (+ `/ansi`, `/styles`) | v2.0.1 | [charmed-glamour](https://github.com/tpscloud/charmedimsure/blob/main/skills/charmed-glamour/SKILL.md) |
| [Log](https://github.com/charmbracelet/log) | `charm.land/log/v2` | v2.0.1 | [charmed-log](https://github.com/tpscloud/charmedimsure/blob/main/skills/charmed-log/SKILL.md) |
| [colorprofile](https://github.com/charmbracelet/colorprofile) | `github.com/charmbracelet/colorprofile` | v0.4.3 | used by the above |
| [Harmonica](https://github.com/charmbracelet/harmonica) | `github.com/charmbracelet/harmonica` | v0.2.0 | spring animation, used by progress |
| [Wish](https://github.com/charmbracelet/wish) | `charm.land/wish/v2` | Not pinned | SSH apps (not covered here) |

The versions above, except Wish, are pinned in the
[example module](https://github.com/tpscloud/charmedimsure/blob/main/examples/go.mod).
Wish is listed for orientation only. It is not a dependency of the example
module, and this repository has no SSH example.

For Bubble Tea, Bubbles, Lip Gloss, Huh, Glamour, and Log, the older
`github.com/charmbracelet/<name>` paths belong to the pre-v2 APIs.
This rule does not apply to colorprofile or Harmonica, which keep their
GitHub module paths. Two major versions can coexist in one binary as
separate modules, but the libraries that share types (Bubble Tea, Bubbles,
Huh, Lip Gloss) have to move together, and a v2 Log or Glamour needs v2
Lip Gloss.

From an existing application module, select the tested core UI versions with:

```sh
go get charm.land/bubbletea/v2@v2.1.0 charm.land/bubbles/v2@v2.2.1 \
       charm.land/lipgloss/v2@v2.0.6 charm.land/huh/v2@v2.0.3 \
       charm.land/glamour/v2@v2.0.1 charm.land/log/v2@v2.0.1
```

## Which library

**A full interactive program (panes, navigation, live data).**  Bubble
Tea for the loop, Bubbles for the widgets, Lip Gloss for layout and
color. Start from the
[compose](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbletea/compose/main.go)
or [fullscreen](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbletea/fullscreen/main.go)
example in charmed-bubbletea.

**A CLI that asks a few questions then does work.**  Huh alone.  Run the
form, then a Huh spinner around the work.  No Bubble Tea code needed.

**Questions inside a larger TUI.**  Huh embedded as a bubble; see the
[embedded example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-huh/embedded/main.go)
in charmed-huh.

**Pretty output from a non interactive command.**  Lip Gloss only:
styles, tables, trees, printed with `lipgloss.Println`.  Add Glamour if the
content is Markdown.  Add Log for the logging.

**A shell script that needs prompts.**  gum (charmed-gum).  Same widgets,
no Go.

**A demo GIF for the README.**  vhs (charmed-vhs).

**Serve the TUI over SSH.**  Wish wraps a Bubble Tea program per
connection; pass `tea.WithEnvironment(session.Environ())` so color
detection uses the client's terminal.

**Animation.**  Harmonica springs; Bubbles' progress bar already uses it.
For custom motion, tick with `tea.Tick` and step the spring on each tick.

## How they fit together

```
ultraviolet, x/ansi, colorprofile      terminal I/O, escape codes, profiles
        │
   Bubble Tea ◄──── Lip Gloss           loop and rendering ◄── styles, layout
        │               │
     Bubbles            │                widgets styled by Lip Gloss
        │               │
       Huh ◄────────────┘                forms built from Bubbles
     Glamour ───► Lip Gloss              Markdown to ANSI, printed via Lip Gloss
       Log ─────► Lip Gloss              styled log lines
```

Data flows as strings. Lip Gloss `Render` returns a string; Bubbles'
`View` returns a string; Glamour `Render` returns a string and an error,
which you must handle. Your top level
`View` wraps the assembled string in `tea.NewView`.  There is exactly one
place where terminal state is declared (fields on `tea.View`) and one
place where color downsampling happens (the Bubble Tea renderer, or the
`lipgloss.Print` family when there is no program).

Dark or light detection follows the same rule.  Inside a program, return
`tea.RequestBackgroundColor` and read `IsDark()` from the message; most
Bubbles style constructors (`DefaultStyles(isDark)`; table and filepicker
take no argument) and every Huh theme take that bool.  Outside a
program, call `lipgloss.HasDarkBackground(os.Stdin, os.Stdout)` once
before printing anything.

## Recipes across libraries

- **Themed app.** One `styles` struct built from `isDark`, passed to
  every bubble's `SetStyles` or `Styles` field and to Huh's `ThemeFunc`.
  Build it in the `BackgroundColorMsg` handler.
- **Markdown help screen.** Glamour renders at the available width, and a
  Bubbles viewport scrolls it. The
  [viewer example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-glamour/viewer/main.go)
  subtracts layout space from `WindowSizeMsg.Width` and clamps the result.
- **Modal dialog.** Lip Gloss layers and compositor over the main view;
  hit test clicks with `Compositor.Hit` (the
  [overlay example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-lipgloss/overlay/main.go)).
- **Long running work with progress.** A `Cmd` that reads from a channel
  the worker writes to; each message updates a Bubbles progress model
  (combine the [subscription pattern](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbletea/subscribe/main.go)
  with the [progress example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbles/progress/main.go)).
- **Debug logging in a TUI.** `tea.LogToFileWith("debug.log", "app",
  log.Default())`, then tail the file. See the
  [logging example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-log/tea/main.go).
- **Tests.** Drive `Update` with constructed messages; run the whole
  program with `tea.WithInput`, `tea.WithOutput`, `tea.WithWindowSize`,
  and `tea.WithColorProfile(colorprofile.Ascii)` for stable strings.

## Migrating from v1

Move Bubble Tea, Bubbles, Huh, and Lip Gloss in one step, since they share
`tea.Model`, `tea.Cmd`, and `lipgloss.Style`; Log and Glamour only need a
matching Lip Gloss major. Start by replacing `github.com/charmbracelet/X`
with `charm.land/X/v2` for bubbletea, bubbles, lipgloss, huh, glamour, and
log; subpackages keep their suffix (`bubbles/v2/list`, `lipgloss/v2/table`).

Read [references/migration.md](references/migration.md) for the
per-library API changes (Bubble Tea keys, mouse, and view fields; Lip Gloss
colors; Bubbles sizes and key maps; Huh themes; Glamour and Log), links to
the pinned upstream upgrade guides, and the test loop to run after each step.

## The CLIs

| Tool | Does |
|------|------|
| [gum](https://github.com/charmbracelet/gum) | prompts, pickers, spinners, styling for scripts |
| [vhs](https://github.com/charmbracelet/vhs) | record terminal demos from tape files |
| [glow](https://github.com/charmbracelet/glow) | read Markdown in the terminal |
| [freeze](https://github.com/charmbracelet/freeze) | screenshots of code and terminal output as PNG or SVG |
| [mods](https://github.com/charmbracelet/mods) | pipe text through an LLM |
| [crush](https://github.com/charmbracelet/crush) | agentic coding assistant in the terminal |
| [soft-serve](https://github.com/charmbracelet/soft-serve) | self hosted git server with a TUI |
| [wishlist](https://github.com/charmbracelet/wishlist) | SSH directory |
| [pop](https://github.com/charmbracelet/pop) | send email from the terminal |
| [skate](https://github.com/charmbracelet/skate) | key value store |
| [sequin](https://github.com/charmbracelet/sequin) | explain escape sequences |

## Where to look when a skill here is not enough

- The [sibling skills](https://github.com/tpscloud/charmedimsure/tree/main/skills)
  and [examples](https://github.com/tpscloud/charmedimsure/tree/main/examples)
  live in the repository; an isolated skill install does not include them.
  To run a linked example from a checkout, enter the
  `examples` directory and run `go run ./charmed-bubbletea/compose` (or the
  corresponding example).
- Module source: `$(go env GOMODCACHE)/charm.land/<lib>/v2@<version>/`.
  `go doc -all charm.land/<lib>/v2` from any module that requires it.
- Upstream examples live in each GitHub repo under `examples/` (not in
  the module zip for most libraries).

## Mistakes that cost time

- Passing v1 and v2 values of the same library across an API boundary. A v1
  `lipgloss.Style` will not satisfy a v2 signature and vice versa.
- Using `@latest` instead of the pinned versions. It selects versions these
  skills have not been checked against.
- Trusting upstream prose over `go doc`. Some READMEs and doc comments still
  show v1 signatures (`View() string` in Huh's README, `tea.NewFrame` in a
  Bubbles doc comment).
