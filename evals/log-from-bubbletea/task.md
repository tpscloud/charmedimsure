# Debug logging from a Bubble Tea program

A Bubble Tea v2 program needs debug logging with Charm Log v2
(`charm.land/log/v2`). Logs must go to `debug.log` while the UI runs, never to
the terminal, at debug level with timestamps. Someone will read `debug.log` with
`grep` and `less`, so it must be plain text with no ANSI escape codes. Log every
key press.

Return the complete `main.go`.
