# Demo GIF in CI with vhs

Write a vhs tape `demo.tape` that records `echo hello` followed by `ls` in Bash at
a fixed size and theme, and a GitHub Actions job that renders it on every push to
`main`. Renders must be reproducible, and nothing may be uploaded anywhere
except as a workflow artifact.

Return the tape and the workflow YAML.
