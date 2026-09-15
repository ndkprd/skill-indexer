# skillstore-install

A tiny CLI, meant to be run via `npx`, that downloads a skill zip produced by
`skillstore` and extracts it into a local skills directory.

## Usage

```bash
npx skillstore-install <zip-url> [--dest <dir>]
```

- `<zip-url>` — the URL of a skill zip, e.g.
  `https://your-site.example/downloads/vue.zip`. Every skill detail page on
  a generated site shows this command pre-filled.
- `--dest` — directory to extract into (default: `./skills`).

The zip's own internal entries are already prefixed with the skill's
directory name (e.g. `vue/SKILL.md`), so extraction reproduces
`<dest>/<skill-name>/...` — no extra wrapping folder is added.

Exits non-zero with a message on stderr if the URL is unreachable, the
response isn't a 2xx, or the downloaded content isn't a valid zip.

## Development

Requires Node.js 18+ (this repo's dev shell uses `bun`, which is
Node-compatible enough to develop and test against, but the shipped
`index.js` itself is plain Node — no bun-only APIs).

```bash
bun install   # or npm install
bun test      # or npm test / node --test index.test.js
```

`index.js` exports `installFromUrl(url, dest)` as a plain async function,
separate from the CLI/argv-parsing wrapper, so it can be tested directly
without spawning a subprocess. See `index.test.js` for examples, including a
local `http` fixture server.

## Publishing

This package has **not** been published to the npm registry yet. Every
Skillstore site's `npx skillstore-install ...` command assumes it has
been — until someone runs `npm publish` from this directory (after
reviewing `package.json`'s name, version, and license), that command will
fail with a "package not found" error for anyone who tries it. Publishing
was explicitly deferred as a follow-up when this package was first built,
not forgotten.
