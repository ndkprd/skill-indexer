# skill-indexer

A static site generator for a "skill marketplace": point it at a directory of
[Claude Code skills](https://docs.claude.com/), each with a `SKILL.md`
frontmatter file, and it produces a searchable, browsable static site — a
card grid with a slide-in detail panel, zip downloads, copy-paste `npx`
install commands, and a Nord-based dark mode toggle. No backend required to
host the output; only the `npx` install step touches the network.

## Quickstart

```bash
go build -o skill-indexer .
./skill-indexer --skill-dir examples/skills --output-dir public
```

This scans `examples/skills/` (6 real skills ship in this repo as
sample/fixture data — see [examples/README.md](examples/README.md)),
skipping and warning on any directory missing a valid `SKILL.md`, and writes
the generated site to `public/`. Serve it with any static file server, e.g.:

```bash
python3 -m http.server -d public 8080
```

Point `--skill-dir` at your own skills directory for real use. See
[examples/](examples/) for a GitLab CI pipeline that publishes to GitLab
Pages, a Docker Compose setup, and a Kubernetes Deployment template — plus
the root `Dockerfile` for a containerized build of the CLI itself.

## CLI

```text
skill-indexer [flags]

Flags:
  -h, --help                help for skill-indexer
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

```text
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
`vue/SKILL.md`, `vue/references/foo.md`), so extracting it anywhere
reproduces a `<name>/...` layout matching the `--skill-dir` convention the
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
npx skills add https://your-site.example/downloads/vue.zip -a claude-code -y
```

This uses [`skills`](https://github.com/vercel-labs/skills) — an existing,
published, actively maintained CLI for the open agent skills ecosystem, not
anything this project ships. It supports installing directly from a zip
download URL, which is exactly what `/downloads/<name>.zip` is, so no
custom installer package is needed. A Project/Global toggle next to the
command adds `-g` for the latter. Project scope installs to
`./.claude/skills/<name>/`; global scope to `~/.claude/skills/<name>/`.
`skills` itself requires Node.js 22.20+.

## Development

```bash
go build ./...
go vet ./...
gofmt -l .
go test ./...
```

### Project layout

```text
cmd/               Cobra CLI wiring (root command, flags, pipeline orchestration)
internal/skill/    SKILL.md frontmatter parsing + directory scanning
internal/site/     HTML rendering (html/template + go:embed), zip archiving,
                   search index generation
examples/          sample skills (used by tests too), GitLab CI/Docker
                   Compose/Kubernetes deployment examples, and its own
                   quickstart README
Dockerfile         containerized build of the CLI (not the generated site)
.agents/plans/     implementation plan(s) for this project
PRODUCT.md         strategic design context (users, purpose, anti-references)
DESIGN.md          visual design system (colors, typography, components)
```

## License

[MIT](LICENSE)
