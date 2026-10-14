# Issue triage script with gum

Write a Bash script that uses gum. It prints a styled banner, then shows
`issues.csv` (columns `id`, `title`, `age`) as an interactive table where the
user picks one row. Next the user picks any number of labels from `bug`,
`docs`, `question`, and `wontfix`, including none. Then the script asks
`Apply labels?` and waits at most ten seconds; silence counts as no. If the
user agrees, it runs `./apply-labels <id> <labels...>` (assume it exists and
can take a while) behind a spinner and shows the command's output, including
when the command fails. Cancelling at any step (Esc or Ctrl+C) must stop the
script before anything is applied. The script's output is often piped into
`tee triage.log`, and the banner must keep its colors there. Print one line
telling the user how to toggle labels.

Return the script.
