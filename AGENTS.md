# AGENTS.md

Instructions for AI coding agents working in this repository.

## What this is

`skill-repo-store` is a Go CLI that generates a static "skill marketplace"
site from a directory of Claude Code skills. See `README.md` for the user-facing
overview, `PRODUCT.md` for who it's for and why, and `DESIGN.md` for the
visual design system. The implementation plan and its resolved FAQ live at
`.agents/plans/001-skill-marketplace-generator.md` — read it before making
structural changes; it records *why* several non-obvious decisions were made.

## Build, test, lint

```bash
go build ./...
go vet ./...
gofmt -l .          # must report nothing
go test ./...
```

Run the generator against the real fixtures to sanity-check end to end:

```bash
go run . --skill-dir skills --output-dir public
```

The `installer/` directory is a separate Node.js package with its own
lifecycle (this shell has `bun` but not `node`/`npm`/`npx` — use `bun`
locally; the shipped `installer/index.js` itself must stay plain
Node-compatible, no bun-only APIs):

```bash
cd installer && bun install && bun test
```

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
  - Zip entries are prefixed with the skill's `DirName` (e.g.
    `vue/SKILL.md`) so extraction reproduces the `--skill-dir` layout — see
    the plan's FAQ for why this matters to the installer.
  - The index and detail page templates are parsed as **two separate**
    `*template.Template` sets (each paired with `layout.html.tmpl`), because
    both define a `content` block under the same name. Parsing them
    together would let one silently clobber the other — don't merge them
    into a single `ParseFS` call without fixing that first.
  - Card badges use a fixed priority order (`version` → `author` →
    `license` → `compatibility` → first remaining key alphabetically),
    capped at 2. The full metadata table on a detail page shows every key,
    sorted alphabetically for deterministic output.
- **Assets are vendored, not CDN-loaded** (except the two Google Fonts
  requests for IBM Plex Sans/Mono, which degrade to the system-font fallback
  stack if unreachable). `internal/site/assets/fuse.min.js` is a vendored
  build, not a runtime `<script src="https://.../fuse.js">` — if it needs
  updating, re-download it and replace the file in place.
- **Design system**: follow `DESIGN.md` (restrained color strategy, one
  indigo/blue accent, flat-by-default elevation with feedback-only hover
  states, `prefers-reduced-motion` respected everywhere). Its Do's/Don'ts
  section is normative for any new UI work in `internal/site/templates` or
  `internal/site/assets`.

## Git workflow

- Work happens on a feature branch named
  `<model-name>/<feature|fix|refactor>/<slug>`, branched from `main`.
- Never merge or push to `main` without the user's **explicit** approval —
  finishing an implementation is not itself approval to merge.
- Prefer new commits over amending; don't rewrite history that's already
  been reported to the user as done.

## Things to not reintroduce

- No `serve` subcommand — explicitly rejected in favor of a single
  generate-only root command (`skill-repo-store --skill-dir ... --output-dir
  ...`). Users run their own static file server.
- No CSS framework, no JS build step for the generated site — it's rendered
  entirely via `html/template` + `go:embed` so the generator stays one
  self-contained Go binary.
- Don't add a `--base-path`-style prefix without checking — the generated
  site currently uses root-relative URLs (`/skills/...`, `/downloads/...`,
  `/assets/...`) and assumes it's hosted at the domain root.
