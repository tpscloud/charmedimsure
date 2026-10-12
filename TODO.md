# TODO

- [ ] Publish the release as `tpscloud/charmedimsure-skill` on `main`. The
  skills link to examples at that path.
- [ ] After publishing, run `npx skills add tpscloud/charmedimsure-skill --list`
  from a clean directory and check that all nine skills appear.
- [ ] Decide whether to publish `tpscloud/charmedimsure` to the Tessl registry.
  The manifest is private; making a plugin public cannot be undone. Once the
  Tessl results are public, point the README's Tessl badges at them instead of
  SECURITY.md.
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
