## [0.2.0] - 2026-09-26

### Features

- Add `--base-url` flag to serve the generated site under a subpath

## [0.1.2] - 2026-09-25

### Miscellaneous Tasks

- Switch Dockerfile runtime image to debian-slim, add explicit non-root user
- Point examples/.gitlab-ci.yml at the project's own published image instead of building in a golang container

## [0.1.1] - 2026-09-16

### Features

- Rebrand to Skill Indexer

### Bug Fixes

- Fix image reference: endekastore -> endekasoft

### Documentation

- Update repo URLs to the real renamed GitLab project
- Strip development-diary narration from docs and AGENTS.md
- Drop AGENTS.md cross-reference from end of README

### Miscellaneous Tasks

- Add MIT license
- Reduce example skills number

## [0.1.0] - 2026-09-15

### Features

- Add Go scaffold, frontmatter parser/scanner, and zip downloads
- Add site rendering, search index, and end-to-end CLI wiring
- Replace detail pages with a slide-in panel, add Nord dark mode, drop body render
- Rebrand to Skillstore; panel width, layout, scope toggle, footer
- Put npx install command above the download button in the panel
- Add Dockerfile, examples/ (gitlab-ci + README), move fixtures under examples/
- Add Docker Compose and Kubernetes deployment examples
- Replace the custom installer/ package with npx skills add

### Documentation

- Add README and AGENTS.md
- Add full project documentation (docs/), close 14 gaps from cold review

### Miscellaneous Tasks

- Initial commit: skills fixtures, project plan, design context
- Integrate to ci/cd
