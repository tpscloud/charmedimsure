# vhs tape settings and actions

Every `Set` option and action command in a vhs 0.12.1 tape, with the rules
for where settings may appear. Read this when a tape needs a setting or key
that the example in [SKILL.md](../SKILL.md) does not show.

## Settings

| Setting | Meaning | Note |
|---------|---------|------|
| `Set Shell "bash"` | shell to drive | bash, zsh, fish, and others |
| `Set Width 1200`, `Set Height 600` | terminal pixel size | |
| `Set Columns 80`, `Set Rows 24` | terminal size in cells | not with `Width` or `Height`; see below |
| `Set FontSize 32`, `Set FontFamily "JetBrains Mono"` | font | family must be installed |
| `Set LetterSpacing 1.0`, `Set LineHeight 1.2` | metrics | |
| `Set Theme "Dracula"` | color theme by name | or inline JSON `{"background": "#...", ...}` |
| `Set Padding 20` | inner padding in pixels | |
| `Set Margin 40`, `Set MarginFill "#000000"` | outer margin and its fill | fill may be an image path |
| `Set WindowBar Colorful` | fake title bar | Rings, RingsRight, Colorful, ColorfulRight |
| `Set WindowBarSize 40`, `Set BorderRadius 8` | bar height, corner radius | |
| `Set TypingSpeed 50ms` | delay between typed characters | override per line with `Type@10ms` |
| `Set Framerate 50`, `Set PlaybackSpeed 1.0` | recording rate and playback | |
| `Set LoopOffset 50%` | GIF loop start frame | |
| `Set WaitTimeout 15s`, `Set WaitPattern /\$/` | defaults for `Wait` | |

`vhs themes` lists the theme names.

## Columns and Rows

`Columns` and `Rows` replace `Width` and `Height` per axis: vhs measures the
font after applying `FontSize`, `FontFamily`, `LetterSpacing`, and
`LineHeight`, then derives the pixel size, including padding, margin, and
window bar. Combining `Width` with `Columns`, or `Height` with `Rows`, fails
the render with `cannot set both ...` in either order. If the measured grid
never matches, the render fails rather than recording at the wrong size.
`vhs validate` checks none of this, and accepts a non numeric value too.

## Where settings apply

Settings apply only before the first command other than `Set`, `Output`, or
`Require`; `Hide` counts as such a command. A later `Set` is ignored with a
warning, except `Set TypingSpeed`, which works anywhere. `Set Shell` and
`Env` take effect wherever they appear.

## Actions

| Command | Effect |
|---------|--------|
| `Type "text"` | type with the current speed; `Type@20ms "fast"` overrides |
| `Enter`, `Tab`, `Space`, `Backspace`, `Delete`, `Insert`, `Escape` | press, optionally `[count]` and `@delay` |
| `Up`, `Down`, `Left`, `Right`, `PageUp`, `PageDown` | arrows and paging, same modifiers |
| `Ctrl+C`, `Ctrl+Alt+Shift+X`, `Alt+key` | chords |
| `Sleep 500ms`, `Sleep 2s` | pause |
| `Wait /regex/`, `Wait+Screen /regex/`, `Wait@10s /regex/` | block until the line (or whole screen) matches, with timeout |
| `Hide`, `Show` | stop and resume capturing frames; setup goes between |
| `Screenshot shot.png` | save a still |
| `Output out.txt` (or `.ascii`) | append the terminal text after each recorded command, for golden-file tests |
| `Copy "text"`, `Paste` | clipboard |
| `Source other.tape` | include another tape |
| `ScrollUp 3`, `ScrollDown 3` | scroll the terminal buffer |
