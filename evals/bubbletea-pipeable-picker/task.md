# Environment picker for a deploy script

Write `main.go` for a Go program using Bubble Tea v2 (`charm.land/bubbletea/v2`).
It shows the environments `dev`, `staging`, and `production` in the terminal's
alternate screen. Up and down (or `k` and `j`) move the highlight, space marks
the highlighted environment, and enter confirms the marked one. On confirm the
program prints the marked environment name on stdout and exits 0. If the user
presses `q`, `esc`, or `ctrl+c`, the program prints nothing and exits 1.

The program is used from a shell script as `env=$(pick-env)`, so stdout must
contain only the result.

Return the complete `main.go`.
