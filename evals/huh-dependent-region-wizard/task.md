# Deployment wizard with dependent choices

Write a Go program with Huh v2 (`charm.land/huh/v2`). One form asks for a cloud
provider (`aws`, `gcp`, `azure`) and then a region; the region choices depend
on the provider (give each provider three regions). The user can go back to the
provider question and change it, and the region choices must follow. A final
confirm asks `Deploy?`. After the form, show a spinner for two seconds of
simulated work, then print `provider region` on stdout.

The program honours the `ACCESSIBLE` environment variable for screen readers.
It is called from scripts as `target=$(deploy-wizard)`, so stdout must carry
only the final line in both modes. Cancelling exits with a non-zero status and
prints nothing on stdout.

Return the complete `main.go`.
