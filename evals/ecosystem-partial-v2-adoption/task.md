# Adding Charm v2 libraries to a v1 application

An existing Go application uses `github.com/charmbracelet/bubbletea` and
`github.com/charmbracelet/lipgloss` (both pre-v2) throughout its TUI. The team
wants two additions without rewriting the TUI now: a setup wizard built with Huh
(a separate `init` subcommand that runs before the TUI starts) and Charm Log for
colored logging. Explain which v2 modules they can add alongside the v1 code,
what would force a wider migration and why, give the exact `go get` command for
the v2 modules you recommend, and list what to watch for.

Return the recommendation and the command.
