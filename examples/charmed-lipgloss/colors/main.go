// examples/charmed-lipgloss/colors/main.go
// Color helpers, gradients, underline styles, hyperlinks, and ranges.
package main

import (
	"strings"

	"charm.land/lipgloss/v2"
)

func main() {
	base := lipgloss.Color("#7D56F4")

	// Derive tints and shades from one brand color.
	lipgloss.Println(
		lipgloss.NewStyle().Foreground(lipgloss.Lighten(base, 0.3)).Render("lighter"),
		lipgloss.NewStyle().Foreground(base).Render("base"),
		lipgloss.NewStyle().Foreground(lipgloss.Darken(base, 0.3)).Render("darker"),
		lipgloss.NewStyle().Foreground(lipgloss.Complementary(base)).Render("complement"),
	)

	// Gradient text: one color per rune from a blend. Index runes, not bytes;
	// ranging over the string directly would sample at UTF-8 byte offsets and
	// skip stops on multibyte input such as "grädient".
	word := "gradient"
	runes := []rune(word)
	stops := lipgloss.Blend1D(len(runes), lipgloss.Color("#FF5F87"), lipgloss.Color("#5FD7FF"))
	var b strings.Builder
	for i, r := range runes {
		b.WriteString(lipgloss.NewStyle().Foreground(stops[i]).Render(string(r)))
	}
	lipgloss.Println(b.String())

	// Gradient borders come built in.
	lipgloss.Println(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForegroundBlend(lipgloss.Color("#FF5F87"), lipgloss.Color("#5FD7FF")).
		Padding(0, 2).
		Render("blended border"))

	// Underline variants and a clickable link (OSC 8, where supported).
	lipgloss.Println(
		lipgloss.NewStyle().UnderlineStyle(lipgloss.UnderlineCurly).UnderlineColor(lipgloss.Red).Render("curly"),
		lipgloss.NewStyle().UnderlineStyle(lipgloss.UnderlineDouble).Render("double"),
		lipgloss.NewStyle().Hyperlink("https://charm.land").Render("charm.land"),
	)

	// Style a terminal-cell range. In this ASCII text, byte offsets match cells.
	s := "match this word here"
	start := strings.Index(s, "word")
	lipgloss.Println(lipgloss.StyleRanges(s,
		lipgloss.NewRange(start, start+len("word"), lipgloss.NewStyle().Reverse(true)),
	))

	// The 16 ANSI colors have names; anything else goes through Color().
	lipgloss.Println(lipgloss.NewStyle().Foreground(lipgloss.BrightGreen).Render("BrightGreen"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render("ANSI 214"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFAA")).Render("hex"))
}
