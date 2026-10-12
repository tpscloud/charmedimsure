# Charm v1 to v2 migration checklist

The per-library API changes for moving the Charm stack from the
`github.com/charmbracelet/<name>` v1 modules to the `charm.land` v2 modules
pinned in this repository. Return to the [main skill](../SKILL.md) for the
library map and import table.

## Order

Move the libraries that exchange types together: Bubble Tea, Bubbles, Huh,
and Lip Gloss share `tea.Model`, `tea.Cmd`, and `lipgloss.Style`, so a
partial upgrade inside that set does not compile.  Log and Glamour only
need a matching Lip Gloss major.  Code that uses none of the shared types
can lag behind on its own schedule.

## Steps

1. Replace import paths.  `github.com/charmbracelet/X` becomes
   `charm.land/X/v2` for bubbletea, bubbles, lipgloss, huh, glamour, log.
   Subpackages keep their suffix (`bubbles/v2/list`,
   `lipgloss/v2/table`).
2. Bubble Tea: `View() string` becomes `View() tea.View`; key-press handlers
   use `tea.KeyPressMsg` instead of the old `tea.KeyMsg` struct (v2 keeps
   `tea.KeyMsg` as an interface); `msg.Type`, `msg.Runes`, `msg.Alt` become
   `Code`, `Text`, `Mod`; `" "` becomes `"space"`; mouse messages split
   into `MouseClickMsg`, `MouseReleaseMsg`, `MouseWheelMsg`,
   `MouseMotionMsg` with `msg.Mouse()`; program options and commands for
   alt screen, mouse, title, cursor, focus, and paste become `tea.View`
   fields; `tea.WindowSize()` becomes `tea.RequestWindowSize`;
   `tea.Sequentially` becomes `tea.Sequence`; `p.Start` becomes `p.Run`.
3. Lip Gloss: `Color` is a function returning `color.Color`;
   `AdaptiveColor` becomes `LightDark(isDark)(light, dark)` (or the
   `compat` package); `CompleteColor` becomes `Complete(profile)(...)`;
   `HasDarkBackground()` takes `(in, out)`; the renderer API is gone;
   print with `lipgloss.Println` outside Bubble Tea; whitespace options
   collapse into `WithWhitespaceStyle`.
4. Bubbles: `Width` and `Height` fields become `SetWidth`, `SetHeight`,
   `Width()`, `Height()`; `DefaultKeyMap` vars become functions;
   `NewModel` aliases are gone; most style constructors take `isDark`
   (table and filepicker do not);
   progress gradient options become `WithColors`, `WithDefaultBlend`,
   `WithScaled`; `spinner.Tick` is a method; `stopwatch.NewWithInterval`
   becomes an option.
5. Huh: themes are `func(isDark bool) *Styles`; field level
   `WithAccessible` is gone (form level only); the `accessibility`
   package is gone; `Input.Password` is deprecated for
   `EchoMode(EchoModePassword)`.
6. Glamour: `WithAutoStyle` and `WithColorProfile` are gone; pick a style
   name and print with `lipgloss.Print`; `Overlined` is gone.
7. Log: `SetColorProfile` takes a `colorprofile.Profile`; styles are v2
   Lip Gloss styles.

## Upstream guides

The pinned modules include detailed upgrade guides:
[Bubble Tea](https://github.com/charmbracelet/bubbletea/blob/v2.1.0/UPGRADE_GUIDE_V2.md),
[Bubbles](https://github.com/charmbracelet/bubbles/blob/v2.2.1/UPGRADE_GUIDE_V2.md),
[Lip Gloss](https://github.com/charmbracelet/lipgloss/blob/v2.0.6/UPGRADE_GUIDE_V2.md),
[Huh](https://github.com/charmbracelet/huh/blob/v2.0.3/UPGRADE_GUIDE_V2.md),
[Glamour](https://github.com/charmbracelet/glamour/blob/v2.0.1/UPGRADE_GUIDE_V2.md), and
[Log](https://github.com/charmbracelet/log/blob/v2.0.1/UPGRADE_GUIDE_V2.md).

## Validation

After each migration step, run `go test ./...` in your application module.
Fix the first relevant compile or test failure, then rerun. If a type mismatch
remains, inspect the imports and resolved versions with `go list -m all` before
changing more code. Finish by checking keys, resizing, and asynchronous updates
in a terminal.
