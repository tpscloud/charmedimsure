# Bubble Tea v2.1.0 API reference

Generated from `go doc -all charm.land/bubbletea/v2` and trimmed to what
application code touches.  Package alias throughout: `tea`.

## Core types

| Symbol | Shape |
|--------|-------|
| `Model` | interface: `Init() Cmd`, `Update(Msg) (Model, Cmd)`, `View() View` |
| `Msg` | `= uv.Event`, effectively `any` |
| `Cmd` | `func() Msg` |
| `View` | struct, see below |
| `Program` | created by `NewProgram(model, opts...)` |

## View fields

| Field | Type | Effect |
|-------|------|--------|
| `Content` | `string` | screen content; set with `NewView(s)` or `v.SetContent(s)` |
| `AltScreen` | `bool` | alternate screen buffer |
| `MouseMode` | `MouseMode` | `MouseModeNone`, `MouseModeCellMotion`, `MouseModeAllMotion` |
| `ReportFocus` | `bool` | deliver `FocusMsg` and `BlurMsg` |
| `DisableBracketedPasteMode` | `bool` | turn off paste bracketing |
| `WindowTitle` | `string` | terminal title |
| `Cursor` | `*Cursor` | nil hides; `NewCursor(x, y)` shows; fields `Color`, `Shape`, `Blink`, `Position` |
| `BackgroundColor` | `color.Color` | terminal background, nil resets |
| `ForegroundColor` | `color.Color` | terminal foreground, nil resets |
| `ProgressBar` | `*ProgressBar` | `NewProgressBar(state, value)`; states `ProgressBarNone`, `ProgressBarDefault`, `ProgressBarError`, `ProgressBarIndeterminate`, `ProgressBarWarning` |
| `KeyboardEnhancements` | struct | `ReportEventTypes`, `ReportAlternateKeys`, `ReportAllKeysAsEscapeCodes`, `ReportAssociatedText` |
| `OnMouse` | `func(MouseMsg) Cmd` | intercept mouse events against last render |
| `ProgramStatus` | `*ProgramStatus` | OSC 7501 status report (`State`, `App`, `Title`, `Message`, `Progress`/`HasProgress`, `Kind`); states `ProgramStateIdle`, `ProgramStateWorking`, `ProgramStateBlocked`, `ProgramStateDone`, `ProgramStateError`, `ProgramStateClear` |

Cursor shapes: `CursorBlock`, `CursorUnderline`, `CursorBar`.

## Program

| Method | Notes |
|--------|-------|
| `Run() (Model, error)` | blocks; returns final model |
| `Send(Msg)` | thread safe; blocks until program started |
| `Quit()` | from outside the program |
| `Kill()` | immediate stop, skips final render |
| `Wait()` | block until shutdown finishes |
| `Println(...)`, `Printf(...)` | print above the UI (not in alt screen) |
| `ReleaseTerminal()`, `RestoreTerminal()` | manual terminal handoff |

Errors from `Run`: `ErrInterrupted`, `ErrProgramKilled`, `ErrProgramPanic`.

## Program options

| Option | Purpose |
|--------|---------|
| `WithContext(ctx)` | cancel the program from a context |
| `WithInput(io.Reader)` | custom input; tests.  `WithInput(nil)` disables input and also skips every terminal query (background color, cursor position, clipboard, version, capability probe), so their reply messages never arrive |
| `WithOutput(io.Writer)` | custom output; tests |
| `WithWindowSize(w, h)` | fixed initial size; tests |
| `WithColorProfile(colorprofile.Profile)` | force a profile; tests |
| `WithEnvironment([]string)` | env for remote sessions (SSH) |
| `WithFPS(n)` | max frame rate, default 60, cap 120 |
| `WithFilter(func(Model, Msg) Msg)` | inspect or drop messages before Update |
| `WithoutRenderer()` | run the loop with no drawing |
| `WithoutSignalHandler()` | handle SIGINT and friends yourself |
| `WithoutSignals()` | ignore OS signals entirely; tests |
| `WithoutCatchPanics()` | let panics propagate with a stack trace |

## Messages delivered to Update

| Message | When | Useful members |
|---------|------|----------------|
| `KeyPressMsg` | key down | `String()`, `Keystroke()`, `Code`, `Text`, `Mod`, `ShiftedCode`, `BaseCode`, `IsRepeat` |
| `KeyReleaseMsg` | key up, only with `ReportEventTypes` | same as above |
| `KeyMsg` | interface over both | `Key() Key` |
| `WindowSizeMsg` | start and every resize | `Width`, `Height` |
| `MouseClickMsg`, `MouseReleaseMsg`, `MouseWheelMsg`, `MouseMotionMsg` | with a mouse mode set | `Mouse()` gives `Mouse{X, Y, Button, Mod}` |
| `MouseMsg` | interface over all four | `Mouse()` |
| `PasteMsg` | bracketed paste | `Content` |
| `PasteStartMsg`, `PasteEndMsg` | paste boundaries | |
| `FocusMsg`, `BlurMsg` | with `ReportFocus` | |
| `BackgroundColorMsg` | reply to `RequestBackgroundColor` | `IsDark()`, `String()` |
| `ForegroundColorMsg`, `CursorColorMsg` | replies to the matching requests | `IsDark()` |
| `ColorProfileMsg` | at start | `Profile` |
| `CursorPositionMsg` | reply to `RequestCursorPosition` | `X`, `Y` |
| `KeyboardEnhancementsMsg` | terminal reports support | `SupportsEventTypes()` etc |
| `TerminalVersionMsg` | reply to `RequestTerminalVersion` | `Name` |
| `CapabilityMsg` | reply to `RequestCapability` | `Content` |
| `ClipboardMsg` | reply to `ReadClipboard` | `Content` |
| `EnvMsg` | at start | `Getenv`, `LookupEnv` |
| `ResumeMsg` | after `Suspend` returns | |
| `QuitMsg`, `SuspendMsg`, `InterruptMsg` | produced by the matching commands | |
| `ModeReportMsg` | DECRPM reply | `Mode`, `Value` |
| `ProgramStatusSupportMsg` | reply to `RequestProgramStatusSupport`; unsupported terminals never reply | |

## Commands

Plain `Cmd` values (return them directly):

`Quit`, `Suspend`, `Interrupt`, `ClearScreen`, `RequestWindowSize`,
`RequestBackgroundColor`, `RequestForegroundColor`, `RequestCursorColor`,
`RequestCursorPosition`, `RequestTerminalVersion`, `ReadClipboard`,
`ReadPrimaryClipboard`, `RequestProgramStatusSupport`.

Constructors that return a `Cmd`:

| Function | Purpose |
|----------|---------|
| `Batch(cmds...)` | run concurrently, nils dropped |
| `Sequence(cmds...)` | run in order |
| `Tick(d, fn)` | one message after d, measured from now |
| `Every(d, fn)` | one message aligned to the system clock |
| `Println(args...)`, `Printf(tmpl, args...)` | print above the UI |
| `SetClipboard(s)`, `SetPrimaryClipboard(s)` | OSC 52 clipboard write |
| `ExecProcess(*exec.Cmd, func(error) Msg)` | hand over the terminal, resume after |
| `Exec(ExecCommand, func(error) Msg)` | same, for a custom runner |
| `RequestCapability(name)` | terminfo query |
| `Raw(any)` | write raw bytes to the terminal |

## Key constants

Navigation: `KeyUp`, `KeyDown`, `KeyLeft`, `KeyRight`, `KeyHome`, `KeyEnd`,
`KeyPgUp`, `KeyPgDown`, `KeyInsert`, `KeyDelete`, `KeyBegin`, `KeyFind`,
`KeySelect`.

Editing and control: `KeyEnter`, `KeyTab`, `KeyBackspace`, `KeyEscape`,
`KeySpace`.

Function keys `KeyF1` through `KeyF63`.  Keypad keys `KeyKp0` through
`KeyKp9`, `KeyKpEnter`, `KeyKpPlus`, and so on.  Lock keys `KeyCapsLock`,
`KeyNumLock`, `KeyScrollLock`, plus `KeyPrintScreen`, `KeyPause`, media keys.

Modifiers: `ModShift`, `ModAlt`, `ModCtrl`, `ModMeta`, `ModSuper`,
`ModHyper`, `ModCapsLock`, `ModNumLock`, `ModScrollLock`.  Test with
`mod.Contains(tea.ModAlt)`.

Mouse buttons: `MouseNone`, `MouseLeft`, `MouseMiddle`, `MouseRight`,
`MouseWheelUp`, `MouseWheelDown`, `MouseWheelLeft`, `MouseWheelRight`,
`MouseBackward`, `MouseForward`, `MouseButton10`, `MouseButton11`.

## Logging helpers

`LogToFile(path, prefix) (*os.File, error)` points the standard logger at a
file.  `LogToFileWith(path, prefix, setter)` does the same for any logger
that has `SetOutput` and `SetPrefix`, such as `charm.land/log/v2`.

## Other

`OpenTTY() (in, out *os.File, err error)` opens the controlling terminal, for
programs whose stdin or stdout is redirected.

## Migrating from v1

Search and replace, in this order:

| v1                                       | v2                                             |
|------------------------------------------|------------------------------------------------|
| `github.com/charmbracelet/bubbletea`     | `charm.land/bubbletea/v2`                      |
| `View() string`                          | `View() tea.View`, wrap with `tea.NewView(s)`  |
| `tea.KeyMsg` (struct)                    | `tea.KeyPressMsg`                              |
| `msg.Type`, `msg.Runes`, `msg.Alt`       | `msg.Code`, `msg.Text`, `msg.Mod`              |
| `case " ":`                              | `case "space":`                                |
| `tea.KeyCtrlC`                           | `msg.String() == "ctrl+c"`                     |
| `tea.MouseMsg` fields                    | `msg.Mouse().X` etc; match concrete types      |
| `tea.MouseButtonLeft`                    | `tea.MouseLeft`                                |
| `tea.WithAltScreen()`, `tea.EnterAltScreen` | `v.AltScreen = true` in View               |
| `tea.WithMouseCellMotion()`              | `v.MouseMode = tea.MouseModeCellMotion`        |
| `tea.SetWindowTitle(s)`                  | `v.WindowTitle = s`                            |
| `tea.HideCursor` / `tea.ShowCursor`      | `v.Cursor = nil` / `v.Cursor = tea.NewCursor(x, y)` |
| `tea.WindowSize()`                       | `tea.RequestWindowSize`                        |
| `tea.Sequentially`                       | `tea.Sequence`                                 |
| `p.Start()`                              | `p.Run()`                                      |

`tea.WithInputTTY()` and `tea.WithANSICompressor()` are gone and need no
replacement. The [v2.1.0 upgrade guide](https://github.com/charmbracelet/bubbletea/blob/v2.1.0/UPGRADE_GUIDE_V2.md)
has the complete list.

## Older v2 releases

These are new in v2.1.0, so a project pinned to an earlier v2 release does
not have them: the `ProgramStatus` view field, the
`RequestProgramStatusSupport` command, and `WithInput(nil)` skipping every
terminal query.
