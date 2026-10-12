# Glamour v2 renderer options

The `TermRenderer` options and render entry points in Glamour v2.0.1. Return
to the [main skill](../SKILL.md) for the render workflow.

## Style selection

Built in style names: `dark`, `light`, `pink`, `dracula`, `tokyo-night`,
`ascii` (no color), `notty`.

- `WithStandardStyle(name)` selects a built in style by name.
  `NewTermRenderer()` does not select a named style by itself.
- `WithStylePath(s)` accepts a name or a path to a JSON file.
- `WithEnvironmentConfig()` and `glamour.RenderWithEnvironmentConfig(in)`
  read the name or path from the `GLAMOUR_STYLE` environment variable.
  These environment-based helpers fall back to `dark` when `GLAMOUR_STYLE`
  is empty.
- `WithStyles(cfg)` takes an `ansi.StyleConfig`; `WithStylesFromJSONFile(path)`
  and `WithStylesFromJSONBytes(b)` load one from JSON. See
  [styles.md](styles.md).

`NewTermRenderer` returns an error for an unknown style name or an
unreadable style path.

## Layout and content

| Option | Effect |
|--------|--------|
| `WithWordWrap(n)` | wrap width; default 80; 0 disables wrapping |
| `WithEmoji()` | `:tada:` becomes the character |
| `WithPreservedNewLines()` | keep newlines from the source |
| `WithBaseURL(u)` | base for relative links |
| `WithChromaFormatter(name)` | formatter for code block highlighting |
| `WithTableWrap(bool)` | tables wrap by default; `false` truncates overflowing cells with an ellipsis |
| `WithInlineTableLinks(bool)` | `true` renders table links inline instead of as a list below the table |

## Entry points

- `glamour.Render(in, style)` is a one shot render with a named style.
- `r.Render(s)` renders a string with a reusable renderer; `RenderBytes`
  accepts a byte slice.
- Buffered I/O: write the Markdown to the renderer, call `Close()` and check
  its error, then read the rendered output. `Write` alone does not render it.
