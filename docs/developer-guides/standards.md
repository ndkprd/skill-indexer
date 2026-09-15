---
title: Standards
description: Code style, logging conventions, and git workflow for Skillstore.
tags: [developer-guide, standards]
---

# Standards

## Go style

- [Google's Go style guide](https://google.github.io/styleguide/go/): `gofmt`-clean,
  short variable names in small scopes, wrapped errors via
  `fmt.Errorf("...: %w", err)`, doc comments on every exported identifier.
- Small, single-purpose functions over long ones — e.g. `internal/site/zip.go`
  splits directory-walking, entry-naming, and file-copying into three
  separate functions rather than one large one.
- Package boundaries are strict: `internal/skill` only parses/scans;
  `internal/site` only renders/zips/indexes. Neither imports `cmd`, and
  neither logs — see below.

## Logging

All logging goes through `github.com/rs/zerolog` to stderr, and lives
**only** in `cmd/root.go`. `internal/skill` and `internal/site` return
errors and warnings; the caller (`cmd`) decides how to report them. This
keeps both packages usable as plain libraries without a logging framework
opinion baked in.

Log fields are structured (`event`, `skill_dir`, `output_dir`,
`valid_count`, `skipped_count`, etc.) rather than freeform interpolated
strings — see the `log.Info()...` / `log.Warn()...` calls in
`cmd/root.go` for the pattern to follow when adding new log lines.

## Front-end conventions

Follow [DESIGN.md](../../DESIGN.md) for anything touching
`internal/site/templates/` or `internal/site/assets/` — its Do's/Don'ts
section is normative (restrained color use, flat-by-default elevation,
`prefers-reduced-motion` support, no rendered Markdown body, no per-skill
pages). See [Pitfalls](./pitfalls.md) for the specific non-obvious
decisions worth knowing before changing this area.

[PRODUCT.md](../../PRODUCT.md) commits the generated site to a WCAG 2.1 AA
accessibility baseline: full keyboard navigation, a trapped focus while
the detail panel is open, no color-only signifiers, and respecting
`prefers-reduced-motion`. `app.js`'s `setInertOutsidePanel`/
`onPanelKeydown` (focus trap and `inert` toggling) and the
`prefers-reduced-motion` handling in both `app.js` and `style.css` exist
specifically to satisfy that commitment — preserve them when touching
panel or motion code.

## Git workflow

- Work happens on a feature branch named
  `<model-name>/<feature|fix|refactor>/<slug>`, branched from `main`.
- Prefer new commits over amending; don't rewrite history that's already
  been reported as done.
- Never merge or push to `main` without explicit approval — finishing an
  implementation is not itself approval to merge.

## Related

- [Testing](./testing.md) — how to verify a change before committing
- [Contributing](./contributing.md) — the practical commit/branch/review flow
