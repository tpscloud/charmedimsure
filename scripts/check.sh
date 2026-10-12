#!/bin/sh
# Check Charm examples and documentation. Skillforge validates frontmatter.
set -eu
root=$(cd "$(dirname "$0")/.." && pwd)
fail=0

if ! command -v go >/dev/null 2>&1; then
  echo "error: Go is not on PATH; install the version required by examples/go.mod" >&2
  exit 1
fi

echo "== vetting examples"
if ! ( cd "$root/examples" && go vet ./... ); then
  echo "error: example vet checks failed (see Go diagnostics above)" >&2
  fail=1
fi

echo "== building examples"
if ! ( cd "$root/examples" && go build ./... ); then
  echo "error: example build failed (see Go diagnostics above)" >&2
  fail=1
fi

echo "== running example tests"
if ! ( cd "$root/examples" && go test ./... ); then
  echo "error: example tests failed (see Go diagnostics above)" >&2
  fail=1
fi

echo "== sweeping markdown for em and en dashes"
if grep -rn --include='*.md' -e "$(printf '\342\200\224')" -e "$(printf '\342\200\223')" "$root" | grep -v '/examples/go.sum'; then
  echo "error: prohibited em or en dash in the files listed above" >&2; fail=1
fi

echo "== sweeping for v1 import paths"
if grep -rn --include='*.md' --include='*.go' -E '"github\.com/charmbracelet/(bubbletea|lipgloss|bubbles|huh|glamour|log)(/[a-z]+)?"' "$root" ; then
  echo "error: v1 Go import string in the files listed above; use the documented v2 path" >&2; fail=1
fi

echo "== checking Go fences name their example file or are marked as excerpts"
# The first line of every Go fence in a skill or reference is either
#   // examples/<skill>/<topic>/main.go        (trimmed from that file)
#   // excerpt ...                             (API sketch, no example file)
# and a named example file must exist.
for f in "$root"/skills/*/SKILL.md "$root"/skills/*/references/*.md; do
  [ -f "$f" ] || continue
  awk -v root="$root" '
    /^```go/ {
      if ((getline line) <= 0) { printf "%s:%d: unterminated Go fence\n", FILENAME, NR; bad = 1; exit }
      if (line ~ /^\/\/ examples\/[^ ]+\.go( |$)/) {
        p = line; sub(/^\/\/ /, "", p); sub(/ .*/, "", p)
        # Open the file instead of shelling out, so a crafted path cannot run a command.
        f = root "/" p
        ok = ((getline tmp < f) >= 0); close(f)
        if (!ok) {
          printf "%s:%d: no such example file %s\n", FILENAME, NR, p; bad = 1
        }
      } else if (line !~ /^\/\/ excerpt/) {
        printf "%s:%d: Go fence must start with // examples/<path> or // excerpt\n", FILENAME, NR; bad = 1
      }
    }
    END { exit bad }' "$f" || fail=1
done

if [ "$fail" -ne 0 ]; then
  echo "Repository checks failed. Fix the reported errors and rerun skillforge check ." >&2
  exit 1
fi
echo "All repository checks passed."
