# Security

## Reporting a problem

Report a security issue in these skills or their examples by email to
[security@tpscloud.services](mailto:security@tpscloud.services) rather than in a
public issue.

## Scan results

skilllint, skilldoctor, SkillSpector, and Cisco Skill Scanner scanned all
nine skills from a clean export of `417556cb141cc2aea495ffe0c60d93f234f45dd3`
on 2026-10-12 UTC. Tessl reviewed the 0.1.2 skill content on 2026-10-12. The
0.1.2 changes trim prose in five skills and add one example program to
charmed-bubbles; the four static analyzers have not been rerun on them yet.
README badges show these results.

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

- [Tessl quality review](https://tessl.io/registry/tpscloud/charmedimsure/0.1.2/quality)
  and [eval results](https://tessl.io/registry/tpscloud/charmedimsure/0.1.2/evals)
  for the published 0.1.2 plugin.
- [Tessl security report](https://tessl.io/registry/tpscloud/charmedimsure/0.1.2/security):
  the registry's security scan of 0.1.2, including Snyk's results.
