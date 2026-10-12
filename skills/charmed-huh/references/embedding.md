# Embedding a form in Bubble Tea

Details for running a `*huh.Form` as a component inside a Bubble Tea v2
program instead of calling `Run`. SKILL.md has the core `Update` and `View`
code; the complete
[embedded example](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-huh/embedded/main.go)
has the model, frame styles, and minimum width.

## Reading the result

Print the result after the program's `Run` returns, from the final model:

```go
// examples/charmed-huh/embedded/main.go
func main() {
    final, err := tea.NewProgram(newModel()).Run()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    // Read the answers from the final model once the terminal is restored.
    if m, ok := final.(model); ok && m.form.State == huh.StateCompleted {
        fmt.Printf("You picked %s at priority %d\n", m.form.GetString("project"), m.form.GetInt("priority"))
    }
}
```

`tea.Printf` prints nothing while the alternate screen is active (its doc
comment says so), so a completion message sent that way from an alt screen
program is lost.

Give fields a `Key` and read answers with `GetString`, `GetInt`, `GetBool`,
or `Get`; or bind `Value` pointers as usual.

## Hooks and navigation

`SubmitCmd` and `CancelCmd` on the form run when it completes or aborts,
which is a cleaner hook than polling `State` when the form is one of several
screens.

Navigation is available as the `*Form` methods `NextField`, `PrevField`,
`NextGroup`, and `PrevGroup`, which return a `tea.Cmd`.  `huh.NextField` and
`huh.PrevField` also exist as package-level `tea.Msg` functions.

## The outer model

The outer model's `View` returns `tea.View`.  The upstream README still
shows `View() string` for the outer model; that is v1.

## Checking an embedded form

Resize the terminal and confirm the form's width follows the frame, then
check the final output after the program exits.
