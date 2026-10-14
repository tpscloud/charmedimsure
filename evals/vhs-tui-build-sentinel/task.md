# Record a TUI demo from source

A Go TUI lives in the current directory (`go build -o ./app .` builds it). It
runs in the alternate screen, shows the text `Ready` once it has started,
toggles a panel on space, and quits on `q`. Write `demo.tape` for vhs 0.12.1
that builds the program, records it starting, pressing space, and quitting, and
saves `demo.gif` plus a PNG of the open panel. The build output must not appear
in the recording. If the build fails, vhs must fail instead of recording an old
binary. Add a GitHub Actions workflow that renders the tape on every pull
request and keeps the GIF as a workflow artifact without uploading it anywhere
else.

Return the tape and the workflow YAML.
