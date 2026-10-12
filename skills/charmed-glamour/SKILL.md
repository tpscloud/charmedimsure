---
name: charmed-glamour
description: "Render Markdown to styled terminal text in Go with Glamour v2 (charm.land/glamour/v2, formerly github.com/charmbracelet/glamour): one shot Render, reusable TermRenderer with word wrap and emoji, built in styles (dark, light, dracula, tokyo-night, pink, ascii), custom StyleConfig in Go or JSON, color downsampling through Lip Gloss, and a scrolling viewer inside Bubble Tea. Use when showing a README, .md file, help text, changelog, or any Markdown in a CLI or TUI, or when asked to pretty-print Markdown in the terminal. Not for editing Markdown or for general styling (charmed-lipgloss)."
license: "MIT"
---

# Glamour v2

Import as `"charm.land/glamour/v2"`; styles live in
`"charm.land/glamour/v2/styles"` and the config types in
`"charm.land/glamour/v2/ansi"`. Pinned and verified against v2.0.1. Go blocks
headed `// examples/...` are trimmed from
[complete programs](https://github.com/tpscloud/charmedimsure/tree/main/examples/charmed-glamour)
that build against the pinned modules; `// excerpt` blocks are API sketches.

If the project's `go.mod` still requires `github.com/charmbracelet/glamour`
(pre-v2), the v2 APIs here do not apply.  Use `go doc` against the
project's version and keep its import paths unless the user asks to
migrate.

## Workflow

1. Pick `"dark"` or `"light"`: outside Bubble Tea call
   `lipgloss.HasDarkBackground(os.Stdin, os.Stdout)`; inside, read
   `tea.BackgroundColorMsg`.
2. Build one renderer with `glamour.NewTermRenderer(glamour.WithStandardStyle(style), glamour.WithWordWrap(n))`
   and reuse it.
3. Render the Markdown with `r.Render(md)` and handle the error.
4. Print with `lipgloss.Print` outside Bubble Tea; inside, hand the string to
   a viewport or the view.
5. Check in a real terminal: light and dark styles, narrow widths, tables,
   code blocks, and links. For wrong colors or unexpected wrapping, see
   Mistakes.

## Rendering

See the complete
[render example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-glamour/render/main.go):

```go
// examples/charmed-glamour/render/main.go
out, err := glamour.Render(doc, "dark") // one shot with a named style
if err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
}
lipgloss.Print(out)

style := "dark"
if !lipgloss.HasDarkBackground(os.Stdin, os.Stdout) {
    style = "light"
}
r, err := glamour.NewTermRenderer(
    glamour.WithStandardStyle(style),
    glamour.WithWordWrap(60),
    glamour.WithEmoji(),
)
if err != nil { // unknown style name or unreadable style path
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
}
out, err = r.Render("Rendered again at 60 columns :tada:")
if err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
}
lipgloss.Print(out)
```

The options used most: `WithWordWrap(n)`
(default 80; 0 disables), `WithEmoji()` (`:tada:` becomes the character),
and `WithStylePath(s)` (a name or a path to a JSON file).

Read [references/options.md](references/options.md) for the other built in
style names, every renderer option, `GLAMOUR_STYLE` environment configuration, `RenderBytes`, and the
buffered `Write` and `Close` API.

## Custom styles

A style is an `ansi.StyleConfig`. Start from a built in config and edit, as
in the
[custom style example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-glamour/custom/main.go):

```go
// examples/charmed-glamour/custom/main.go
cfg := styles.DraculaStyleConfig // a value copy; safe to edit

cfg.Document.Margin = new(uint(0)) // new(value) needs Go 1.26 or later
cfg.H1.Prefix = ""
cfg.H1.BackgroundColor = new("#5A2D82")
cfg.H1.Color = new("#FFFFFF")
cfg.Code.StylePrimitive.Color = new("#C4A7F5")
cfg.Link.Underline = new(true)

r, err := glamour.NewTermRenderer(
    glamour.WithStyles(cfg),
    glamour.WithWordWrap(70),
)
if err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
}
```

Read [references/styles.md](references/styles.md) for the full list of
StyleConfig blocks and primitives, the exported configs, and the JSON style
format.

## A Markdown viewer in Bubble Tea

Render to a string, hand it to a viewport, and re-render when the width or
the theme changes. The complete
[viewer example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-glamour/viewer/main.go)
defines `model`, `minWrap`, and the message handlers excerpted below.
From a repository checkout, run it in a terminal:

```sh
cd examples
go run ./charmed-glamour/viewer ../README.md
```

```go
// examples/charmed-glamour/viewer/main.go
func (m model) Init() tea.Cmd { return tea.RequestBackgroundColor }

func (m model) render() (string, error) {
    style := "light"
    if m.isDark {
        style = "dark"
    }
    r, err := glamour.NewTermRenderer(
        glamour.WithStandardStyle(style),
        glamour.WithWordWrap(max(m.width-2, minWrap)),
    )
    if err != nil {
        return "", err
    }
    return r.Render(m.md)
}

// rerender is called from both message handlers below.
func (m *model) rerender() {
    out, err := m.render()
    if err != nil {
        out = "render error: " + err.Error()
    }
    m.vp.SetContent(out)
}

case tea.BackgroundColorMsg:
    m.isDark = msg.IsDark()
    if m.ready {            // the size may have arrived first
        m.rerender()
    }
case tea.WindowSizeMsg:
    m.width = msg.Width
    // (re)size the viewport, then:
    m.rerender()
```

`WindowSizeMsg` and `BackgroundColorMsg` arrive in either order, so both
handlers re-render; otherwise a dark terminal keeps the light rendering
until the next resize.  `minWrap` keeps the wrap width positive on tiny
terminals.  Rendering is not free on
large documents; do it on resize, theme, and content change, not on every
frame.

## Migrating from v1

| v1                                     | v2                                                   |
|----------------------------------------|------------------------------------------------------|
| `github.com/charmbracelet/glamour`     | `charm.land/glamour/v2` (and `/ansi`, `/styles`)     |
| `glamour.WithAutoStyle()`              | removed; pass `WithStandardStyle` (see Mistakes)     |
| `glamour.WithColorProfile(p)`          | removed; print with `lipgloss.Print`                 |
| `StylePrimitive.Overlined`             | removed                                              |

Custom `ansi.MarginWriter` users must call `Close()`. See the
[v2.0.1 upgrade guide](https://github.com/charmbracelet/glamour/blob/v2.0.1/UPGRADE_GUIDE_V2.md)
for migration details. Its background-detection example omits the input and
output arguments required by the pinned Lip Gloss version; use the call in
the render example.

After migrating, run your application's tests, fix failures, and rerun. Then
do the terminal check from step 5 of the workflow. For a viewer, also resize
the terminal and confirm that scrolling and background-color updates still
work.

## Mistakes that cost time

- Expecting the style to follow the terminal theme. Glamour v2 has no auto
  style and does not detect the terminal's background or color support, and
  `NewTermRenderer()` does not select a named style by itself. Detect light
  or dark yourself and pass `WithStandardStyle` explicitly, even where the
  upgrade guide suggests an implicit dark style.
- Printing with `fmt` outside Bubble Tea. Only the `lipgloss.Print` family
  downsamples colors to what the terminal can show, so `fmt.Print` gives
  wrong colors on a 256 color terminal. Inside Bubble Tea the program's
  renderer downsamples for you.
- Rendering a wide table or long code block with `WithWordWrap` smaller
  than the content and being surprised by the wrapping. For tables, see
  `WithTableWrap` in [references/options.md](references/options.md).
- Re-creating a renderer for every call in a hot path.  Build once, reuse.
