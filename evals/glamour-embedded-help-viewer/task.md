# Markdown help screen inside an existing TUI

An existing Bubble Tea v2 application needs a help screen component that renders
`HELP.md` with Glamour v2 (`charm.land/glamour/v2`) inside a scrollable Bubbles
viewport. The component is a bubble inside an app that already owns the
terminal, so it cannot touch stdin or stdout itself. The rendering must use a
style that suits the terminal's light or dark background and must wrap to the
window width, re-flowing when the window is resized. `q` closes the app.

Return `main.go` with the component and a minimal program that runs it.
