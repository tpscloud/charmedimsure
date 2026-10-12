---
name: charmed-log
description: "Structured, leveled, styled logging in Go with Charm Log v2 (charm.land/log/v2, the charmbracelet logger): package level and custom loggers, key value pairs, levels, text, JSON, and logfmt formatters, colored log output with Lip Gloss v2 styles, slog handler and standard log adapters, and logging from inside a Bubble Tea program to a debug log file. Use when adding, configuring, or restyling Charm Log output in a Go CLI or TUI, or migrating from github.com/charmbracelet/log. Not for other Go loggers such as zap or zerolog, or for printing UI (charmed-lipgloss)."
license: "MIT"
---

# Charm Log v2

Add structured, leveled, colored logging to a Go CLI or TUI with Charm Log
v2, including JSON and logfmt output, slog adapters, and plain-text log files
from inside Bubble Tea.

Import as `"charm.land/log/v2"`. Go blocks are trimmed from
[complete programs](https://github.com/tpscloud/charmedimsure/tree/main/examples/charmed-log)
that build against the pinned modules.

If the project's `go.mod` still requires `github.com/charmbracelet/log`
(pre-v2), the v2 APIs here do not apply.  Use `go doc` against the
project's version and keep its import paths unless the user asks to
migrate.

## Workflow

1. Pick the logger: the package level `log` functions for a small CLI, or
   `log.NewWithOptions(w, log.Options{...})` for a dedicated one.
2. Set the level (`log.ParseLevel` for config text), the formatter, and the
   writer when you create the logger: the color profile comes from that
   writer.
3. Log a message followed by key value pairs in even numbers; use `With` for
   pairs every line should carry.
4. Inside a Bubble Tea program, send logs to a file with `tea.LogToFileWith`.
   For code that logs through `log/slog`, hand it `slog.New(logger)`.
5. Run it and check level filtering, the output stream, and structured output
   with its consumer; in a TUI, confirm logs reach the file and not the
   screen. On a failure, compare against Mistakes below.

## Basics

The package level logger writes to stderr at info level with timestamps
(see the [basic example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-log/basic/main.go)):

```go
// examples/charmed-log/basic/main.go
log.Info("starting", "project", "site")
log.Debug("hidden at info level")
log.Warn("disk almost full", "free", "1.2G")
log.Error("check failed", "err", errors.New("invalid configuration"))
log.Print("no level prefix, always printed")

log.SetLevel(log.DebugLevel)
log.SetReportCaller(true)
```

`Fatal` logs and calls `os.Exit(1)`.  `Print` ignores the level. A dedicated
logger:

```go
// examples/charmed-log/basic/main.go
logger := log.NewWithOptions(os.Stderr, log.Options{
    ReportTimestamp: true,
    TimeFormat:      time.Kitchen,
    Prefix:          "tasks",
    Level:           log.InfoLevel,
})

projectLog := logger.With("project", "site", "task", "build")   // child with fixed pairs
projectLog.Info("completed")
logger.Infof("%d of %d done", 3, 6)
```

`log.New(w)` has no timestamp, unlike the package logger; call
`SetReportTimestamp(true)` for the same look.

Read [references/api.md](references/api.md) for the full setter list,
levels, caller and time options, context helpers, and every `Styles` field.

## Formatters

The [structured example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-log/structured/main.go)
includes formatters and custom styles. Run it from a repository checkout:

```sh
cd examples
go run ./charmed-log/structured
```

```go
// examples/charmed-log/structured/main.go
jsonLogger := log.NewWithOptions(os.Stdout, log.Options{Formatter: log.JSONFormatter, ReportTimestamp: true})
jsonLogger.Info("json line", "project", "api", "ok", true)

fmtLogger := log.New(os.Stdout)
fmtLogger.SetFormatter(log.LogfmtFormatter)
```

`TextFormatter` (default, colored, human oriented), `JSONFormatter`, and
`LogfmtFormatter`.  Colors are chosen from the output writer's profile,
detected once when the logger is created: a logger built on a file writes
plain text.  Force a profile with
`SetColorProfile(colorprofile.TrueColor)` from
`github.com/charmbracelet/colorprofile`.

## Styles

`log.DefaultStyles()` returns a `*Styles` of Lip Gloss v2 values with a style
per element, a `Levels` map, and `Keys` and `Values` maps for per field
overrides:

```go
// examples/charmed-log/structured/main.go
styles := log.DefaultStyles()
styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
    SetString("ERROR").
    Padding(0, 1).
    Background(lipgloss.Color("204")).
    Foreground(lipgloss.Color("0"))
styles.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
styles.Values["err"] = lipgloss.NewStyle().Bold(true)

styled := log.New(os.Stderr)
styled.SetStyles(styles)
```

`SetString` on the level style changes the label text.

## slog and the standard library

A `*log.Logger` is a `slog.Handler`, so libraries that log through
`log/slog` can be pointed at it.  It also produces a standard `*log.Logger`
for APIs that demand one (see the [adapter example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-log/slog/main.go)):

```go
// examples/charmed-log/slog/main.go
handler := log.NewWithOptions(os.Stderr, log.Options{ReportTimestamp: true, Level: log.DebugLevel})

sl := slog.New(handler)
sl.Info("from slog", "project", "worker")
// Only code that receives sl logs through Charm. Package level calls such
// as slog.Info still use slog's own default until slog.SetDefault(sl).

std := handler.StandardLog(log.StandardLogOptions{ForceLevel: log.WarnLevel})
std.Println("from the standard logger, forced to warn")
```

`slog.New(handler)` builds one logger; pass it to the code that should use
it.  Redirecting every package level `slog` call in the process is a
separate, global decision: `slog.SetDefault(sl)`, which also captures the
standard library `log` package.  Do that in `main`, not in a library.

## Inside Bubble Tea

The renderer owns the terminal, so logs must not go to stdout or stderr
while the program runs.  Any other writer works (a file, an in memory
buffer the app shows in a pane, a remote sink); a file is the usual choice
because you can tail it.  Bubble Tea's `LogToFileWith` accepts the Charm
logger because it implements `SetOutput` and `SetPrefix`
(see the complete [Bubble Tea example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-log/tea/main.go)):

```go
// examples/charmed-log/tea/main.go
f, err := tea.LogToFileWith("debug.log", "tui", log.Default())
if err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
}
defer f.Close()
// log.Default() detected its profile for stderr; match it to the file.
log.SetColorProfile(colorprofile.Detect(f, os.Environ()))
log.SetLevel(log.DebugLevel)
log.SetReportTimestamp(true)

// Inside Update, where msg is a tea.Msg:
if k, ok := msg.(tea.KeyPressMsg); ok {
    log.Debug("key", "key", k.String())
}
```

The program creates `debug.log` in its working directory. Run
`tail -f debug.log` from that directory in a second terminal.

## Migrating from v1

| v1                                          | v2                                             |
|---------------------------------------------|------------------------------------------------|
| `github.com/charmbracelet/log`              | `charm.land/log/v2`                            |
| `SetColorProfile(termenv.TrueColor)`        | `SetColorProfile(colorprofile.TrueColor)`      |
| `Styles` fields of v1 `lipgloss.Style`      | `charm.land/lipgloss/v2` styles                |

These skills are pinned to and verified against v2.0.1. The
[v2.0.1 upgrade guide](https://github.com/charmbracelet/log/blob/v2.0.1/UPGRADE_GUIDE_V2.md)
covers the changes. `termenv` and `go-osc52` are no longer dependencies.
After migrating, run your application's tests, fix failures, and rerun, then
do the checks in step 5 of the workflow.

## Mistakes that cost time

- Odd numbers of key value arguments.  The last key prints with
  `ErrMissingValue` as its value.
- `Fatal` inside library code.  It exits the process; return the error.
- Logging to stderr from inside a TUI and getting garbage on screen.  Use
  the file approach above.
- Calling `SetOutput(file)` on a logger created for a terminal. It swaps the
  writer but keeps the old color profile, so the file fills with ANSI
  escapes. Create the logger on the file, or set the profile with
  `SetColorProfile` as in the Bubble Tea example.
- Mixing v1 and v2 Lip Gloss styles when customizing `Styles`.  The types
  are different packages and the compile error names both.
