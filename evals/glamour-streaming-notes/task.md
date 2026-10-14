# Live Markdown log in a TUI

Write a Bubble Tea v2 program that uses Glamour v2 (`charm.land/glamour/v2`)
to show a Markdown document in a scrollable Bubbles viewport while new Markdown
sections are appended once a second (simulate ten appends, each with a heading,
a paragraph, and a short code block). The view follows the newest content
unless the user has scrolled up, in which case it stays where they are until
they scroll back to the bottom. Rendering must suit the terminal's light or
dark background, wrap to the window width, and keep working when the window is
very narrow (ten columns). `q` quits.

Return the complete `main.go`.
