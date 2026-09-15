---
title: Quick Start
description: Generate and view a Skill Indexer site in under 5 minutes.
tags: [quick-start, cli]
---

# Quick Start

## Prerequisites

- Go 1.27 or later
- A directory of skills, each a folder containing a `SKILL.md` file with
  YAML frontmatter (`name` and `description` required). No skills of your
  own yet? The repo ships sample skills at `examples/skills/`.

## 1. Clone and build

```bash
git clone https://gitlab.com/endekasoft/skillstore.git
cd skill-indexer
go build -o skill-indexer .
```

## 2. Generate the site

```bash
./skill-indexer --skill-dir examples/skills --output-dir public
```

There is no other configuration to set — those two flags are the entire
config surface.

## 3. Serve it

```bash
python3 -m http.server -d public 8080
```

## You should see

Opening `http://localhost:8080` shows a card grid of skills with a search
box centered in the header. Clicking a card slides in a detail panel from
the right with the skill's metadata, an `npx` install command, and a
download button. The CLI itself prints structured log lines to stderr as it
runs, ending with something like:

```text
INF generation complete event=generate_done output_dir=public skill_count=6
```

## Next steps

- Point `--skill-dir` at your own skills directory instead of
  `examples/skills`.
- Publish the output somewhere permanent — see [Deployment](./operator-guides/deployment.md)
  for a Docker option and a GitLab Pages CI pipeline.
- Changing the code? See [Development Setup](./developer-guides/development-setup.md).

## Related

- [Operator Guides](./operator-guides/index.md) — full installation and deployment options
- [Developer Guides](./developer-guides/index.md) — for contributing to Skill Indexer itself
