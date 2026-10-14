# Name prompt with the terminal's real cursor

Write a full screen Bubble Tea v2 program (`charm.land/bubbletea/v2`) with a
Bubbles text input inside a bordered, padded panel that is centered in the
window. The terminal's own cursor (the one the terminal draws, with the shape
and blink the user configured) must sit at the typing position inside the
input; the program must not draw a simulated cursor. Pressing tab moves focus
between the input and an `OK` button; while the button has focus no cursor is
shown. Enter on the button, or `esc`, quits; the entered name is printed after
the program ends.

Return the complete `main.go`.
