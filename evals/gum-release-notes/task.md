# Release-notes prompt with gum

Write a Bash script that uses gum to ask for multi-line release notes, then asks
"Publish these notes?" and, only if the user agrees, writes them to `NOTES.md`.
If the user cancels at any prompt, the script must exit without writing
anything. Print one line telling the user which keys submit the notes and which
key adds a new line.

Return the script.
