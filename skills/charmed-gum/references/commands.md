# gum subcommand reference

Flags, keys, and examples for each gum 2.0.2 subcommand. Read the section for
the command you are using; confirm a flag with `gum <cmd> --help` on the
installed version. Shared behavior (stdout and stderr, exit codes, env
defaults, color stripping) is in [SKILL.md](../SKILL.md#how-gum-behaves).

## choose

```sh
# single choice
format=$(gum choose --header "Format" JSON YAML TOML) || exit 1

# multiple, with everything preselected and ordered output
tasks=$(printf '%s\n' "List files" "Show size" | \
  gum choose --no-limit --selected='*' --ordered --header "Tasks") || exit 1

# labels that differ from values: label:value with --label-delimiter
gum choose --label-delimiter=: "Current directory:." "Parent directory:.."

# in a loop
items=$(gum choose --no-limit a b c) || exit 1
[[ -n $items ]] || exit 0   # --no-limit allows confirming an empty selection
printf '%s\n' "$items" | while IFS= read -r item; do echo "got $item"; done
```

Flags: `--limit N` (default 1), `--no-limit`, `--selected a,b` or `'*'`,
`--ordered`, `--height 10`, `--cursor "> "`, `--cursor-prefix`,
`--selected-prefix`, `--unselected-prefix` (shown only when limit is not
1), `--header`, `--select-if-one`, `--input-delimiter`,
`--output-delimiter`, `--label-delimiter`, `--timeout`, `--show-help`.

Keys: arrows or j/k, x or tab toggles, a selects all, enter confirms. Space
and Ctrl+Space do nothing (Bubble Tea v2 reports `space`, the binding expects
`" "`). There is no filter input; use `gum filter` when searching is needed.

## filter

```sh
file=$(find . -type f -name '*.md' | gum filter --placeholder "Which Markdown file?") || exit 1
branch=$(git branch --format='%(refname:short)' | gum filter --limit 1 --header "Branch") || exit 1
gum filter --no-limit --selected='*' < list.txt
```

Flags: `--limit`, `--no-limit`, `--selected`, `--placeholder`, `--prompt`,
`--header`, `--indicator`, `--selected-prefix`, `--unselected-prefix`,
`--width`, `--height`, `--value` (initial query), `--reverse`,
`--fuzzy` / `--no-fuzzy` (case-insensitive substring match), `--fuzzy-sort`,
`--strip-ansi`, `--timeout`.

Keys: type to search, up/down or Ctrl+N/Ctrl+P move, tab toggles with
`--no-limit` (Ctrl+Space does not), Ctrl+A selects all, enter confirms. With
`--no-limit` and nothing toggled, enter returns the highlighted line; with no
match it returns nothing and exits 0.

## confirm

```sh
gum confirm --default --affirmative "Run" --negative "Skip" "Ready?"
gum confirm --timeout 10s --default "Auto continue in 10s?"   # default wins on timeout
```

`--show-output` echoes the prompt and answer to stdout for logs. Keys:
left/right, h/l, or tab move between the buttons, y and n answer at once,
enter submits the highlighted one.

## input

```sh
name=$(gum input --placeholder "project name" --header "Name") || exit 1
password=$(gum input --password --placeholder "Password") || exit 1
path=$(gum input --value "$PWD" --width 60) || exit 1
```

Flags: `--placeholder`, `--prompt "> "`, `--value`, `--password`,
`--char-limit 400` (0 unlimited), `--width`, `--header`,
`--cursor.mode blink|static|hide`, `--timeout`. The initial value may come
from stdin.

## write

```sh
msg=$(gum write --header "Commit message" --placeholder "what changed?" --char-limit 500) || exit 1
body=$(gum write --show-line-numbers --height 15 < draft.md) || exit 1
```

Flags: `--placeholder`, `--header`, `--prompt`, `--width`, `--height 5`,
`--char-limit`, `--max-lines`, `--show-line-numbers`, `--show-cursor-line`,
`--value`. Enter submits and Ctrl+J inserts a newline; Ctrl+E opens
`$EDITOR`. Esc exits 1 without output and Ctrl+c exits 130, so either one
discards the text.

## file

```sh
file=$(gum file . --file --height 15) || exit 1
dir=$(gum file . --directory --file=false --all) || exit 1
```

Flags: `--file`, `--directory` (what may be selected), `--all` (hidden
files), `--permissions`, `--size`, `--cursor`, `--height 10`, `--header`,
`--timeout`. The argument is the starting directory; the output is an
absolute path. Enter or l opens a directory in a file picker; Esc or q exits
1. The picker selects a listed entry, not the starting directory itself, so a
directory picker needs at least one subdirectory to choose from (otherwise
Enter does nothing). `--all` also lists hidden directories such as `.git`.

## spin

```sh
gum spin --spinner dot --show-output --title "Listing files" -- find . -type f
gum spin --show-error -- make build      # requires a Makefile with a build target
```

Spinners: `line`, `dot`, `minidot`, `jump`, `pulse`, `points`, `globe`,
`moon`, `monkey`, `meter`, `hamburger`. `--show-stdout` and `--show-stderr`
pick one stream; `--align left|right` places the spinner. Output handling and
`--timeout` are covered in [SKILL.md](../SKILL.md#spin).

## style

`style` wraps text in Lip Gloss formatting:

```sh
gum style --foreground 212 --border double --border-foreground 57 \
  --padding "1 3" --margin "1 0" --align center --width 50 \
  "Project tasks" "Choose a local task"

# environment defaults for a whole script
export BORDER=rounded BORDER_FOREGROUND=63 PADDING="0 1"
gum style "uses the env defaults"
```

Flags: `--foreground`, `--background`, `--border none|hidden|normal|
rounded|thick|double`, `--border-foreground`, `--border-background`,
`--align left|center|right`, `--width`, `--height`, `--padding`,
`--margin`, `--bold`, `--faint`, `--italic`, `--strikethrough`,
`--underline`, `--trim`. Each argument is one line.

## join

`join` composes rendered blocks:

```sh
left=$(gum style --border normal --padding "0 1" "left")
right=$(gum style --border normal --padding "0 1" "right")
gum join --horizontal --align top "$left" "$right"
gum join --vertical "$header" "$body"
```

## format

`format` renders Markdown, code, templates, or emoji:

```sh
gum format < README.md                              # markdown, theme pink by default
gum format --theme dracula -- "# Title" "- item"   # -- stops "- item" parsing as a flag
gum format --type code --language go < main.go
gum format --type template '{{ Bold "bold" }} and {{ Color "212" "" "pink" }}'
gum format --type emoji "ship it :rocket:"
```

## table

```sh
# interactive: pick a row from CSV, print it (or one column) to stdout.
# Without --columns the first CSV row is the header.
name=$(gum table --return-column 1 < status.csv) || exit 1
[[ -n $name ]] || exit 1   # Esc or q exits 0 with no selection

# rows.csv has no header row; --columns supplies one
row=$(gum table --columns Task,Status,Duration --widths 12,10,10 < rows.csv) || exit 1

# static print
gum table --print --border rounded --file status.csv
```

Flags: `--separator ","`, `--columns`, `--widths`, `--height`, `--print`,
`--file`, `--border`, `--return-column N`, `--lazy-quotes`,
`--fields-per-record`, `--timeout`.

## pager

```sh
gum pager < CHANGELOG.md
gum pager --show-line-numbers --soft-wrap "$(git log --oneline)"
```

`/` searches inside the pager; `--timeout` exits on its own.

## log

Structured log lines in the Charm Log style, for scripts that want the same
look as the Go tools:

```sh
gum log --structured --level info "task completed" task "List files"
gum log --structured --level error --time kitchen "task failed" err "file not found"
gum log --structured --formatter json --level warn "empty directory"
gum log --min-level warn --level debug "hidden"
gum log --file tasks.log --level info "written to a file"
```

Levels: `debug`, `info`, `warn`, `error`, `fatal` (fatal exits 1).
Formatters: `text`, `json`, `logfmt`. `--time` takes a Go layout or a name
(`kitchen`, `rfc822`, `ansic`). `--prefix`, `-f` for printf formatting
(arguments are strings, so use `%s`).
