# Searchable file list with a background scan

Using Bubbles v2 (`charm.land/bubbles/v2`) inside a Bubble Tea v2 program, show
a list of twelve Markdown file names. The user can filter the list by typing `/`
and a search term, and confirm the highlighted entry with enter, which prints
the file name and quits. `ctrl+c` quits without printing. While a simulated scan
runs for the first three seconds, a spinner animates next to the list title;
afterwards the title says `scan complete`.

Return the complete `main.go`.
