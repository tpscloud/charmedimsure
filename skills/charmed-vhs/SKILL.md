---
name: charmed-vhs
description: "Record terminal demos as a terminal GIF, MP4, or WebM demo video from a .tape text file with vhs (Charm, charmbracelet): tape syntax, Output, Require, Set options (size, font, theme, typing speed, window bar, margin), typing and key commands, Sleep and Wait, Hide and Show, Screenshot, Source, and CI usage. Use when documenting a CLI or TUI with an animated demo, screencast, terminal recording, or screen recording for a README, or asserting its output in tests with golden files. Not for writing the program itself."
license: "MIT"
---

# vhs

Record terminal demos as GIF, MP4, or WebM from a `.tape` script with vhs,
and check terminal output in tests with golden files, without uploading
anything by accident.

The documented baseline is [vhs 0.12.1](https://github.com/charmbracelet/vhs/tree/v0.12.1);
check `vhs --version` and validate with the version you will use to render.
VHS needs `ttyd` and `ffmpeg` on its process PATH; package-manager
installations may provide them through a wrapper.

## Workflow

1. Write a tape with output paths, fixed settings, and commands that work from
   its intended working directory. Put every `Set` at the top.
2. Check that the tape exists and is readable, then run `vhs validate demo.tape`.
3. Check that `VHS_PUBLISH` is unset, then render with `vhs demo.tape`. If a
   command fails or a wait times out, fix the command or expected text and
   rerun validation and rendering.
4. Inspect the generated animation for readable text, complete output, and the
   intended final frame. Review its contents before sharing it.

## Commands

```sh
vhs demo.tape                  # render the tape's Output files
vhs demo                       # same; the render command appends .tape
vhs new demo.tape              # write a documented starter tape
vhs validate ./*.tape          # parse without rendering
vhs record > demo.tape         # capture your own session as a tape (-s shell)
vhs themes                     # list available theme names
vhs publish demo.gif           # upload to vhs.charm.sh and print a URL
vhs -o out.gif demo.tape       # replace the GIF path; other Output lines still render
```

The render command appends `.tape` to a file argument that lacks it, so a tape
saved without that extension does not open. `validate` reads its arguments as
given. `-q` silences progress. `vhs serve` runs the SSH server used by
`ssh vhs.charm.sh`. `vhs publish` uploads the file to a public host. A plain
render also uploads the GIF when `VHS_PUBLISH=true` is set in the environment
or `-p`/`--publish` is passed. Check for that variable before rendering
(`env -u VHS_PUBLISH vhs demo.tape` clears it), and publish only when the
user asks for a shareable URL.

## A tape

Save this Bash-only example as `demo.tape` in a writable directory. It prints
a local task summary; it needs no project binary or fixture.

```
# demo.tape
Output demo.gif
Output demo.mp4

Require bash

Set Shell "bash"
Set FontSize 24
Set Width 1200
Set Height 640
Set Theme "Catppuccin Mocha"
Set Padding 20
Set TypingSpeed 40ms
Set WindowBar Colorful
Set BorderRadius 8
Set MarginFill "#1e1e2e"
Set Margin 40

Hide
Type "export PS1='$ '" Enter
Type "clear" Enter
Show

Type "echo 'project: site'; echo 'task: build'; echo rea''dy" Enter
Wait+Screen /ready/
Sleep 2s
```

Commands run top to bottom in a real shell inside a headless terminal.
`Output` may appear more than once and the extension picks the format.
`Require` fails fast if a program is missing. `Sleep` after every command so
output is visible, `Wait` on a known completion string rather than guessing a
sleep, and end with `Sleep 2s` so the last frame lingers before the loop.

Read [references/tape-syntax.md](references/tape-syntax.md) for every `Set`
option (fonts, themes, `Columns` and `Rows`, window bar, framerate, loop
offset, wait defaults), the rules for where settings apply, and the full
action list (keys and chords, `Wait` forms, `Screenshot`, text `Output`,
clipboard, `Source`, scrolling).

## Recording a TUI

This tape records the repository's
[Lip Gloss overlay example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-lipgloss/overlay/main.go).
Save it as `tui.tape` in a checkout's `examples` directory and run VHS there.
It builds `./vhs-app` and writes `tui.gif` and `tui.png`; keep those generated
files out of commits. An isolated skill install does not include the Go example.

```
Output tui.gif
Require go
Set Shell "bash"
Set Width 1000
Set Height 600
Set WaitTimeout 60s

Hide
Type "export PS1='$ '" Enter
Type "go build -o ./vhs-app ./charmed-lipgloss/overlay && echo build-''ok" Enter
Wait+Screen /build-ok/
Type "clear" Enter
Show

Type "./vhs-app" Enter
Wait+Screen /space toggles the dialog/
Space
Wait+Screen /a dialog on its own layer/
Sleep 1s
Screenshot tui.png
Space
Sleep 1s
Type "q"
Sleep 1s
```

`Wait+Screen` is the tool for alternate screen apps, since the prompt never
appears there. Keep `Width` and `Height` modest so text stays legible in a
README.

Wait for a success sentinel after the build, not for the prompt: a failed
`go build` returns to the prompt too, and with a fixed output path a stale
binary from an earlier run would then be recorded. With `&& echo build-''ok`
and `Wait+Screen /build-ok/`, the shell prints the sentinel only on success.
The split quoting keeps `build-ok` out of the echoed command itself, which
would otherwise satisfy a screen-wide match before the build finishes. A
failed build makes the `Wait` time out and vhs exits non zero instead of
recording the wrong program. Build into a path under the project
(`./vhs-app`, which you should ignore in Git) rather than a shared `/tmp` name
another run could have left behind.

## CI and tests

Read [references/ci.md](references/ci.md) before adding vhs to a pipeline: it
covers guarding `vhs validate` against missing files, the pinned
`charmbracelet/vhs-action` step, and what to pin for golden-file and
screenshot fixtures.

## Mistakes that cost time

- Forgetting `Output`. The tape runs and produces nothing.
- A `Set` after `Hide` or the first `Type`. It is ignored with a warning; move
  settings to the top.
- No `Sleep` after `Enter`. The command's output never makes it into a frame.
- Waiting on the prompt inside an alt screen program. Use `Wait+Screen` with
  text the program draws.
- Recording with a font that is not installed. vhs falls back silently and
  the demo looks different on every machine.

## Older versions

- 0.12.0 added the `.tape` suffix on render arguments; earlier versions opened
  a tape saved without that extension.
- 0.12.0 added `Set Columns` and `Set Rows`.
- Before 0.12.1, a failed ffmpeg encode was logged and vhs exited zero.
