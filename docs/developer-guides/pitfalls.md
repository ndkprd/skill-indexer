---
title: Pitfalls
description: Known gotchas and non-obvious decisions worth knowing before changing Skillstore.
tags: [developer-guide, pitfalls]
---

# Pitfalls

Things that look like bugs or missing features on first read, but are
deliberate — plus a couple of real footguns.

## There is exactly one rendered page

`internal/site/templates/index.html.tmpl` is the only template. There are
**no per-skill HTML pages**. Skill detail is populated entirely
client-side, from `search-index.json`, into a slide-in panel (`app.js`).
This was a deliberate architecture change (an earlier version rendered a
`skills/<name>.html` per skill) — don't reintroduce that without removing
this design again. `search-index.json`'s `zipPath` and `metadata` fields
exist specifically so the panel never needs a second fetch or a
server-side route per skill.

## Zip entries are prefixed with the skill's directory name

`ZipSkillDir` writes every entry as `<DirName>/<relative path>` (e.g.
`vue/SKILL.md`), not bare `SKILL.md` at the archive root. This matters
because the `npx` install command extracts straight into `--dest` with no
extra wrapping folder — if zip entries were flat, `--dest ./skills` would
dump files directly into `./skills/` instead of `./skills/vue/`,
colliding with every other installed skill. `zip_test.go` explicitly
asserts against the flat form to catch a regression here.

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
anywhere in `internal/site`. This was also a deliberate removal (an
earlier version rendered it as sanitized Markdown via `goldmark` on the
detail page) — the frontmatter `description` is considered sufficient.
Don't re-add a Markdown render step without confirming that's actually
wanted again; the `goldmark` dependency was removed along with it.

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

## The install command's URL is computed in the browser, not baked in

`app.js` builds the `npx` command from `window.location.origin` at panel-open
time — the Go side has no concept of the site's eventual public URL and
never needs one. If you're debugging a "wrong URL" report, the cause is
almost always how the site is being _served_ (a proxy rewriting the
origin, or the page opened via `file://`), not the generator.

## Related

- [Architecture](./architecture.md) — where each of these decisions sits in the overall design
- [Testing](./testing.md) — which tests exist specifically to catch a regression on the above
