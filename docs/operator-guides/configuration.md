---
title: Configuration
description: Full reference for the Skill Indexer CLI's flags and the SKILL.md format it reads.
tags: [operator-guide, configuration, reference]
---

# Configuration

Skill Indexer has no environment variables and no config file. Everything is
either a CLI flag or a convention about the shape of a `SKILL.md` file.

## CLI flags

| Flag               | Type               | Default  | Description                                                                                                                                                                                                        |
| ------------------ | ------------------ | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `--skill-dir`      | string             | `skills` | Directory containing skill subdirectories to scan. Each immediate subdirectory must contain a `SKILL.md`.                                                                                                          |
| `--output-dir`     | string             | `public` | Directory to write the generated static site into. **Cleared and recreated on every run** — anything already there is deleted first.                                                                               |
| `--repository-url` | string             | —        | Git URL of the repository holding the skills. Switches install commands to `npx skills add <url> --skill <name>` and adds an "Install everything" banner. See [Repository install mode](#repository-install-mode). |
| `--site-name`      | string             | —        | Site title shown in the header in place of the `skill-indexer` wordmark; also used as the page `<title>`.                                                                                                          |
| `--extra-css`      | string, repeatable | —        | CSS file copied to `assets/extra/` and loaded after the built-in stylesheet. See [Theming](#theming).                                                                                                              |
| `--extra-js`       | string, repeatable | —        | JS file copied to `assets/extra/` and loaded after `app.js`.                                                                                                                                                       |
| `--help` / `-h`    | flag               | —        | Prints usage text.                                                                                                                                                                                                 |

Example:

```bash
skill-indexer --skill-dir path/to/skills --output-dir dist
```

> **Warning:** `--output-dir` is fully wiped (`os.RemoveAll`) before
> generation starts, every run. Never point it at a directory containing
> anything you haven't already committed or backed up elsewhere.

## The `SKILL.md` format

Every skill is a directory containing a `SKILL.md` file with YAML
frontmatter, delimited by `---` lines, followed by an optional Markdown
body (the body is parsed but never rendered anywhere in the generated
site — only the frontmatter `description` is shown to visitors).

```markdown
---
name: my-skill
description: What this skill does and when to use it.
metadata:
  author: your-name
  version: "1.0.0"
license: MIT
---

# My Skill

Anything here is parsed but not displayed in the generated site.
```

| Field             | Required | Notes                                                                                                                                                                                                                                                              |
| ----------------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `name`            | Yes      | Non-empty string. Shown as the card title and panel headline.                                                                                                                                                                                                      |
| `description`     | Yes      | Non-empty string. Shown on the card and in the panel — this is the only prose visitors see.                                                                                                                                                                        |
| `metadata`        | No       | A nested map of any keys you like (`author`, `version`, `source`, ...). Every entry is shown in the panel's metadata table.                                                                                                                                        |
| _(anything else)_ | No       | Any other top-level key — `license`, `compatibility`, `version`, `author`, `dependencies`, `allowed-tools`, etc. — is folded into the same metadata set as `metadata`, with no schema and no type conversion. Strings, lists, and booleans are all accepted as-is. |

A skill directory can contain any other files alongside `SKILL.md`
(`references/`, images, license files, ...) — the whole directory is zipped
for download and `npx` install, not just the frontmatter file.

> **Nota:** if the same key appears both nested under `metadata:` and as a
> top-level key (e.g. `metadata: { version: "1.0" }` **and** a separate
> top-level `version: "2.0"`), the top-level value silently wins — there's
> no warning about the collision. Avoid defining the same key in both
> places.

### What makes a skill invalid

A subdirectory of `--skill-dir` is skipped (with a logged warning, not a
failed run) when:

- It has no `SKILL.md` file at all, or
- Its `SKILL.md` is missing `name` or `description`, either has an empty
  value, or the file's YAML frontmatter fails to parse.

A stray non-directory entry directly inside `--skill-dir` (a stray file
like a `README.md` or `.DS_Store`) is also skipped, but silently — no
warning is logged for it, unlike the cases above.

> **Nota (known limitation):** a `metadata:` key whose value isn't a YAML
> mapping (a plain string, a list, `null`, ...) does **not** make the skill
> invalid. The skill still parses successfully, but the entire `metadata:`
> block is silently dropped — nothing from it reaches the generated site,
> and no warning is logged. Keep `metadata:` a plain key-value mapping.

See [Troubleshooting](./troubleshooting.md) if skills you expect to see are
missing from the output.

## Repository install mode

By default each panel's command installs from the generated zip URL. With
`--repository-url https://host/group/skills.git` it becomes:

```bash
npx skills add https://host/group/skills.git --skill <name> -a claude-code -y
```

and a banner at the top of the page offers `--skill '*'` to install every
skill at once. `<name>` is the skill's frontmatter `name`. The repository
must contain the skills in a layout the `skills` CLI can discover — that is
the operator's responsibility; the generator does not check it. The URL is
passed through verbatim (trailing `/` trimmed). The **Download .zip** and
**Download .skill** buttons stay available in this mode.

## Downloads

Each panel offers **Download .zip** and **Download .skill**. A `.skill` file
is the same archive with a different extension, so both links point at the
one `downloads/<name>.zip` on disk and only the browser's saved file name
differs (via the `download` attribute). That attribute is ignored if the
zips are served from a different origin than the page; the file then saves
under its `.zip` name.

## Theming

`--extra-css` and `--extra-js` (repeatable) copy files into `assets/extra/`
and link them after the built-in `style.css` / `app.js`. Each must end in
`.css` / `.js`, and basenames must be unique across both lists. To re-skin
the site, override the `--color-*`, `--font-*` and `--radius-*` custom
properties from `style.css` in all three places it defines them: `:root`,
`:root[data-theme="dark"]`, and the `prefers-color-scheme: dark` block.

```bash
skill-indexer --extra-css my-theme.css --site-name "Acme Skills"
```

## Site branding

The generated site's footer ("Skill Indexer vX.Y.Z | by ndkprd", linking to
`gitlab.com/endekasoft/skill-indexer` and `gitlab.com/endekasoft`) is
currently **hardcoded** in `internal/site/templates/index.html.tmpl` —
there is no flag or config option to change or remove it (`--site-name` and
`--extra-css` change the header and colors, not the footer). Anyone
generating their own site with this version of Skill Indexer will ship that
same attribution. If you need different branding, you currently have to
edit the template source directly before building.

## Related

- [Installation](./installation.md) — getting the CLI built and runnable
- [Troubleshooting](./troubleshooting.md) — diagnosing skipped or missing skills
