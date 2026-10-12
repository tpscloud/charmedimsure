// examples/charmed-glamour/render/main.go
// Render markdown with a built in style, then print through Lip Gloss so
// colors match the terminal. Glamour itself is pure and knows nothing
// about the terminal.
package main

import (
	"fmt"
	"os"

	"charm.land/glamour/v2"
	"charm.land/lipgloss/v2"
)

const doc = `# Local project

Status and notes for a local development project.

## Changes

- Fixed *task list* sorting
- Updated the **setup guide**

| Component | Status | Next step |
|-----------|--------|-----------|
| CLI       | ready  | review help |
| Docs      | draft  | check links |
| Tests     | pass   | review coverage |

` + "```sh\nls -la .\n```\n"

func main() {
	// One shot with a named style: "dark" (default), "light", "pink",
	// "dracula", "tokyo-night", "ascii", "notty".
	out, err := glamour.Render(doc, "dark")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	lipgloss.Print(out)

	// A reusable renderer with options. Pick light or dark yourself.
	style := "dark"
	if !lipgloss.HasDarkBackground(os.Stdin, os.Stdout) {
		style = "light"
	}
	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle(style),
		glamour.WithWordWrap(60),
		glamour.WithEmoji(),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	out, err = r.Render("Rendered again at 60 columns :tada:")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	lipgloss.Print(out)
}
