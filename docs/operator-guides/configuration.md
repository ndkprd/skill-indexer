---
title: Configuration
description: Full reference for the Skill Indexer CLI's flags and the SKILL.md format it reads.
tags: [operator-guide, configuration, reference]
---

# Configuration

Skill Indexer has no environment variables and no config file. Everything is
either a CLI flag or a convention about the shape of a `SKILL.md` file.

## CLI flags

| Flag            | Type   | Default  | Description                                                                                                                          |
| --------------- | ------ | -------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| `--skill-dir`   | string | `skills` | Directory containing skill subdirectories to scan. Each immediate subdirectory must contain a `SKILL.md`.                            |
| `--output-dir`  | string | `public` | Directory to write the generated static site into. **Cleared and recreated on every run** — anything already there is deleted first. |
| `--help` / `-h` | flag   | —        | Prints usage text.                                                                                                                   |

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

## Site branding

The generated site's footer ("Skill Indexer vX.Y.Z | by ndkprd", linking to
`gitlab.com/endekasoft/skill-indexer` and `gitlab.com/endekasoft`) is
currently **hardcoded** in `internal/site/templates/index.html.tmpl` —
there is no flag or config option to change or remove it. Anyone
generating their own site with this version of Skill Indexer will ship that
same attribution. If you need different branding, you currently have to
edit the template source directly before building.

## Related

- [Installation](./installation.md) — getting the CLI built and runnable
- [Troubleshooting](./troubleshooting.md) — diagnosing skipped or missing skills
