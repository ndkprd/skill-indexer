# skill-repo-store

A static site generator for a "skill marketplace": point it at a directory of
[Claude Code skills](https://docs.claude.com/), each with a `SKILL.md`
frontmatter file, and it produces a searchable, browsable static site — a
card grid with a slide-in detail panel, zip downloads, copy-paste `npx`
install commands, and a Nord-based dark mode toggle. No backend required to
host the output; only the `npx` install step touches the network.

## Quickstart

```bash
go build -o skill-repo-store .
./skill-repo-store --skill-dir skills --output-dir public
```

This scans `skills/` (46 real skills ship in this repo as sample/fixture
data), skipping and warning on any directory missing a valid `SKILL.md`, and
writes the generated site to `public/`. Serve it with any static file
server, e.g.:

```bash
python3 -m http.server -d public 8080
```

## CLI

```
skill-repo-store [flags]

Flags:
  -h, --help                help for skill-repo-store
      --output-dir string   directory to write the generated static site into (default "public")
      --skill-dir string    directory containing skill subdirectories to scan (default "skills")
```

Each subdirectory of `--skill-dir` must contain a `SKILL.md` with YAML
frontmatter (`name` and `description` are required; everything else —
`metadata`, `license`, `compatibility`, `version`, `author`, ... — is folded
into a free-form metadata map and shown in the skill's detail panel). The
rendered SKILL.md body is not shown anywhere in the UI — the frontmatter
`description` is the only prose surfaced to visitors, by design.

## Output layout

```
public/
  index.html              the single page: header, card grid, detail panel
  downloads/<name>.zip    each skill's full directory, zipped
  search-index.json       drives both client-side search and the detail
                           panel (Fuse.js reads it; so does app.js)
  assets/                 style.css, app.js, vendored fuse.min.js
```

There are no per-skill HTML pages. Clicking a card opens a slide-in panel
(populated client-side from `search-index.json`) rather than navigating
away; the panel is deep-linkable via `#<skill-name>` in the URL and closes
on Escape, backdrop click, or browser back.

Each zip's entries are prefixed with the skill's own directory name (e.g.
`vue/SKILL.md`, `vue/references/foo.md`), so extracting it under
`--dest ./skills` reproduces the same `./skills/<name>/...` layout the
generator itself expects as input.

## Dark mode

A theme toggle in the header switches between the default light theme and a
[Nord](https://www.nordtheme.com/)-based dark theme. The choice persists via
`localStorage`; absent an explicit choice, the site follows the OS's
`prefers-color-scheme`. All theming is pure CSS custom properties in
`internal/site/assets/style.css` — no separate dark asset build.

## Installing a skill via npx

Every skill's detail panel shows a command like:

```bash
npx skill-repo-store-install https://your-site.example/downloads/vue.zip --dest ./skills
```

The `installer/` directory is that standalone npm package
(`skill-repo-store-install`). It fetches a zip by URL and extracts it into
`--dest` (default `./skills`). See [installer/README.md](installer/README.md)
for details.

## Development

```bash
go build ./...
go vet ./...
gofmt -l .
go test ./...
```

The `installer/` package is a separate Node.js project:

```bash
cd installer
bun install   # or npm install
bun test      # or npm test
```

### Project layout

```
cmd/               Cobra CLI wiring (root command, flags, pipeline orchestration)
internal/skill/    SKILL.md frontmatter parsing + directory scanning
internal/site/     HTML rendering (html/template + go:embed), zip archiving,
                   search index generation
installer/         standalone npx-installable Node package
skills/            sample/fixture skills used for manual runs and tests
.agents/plans/     implementation plan(s) for this project
PRODUCT.md         strategic design context (users, purpose, anti-references)
DESIGN.md          visual design system (colors, typography, components)
```

See `AGENTS.md` for conventions to follow when working in this codebase.
