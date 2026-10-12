// examples/charmed-lipgloss/adaptive/main.go
// Standalone programs must ask the terminal whether the background is dark
// and which color profile it supports. Inside Bubble Tea, use
// tea.RequestBackgroundColor instead (see charmed-bubbletea).
package main

import (
	"os"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
)

func main() {
	// Query once at startup. This writes and reads an escape sequence on the
	// terminal, so do it before any other output.
	isDark := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	lightDark := lipgloss.LightDark(isDark)

	// lightDark(light, dark) picks one color.Color for the current background.
	fg := lightDark(lipgloss.Color("#1A1A1A"), lipgloss.Color("#F0F0F0"))
	accent := lightDark(lipgloss.Color("#5A2D82"), lipgloss.Color("#C4A7F5"))

	// Complete picks per profile when you have hand tuned values for each.
	profile := colorprofile.Detect(os.Stdout, os.Environ())
	complete := lipgloss.Complete(profile)
	brand := complete(lipgloss.Color("5"), lipgloss.Color("170"), lipgloss.Color("#D75FEE"))

	text := lipgloss.NewStyle().Foreground(fg)
	head := lipgloss.NewStyle().Bold(true).Foreground(accent)
	logo := lipgloss.NewStyle().Foreground(brand).Bold(true)

	lipgloss.Println(head.Render("adaptive"), text.Render("dark background:"), isDark)
	lipgloss.Println(logo.Render("profile:"), profile)
}
