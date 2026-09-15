# Product

## Register

product

## Users

Developers already using (or evaluating) Claude Code, arriving in one of two modes:

- **Directed**: they know roughly what they need ("a skill for Vue testing", "gitlab cli
  commands") and land via search or a direct link. They want to confirm fit and get an
  install command as fast as possible.
- **Exploratory**: they have no specific need yet — they're surveying what skills exist,
  evaluating whether to adopt the skills catalog at all, or looking for ideas. They need
  an inviting browse path, not just a search box.

Both personas need the same underlying surface to hold up: fast when they know what they
want, legible and browsable when they don't.

## Product Purpose

A static, generated marketplace for Claude Code skills. Built from a `skills/` directory
of skill definitions (`SKILL.md` frontmatter + body), it renders a searchable card grid,
a detail page per skill (full description, metadata, rendered body), a zip download, and
a copy-paste `npx` install command. No backend — browsing, search, and download all work
from static files; only the `npx` install step touches the network.

Success is a developer going from landing on the site to a copied install command for the
skill they came for in well under a minute — and, secondarily, a catalog that holds up to
open browsing so exploratory visitors don't bounce immediately.

## Brand Personality

Utilitarian, dense, trustworthy. Modeled on the feel of the npm registry and the VS Code
Marketplace: function over flourish, information visible at a glance rather than hidden
behind clicks, monospace accents for code and install commands, near-zero decorative
color. This should read as established developer infrastructure, not a marketing site —
the skill descriptions and rendered `SKILL.md` content are the pitch; nothing is layered
on top to sell it.

## Anti-references

- **Generic SaaS marketing template.** No hero-metric blocks, no gradient-text headlines,
  no stock-illustration energy, no "trusted by" logo strips. This site isn't selling
  anything.
- **Enterprise admin-dashboard bloat.** No heavy nested sidebar navigation, no dense
  data-table chrome, no dashboard-style stat tiles. This is a browse/search surface, not
  an app shell.
- **Generic AI-slop aesthetic.** No default purple-gradient-on-dark, no glassmorphism, no
  side-stripe accent borders on cards — the tells that make an interface read as
  AI-generated on sight.

## Design Principles

- **Speed to install beats delight.** Every design decision should shorten the path from
  landing to a copied, working install command. Never trade friction for polish.
- **Density with clarity.** Surface real metadata (name, description, key frontmatter
  fields) up front instead of hiding it behind clicks — but never at the cost of
  scanability. Dense is fine; cluttered is not.
- **Earn trust through restraint.** A developer tool signals reliability by looking plain
  and functional, not by trying to impress. Spend visual energy on legibility and
  hierarchy, not decoration.
- **Support two reading speeds.** The same surface must work for someone who knows
  exactly what they want (search, scan, install) and someone who's just browsing to see
  what's possible.
- **Show, don't market.** Skill descriptions and rendered `SKILL.md` bodies are the
  pitch. No persuasive copy layered on top of what the skill author actually wrote.

## Accessibility & Inclusion

WCAG 2.1 AA baseline. Full keyboard navigation across the search box, card grid, and
detail page links/actions. Sufficient color contrast in all states. Respect
`prefers-reduced-motion` — disable or simplify hover/transition animation for users who
request it. No color-only signifiers: any status or metadata badge must carry a text or
icon cue, not hue alone.
