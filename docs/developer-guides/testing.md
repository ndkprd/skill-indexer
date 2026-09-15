---
title: Testing
description: How Skillstore's test suite is structured and how to run it.
tags: [developer-guide, testing]
---

# Testing

## Running the Go tests

```bash
go test ./...
```

There's no mocking framework and no test database — every test either
uses in-memory fixtures or reads real files from `examples/skills/`.

## What's covered

| File                                | Covers                                                                                                                                                                                                                                                                                                     |
| ----------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `internal/skill/parse_test.go`      | Frontmatter parsing against real fixtures: nested `metadata:` maps, top-level `license`/`compatibility`, 4-space indentation, a minimal name+description-only skill, a top-level list value (`dependencies`), a top-level string value (`allowed-tools`), plus malformed-YAML and missing-file error cases |
| `internal/skill/scan_test.go`       | `ScanDir` against the full `examples/skills/` fixture set (expects zero warnings, and a skill count read from the directory itself — not hardcoded, see below) and against a constructed temp directory with one deliberately broken skill                                                                 |
| `internal/site/zip_test.go`         | `ZipSkillDir` produces entries prefixed with the skill's directory name (`vue/SKILL.md`, not bare `SKILL.md`)                                                                                                                                                                                              |
| `internal/site/render_test.go`      | `Render` against synthetic in-memory skills: `index.html` contains each skill's name, a `#<dirname>` panel-trigger link, the detail panel skeleton, and the theme toggle; all three static assets get written                                                                                              |
| `internal/site/searchindex_test.go` | `BuildSearchIndex` produces valid JSON with one entry per skill, including the computed `zipPath`                                                                                                                                                                                                          |
| `installer/index.test.js`           | `installFromUrl` against a local `http` fixture server serving a zip built on the fly, plus `parseArgs` edge cases and failure paths (404, corrupt zip, unreachable host)                                                                                                                                  |

## Why the fixture-count assertion reads the directory instead of a number

`scan_test.go`'s "real fixture directory" case counts `examples/skills/`'s
own subdirectories and asserts `ScanDir` returns that many skills, rather
than asserting a literal `33`. The fixture set has already changed size
once (46 → 33, when the ASDP-internal skills were dropped from the public
example set) and silently broke a hardcoded assertion. If you add or
remove a fixture skill, this test keeps passing without an edit.

## Installer tests

```bash
cd installer
bun test   # or: node --test index.test.js
```

Uses only Node's built-in `node:test` + `node:assert` — no test framework
dependency. See [Development Setup](./development-setup.md) for the
bun/Node distinction in this repo's dev shell.

## Related

- [Development Setup](./development-setup.md) — building and running everything locally
- [Pitfalls](./pitfalls.md) — decisions the tests above are specifically guarding
