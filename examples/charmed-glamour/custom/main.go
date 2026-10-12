// examples/charmed-glamour/custom/main.go
// Start from a built in StyleConfig and change a few things.
package main

import (
	"fmt"
	"os"

	"charm.land/glamour/v2"
	"charm.land/glamour/v2/ansi"
	"charm.land/glamour/v2/styles"
	"charm.land/lipgloss/v2"
)

func main() {
	cfg := styles.DraculaStyleConfig // a value copy; safe to edit

	cfg.Document.Margin = new(uint(0))
	cfg.H1.Prefix = "" // drop the "# " style prefix on top level headings
	cfg.H1.BackgroundColor = new("#5A2D82")
	cfg.H1.Color = new("#FFFFFF")
	cfg.Code.StylePrimitive.Color = new("#C4A7F5")
	cfg.Link.Underline = new(true)

	r, err := glamour.NewTermRenderer(
		glamour.WithStyles(cfg),
		glamour.WithWordWrap(70),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	out, err := r.Render("# Custom style\n\nInline `code` and a [link](https://charm.land).\n")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	lipgloss.Print(out)

	// Styles can also come from JSON, the same format as the styles/ gallery.
	_ = ansi.StyleConfig{} // WithStylesFromJSONFile("my.json") or WithStylesFromJSONBytes(b)
}
