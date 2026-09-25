# AGENTS.md

Instructions for AI coding agents working in this repository.

## What this is

`skill-indexer` is a Go CLI that generates a static "skill marketplace"
site from a directory of Claude Code skills. See `README.md` for the user-facing
overview, `PRODUCT.md` for who it's for and why, and `DESIGN.md` for the
visual design system. The implementation plan and its resolved FAQ live at
`.agents/plans/001-skill-marketplace-generator.md` — read it before making
structural changes; it records _why_ several non-obvious decisions were made.

## Build, test, lint

```bash
go build ./...
go vet ./...
gofmt -l .          # must report nothing
go test ./...
```

Run the generator against the real fixtures to sanity-check end to end:

```bash
go run . --skill-dir examples/skills --output-dir public
```

Sample skills live at `examples/skills/`, alongside the
`examples/.gitlab-ci.yml` and `examples/README.md` quickstarts. All test
fixture paths in `internal/skill/*_test.go` and `internal/site/zip_test.go`
point at `../../examples/skills/...`. `internal/skill/scan_test.go`'s count
assertion reads the directory itself rather than hardcoding a number, so
it tolerates the fixture set changing size. The _named_ fixtures other
tests reference (`parse_test.go`'s specific skills for each frontmatter
shape, `zip_test.go`'s skill with a `references/` subdirectory) are
**not** self-adjusting the same way — if you prune `examples/skills/`,
check those tests still point at skills that exist and still have the
shape each subtest needs.

## Conventions

- **Go style**: Google Go style guide (gofmt-clean, wrapped errors via
  `fmt.Errorf("...: %w", err)`, short names in small scopes, doc comments on
  exported identifiers, small single-purpose functions over long ones).
- **Logging**: `github.com/rs/zerolog` to stderr, structured fields
  (`event`, plus whatever's relevant — `skill_dir`, `skill`, counts). Logging
  lives in `cmd/`, not in `internal/skill` or `internal/site` — those
  packages return errors/warnings and let the caller decide how to report
  them.
- **`internal/skill`** owns the `Skill` type and everything about turning a
  `SKILL.md` file into one. `ParseFrontmatter` requires `name`/`description`
  and folds every other top-level frontmatter key — nested `metadata:` map
  included — into `Skill.Metadata` as-is, with no schema and no
  normalization. `ScanDir` never fails the whole scan on one bad skill; it
  warns and skips.
- **`internal/site`** owns everything about turning `[]*skill.Skill` into
  the static site: HTML rendering (`render.go`), zip archiving (`zip.go`),
  and the search index (`searchindex.go`).
  - There is exactly **one** rendered page (`templates/index.html.tmpl`) —
    no per-skill pages. Skill detail lives entirely client-side, in a
    slide-in panel populated from `search-index.json` (see `app.js`). Don't
    reintroduce a `skills/<name>.html` template without removing this again;
    `search-index.json`'s `zipPath`/`metadata` fields exist specifically so
    the panel never needs a second fetch or a server route per skill.
  - Zip entries are prefixed with the skill's `DirName` (e.g.
    `vue/SKILL.md`) so extraction reproduces the `--skill-dir` layout — see
    the plan's FAQ for why this matters, and note that the third-party
    `skills` CLI (see below) also relies on this when installing directly
    from a zip download URL.
  - Card badges use a fixed priority order (`version` → `author` →
    `license` → `compatibility` → first remaining key alphabetically),
    capped at 2. The panel's metadata table shows every key, sorted
    alphabetically client-side (`app.js`'s `formatMetadataValue` mirrors
    `render.go`'s Go-side equivalent — keep both in sync if the format
    logic changes).
  - The `SKILL.md` body is parsed into `Skill.Body` but deliberately never
    rendered anywhere in the UI — the frontmatter `description` is
    considered sufficient. Don't re-add a markdown-to-HTML render step
    without confirming that's actually wanted again.
  - The install command (`app.js`'s `updateInstallCommand`) shells out to
    `npx skills add <zip-url> -a claude-code -y[-g]` — the third-party
    [`vercel-labs/skills`](https://github.com/vercel-labs/skills) CLI, not
    a package this project ships or maintains. `skills` requires Node.js
    22.20+, and its project-scope default is `./.claude/skills/<name>/`,
    not the bare `./skills/<name>/` this project's own `--skill-dir`
    default uses. If you're tempted to bundle a custom installer instead,
    check `skills` actually stopped working first — see
    [Pitfalls](docs/developer-guides/pitfalls.md).
- **Assets are vendored, not CDN-loaded** (except the two Google Fonts
  requests for IBM Plex Sans/Mono, which degrade to the system-font fallback
  stack if unreachable). `internal/site/assets/fuse.min.js` is a vendored
  build, not a runtime `<script src="https://.../fuse.js">` — if it needs
  updating, re-download it and replace the file in place.
- **Theming**: light is the restrained indigo/blue system from `DESIGN.md`;
  dark is Nord (`https://www.nordtheme.com/`) — Polar Night surfaces, Snow
  Storm text, Frost accent — applied via `:root[data-theme="dark"]` and
  mirrored under `@media (prefers-color-scheme: dark) { :root:not([data-theme="light"]) { ... } }`
  so an explicit choice (persisted to `localStorage`, toggled in the header)
  always wins over the OS default. Keep both blocks' variable sets
  identical when adding new tokens. An inline `<script>` at the top of
  `<head>` applies any stored `data-theme` before first paint to avoid a
  flash of the wrong theme — don't move theme-affecting CSS above that
  script's effect, and don't remove the script assuming the FOUC doesn't
  matter.
- **Design system**: follow `DESIGN.md` (flat-by-default elevation with
  feedback-only hover states, `prefers-reduced-motion` respected
  everywhere, focus trapping + `inert` on background content while the
  panel is open). Its Do's/Don'ts section is normative for any new UI work
  in `internal/site/templates` or `internal/site/assets`.

## Git workflow

- Work happens on a feature branch named
  `<model-name>/<feature|fix|refactor>/<slug>`, branched from `main`.
- Never merge or push to `main` without the user's **explicit** approval —
  finishing an implementation is not itself approval to merge.
- Prefer new commits over amending; don't rewrite history that's already
  been reported to the user as done.
- Commit messages follow [Conventional Commits](https://www.conventionalcommits.org/)
  (`<type>(<scope>): <description>`, e.g. `fix(docker): switch runtime image
  to debian-slim`) — the changelog is generated from commit history via
  git-cliff, which filters out non-conventional commits.

## Things to not reintroduce

- No `serve` subcommand — explicitly rejected in favor of a single
  generate-only root command (`skill-indexer --skill-dir ... --output-dir
...`). Users run their own static file server.
- No CSS framework, no JS build step for the generated site — it's rendered
  entirely via `html/template` + `go:embed` so the generator stays one
  self-contained Go binary.
- Don't add a `--base-path`-style prefix without checking — the generated
  site currently uses root-relative URLs (`/downloads/...`, `/assets/...`,
  `#<skill-name>` hashes) and assumes it's hosted at the domain root.
- The `Dockerfile` builds and ships only the CLI binary (distroless,
  non-root) — it never bakes in a `skills/` directory or the generated
  site. Don't add either without a real reason; users mount their own
  skills dir and output dir as volumes.
