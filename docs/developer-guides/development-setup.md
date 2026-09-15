---
title: Development Setup
description: Setting up a local environment to work on Skillstore itself.
tags: [developer-guide, setup]
---

# Development Setup

## Prerequisites

- Go 1.27+

Node.js is only needed by end users running the generated site's `npx`
install command (see [Stack](./stack.md)) — not for building, testing, or
running the CLI itself.

## Clone and build

```bash
git clone https://gitlab.com/endekasoft/skillstore.git
cd skillstore
go build ./...
```

## Run against the bundled fixtures

```bash
go run . --skill-dir examples/skills --output-dir public
```

`examples/skills/` (33 real skills) is both the manual-testing fixture set
and what `go test` reads — see [Testing](./testing.md).

## Run the checks

```bash
go build ./...
go vet ./...
gofmt -l .          # must print nothing
go test ./...
```

## Iterating on the generated site's HTML/CSS/JS

Templates and assets live under `internal/site/templates/` and
`internal/site/assets/`, embedded into the binary via `go:embed`. There is
no separate build step for them — edit the files, then re-run
`go run . --skill-dir examples/skills --output-dir public` and reload
`public/index.html` in a browser (served over real HTTP; see
[Troubleshooting](../operator-guides/troubleshooting.md) for why
`file://` breaks search and the panel).

## Related

- [Testing](./testing.md) — what the test suite covers and how it's structured
- [Standards](./standards.md) — code style and git workflow conventions
