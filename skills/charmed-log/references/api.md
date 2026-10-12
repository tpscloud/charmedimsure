# Charm Log v2 logger API

Constructors, setters, levels, and the `Styles` fields in Charm Log v2.0.1.
Return to the [main skill](../SKILL.md) for the workflow and examples.

## Constructors

- The package level logger (`log.Default()`, used by `log.Info` and the other
  package functions) writes to stderr at info level with timestamps; it is
  built with `ReportTimestamp: true`.
- `log.NewWithOptions(w, log.Options{...})` sets everything at once.
- `log.New(w)` is `NewWithOptions(w, Options{})`: info level, no prefix, and
  no timestamp. Call `SetReportTimestamp(true)` on a custom logger if you
  want the same look as the package logger.

## Setters

Setters exist for everything in `Options`:

- `SetLevel`, `SetPrefix`, `SetReportTimestamp`, `SetReportCaller`,
  `SetTimeFormat`, `SetCallerOffset`.
- `SetTimeFunction`: `log.NowUTC` for UTC stamps.
- `SetCallerFormatter`: `ShortCallerFormatter` or `LongCallerFormatter`.
- `SetOutput`: swaps the writer but keeps the color profile detected when
  the logger was created.
- `SetFormatter`: `TextFormatter`, `JSONFormatter`, or `LogfmtFormatter`.
- `SetStyles`: takes a `*Styles` (below).
- `SetColorProfile`: takes a `colorprofile.Profile` from
  `github.com/charmbracelet/colorprofile`.

`WithPrefix(p)` returns a child with a new prefix; `With(pairs...)` returns a
child that carries fixed key value pairs. `Helper()` marks a function so the
caller shown is its caller.

## Levels

`DebugLevel`, `InfoLevel`, `WarnLevel`, `ErrorLevel`, `FatalLevel`. `Fatal`
logs and calls `os.Exit(1)`. `Print` ignores the level. `ParseLevel("warn")`
turns config text into a level. Every level function takes a message (any
value) followed by key value pairs; the `f` variants (`Infof` and so on)
take a format string instead.

## Context

Pass a logger through a context with `log.WithContext(ctx, logger)` and
recover it with `log.FromContext(ctx)` (falls back to the default logger).

## Styles

`log.DefaultStyles()` returns a `*Styles` of Lip Gloss v2 values:

| Field | Holds |
|-------|-------|
| `Timestamp`, `Caller`, `Prefix`, `Message`, `Key`, `Value`, `Separator` | one style each |
| `Levels` | map by level; `SetString` on a level style changes the label text |
| `Keys`, `Values` | maps by key name for per field overrides |
