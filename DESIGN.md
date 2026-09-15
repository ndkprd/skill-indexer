---
name: Skillstore
description: A static marketplace for Claude Code skills — browse, search, install.
colors:
  bg-light: "oklch(98.3% 0.004 265)"
  surface-light: "oklch(99.3% 0.003 265)"
  surface-sunken-light: "oklch(96% 0.006 265)"
  border-light: "oklch(87% 0.008 265)"
  border-strong-light: "oklch(78% 0.012 265)"
  text-light: "oklch(22% 0.014 265)"
  text-secondary-light: "oklch(46% 0.012 265)"
  text-tertiary-light: "oklch(60% 0.01 265)"
  accent-light: "oklch(47% 0.15 265)"
  accent-hover-light: "oklch(40% 0.16 265)"
  success-light: "oklch(55% 0.14 150)"
  bg-dark: "#2e3440"
  surface-dark: "#3b4252"
  surface-sunken-dark: "#2e3440"
  border-dark: "#434c5e"
  border-strong-dark: "#4c566a"
  text-dark: "#eceff4"
  text-secondary-dark: "#d8dee9"
  accent-dark: "#88c0d0"
  accent-hover-dark: "#8fbcbb"
  success-dark: "#a3be8c"
typography:
  display:
    fontFamily: "IBM Plex Sans, -apple-system, BlinkMacSystemFont, Segoe UI, Roboto, sans-serif"
    fontSize: "1.75rem"
    fontWeight: 600
    lineHeight: 1.3
  body:
    fontFamily: "IBM Plex Sans, -apple-system, BlinkMacSystemFont, Segoe UI, Roboto, sans-serif"
    fontSize: "1rem"
    fontWeight: 400
    lineHeight: 1.5
  label:
    fontFamily: "IBM Plex Sans, -apple-system, BlinkMacSystemFont, Segoe UI, Roboto, sans-serif"
    fontSize: "0.75rem"
    fontWeight: 500
    letterSpacing: "0.06em"
  mono:
    fontFamily: "IBM Plex Mono, ui-monospace, SFMono-Regular, SF Mono, Menlo, Consolas, Liberation Mono, monospace"
    fontSize: "0.875rem"
    fontWeight: 400
rounded:
  sm: "6px"
  md: "10px"
  pill: "999px"
spacing:
  1: "4px"
  2: "8px"
  3: "12px"
  4: "16px"
  6: "24px"
  8: "32px"
  12: "48px"
  16: "64px"
components:
  card:
    backgroundColor: "{colors.surface-light}"
    rounded: "{rounded.md}"
    padding: "{spacing.4}"
  badge:
    backgroundColor: "{colors.surface-sunken-light}"
    textColor: "{colors.text-secondary-light}"
    rounded: "{rounded.pill}"
  button-primary:
    backgroundColor: "{colors.accent-light}"
    textColor: "{colors.bg-light}"
    rounded: "{rounded.sm}"
    padding: "8px 16px"
  button-primary-hover:
    backgroundColor: "{colors.accent-hover-light}"
---

# Design System: Skillstore

## 1. Overview

**Creative North Star: "The Registry Desk"**

Dense, plain, and fast: a page that behaves like developer infrastructure, not a
pitch. The reference lineage is npm's package cards, VS Code Marketplace's
extension grid, and GitHub's file browser — surfaces that put real metadata in
front of you immediately and never ask you to scroll past a hero section to find
it. Color stays restrained (tinted neutrals, one small accent used sparingly) so
the accent still means something when it appears. Motion is present but modest:
cards respond to hover, search results crossfade, a copy button confirms itself,
a detail panel slides in from the right — nothing choreographed, nothing that
makes the visitor wait to act.

This system explicitly rejects the Vercel/Linear-style animated SaaS landing
page: no scroll-triggered reveals, no gradient hero text, no oversized marketing
typography standing between the visitor and the card grid. It also rejects
showing more than a visitor asked for: the card grid and its detail panel show
name, description, and metadata — never the full skill body. The frontmatter
`description` is treated as sufficient; if it isn't, that's a content problem to
fix upstream, not something this UI should compensate for with a wall of
rendered Markdown.

**Key Characteristics:**
- Information visible at a glance, not behind a click or a hover
- One accent color, spent deliberately, not decoratively
- Monospace as the visual signal for "this is literal and copyable"
- Motion confirms actions; it never performs for its own sake
- Detail is a panel, not a destination — opening a skill never leaves the grid

## 2. Colors

Neutral-dominant surface with a single accent doing all the color work in the
system. Two committed themes, not a token-by-token dark mode afterthought:
**light** (the restrained indigo system below) and **dark** (Nord — Polar
Night surfaces, Snow Storm text, Frost accent). A header toggle switches
between them explicitly; absent a stored choice, the OS's
`prefers-color-scheme` decides.

### Primary
- **Deep indigo/blue accent** (light: `oklch(47% 0.15 265)`; dark: Nord Frost
  `#88c0d0`): links, the install-command CTA, active search/filter state.
  Nowhere else.

### Neutral — Light
- **Warm-tinted off-white** (`oklch(98.3% 0.004 265)`): base background.
- **Warm-tinted near-black** (`oklch(22% 0.014 265)`): primary text.
- **Mid-tone neutrals** (`oklch(46–87% 0.008–0.012 265)`): borders, dividers,
  secondary/metadata text.

### Neutral — Dark (Nord)
- **Polar Night** (`#2e3440` bg, `#3b4252` surface, `#434c5e` /
  `#4c566a` borders): the same elevation logic as light — sunken elements
  (code blocks, the install command) drop to `#2e3440`, the same value as the
  page background, so they read as recessed beneath whatever surface holds
  them.
- **Snow Storm** (`#eceff4` primary text, `#d8dee9` secondary text): light
  text on dark needs to stay light — Nord's own foreground family, not a
  darkened neutral, which is why dark-mode text doesn't mirror light mode's
  mid-tone-gray approach.

### Named Rules
**The One Accent Rule.** The accent appears on ≤10% of any given screen. It
marks exactly three things: a link, the install CTA, and the active search
state. If a fourth use case shows up, reconsider before reaching for the
accent again.

**The Sunken-Equals-Background Rule (dark).** In dark mode, a "sunken"
surface (code blocks, the install command) is the *same* color as the page
background, not a separately-invented darker shade — depth reads from what's
lighter (the surface around it), not from a bespoke deeper black.

## 3. Typography

**Pairing:** IBM Plex Sans (display, body, labels) + IBM Plex Mono (anything
copyable). A superfamily pairing on purpose — same technical heritage, no
visual tension between the two, loaded from Google Fonts with a system-font
fallback stack so the page still reads correctly if the request fails.

**Character:** A clean technical sans carries headings, body copy, and
descriptions; monospace is reserved for anything the visitor might copy or that
represents a literal identifier — skill names, install commands, metadata
values. The switch to mono is itself a signal: "this is exact, paste it as-is."

### Hierarchy
- **Display** (600, 1.75rem, 1.3): the panel's skill-name headline only.
- **Title** (500, 1rem, mono): card titles in the grid.
- **Body** (400, 1rem, 1.5): descriptions, in cards and in the panel.
- **Label** (500, 0.75rem, uppercase, 0.06em tracking): metadata table
  headers, install-command label.
- **Mono** (400, 0.875rem or 0.75rem): skill directory names, install
  commands, zip filenames, every metadata table value.

### Named Rules
**The Mono-Means-Copyable Rule.** If a string on the page is something a
developer would paste into a terminal or a config file, it renders in
monospace. Nothing else does.

## 4. Elevation

Flat by default in both themes. Depth is not a resting-state decoration here;
it shows up only as direct feedback to interaction — a card lifting slightly
on hover, a focus ring appearing on keyboard navigation, the detail panel
sliding in with a real drop shadow because it is a genuine overlay above the
page. No ambient shadows sit under static content.

### Shadow Vocabulary
- **Hover shadow** (`0 2px 10px` accent-tinted at low alpha, light; `0 2px
  10px rgb(0 0 0 / 0.3)`, dark): cards and other resting elements, only on
  hover/focus.
- **Panel shadow** (`-8px 0 32px` at higher alpha than the hover shadow): the
  one place a shadow exists without user interaction driving it, because the
  panel is a legitimate elevated overlay, not resting content.

### Named Rules
**The Feedback-Only Rule.** Shadows and elevation changes are a response to
state (hover, focus, active, or "this is an overlay above the page"), never a
permanent property of resting, in-flow content.

## 5. Components

### Cards
- **Shape:** 10px radius, 1px border.
- **Default:** surface background, border in the neutral border color.
- **Hover / Focus:** border shifts to the accent color, lifts 2px
  (`translateY`), hover shadow appears. Never a shadow at rest.
- **Content:** name (mono, 500) + up to 2 metadata badges on one row,
  2-line-clamped description below, dirname (mono, tertiary) as a footer
  line. The whole card is a single link — no nested interactive elements.

### Badges
- **Style:** pill shape, sunken background, 1px border, mono type, tertiary
  text color. Up to 2 per card, chosen by priority (`version` → `author` →
  `license` → `compatibility` → first remaining key alphabetically); every
  metadata key is shown in the panel's full table, sorted alphabetically.

### Detail Panel
- **Shape:** fixed to the right edge of the viewport, full height, `min(440px,
  100vw)` wide (full-width takeover under ~768px). Slides in via `transform:
  translateX()`, backdrop fades via opacity — both respect
  `prefers-reduced-motion`.
- **Behavior:** opens on card click, closes on the close button, a backdrop
  click, Escape, or browser back (its open state lives in the URL hash, so
  it's a real, shareable, back-button-able navigation, not just a JS overlay
  toggle). Background content gets `inert` while the panel is open; focus is
  trapped inside it and returns to the triggering card on close.
- **Content:** name (display) + mono dirname, description, metadata table
  (omitted entirely when a skill has no extra metadata — never rendered
  empty), then the action area: a full-width primary "Download .zip" button
  and the install-command block. No rendered skill body, ever — see Overview.

### Buttons
- **Primary:** accent background, background-colored text, 6px radius. Used
  once per panel (Download).
- **Ghost:** transparent/bg background, bordered, used for the copy button.
- **Copy button states:** idle → "Copy"; after a successful clipboard write,
  "Copied" with the border/text color swapped to the success color for
  ~1.6s, then reverts.

### Theme Toggle
- **Style:** icon-only button in the header, sun/moon SVG (inline, no icon
  font), bordered like other secondary controls. Reflects the *effective*
  theme (explicit choice, or OS preference if none stored) at all times,
  including live updates if the OS preference changes mid-session.

## 6. Do's and Don'ts

### Do:
- **Do** keep the accent color to links, the install CTA, and active search
  state — nowhere else.
- **Do** render skill names, install commands, and any copyable value in
  monospace.
- **Do** show real metadata (name, description, key frontmatter fields) on the
  card itself, not behind a click.
- **Do** keep the detail panel's sunken surfaces the same color as the page
  background in dark mode — don't invent a separate darker shade.
- **Do** respect `prefers-reduced-motion`: disable hover lift, panel slide,
  and crossfade transitions when requested.
- **Do** keep every interactive element keyboard-reachable with a visible
  focus state, and trap focus inside the panel while it's open (WCAG 2.1 AA
  baseline, per PRODUCT.md).

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
- **Don't** render the skill's Markdown body anywhere. The description is the
  pitch; a detail panel is not a documentation reader.
- **Don't** use more than one accent hue per theme, or let the accent creep
  past ~10% of any screen's surface.
- **Don't** treat dark mode as an inverted light mode — Nord's own surface
  and text families, not a programmatic color-flip of the light tokens.
