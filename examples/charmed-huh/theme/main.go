// examples/charmed-huh/theme/main.go
// Themes are functions of isDark. Start from a built in one and override.
package main

import (
	"fmt"
	"os"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

func myTheme(isDark bool) *huh.Styles {
	s := huh.ThemeCharm(isDark)
	ld := lipgloss.LightDark(isDark)
	accent := ld(lipgloss.Color("#5A2D82"), lipgloss.Color("#C4A7F5"))

	s.Focused.Title = s.Focused.Title.Foreground(accent).Bold(true)
	s.Focused.SelectedOption = s.Focused.SelectedOption.Foreground(accent)
	s.Focused.MultiSelectSelector = s.Focused.MultiSelectSelector.Foreground(accent)
	s.Focused.FocusedButton = s.Focused.FocusedButton.Background(accent)
	return s
}

func main() {
	// huh asks the terminal for its background color itself and calls the
	// ThemeFunc with the answer. The built in themes have the same shape:
	// huh.ThemeDracula, huh.ThemeCatppuccin, huh.ThemeBase16, huh.ThemeBase.
	var ok bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().Title("Themed confirm").Value(&ok),
		),
	).WithTheme(huh.ThemeFunc(myTheme))

	// NewForm already turns accessible mode on for TERM=dumb. Only override
	// that when ACCESSIBLE is actually set; calling WithAccessible(false)
	// unconditionally would switch it back off.
	if v, set := os.LookupEnv("ACCESSIBLE"); set {
		form = form.WithAccessible(v != "" && v != "0" && v != "false")
	}

	err := form.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("ok:", ok)
}
