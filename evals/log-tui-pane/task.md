# Colored log pane inside a TUI

A Bubble Tea v2 application runs three background workers that log through
Charm Log v2 (`charm.land/log/v2`) from their own goroutines, a line every few
hundred milliseconds with the worker name as a key value pair. The application
shows the last twenty log lines in a pane at the bottom of the screen, with the
same level colors the logger uses on a terminal, and also appends every line
as plain text to `app.log`. The terminal itself must never receive log output
directly. `q` quits.

Return the complete `main.go`.
