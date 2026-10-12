# Security

## Reporting a problem

Report a security issue in these skills or their examples by email to
[security@tpscloud.services](mailto:security@tpscloud.services) rather than in a
public issue.

## Scan results

All five analyzers scanned all nine skills of release 0.1.2 on 2026-10-12
UTC. The static analyzers scanned a clean export of
`60920aa6b476cf871e8b20f57a180d7f69e2805d`; Tessl reviewed the same skill
content. Later commits change no skill directory.

| Skill | <a href="https://github.com/bueti/skilllint" target="_blank" rel="noopener noreferrer">skilllint</a><br>0.3.0 | <a href="https://github.com/xyiqq/skilldoctor" target="_blank" rel="noopener noreferrer">skilldoctor</a><br>0.2.4 | <a href="https://github.com/NVIDIA/SkillSpector" target="_blank" rel="noopener noreferrer">SkillSpector</a><br>2.12.0 | <a href="https://github.com/cisco-ai-defense/skill-scanner" target="_blank" rel="noopener noreferrer">Cisco Skill Scanner</a><br>2.2.2 | <a href="https://docs.tessl.io/improving-your-skills/reviewing-skills" target="_blank" rel="noopener noreferrer">Tessl</a><br>0.114.0 |
|---|---|---|---|---|---|
| [charmed-bubbles](skills/charmed-bubbles/SKILL.md) | pass | pass | pass | pass | 100/100 |
| [charmed-bubbletea](skills/charmed-bubbletea/SKILL.md) | pass | pass | pass | pass | 100/100 |
| [charmed-ecosystem](skills/charmed-ecosystem/SKILL.md) | pass | pass | pass | pass | 100/100 |
| [charmed-glamour](skills/charmed-glamour/SKILL.md) | pass | pass | pass | pass | 100/100 |
| [charmed-gum](skills/charmed-gum/SKILL.md) | pass | pass | pass | pass | 100/100 |
| [charmed-huh](skills/charmed-huh/SKILL.md) | pass | pass | pass | pass | 100/100 |
| [charmed-lipgloss](skills/charmed-lipgloss/SKILL.md) | pass | pass | pass | pass | 100/100 |
| [charmed-log](skills/charmed-log/SKILL.md) | pass | pass | pass | pass | 100/100 |
| [charmed-vhs](skills/charmed-vhs/SKILL.md) | pass | pass | pass | pass | 100/100 |

SkillSpector passes reflect the reviewed false positives listed below.
SkillSpector and Cisco used static analysis. Tessl ran its hosted quality
review.

### Accepted false positives

- **SkillSpector `EA3` (all skills):** matches `NOT LIMITED TO` in the
  [MIT warranty disclaimer](skills/charmed-gum/LICENSES/MIT.txt). This is license
  text about warranty coverage, not an instruction to the agent. The exception
  covers this exact license passage only.
- **SkillSpector `AS3` (charmed-ecosystem):** flags six GitHub links in the
  [library table](skills/charmed-ecosystem/SKILL.md#import-paths) to
  this repository's Bubble Tea, Bubbles, Lip Gloss, Huh, Glamour, and Log skills.
  These are documentation links used to select a library skill. They do not
  enumerate installed skills or inspect private agent configuration. The
  exception covers these six links only.

License text and links are unchanged; the exceptions only change how the
results are shown.

These checks cover their stated tools and scope, not every possible behavior.

## Independent reports

skills.sh runs three security audits on each published skill, and Tessl
reviews each skill's quality. All nine skills passed every audit and scored
100/100 on 2026-10-12. Each result links to that skill's live report; the
README badges link to the same pages.

| Skill | Gen Agent Trust Hub | Socket | Snyk | Tessl quality |
|---|---|---|---|---|
| [charmed-bubbles](skills/charmed-bubbles/SKILL.md) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-bubbles/security/agent-trust-hub) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-bubbles/security/socket) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-bubbles/security/snyk) | [100/100](https://tessl.io/registry/tpscloud/charmedimsure/quality/charmed-bubbles) |
| [charmed-bubbletea](skills/charmed-bubbletea/SKILL.md) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-bubbletea/security/agent-trust-hub) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-bubbletea/security/socket) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-bubbletea/security/snyk) | [100/100](https://tessl.io/registry/tpscloud/charmedimsure/quality/charmed-bubbletea) |
| [charmed-ecosystem](skills/charmed-ecosystem/SKILL.md) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-ecosystem/security/agent-trust-hub) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-ecosystem/security/socket) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-ecosystem/security/snyk) | [100/100](https://tessl.io/registry/tpscloud/charmedimsure/quality/charmed-ecosystem) |
| [charmed-glamour](skills/charmed-glamour/SKILL.md) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-glamour/security/agent-trust-hub) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-glamour/security/socket) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-glamour/security/snyk) | [100/100](https://tessl.io/registry/tpscloud/charmedimsure/quality/charmed-glamour) |
| [charmed-gum](skills/charmed-gum/SKILL.md) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-gum/security/agent-trust-hub) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-gum/security/socket) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-gum/security/snyk) | [100/100](https://tessl.io/registry/tpscloud/charmedimsure/quality/charmed-gum) |
| [charmed-huh](skills/charmed-huh/SKILL.md) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-huh/security/agent-trust-hub) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-huh/security/socket) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-huh/security/snyk) | [100/100](https://tessl.io/registry/tpscloud/charmedimsure/quality/charmed-huh) |
| [charmed-lipgloss](skills/charmed-lipgloss/SKILL.md) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-lipgloss/security/agent-trust-hub) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-lipgloss/security/socket) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-lipgloss/security/snyk) | [100/100](https://tessl.io/registry/tpscloud/charmedimsure/quality/charmed-lipgloss) |
| [charmed-log](skills/charmed-log/SKILL.md) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-log/security/agent-trust-hub) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-log/security/socket) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-log/security/snyk) | [100/100](https://tessl.io/registry/tpscloud/charmedimsure/quality/charmed-log) |
| [charmed-vhs](skills/charmed-vhs/SKILL.md) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-vhs/security/agent-trust-hub) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-vhs/security/socket) | [pass](https://www.skills.sh/tpscloud/charmedimsure/charmed-vhs/security/snyk) | [100/100](https://tessl.io/registry/tpscloud/charmedimsure/quality/charmed-vhs) |

For the whole plugin, Tessl also publishes
[eval results](https://tessl.io/registry/tpscloud/charmedimsure/0.1.2/evals)
and a [security report](https://tessl.io/registry/tpscloud/charmedimsure/0.1.2/security)
for 0.1.2.
