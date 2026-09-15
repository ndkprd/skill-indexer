---
title: Testing
description: How Skill Indexer's test suite is structured and how to run it.
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

| File                                | Covers                                                                                                                                                                                                                                                                                                                                                                       |
| ----------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `internal/skill/parse_test.go`      | Frontmatter parsing: a nested `metadata:` map, a minimal name+description-only skill, and a top-level list/string value (`dependencies`/`license`) all against real fixtures; top-level `license`/`compatibility` with 4-space metadata indentation against a synthetic fixture (no single real skill combines that shape); plus malformed-YAML and missing-file error cases |
| `internal/skill/scan_test.go`       | `ScanDir` against the full `examples/skills/` fixture set (expects zero warnings, and a skill count read from the directory itself — not hardcoded, see below) and against a constructed temp directory with one deliberately broken skill                                                                                                                                   |
| `internal/site/zip_test.go`         | `ZipSkillDir` produces entries prefixed with the skill's directory name (`skill-creator/SKILL.md`, not bare `SKILL.md`)                                                                                                                                                                                                                                                      |
| `internal/site/render_test.go`      | `Render` against synthetic in-memory skills: `index.html` contains each skill's name, a `#<dirname>` panel-trigger link, the detail panel skeleton, and the theme toggle; all three static assets get written                                                                                                                                                                |
| `internal/site/searchindex_test.go` | `BuildSearchIndex` produces valid JSON with one entry per skill, including the computed `zipPath`                                                                                                                                                                                                                                                                            |

## Why the fixture-count assertion reads the directory instead of a number

`scan_test.go`'s "real fixture directory" case counts `examples/skills/`'s
own subdirectories and asserts `ScanDir` returns that many skills, rather
than asserting a literal number. The fixture set has already changed size
twice (46 → 33, when the ASDP-internal skills were dropped from the public
example set; then 33 → 6, trimming further) and would have silently broken
a hardcoded assertion each time. If you add or remove a fixture skill,
this particular test keeps passing without an edit — but see
[Pitfalls](./pitfalls.md) and `AGENTS.md`: the _named_ fixtures other
tests reference aren't self-adjusting the same way.

## Related

- [Development Setup](./development-setup.md) — building and running everything locally
- [Pitfalls](./pitfalls.md) — decisions the tests above are specifically guarding
