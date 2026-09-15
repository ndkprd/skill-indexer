---
title: Development Setup
description: Setting up a local environment to work on Skillstore itself.
tags: [developer-guide, setup]
---

# Development Setup

## Prerequisites

- Go 1.27+
- Node.js 18+ or [Bun](https://bun.sh/) for the `installer/` package (this
  repo's own dev environment uses Bun as a Node-compatible runtime — see
  the note below)

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

## Working on the installer package

`installer/` is a separate Node.js project with its own lifecycle:

```bash
cd installer
bun install && bun test   # or: npm install && npm test
```

> **Note:** this repo's shell has `bun` but not `node`/`npm`/`npx`
> installed directly. Bun is Node-compatible enough to develop and test
> against, but `installer/index.js` itself must stay plain-Node-compatible
> (CommonJS, global `fetch`, no bun-only APIs) since end users invoke it
> via real `npx`.

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
