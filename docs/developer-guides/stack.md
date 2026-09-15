---
title: Stack
description: Languages, frameworks, and libraries Skillstore is built on, and why.
tags: [developer-guide, stack]
---

# Stack

## CLI (Go)

| Component     | Choice                                                    | Why                                                                                                                       |
| ------------- | --------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------- |
| Language      | Go 1.27                                                   | Compiles to a single static binary — no runtime install needed to distribute the generator.                               |
| CLI framework | [`spf13/cobra`](https://github.com/spf13/cobra)           | Standard, well-known flag/command handling for a generate-only root command.                                              |
| Logging       | [`rs/zerolog`](https://github.com/rs/zerolog)             | Structured logging to stderr; used only in `cmd`, never in the library packages.                                          |
| YAML          | [`gopkg.in/yaml.v3`](https://pkg.go.dev/gopkg.in/yaml.v3) | Decodes frontmatter into `map[string]any` natively — no custom indentation handling needed for nested `metadata:` blocks. |
| Templating    | stdlib `html/template` + `go:embed`                       | Keeps the whole generator a single binary: templates and static assets are compiled in, not read from disk at runtime.    |

No web framework, no database, no ORM — the CLI's only job is reading
files and writing files.

## Generated site (client-side)

| Component | Choice                                                                                                             | Why                                                                                                                                               |
| --------- | ------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| Search    | [Fuse.js](https://fusejs.io/) v7, vendored into `internal/site/assets/fuse.min.js`                                 | Fuzzy client-side search with zero server involvement; vendored rather than CDN-loaded so the site works without a third-party script dependency. |
| Fonts     | IBM Plex Sans (UI) + IBM Plex Mono (anything copyable), loaded from Google Fonts with a system-font fallback stack | A superfamily pairing with real technical character — deliberately not Inter/Roboto/Arial. Falls back gracefully if the font request fails.       |
| JS        | Plain vanilla JS (`app.js`), no framework, no build step                                                           | The generated site is meant to be a handful of static files; a framework/bundler would add a build step the Go binary doesn't otherwise need.     |
| CSS       | Plain CSS custom properties, no framework                                                                          | See [DESIGN.md](../../DESIGN.md) for the full token system (light + Nord dark themes).                                                            |

## Installer (Node)

| Component      | Choice                                             | Why                                                                                                                                       |
| -------------- | -------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------- |
| Runtime        | Node.js 18+ (via `npx`)                            | The install command is meant to be copy-pasted by developers who already have Node for other tooling.                                     |
| Zip extraction | [`adm-zip`](https://www.npmjs.com/package/adm-zip) | Pure JS, no native bindings — keeps `npx` invocation fast with no compile step.                                                           |
| Test runner    | Node's built-in `node:test`                        | No extra dependency just for tests; run via `bun test` in this repo's dev shell (Node-compatible), but the shipped code stays plain-Node. |

## Related

- [Architecture](./architecture.md) — how these pieces fit together
- [Development Setup](./development-setup.md) — running everything locally
