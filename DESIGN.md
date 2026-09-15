<!-- SEED: re-run /impeccable document once there's code to capture the actual tokens and components. -->

---
name: Skill Repo Store
description: A static marketplace for Claude Code skills — browse, search, install.
---

# Design System: Skill Repo Store

## 1. Overview

**Creative North Star: "The Registry Desk"**

Dense, plain, and fast: a page that behaves like developer infrastructure, not a
pitch. The reference lineage is npm's package cards, VS Code Marketplace's
extension grid, and GitHub's file browser — surfaces that put real metadata in
front of you immediately and never ask you to scroll past a hero section to find
it. Color stays restrained (tinted neutrals, one small indigo/blue accent used
sparingly) so the accent still means something when it appears. Motion is present
but modest: cards respond to hover, search results crossfade, a copy button
confirms itself — nothing choreographed, nothing that makes the visitor wait to
act.

This system explicitly rejects the Vercel/Linear-style animated SaaS landing
page: no scroll-triggered reveals, no gradient hero text, no oversized marketing
typography standing between the visitor and the card grid.

**Key Characteristics:**
- Information visible at a glance, not behind a click or a hover
- One accent color, spent deliberately, not decoratively
- Monospace as the visual signal for "this is literal and copyable"
- Motion confirms actions; it never performs for its own sake

## 2. Colors

Neutral-dominant surface with a single accent doing all the color work in the
system.

### Primary
- **Deep indigo/blue accent** (`[to be resolved during implementation]`): links,
  the install-command CTA, active search/filter state. Nowhere else.

### Neutral
- **Warm-tinted off-white** (`[to be resolved during implementation]`): base
  background.
- **Warm-tinted near-black** (`[to be resolved during implementation]`): primary
  text.
- **Mid-tone neutrals** (`[to be resolved during implementation]`): borders,
  dividers, secondary/metadata text, disabled states.

### Named Rules
**The One Accent Rule.** The indigo/blue accent appears on ≤10% of any given
screen. It marks exactly three things: a link, the install CTA, and the active
search state. If a fourth use case shows up, reconsider before reaching for the
accent again.

## 3. Typography

**Direction:** Display + mono pairing — `[font pairing to be chosen at
implementation]`.

**Character:** A clean technical sans carries headings, body copy, and
descriptions; monospace is reserved for anything the visitor might copy or that
represents a literal identifier — skill names, install commands, metadata
badges (version, license). The switch to mono is itself a signal: "this is
exact, paste it as-is."

### Hierarchy
- **Display**: page title / site heading only.
- **Headline**: skill name on detail pages.
- **Title**: card titles in the grid, section headers.
- **Body**: descriptions, rendered `SKILL.md` body content. Cap line length at
  65–75ch for the rendered body.
- **Label**: metadata badges, filter chips, table headers — likely uppercase,
  small, wide letter-spacing.
- **Mono**: skill directory names, install commands, zip filenames, frontmatter
  values in the metadata table.

### Named Rules
**The Mono-Means-Copyable Rule.** If a string on the page is something a
developer would paste into a terminal or a config file, it renders in
monospace. Nothing else does.

## 4. Elevation

Flat by default. Depth is not a resting-state decoration here; it shows up only
as direct feedback to interaction — a card lifting slightly on hover, a focus
ring appearing on keyboard navigation. No ambient shadows sit under static
content.

### Named Rules
**The Feedback-Only Rule.** Shadows and elevation changes are a response to
state (hover, focus, active), never a permanent property of a resting element.

## 6. Do's and Don'ts

### Do:
- **Do** keep the accent color to links, the install CTA, and active search
  state — nowhere else.
- **Do** render skill names, install commands, and any copyable value in
  monospace.
- **Do** show real metadata (name, description, key frontmatter fields) on the
  card itself, not behind a click.
- **Do** respect `prefers-reduced-motion`: disable hover lift and crossfade
  transitions when requested.
- **Do** keep every interactive element keyboard-reachable with a visible focus
  state (WCAG 2.1 AA baseline, per PRODUCT.md).

### Don't:
- **Don't** build a "generic SaaS marketing template" — no hero-metric blocks,
  no gradient-text headlines, no stock-illustration energy, no "trusted by"
  logo strips (per PRODUCT.md's anti-references).
- **Don't** build "enterprise admin-dashboard bloat" — no heavy nested sidebar
  navigation, no dense data-table chrome, no dashboard-style stat tiles (per
  PRODUCT.md's anti-references).
- **Don't** reach for the "generic AI-slop aesthetic" — no default
  purple-gradient-on-dark, no glassmorphism, no side-stripe accent borders on
  cards (per PRODUCT.md's anti-references).
- **Don't** build a Vercel/Linear-style animated SaaS landing page: no
  scroll-triggered reveals, no oversized gradient hero typography, no
  choreographed entrance sequences.
- **Don't** use more than one accent hue, or let the accent creep past ~10% of
  any screen's surface.
