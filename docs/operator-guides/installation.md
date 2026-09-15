---
title: Installation
description: System requirements and install procedure for the Skill Indexer CLI.
tags: [operator-guide, installation]
---

# Installation

## System requirements

| Requirement         | Version               | Needed for                                                                                                   |
| ------------------- | --------------------- | ------------------------------------------------------------------------------------------------------------ |
| Go                  | 1.27+                 | Building the CLI from source                                                                                 |
| Docker              | any recent version    | Building/running the container image instead                                                                 |
| Docker Compose      | v2 (`docker compose`) | Only for the `examples/docker-compose.yaml` demo                                                             |
| kubectl + a cluster | any recent version    | Only for deploying `examples/kubernetes.yaml`                                                                |
| Node.js             | 22.20+                | Only for end users running the generated `npx` install command — not required to build or run the CLI itself |

Skill Indexer itself has no runtime dependencies beyond the compiled binary:
no database, no external services, no environment-specific configuration.

## Option A: build from source

```bash
git clone https://gitlab.com/endekasoft/skillstore.git
cd skill-indexer
go build -o skill-indexer .
```

This produces a single static binary, `skill-indexer`, in the repository
root. Move it onto your `PATH` if you want to run it from anywhere:

```bash
mv skill-indexer /usr/local/bin/
```

## Option B: build a container image

```bash
docker build -t skill-indexer .
```

The image is based on `gcr.io/distroless/static-debian12:nonroot` and
contains only the compiled binary — no shell, no package manager, running
as a non-root user. See [Deployment](./deployment.md) for how to run it.

## Verify the install

```bash
./skill-indexer --help
```

should print usage text listing `--skill-dir` and `--output-dir`. To verify
end to end, generate the bundled sample site:

```bash
./skill-indexer --skill-dir examples/skills --output-dir /tmp/skill-indexer-check
```

and confirm `/tmp/skill-indexer-check/index.html` and
`/tmp/skill-indexer-check/search-index.json` were created, with no `WRN`
(warning) lines in the CLI's log output.

## Related

- [Configuration](./configuration.md) — the flags and the `SKILL.md` format they operate on
- [Deployment](./deployment.md) — running this in Docker, Compose, Kubernetes, or CI
