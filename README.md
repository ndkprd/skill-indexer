# skill-repo-store

A static site generator for a "skill marketplace": point it at a directory of
[Claude Code skills](https://docs.claude.com/), each with a `SKILL.md`
frontmatter file, and it produces a searchable, browsable static site —
card grid, per-skill detail pages, zip downloads, and copy-paste `npx`
install commands. No backend required to host the output; only the `npx`
install step touches the network.

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
into a free-form metadata map and shown on the skill's detail page).

## Output layout

```
public/
  index.html              card grid + search
  skills/<name>.html      one detail page per skill
  downloads/<name>.zip    the skill's full directory, zipped
  search-index.json       client-side search index (Fuse.js)
  assets/                 style.css, app.js, vendored fuse.min.js
```

Each zip's entries are prefixed with the skill's own directory name (e.g.
`vue/SKILL.md`, `vue/references/foo.md`), so extracting it under
`--dest ./skills` reproduces the same `./skills/<name>/...` layout the
generator itself expects as input.

## Installing a skill via npx

Every skill detail page shows a command like:

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
