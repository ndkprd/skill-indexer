---
title: Skill Indexer
description: A static site generator that turns a directory of Claude Code skills into a searchable marketplace.
tags: [overview, skill-indexer]
---

# Skill Indexer

Skill Indexer is a Go CLI that scans a directory of [Claude Code
skills](https://docs.claude.com/) — each a folder with a `SKILL.md`
frontmatter file — and generates a static, searchable marketplace site: a
card grid with a slide-in detail panel, zip downloads, copy-paste `npx`
install commands, and a Nord-based dark mode toggle. The output is plain
static files; no backend is required to host it.

## Who this is for

- **Site visitors** — developers browsing a generated marketplace to find
  and install a skill. See [User Guides](./user-guides/index.md).
- **Operators** — whoever builds the CLI, runs it against a `skills/`
  directory, and publishes the result (locally, via Docker Compose, on
  Kubernetes, or via CI). See [Operator Guides](./operator-guides/index.md).
- **Contributors** — anyone changing Skill Indexer itself. See [Developer
  Guides](./developer-guides/index.md).

## In This Section

- [Quick Start](./quick-start.md) — generate and view a site in under 5
  minutes
- [User Guides](./user-guides/index.md) — using the generated marketplace
  site
- [Operator Guides](./operator-guides/index.md) — installing, configuring,
  and deploying Skill Indexer
- [Developer Guides](./developer-guides/index.md) — architecture, stack,
  and contributing

## Related

- [gitlab.com/endekasoft/skillstore](https://gitlab.com/endekasoft/skillstore) — source
- [examples/](../examples/README.md) — a runnable quickstart and a
  copy-pasteable GitLab CI pipeline
- [PRODUCT.md](../PRODUCT.md) — the strategic brief this project is built
  against: personas, purpose, and its WCAG 2.1 AA accessibility commitment
- [DESIGN.md](../DESIGN.md) — the visual design system implementing that brief
