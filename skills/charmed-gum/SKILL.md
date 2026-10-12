---
name: charmed-gum
description: "Write interactive shell scripts with the gum CLI (Charm, charmbracelet): choose, filter, confirm, input, write, file, spin, table, pager, style, join, format, and log, with their flags, exit codes, environment variable defaults, and Bash examples. Use when a Bash or shell script needs a prompt, menu, picker, fuzzy finder, yes or no question, text input, spinner, CSV table, pager, styled output with borders and terminal colors, rendered Markdown, or log lines, without writing Go. Not for Go programs (charmed-huh, charmed-bubbletea)."
license: "MIT"
---

# gum

Written from gum 2.0.2 (its source, `gum --help`, and each subcommand's help),
the baseline for this skill; run `gum --version` first (2.x prints
`gum version 2.0.2`) and confirm a flag with `gum <cmd> --help` on the
installed version before relying on it.

If the request names a subcommand, go straight to it. Full flags and keys for
every subcommand are in [references/commands.md](references/commands.md).

## Workflow

1. Check the installed version with `gum --version`, or guard the script with
   `gum version-check '>=2.0.0' || exit 1`, which exits 1 outside the range.
2. Pick the command: `choose` for a short list, `filter` past a couple dozen
   items, `confirm` for yes/no, `input` or `write` for text, `file` for a
   path, `spin` around a slow command, `style`, `join`, `format`, `table`,
   `pager`, or `log` for output.
3. Capture the result with `$(...)` and chain `|| exit 1` so an abort stops
   the script; then check for an empty result where the command allows one.
4. Run `bash -n script.sh` to catch syntax errors before running it.
5. Run the script in a terminal and exercise the edges: Ctrl+c, Esc, a
   declined confirm, and an empty selection. If a prompt does not show or a
   value is wrong, check the rules in [How gum behaves](#how-gum-behaves)
   and the command's section in [references/commands.md](references/commands.md).

## How gum behaves

- Prompts draw on stderr; the result goes to stdout. Capture with `$(...)`
  and the UI still shows.
- Ctrl+c aborts with exit status 130, and Esc aborts with a nonzero status in
  every prompt except `table` (in `filter`, the first Esc only leaves the
  search box). The exit status is the reliable signal: some commands can
  still print a preselected value on abort, so test `$?` (or chain
  `|| exit 1`) before using the captured variable. `gum table` is the
  exception: Esc or q exits 0 with an empty line, so check the captured row
  as well.
- Many flags expose environment defaults such as `GUM_CHOOSE_CURSOR` and
  `GUM_INPUT_PLACEHOLDER`; check the subcommand help for the exact names.
  Selection flags such as `choose --limit` and `--no-limit` have no env
  binding. `gum style` uses bare names (`FOREGROUND`, `BORDER`, `PADDING`).
- Style flags share a shape: `--<part>.foreground`, `--<part>.background`,
  and often `--padding "0 1"`. Colors are ANSI indexes or hex.
- `style`, `join`, `format`, and `table --print` write through a color
  profile writer: when stdout is not a terminal (a pipe, a file, or `$(...)`)
  the ANSI styling is stripped and the borders and layout stay.
  `CLICOLOR_FORCE=1` keeps the colors and `NO_COLOR=1` removes them.
- `--timeout 30s` ends an interactive command after the duration. `confirm`
  then returns its default; the other prompts print "timed out" and exit 124.
- Options arrive as arguments or one per line on stdin. Output is one
  selection per line; change with `--output-delimiter`.

## Prompts

```sh
# pick one, or several (x or tab toggles, a selects all)
format=$(gum choose --header "Format" JSON YAML TOML) || exit 1
tasks=$(printf '%s\n' "List files" "Show size" | \
  gum choose --no-limit --selected='*' --ordered --header "Tasks") || exit 1
[[ -n $tasks ]] || exit 0   # --no-limit allows confirming an empty selection

# fuzzy search over a long list
file=$(find . -type f -name '*.md' | gum filter --placeholder "Which Markdown file?") || exit 1

# text
name=$(gum input --placeholder "project name" --header "Name") || exit 1
msg=$(gum write --header "Commit message" --placeholder "what changed?" --char-limit 500) || exit 1

# a directory (needs at least one subdirectory to choose from)
dir=$(gum file . --directory --file=false --all) || exit 1

# yes or no
if gum confirm "List files in the current directory?"; then
  ls -la .
fi
```

`confirm` exits 0 for the affirmative and 1 for the negative or Esc. Ctrl+c
returns a nonzero status. An `if` branch runs only on success. In `write`,
Enter submits and Ctrl+J inserts a newline.

## spin

```sh
gum spin --spinner dot --show-output --title "Listing files" -- find . -type f
gum spin --show-output --title "Measuring directory" -- du -sh .
```

The command's exit status is gum's exit status. `--show-output` streams both
stdout and stderr. With `--show-output` alone, a failing command's output is
discarded; add `--show-error` to print it. These flags apply only when stdout
is a terminal. When stdout is a pipe, a file, or `$(...)`, gum hands it to the
command directly, so its stdout is captured with or without `--show-output`,
even when it fails. `--timeout` makes gum stop waiting and exit 124, but the
command keeps running; wrap it in `timeout(1)` when it must stop.

## Output commands

```sh
gum style --foreground 212 --border double --padding "1 3" "Project tasks"
gum join --horizontal --align top "$left" "$right"
gum format < README.md
gum table --print --border rounded --file status.csv
gum pager < CHANGELOG.md
gum log --structured --level info "task completed" task "List files"
```

Read [references/commands.md](references/commands.md) for the style flags and
env defaults, `join` layouts, `format` types (code, template, emoji) and
themes, interactive `table` row selection, `pager` search, and `log` levels,
formatters, and time layouts.

## Putting it together

```sh
#!/usr/bin/env bash
set -euo pipefail

gum style --border rounded --padding "0 2" --foreground 212 "Project tasks"

tasks=$(gum choose --no-limit --selected='*' --ordered --header "Tasks" \
  "List files" "Show size") || exit 1
[[ -n $tasks ]] || { gum log --level warn "No tasks selected"; exit 0; }
dir=$(gum file . --directory --file=false --all) || exit 1

gum confirm "Run the selected tasks in $dir?" || exit 1

while IFS= read -r task; do
  case "$task" in
    "List files") gum spin --show-output --show-error --title "$task" -- ls -la "$dir" ;;
    "Show size") gum spin --show-output --show-error --title "$task" -- du -sh "$dir" ;;
    *) printf 'Unknown task: %s\n' "$task" >&2; exit 1 ;;
  esac
  gum log --structured --level info "task completed" task "$task" directory "$dir"
done <<<"$tasks"
```

The examples use Bash syntax. The complete recipe was exercised with gum
2.0.2, including cancellation, a declined confirm, and an empty selection. It
uses `set -euo pipefail` and a here string. For fish, adapt the assignments,
loops, and error handling to fish syntax. A complete fish recipe is not
included; fish and Zsh versions have not been tested.

## Mistakes that cost time

- Testing `[ -n "$choice" ]` but not the exit status. An aborted prompt can
  still print a preselected value in some commands; check `$?`.
- Redirecting stderr to a file while a prompt runs. The UI vanishes, and in
  2.0.2 `gum input` without `--width` panics and exits 124 with "timed out".
- Passing options with spaces as one argument to `choose`. Quote each option,
  or feed them on stdin one per line.
- Forgetting `--no-limit` and then wondering why only one item comes back.
- Telling users to press space in `gum choose` or Ctrl+Space in `gum filter`.
  Use x or tab in `choose` and tab in `filter`.
- Using `gum spin` for a command that needs the TTY (an editor, a login
  prompt). spin captures the command's output; run those directly.

## Older versions

- Install: package managers ship 2.0.2 (`brew install gum`,
  `nix-env -iA nixpkgs.gum`), and `go install charm.land/gum/v2@latest`
  builds the current release. The old path `github.com/charmbracelet/gum`
  stops at 0.17.0.
- `gum --version` on 0.17.0 printed the version with a leading `v`; 2.x
  prints it without one (`gum version 2.0.2`).
- Flags, defaults, and env names match 0.17.0. The move to Bubble Tea v2
  changed some keys and output handling:
  - Space and Ctrl+Space toggled in `choose` in 0.17.0 and do nothing in 2.x.
  - Ctrl+Space toggled in `filter` in 0.17.0; 2.x uses tab.
  - In 0.17.0, `format` kept its colors when piped; 2.x strips them like
    `style` and `join`.
