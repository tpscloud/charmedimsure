// examples/charmed-lipgloss/basics/main.go
// Styles are values. Build them with chained setters, then Render.
// Outside Bubble Tea, print through lipgloss so colors are downsampled to
// what the terminal supports.
package main

import (
	"charm.land/lipgloss/v2"
)

func main() {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	card := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Width(40)

	body := lipgloss.NewStyle().Italic(true).Faint(true)

	// Styles are values, so deriving one never mutates the original.
	warning := body.Faint(false).Italic(false).Foreground(lipgloss.Yellow)

	content := title.Render("Lip Gloss v2") + "\n\n" +
		body.Render("Render returns a string with ANSI codes. Print it with lipgloss.Println so it is downsampled for this terminal.") + "\n\n" +
		warning.Render("fmt.Println would emit true color codes regardless of terminal support.")

	lipgloss.Println(card.Render(content))
}
