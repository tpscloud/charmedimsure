// examples/charmed-lipgloss/layout/main.go
// Measure, join, and place blocks. Width and Height count cells, not bytes.
package main

import (
	"fmt"

	"charm.land/lipgloss/v2"
)

func main() {
	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(0, 1)

	left := box.Width(20).Render("left column\nwith two lines")
	right := box.Width(24).Render("right column")

	// Join side by side, aligning the tops.
	row := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	// Stack vertically, centering the narrower block.
	footer := lipgloss.NewStyle().Faint(true).Render("joined with JoinVertical")
	page := lipgloss.JoinVertical(lipgloss.Center, row, footer)

	// Place inside a fixed area (what you would do with WindowSizeMsg values).
	w, h := lipgloss.Size(page)
	fmt.Printf("page is %d cells wide and %d tall\n", w, h)
	framed := lipgloss.Place(w+10, h+4, lipgloss.Center, lipgloss.Center, page,
		lipgloss.WithWhitespaceChars("."),
		lipgloss.WithWhitespaceStyle(lipgloss.NewStyle().Faint(true)),
	)
	lipgloss.Println(framed)

	// Alignment inside a fixed width block.
	centered := lipgloss.NewStyle().Width(30).Align(lipgloss.Center).
		Border(lipgloss.DoubleBorder()).Render("centered text")
	lipgloss.Println(centered)

	// Wrap long text to a width while keeping any styling intact.
	long := "Lip Gloss wraps at word boundaries and preserves ANSI sequences across the break."
	lipgloss.Println(lipgloss.Wrap(long, 30, ""))
}
