// examples/charmed-lipgloss/layout/box_test.go
// Pins the box model described in charmed-lipgloss: Width includes padding
// and border and excludes margins.
package main

import (
	"testing"

	"charm.land/lipgloss/v2"
)

func TestWidthIncludesPaddingAndBorderNotMargin(t *testing.T) {
	s := lipgloss.NewStyle().Width(10).Padding(0, 1).Border(lipgloss.NormalBorder())

	if got := lipgloss.Width(s.Render("hello")); got != 10 {
		t.Fatalf("Width(10) with padding and border rendered %d cells, want 10", got)
	}
	if got := lipgloss.Width(s.Margin(0, 1).Render("hello")); got != 12 {
		t.Fatalf("adding Margin(0, 1) rendered %d cells, want 12", got)
	}
	if h, _ := s.GetFrameSize(); h != 4 {
		t.Fatalf("GetFrameSize horizontal = %d, want 4 (padding 2 + border 2)", h)
	}
	if got := lipgloss.Height(lipgloss.NewStyle().Height(3).Border(lipgloss.NormalBorder()).Render("x")); got != 3 {
		t.Fatalf("Height(3) with border rendered %d lines, want 3", got)
	}
}
