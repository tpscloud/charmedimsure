# TODO

- [ ] Publish `tpscloud/charmedimsure` to Tessl: add a workspace API key as the
  `TESSL_TOKEN` repository secret, then run the "Publish to Tessl" workflow.
  Publishing through the action hides the community copies Tessl already
  indexed from this repository. Once the results are public, point the
  README's Tessl badges at them instead of SECURITY.md.
- [ ] Plan how skills follow upstream releases. Decide how a new Charm library,
  gum, or vhs release is detected (for example a scheduled
  `go list -m -u` and release check), who reviews the API diff, and the order:
  bump `examples/go.mod` or the CLI baseline, fix the examples, re-verify the
  affected skill text, rescan, then tag a release.
- [ ] Run every scanner in CI. Today skilllint, skilldoctor, SkillSpector,
  Cisco Skill Scanner, and Tessl run by hand; add them to Skillforge as
  pinned, opt-in steps with recorded results, so each push or release is
  scanned the same way and SECURITY.md and the badges come from CI output.
  Tessl uploads the skill and costs credits, so run it on release tags only.
