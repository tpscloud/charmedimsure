# vhs in CI and tests

How to validate and render tapes in CI with vhs 0.12.1, and how to use tapes
as golden-file or screenshot fixtures. Read this before adding vhs to a
pipeline or a test suite.

## Validating

Pass existing readable files to `vhs validate`: v0.12.1 silently skips
unreadable or missing paths, and does not expand a quoted glob itself. Have
CI fail if its file list is empty or a tape is unreadable before calling the
validator.

```sh
set -euo pipefail
shopt -s nullglob
tapes=(./*.tape)
(( ${#tapes[@]} > 0 )) || { echo "no tapes found" >&2; exit 1; }
for t in "${tapes[@]}"; do
  [[ -r $t ]] || { echo "unreadable: $t" >&2; exit 1; }
done
vhs validate "${tapes[@]}"
```

## Rendering

Render in a job that has ffmpeg and ttyd. A failed ffmpeg encode makes vhs
exit non zero. The [official action](https://github.com/charmbracelet/vhs-action):

```yaml
- uses: charmbracelet/vhs-action@f6d7db07a432fcd3b06772628d36a57f96d95dcf # v2.1.1
  with:
    version: v0.12.1
    path: demo.tape
```

The action's `version` defaults to `latest`, and it always installs the
latest ttyd and ffmpeg. Pinning those needs your own install step or a
container image.

## Golden files and screenshots

For text assertions, add `Output golden.txt` and diff it against a committed
copy. For deterministic screenshots in tests, `Set Theme` to a fixed theme,
`Set FontFamily` to a font you ship, `Hide` the setup, and compare the
`Screenshot` output. Theme and font are not enough on their own: pin the vhs,
ttyd, and ffmpeg versions in CI, set `Width` and `Height` explicitly, feed the
program fixed input data, and `Wait` on drawn text rather than `Sleep` so
timing differences do not change the frame. Treat all of those as inputs to
the fixture.
