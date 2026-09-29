## [0.4.1] - 2026-09-29

### Bug Fixes

- *(site)* Bump footer version const to 0.4.1

## [0.4.0] - 2026-09-29

### Features

- *(site)* Add list view with tile/list icon toggle
- *(site)* Polish header, list rows, card badges and panel width

### Refactor

- *(site)* Slim single-row sticky header, scrolling install-everything block

### Miscellaneous Tasks

- Untrack built binary and gitignore it

## [0.3.0] - 2026-09-29

### Features

- *(cli)* Add --repository-url, --site-name, --extra-css and --extra-js flags
- *(site)* Repo-mode install, install-everything banner, site title, .skill download

### Documentation

- *(plan)* Add plan 002 for repo-url install, site name, theming
- Document repository-url, site-name, extra assets and .skill download
- Fix markdown lint in config table and AGENTS
- *(plan)* Mark plan 002 tasks done
- Remove internal references from plan 002, AGENTS and tests
- *(examples)* Add sample dracula and monokai themes
- Remove internal references from plan 001 and gitlab-cli example

### Testing

- *(site)* Cover repo URL, site name, extras and .skill link

### Miscellaneous Tasks

- Gitignore and dockerignore local themes dir

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
