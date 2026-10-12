# Accessible mode

How Huh v2.0.3 decides on accessible mode, and where it behaves differently
from the TUI.

`form.WithAccessible(true)` swaps the TUI for plain line based prompts that
screen readers handle.

## Turning it on without losing the default

`NewForm` switches accessible mode on when `TERM` is `dumb`, and any later
`WithAccessible` call replaces that decision.  So
`WithAccessible(os.Getenv("ACCESSIBLE") != "")` turns it back off on a dumb
terminal whenever the variable is unset.  Honor the common `ACCESSIBLE`
convention without losing the default by only calling `WithAccessible` when
the variable is present, as the
[theme example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-huh/theme/main.go)
does with `os.LookupEnv`.

## Differences from the TUI

- Accessible prompts default to stdout, while the TUI renders to stderr.
  Use `WithOutput(os.Stderr)` to keep stdout for results in both modes.
- Timeouts are not supported (`ErrTimeoutUnsupported`).
- The accessible path does not use the context passed to `RunWithContext`
  and does not skip hidden groups.  Check accessible behavior separately
  when using conditional pages.
- `Run` discards field errors and returns nil: on EOF a prompt keeps its
  default value and skips validation, and `ErrUserAborted` is never
  returned.  Validate the bound values after `Run`.

The spinner has its own `WithAccessible`, mirroring the form.
