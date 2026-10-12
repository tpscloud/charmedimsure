# Sign-up form that also works for screen readers

Write a Go program with Huh v2 (`charm.land/huh/v2`) that asks for a project
name (required) and a deployment region (one of three). It must honour an
`ACCESSIBLE` environment variable for screen-reader users. After the form, print
the chosen values, or print `cancelled` and exit 1 if the user cancelled. An
empty project name must never be accepted, in either mode.

Return the complete `main.go`.
