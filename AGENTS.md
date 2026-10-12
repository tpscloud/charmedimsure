# charmedimsure

Agent skills for the Charm terminal stack, written from the module source
rather than from memory.  Every Go snippet trimmed from a program has a
matching file under `examples/` that must compile.  If the two ever disagree,
the example is right and the skill text is wrong.

## What is in here

```
skills/<name>/SKILL.md        the skill (frontmatter plus body)
skills/<name>/references/     optional deep dives, loaded on demand
examples/<skill>/<topic>/     one main.go per snippet, all in one Go module
scripts/check.sh              builds every example, sweeps prose for banned text
.claude-plugin/plugin.json    lets the repo install as a Claude Code plugin
.tessl-plugin/plugin.json     Tessl registry manifest
skills.sh.json                skill groupings for the skills.sh page
```

Skills and the module versions they describe:

| Skill             | Covers                                 | Pinned version |
|-------------------|----------------------------------------|----------------|
| charmed-ecosystem | which library for what, v1 to v2 moves | all below      |
| charmed-bubbletea | charm.land/bubbletea/v2                | v2.1.0         |
| charmed-bubbles   | charm.land/bubbles/v2                  | v2.2.1         |
| charmed-lipgloss  | charm.land/lipgloss/v2                 | v2.0.6         |
| charmed-huh       | charm.land/huh/v2                      | v2.0.3         |
| charmed-log       | charm.land/log/v2                      | v2.0.1         |
| charmed-glamour   | charm.land/glamour/v2                  | v2.0.1         |
| charmed-gum       | gum CLI                                | 2.0.2          |
| charmed-vhs       | vhs CLI and tape files                 | 0.12.1         |

The pins live in `examples/go.mod`.  Bump them there first, run the check, and
only then edit prose.

## Rules for editing a skill

1. Verify before you write.  Use `go doc` against the pinned module in
   `examples/`, or read the module source in the Go module cache.  Do not
   trust a README, a blog post, or an older skill; several upstream READMEs
   still show v1 signatures.
2. Every code block that is Go is either trimmed from a file under
   `examples/` that builds, or a short API sketch.  The first line of the
   block says which: `// examples/<skill>/<topic>/main.go` for the former
   (the check verifies the file exists), `// excerpt: ...` for the latter.
   Prefer a file; use an excerpt only for a signature or field list that
   no example would sensibly contain.  When a trimmed block and its file
   disagree, fix the block.  Shell and tape snippets are exempt from the
   build but not from review.
3. Keep SKILL.md short: a workflow, the core API, and the mistakes, ideally
   under 250 lines.  Put per-topic detail, full option tables, and migration
   tables in `references/`, one level deep, each linked from SKILL.md with a
   line saying what it holds.
4. Frontmatter carries only `name`, `description`, and `license`.
   Visibility on a given machine is a settings concern
   (`skillOverrides` in settings.json), not something baked into the file, so
   the skills stay portable.
5. Descriptions must say when to use the skill and when not to, in one or two
   sentences, and stay under 1024 characters.  Double quote the description:
   a plain YAML scalar containing `: ` does not parse. Skillforge runs a
   real YAML parser over the frontmatter.
6. No em dashes, no en dashes, anywhere: prose, code comments, tables,
   commit messages.  Use a period, comma, colon, or parentheses instead.
   Box drawing characters in diagrams are fine.
7. Write plainly.  No triplet punchlines, no "it's not X, it's Y" pivots, no
   inflated significance.  Say what the API does and what goes wrong when it
   is misused.

## Checking your work

```sh
skillforge validate .
skillforge check .
```

[Skillforge](https://github.com/tpscloud/skillforge-ci) validates frontmatter
with the policy in `.skill-quality.json`: the name matches the directory, the
description is present and within the length limit, and no other keys appear.
`check` runs `scripts/check.sh`, which builds, vets, and tests every Go example,
sweeps Markdown for dashes and v1 imports, and checks that every Go fence in
`SKILL.md` and `references/` names an existing example or is marked as an
excerpt. The script needs Go and standard shell tools.

Run both before every commit. Install Skillforge from its
[releases](https://github.com/tpscloud/skillforge-ci/releases) or build it from
`.skill-quality/` (`go -C .skill-quality build -o /tmp/skillforge .`). CI builds
it from that vendored copy. Do not edit `.skill-quality/` by hand; update it
with `skillforge install . --update` from a newer Skillforge release.

The gum and vhs versions in the table are documentation baselines; the check
does not pin or run those binaries. Run changed shell snippets through
`bash -n` and changed tapes through `vhs validate` with the baseline version.

## Commits

One skill per commit.  Subject line names the skill and what changed, for
example `charmed-huh: add multiselect filtering section`.  Do not bundle
unrelated skills into one commit.

## Adding a skill

1. `mkdir skills/<name>` and write `SKILL.md` with the frontmatter above.
2. Add `examples/<name>/...` for every Go block.
3. Add a row to the table in this file, an entry in README.md, the skill to
   `.skill-quality.json`, `skills.sh.json`, and `.tessl-plugin/plugin.json`,
   and a `LICENSES/MIT.txt` and `NOTICE.md` like the other skills.
4. Run validation and the configured checks, then commit.
