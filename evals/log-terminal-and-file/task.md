# CLI logging to the terminal and a file

A Go command line tool logs with Charm Log v2 (`charm.land/log/v2`). Every log
line must appear on stderr and also be appended to `app.log`. On a terminal the
stderr output is colored; the file is always plain text that `grep` can search.
A `--verbose` flag switches both destinations from info to debug level, and
every line carries a timestamp. Log at least one message with key value pairs.

Return the complete `main.go`.
