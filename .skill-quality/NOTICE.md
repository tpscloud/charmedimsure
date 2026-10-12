# License and dependencies

Skillforge CI is Copyright 2026 John Deatherage (@coversheets), licensed under MIT. This includes
the Go implementation, tests, workflow templates, and authored documentation.
The installed bundle carries this notice and the full license.

The implementation supports the configuration and report formats used by the
one-skill-to-rule-them-all toolkit. It does not bundle that toolkit's Python
code or the downloaded skills.sh JSON schema.

Go resolves the dependencies pinned in go.mod and go.sum. Their license texts
are retained in licenses/ for distributions of the compiled binary:

| Dependency | Version | License |
| --- | --- | --- |
| github.com/spf13/cobra | v1.10.2 | Apache-2.0 |
| github.com/spf13/pflag | v1.0.9 | BSD-3-Clause |
| github.com/inconshreveable/mousetrap | v1.1.0 | Apache-2.0 |
| github.com/yuin/goldmark | v1.8.6 | MIT |
| gopkg.in/yaml.v3 | v3.0.1 | MIT and Apache-2.0 notices |

The Go runtime and standard library retain their BSD-style license in
licenses/GO.txt.

Gitleaks v8.30.1 is downloaded separately for each scan, verified against a
pinned archive digest, and removed afterward. Its MIT license is retained in
licenses/GITLEAKS.txt. GitHub Actions run separately under their own terms;
their implementations are not copied into this repository.
