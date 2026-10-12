---
name: charmed-lipgloss
description: "Style and lay out terminal text in Go with Lip Gloss v2 (charm.land/lipgloss/v2, charmbracelet): ANSI colors and terminal colors, borders and boxes, padding, alignment, joining and placing blocks, light or dark adaptation, gradients, static tables, trees, lists, and layered overlays. Use when rendering styled strings or pretty terminal output for a TUI or CLI, drawing a bordered box, building a layout from measured blocks, or migrating v1 lipgloss code. Not for the event loop (charmed-bubbletea) or interactive widgets (charmed-bubbles)."
license: "MIT"
---

# Lip Gloss v2

Style and lay out terminal text in Go with Lip Gloss v2: colors that adapt to
the background, borders, padding, alignment, joined and layered blocks,
tables, trees, and lists, measured in terminal cells.

Import as `"charm.land/lipgloss/v2"`; subpackages are `table`, `tree`, `list`,
and `compat`. Pinned and verified against v2.0.6. Go blocks headed
`// examples/...` are trimmed from
[complete programs](https://github.com/tpscloud/charmedimsure/tree/main/examples/charmed-lipgloss)
that build against the pinned modules; `// excerpt` blocks are API sketches.

If the project's `go.mod` still requires `github.com/charmbracelet/lipgloss`
(pre-v2), the v2 APIs here do not apply.  Use `go doc` against the
project's version and keep its import paths unless the user asks to
migrate.

## Workflow

1. Standalone, decide light or dark once with `HasDarkBackground` before
   printing anything, and build a `LightDark` picker. Inside Bubble Tea, use
   `tea.BackgroundColorMsg` instead (see Light and dark backgrounds).
2. Define styles as values; derive variants by chaining setters.
3. Render blocks and measure them with `lipgloss.Width`, `Height`, or `Size`.
4. Compose with `JoinHorizontal`, `JoinVertical`, and `Place`, or with layers
   for overlays; use the subpackages for tables, trees, and lists.
5. Print with `lipgloss.Println` (standalone) or return the string in a
   `tea.View`.
6. Check in a real terminal and compare `lipgloss.Size` with the width you
   expected. Wrong colors usually mean output went through `fmt`; a box off by
   its frame usually means `GetFrameSize` was subtracted from `Width`. Fix and
   rerun.

## Mental model

A `Style` is a plain value.  Setters return a new value, so you build
styles by chaining and derive variants without mutating the base.  `Render`
turns a style plus text into a string carrying ANSI escape codes.  Nothing
else happens until that string reaches a terminal.

`Render` always emits full color; downsampling to what the terminal
supports happens at output.  Inside Bubble Tea v2 that is automatic.
Standalone, print through Lip Gloss
(see the [basics example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-lipgloss/basics/main.go)):

```go
// examples/charmed-lipgloss/basics/main.go
title := lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FAFAFA")).
    Background(lipgloss.Color("#7D56F4")).
    Padding(0, 1)

card := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("63")).
    Padding(1, 2).
    Width(40)

body := lipgloss.NewStyle().Italic(true).Faint(true)
warning := body.Faint(false).Italic(false).Foreground(lipgloss.Yellow) // body unchanged

lipgloss.Println(card.Render(title.Render("Lip Gloss v2") + "\n\n" + body.Render("...")))
```

`lipgloss.Print`, `Println`, `Printf`, `Fprint*`, and `Sprint*` all
downsample.  `fmt.Println` does not.  The default writer targets stdout;
for stderr use `lipgloss.Fprintln(os.Stderr, s)` or reassign
`lipgloss.Writer`.

## Colors

`lipgloss.Color(s)` is a function returning `image/color.Color`.  It accepts
hex (`"#7D56F4"`) or an ANSI index (`"212"`, `"5"`).  The sixteen basic
colors have names: `lipgloss.Red`, `lipgloss.BrightGreen`, and so on.  Any
`color.Color` works, so `color.RGBA{...}` from the standard library is
fine too.

Style setters take `color.Color`: `Foreground`, `Background`,
`BorderForeground`, `BorderBackground`, `UnderlineColor`, `MarginBackground`.

Helpers from the [colors example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-lipgloss/colors/main.go)
(the gradient uses `strings.Builder`).  `lipgloss.Alpha(c, 0.5)` also
exists and sets a color's alpha on a 0 to 1 scale:

```go
// examples/charmed-lipgloss/colors/main.go
base := lipgloss.Color("#7D56F4")
lipgloss.Lighten(base, 0.3)          // tint by 30 percent
lipgloss.Darken(base, 0.3)
lipgloss.Complementary(base)

// One color per rune from a two stop gradient. Index runes, not bytes:
// ranging over the string would sample at UTF-8 byte offsets.
word := "gradient"
runes := []rune(word)
stops := lipgloss.Blend1D(len(runes), lipgloss.Color("#FF5F87"), lipgloss.Color("#5FD7FF"))
var b strings.Builder
for i, r := range runes {
    b.WriteString(lipgloss.NewStyle().Foreground(stops[i]).Render(string(r)))
}

// Gradient borders are a style setting.
lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).
    BorderForegroundBlend(lipgloss.Color("#FF5F87"), lipgloss.Color("#5FD7FF"))
```

`Blend2D(width, height, angle, stops...)` fills a rectangle.  The gradient
loop above colors per rune; wide characters and combining marks occupy
cells differently, so for cell accurate gradients split by grapheme and
measure with `lipgloss.Width`.

### Light and dark backgrounds

Decide once whether the background is dark, then pick colors with
`LightDark`.  Standalone, as in the
[adaptive example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-lipgloss/adaptive/main.go):

```go
// examples/charmed-lipgloss/adaptive/main.go
isDark := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)   // queries the terminal
lightDark := lipgloss.LightDark(isDark)
fg := lightDark(lipgloss.Color("#1A1A1A"), lipgloss.Color("#F0F0F0")) // (light, dark)

profile := colorprofile.Detect(os.Stdout, os.Environ())
complete := lipgloss.Complete(profile)
brand := complete(lipgloss.Color("5"), lipgloss.Color("170"), lipgloss.Color("#D75FEE")) // (ansi, ansi256, truecolor)
```

Call `HasDarkBackground` before printing anything else; it exchanges an
escape sequence with the terminal.  Inside Bubble Tea, return
`tea.RequestBackgroundColor` from `Init` and read `msg.IsDark()` from the
`tea.BackgroundColorMsg` (see charmed-bubbletea).

## Text decoration

`Bold`, `Italic`, `Faint`, `Blink`, `Reverse`, `Strikethrough`,
`Underline(bool)`, `UnderlineStyle(lipgloss.UnderlineCurly)` (also
`UnderlineSingle`, `UnderlineDouble`, `UnderlineDotted`, `UnderlineDashed`,
`UnderlineNone`), and `UnderlineColor`.  `Hyperlink(url)` emits an OSC 8
link where the terminal supports it.  `Transform(func(string) string)` runs
on the text at render time (uppercase, truncation).

Style part of a string you already have. `StyleRanges` uses terminal-cell
offsets; `StyleRunes` uses rune indexes. The colors example uses ASCII, where
byte and cell offsets coincide:

```go
// excerpt: StyleRanges appears in examples/charmed-lipgloss/colors/main.go; StyleRunes has no example
lipgloss.StyleRanges(s, lipgloss.NewRange(start, end, lipgloss.NewStyle().Reverse(true)))
lipgloss.StyleRunes(s, []int{0, 3, 5}, matchedStyle, unmatchedStyle)
```

## Box model

Order from the inside out: content, padding, border, margin.

- `Padding(v...)` and `Margin(v...)` take one to four ints in CSS order
  (all; vertical horizontal; top horizontal bottom; top right bottom left).
  Single sides: `PaddingLeft`, `MarginTop`, and so on.  `PaddingChar` and
  `MarginChar` change the fill rune.
- `Width(n)` and `Height(n)` set the size of the block before margins:
  padding and border are inside that number, margins are outside it.  Text
  wraps at the width minus horizontal padding and border.  Measured with
  the pinned module: `Width(10).Padding(0, 1).Border(NormalBorder())`
  renders exactly 10 cells wide; adding `Margin(0, 1)` makes it 12.
  `MaxWidth` and `MaxHeight` truncate instead.
- `Align(pos)` or `AlignHorizontal` and `AlignVertical` position text inside
  a fixed box.  Positions are floats from 0 to 1 with names `Left`,
  `Center`, `Right`, `Top`, `Bottom`.
- `Border(b, sides...)` sets the border and optionally which sides.  Built
  in borders: `NormalBorder`, `RoundedBorder`, `ThickBorder`,
  `DoubleBorder`, `BlockBorder`, `OuterHalfBlockBorder`,
  `InnerHalfBlockBorder`, `HiddenBorder` (takes space, draws nothing),
  `ASCIIBorder`, `MarkdownBorder`.  Per side toggles `BorderTop(bool)`
  etc, per side colors `BorderTopForeground` etc.
- `Inline(true)` strips newlines and disables padding, borders, and
  margins, for a single line.
- `Inherit(other)` copies unset properties from another style. Margins,
  padding, hyperlinks, and the underlying text are not inherited; an
  inherited background also becomes the margin background unless one is
  set.
- `SetString(s)` attaches default text so `style.String()` renders it.

Every setter has a `Get*` and an `Unset*`.  `GetFrameSize()` returns the
horizontal and vertical space taken by padding, border, and margin together.
Use it to find the content area inside a decorated box, or to size a child
you render yourself inside a parent's frame.  A box meant to fill `w` cells
is `Width(w - GetHorizontalMargins())`, and the text inside it gets
`w - GetHorizontalFrameSize()` cells.

## Measuring, joining, placing

Measure with `lipgloss.Width`, `Height`, or `Size`, which count terminal
cells and ignore escape codes
(see the [layout example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-lipgloss/layout/main.go)):

```go
// examples/charmed-lipgloss/layout/main.go
box := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 1)
left := box.Width(20).Render("left column\nwith two lines")
right := box.Width(24).Render("right column")

row := lipgloss.JoinHorizontal(lipgloss.Top, left, right)       // side by side
footer := lipgloss.NewStyle().Faint(true).Render("joined with JoinVertical")
page := lipgloss.JoinVertical(lipgloss.Center, row, footer)     // stacked

w, h := lipgloss.Size(page)
framed := lipgloss.Place(w+10, h+4, lipgloss.Center, lipgloss.Center, page,
    lipgloss.WithWhitespaceChars("."),
    lipgloss.WithWhitespaceStyle(lipgloss.NewStyle().Faint(true)),
)

centered := lipgloss.NewStyle().Width(30).Align(lipgloss.Center).Render("centered text")
long := "Lip Gloss wraps at word boundaries and preserves ANSI sequences across the break."
lipgloss.Println(lipgloss.Wrap(long, 30, "")) // word wrap, preserves ANSI and links
```

`JoinHorizontal` takes a vertical alignment for blocks of different height;
`JoinVertical` takes a horizontal alignment.  `Place`, `PlaceHorizontal`,
and `PlaceVertical` pad a block out to a size; in Bubble Tea, feed them the
values from `tea.WindowSizeMsg`.

## Overlays, tables, trees, and lists

- Read [references/overlays.md](references/overlays.md) for layers, the
  compositor, Z order, and mouse hit testing (a modal over a full screen
  view).
- Read [references/tables.md](references/tables.md) for static tables:
  borders, `StyleFunc` and the header row index, data sources, and windowing.
- Read [references/trees-and-lists.md](references/trees-and-lists.md) for
  trees, nested lists, enumerators, and per item styling.

## Migrating from v1

Read [references/migration.md](references/migration.md) for the v1 to v2
mapping table (renderer removal, `Color` as a function, `AdaptiveColor` and
the `compat` package, whitespace options) and the steps to check a migrated
program.

## Mistakes that cost time

- Measuring with `len(s)` or `utf8.RuneCountInString`.  Use
  `lipgloss.Width`; escape codes and wide runes both break the byte count.
- Calling `HasDarkBackground` inside a running Bubble Tea program.  It
  fights the renderer for the terminal.  Use the message instead.
- Printing styled output with `fmt` in a standalone tool and then reporting
  broken colors on a 256 color terminal.
- Forgetting that `Width` includes padding and border but not margin.
  Subtracting `GetFrameSize` from a number you then pass to `Width` shrinks
  the box twice; subtract only the margins (see Box model).
- Expecting `Inherit` to copy padding or margins.  It does not.
- Passing a v1 `lipgloss.Style` to a v2 library (or the reverse).  The
  types are different packages; the compiler error names both paths.
