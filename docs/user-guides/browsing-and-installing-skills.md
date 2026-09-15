---
title: Browsing and Installing Skills
description: How to search, inspect, and install a skill from a generated Skillstore site.
tags: [user-guide, search, install]
---

# Browsing and Installing Skills

## Finding a skill

1. Open the site. Every available skill shows as a card: a name, a short
   description, and up to two badges (for example a version number or
   author) when that information exists. Cards are sorted alphabetically
   by name (case-sensitive, so all-capitals names sort before lowercase
   ones).
2. Type into the search box in the header. Results filter as you type,
   matching against both the skill's name and its description.
3. If nothing matches, the page tells you: "No skills match '&lt;your
   search&gt;'. Try a different search term." Clear the box to see every
   skill again.

## Viewing a skill's details

1. Click any card. A panel slides in from the right with the skill's full
   name, its directory name, description, and — if the skill has any —
   a table of every other frontmatter field (license, version, author,
   compatibility, or anything else its author included).
2. Close the panel with the **×** button, by clicking outside it, by
   pressing `Escape`, or with your browser's back button. Each
   skill's panel has its own shareable URL (the page address gets a
   `#skill-name` suffix while open) — copy it from your address bar to
   send someone directly to that skill.

## Installing a skill

Each panel offers two ways to get a skill:

- **`npx` install** (top of the panel): copy the pre-filled command with
  the **Copy** button and run it in your terminal. It downloads the
  skill and extracts it for you, using
  [`skills`](https://github.com/vercel-labs/skills), an existing
  third-party package manager for Claude Code and other coding agents'
  skills — requires Node.js 22.20+.
- **Download .zip**: downloads the skill's folder as a zip archive to
  extract yourself. Works regardless of what's installed on your machine.

### Choosing where it installs

The `npx` command has a **Project** / **Global** toggle above it:

- **Project** (default) installs into `./.claude/skills/` in your current
  directory.
- **Global** installs into `~/.claude/skills/`, available to every project
  on your machine.

Your choice is remembered while you keep browsing, so you don't have to
reselect it for every skill.

## Dark mode

Click the sun/moon icon in the header to switch between light and dark
(Nord-based) themes. Your choice is remembered on this device; until you
choose one, the site follows your operating system's light/dark setting
automatically.

## If something doesn't work

- **"Could not load skill data" banner, search box disabled**: the site
  couldn't fetch its search index (usually a network hiccup, or the site
  wasn't served correctly — see [Troubleshooting](../operator-guides/troubleshooting.md)
  if you run the site yourself). Reload the page.
- **Clicking a card does nothing**: this points to the same cause as
  above — the page failed to load its skill data. Check for the banner
  described above.
- **The `npx` command doesn't work**: make sure you have Node.js 22.20+
  installed (`npx` ships with it) — this is a real requirement of the
  `skills` package the command uses, not something specific to this site.
- **Clicking Copy does nothing** (no "Copied" confirmation): your browser
  is blocking clipboard access, typically because the site isn't served
  over HTTPS (or `localhost`). Select the command text in the box and
  copy it manually instead.

## Related

- [Quick Start](../quick-start.md) — generate a site to try this on
- [Troubleshooting](../operator-guides/troubleshooting.md) — fixing a site that isn't working correctly
