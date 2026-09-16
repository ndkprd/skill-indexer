---
title: Pitfalls
description: Known gotchas and non-obvious decisions worth knowing before changing Skill Indexer.
tags: [developer-guide, pitfalls]
---

# Pitfalls

Things that look like bugs or missing features on first read, but are
deliberate — plus a couple of real footguns.

## There is exactly one rendered page

`internal/site/templates/index.html.tmpl` is the only template. There are
**no per-skill HTML pages**. Skill detail is populated entirely
client-side, from `search-index.json`, into a slide-in panel (`app.js`).
`search-index.json`'s `zipPath` and `metadata` fields exist specifically
so the panel never needs a second fetch or a server-side route per skill
— don't add a `skills/<name>.html` template without removing this design
first.

## Zip entries are prefixed with the skill's directory name

`ZipSkillDir` writes every entry as `<DirName>/<relative path>` (e.g.
`vue/SKILL.md`), not bare `SKILL.md` at the archive root. This matters
because the `npx skills add <zip-url>` install command (see below) expects
a zip's own contents to already carry the right top-level folder name — if
zip entries were flat, extraction would dump files directly into the
install target instead of a `vue/` subfolder within it, colliding with
every other installed skill. `zip_test.go` explicitly asserts against the
flat form to catch a regression here.

## Frontmatter metadata has no schema, and that's intentional

`ParseFrontmatter` folds _every_ top-level frontmatter key it doesn't
recognize (`license`, `compatibility`, `version`, `author`,
`dependencies`, `allowed-tools`, anything) into `Skill.Metadata`, with no
type coercion and no normalization. Real skills in `examples/skills/`
disagree on where things go (some nest everything under `metadata:`, some
put `license`/`compatibility` at the top level, some use 4-space
indentation) — the parser accommodates all of it rather than picking one
"correct" shape and rejecting the rest.

Two sharp edges follow from this in `skillFromFields` (`internal/skill/parse.go`):

- The nested-`metadata:` loop runs first and the top-level-key loop runs
  second with no "already set" guard, so a key defined in both places
  silently takes the top-level value — there's no error and no warning.
- `coerceMap` returns `nil` for a `metadata:` value that isn't a YAML
  mapping, and the code ranges over that `nil` as a no-op. A malformed
  `metadata:` block doesn't fail parsing; it just silently disappears.
  If you're touching this function, consider whether either of these
  should become a hard error instead of silent behavior.

## Card badges and the panel's metadata table format values independently, in two languages

`internal/site/render.go`'s `pickBadges`/`formatMetadataValue` (Go, used
for card badges) and `app.js`'s `formatMetadataValue` (JS, used for the
panel's metadata table) implement the same "turn an arbitrary YAML-decoded
value into a display string" logic separately, by hand, in two languages.
They're meant to stay behaviorally identical (join lists with `", "`,
`String()`/`fmt.Sprint()` everything else) — if you change how one
formats a value, change the other the same way, or the badges on a card
and the same field's row in its detail panel will disagree.

## The `SKILL.md` body is parsed but never rendered

`Skill.Body` is populated by `ParseFrontmatter` and then never used
anywhere in `internal/site` — the frontmatter `description` is considered
sufficient on its own. Don't add a Markdown render step for it without
confirming that's actually wanted; there's no `goldmark`-style dependency
in `go.mod` to reach for, so adding one back is a real decision, not a
revert.

## A skill directory missing `SKILL.md` never fails the whole run

`ScanDir` collects parse failures as warnings and skips that skill,
rather than returning an error. `runGenerate` in `cmd/root.go` logs each
one and keeps going. If you're adding a new failure mode to
`ParseFrontmatter`, make sure it surfaces as an `error` return (which
`ScanDir` already turns into a skipped-with-warning skill) rather than a
panic — a panic here would take down the entire generate run over one bad
file.

## `--output-dir` is destroyed on every run

`runGenerate` calls `os.RemoveAll(outputDir)` before writing anything.
This is intentional (no stale files from skills removed since a previous
run) but has no confirmation prompt and no dry-run flag. Never point
`--output-dir` at something you didn't mean to delete.

## The `npx` install command depends on a third-party CLI this project doesn't own

`app.js`'s `updateInstallCommand` generates
`npx skills add <zip-url> -a claude-code -y[-g]`, wrapping
[`vercel-labs/skills`](https://github.com/vercel-labs/skills). This
project builds and ships nothing to make that command work — it relies on
`skills` continuing to support installing directly from a zip download URL
and continuing to accept `-a`/`-y`/`-g` the way it does today. If that
package's behavior changes incompatibly, the fix is in `app.js`'s command
string, not in any Go code. Don't reach for a bundled installer package as
the fix without confirming `skills` actually stopped working.

Two things worth knowing about `skills` itself:

- It requires **Node.js 22.20+**.
- Its project-scope default is `./.claude/skills/<name>/` (agent-
  namespaced), not the bare `./skills/<name>/` this project's own
  `--skill-dir` default uses. Its global scope (`~/.claude/skills/`)
  matches ours exactly.

## The install command's URL is computed in the browser, not baked in

`app.js` builds the `npx` command from `window.location.origin` at panel-open
time — the Go side has no concept of the site's eventual public URL and
never needs one. If you're debugging a "wrong URL" report, the cause is
almost always how the site is being _served_ (a proxy rewriting the
origin, or the page opened via `file://`), not the generator.

## Named test fixtures don't self-adjust if `examples/skills/` changes

`internal/skill/scan_test.go`'s skill-count assertion reads
`examples/skills/`'s own directory listing rather than a hardcoded number
(see [Testing](./testing.md)), so it tolerates the fixture set changing
size. But `parse_test.go` and `zip_test.go` reference specific _named_
skills to test specific frontmatter shapes (a nested `metadata:` map, a
minimal skill, a `references/` subdirectory, ...) — those aren't
self-adjusting. If a skill those tests depend on is ever removed from
`examples/skills/`, the tests fail with a plain "no such file or
directory", not a helpful message pointing at the real cause. If you
prune `examples/skills/`, run `go test ./...` afterward and pick new
named fixtures for whatever shape each failing subtest needs.

## Related

- [Architecture](./architecture.md) — where each of these decisions sits in the overall design
- [Testing](./testing.md) — which tests exist specifically to catch a regression on the above
