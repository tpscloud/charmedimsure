# Testing and debugging Bubble Tea v2

These examples use Bubble Tea v2.1.0. Return to the
[main skill](../SKILL.md) for the program loop and command patterns.

## Testing

`View` and `Update` are plain methods, so unit tests construct messages and
call them.  For end to end runs, give the program an `io.Reader` of
keystrokes, a buffer for output, and a fixed size
([examples/charmed-bubbletea/testing/main_test.go](https://github.com/tpscloud/charmedimsure/blob/main/examples/charmed-bubbletea/testing/main_test.go)):

```go
// examples/charmed-bubbletea/testing/main_test.go
func TestCursorMovesAndSelects(t *testing.T) {
    var m tea.Model = model{items: []string{"a", "b", "c"}}
    m, _ = m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
    m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
    m, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
    if got := m.(model).chosen; got != "c" {
        t.Fatalf("chosen = %q, want c", got)
    }
    if _, ok := cmd().(tea.QuitMsg); !ok {
        t.Fatal("expected the returned command to produce a QuitMsg")
    }
}

func TestProgramHeadless(t *testing.T) {
    in := strings.NewReader("j\r")
    var out bytes.Buffer
    p := tea.NewProgram(model{items: []string{"a", "b", "c"}},
        tea.WithInput(in),
        tea.WithOutput(&out),
        tea.WithWindowSize(40, 10),
        tea.WithoutSignals(),
    )
    final, err := p.Run()
    // ...
}
```

`tea.WithColorProfile(colorprofile.Ascii)` makes rendered output free of
color codes when you want to compare strings.

## Debugging

Standard output belongs to the renderer while the program runs.  Log to a
file instead:

```go
// excerpt: tea.LogToFile with the standard library logger; charmed-log's tea example uses LogToFileWith
f, err := tea.LogToFile("debug.log", "debug")
if err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
}
defer f.Close()
log.Println("this goes to debug.log")
```

Then `tail -f debug.log` in another terminal.  Bubble Tea recovers panics by
default and restores the terminal; pass `tea.WithoutCatchPanics()` to get a
raw stack trace while developing.

