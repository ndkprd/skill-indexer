---
title: Contributing
description: How to make and submit a change to Skill Indexer.
tags: [developer-guide, contributing]
---

# Contributing

## Making a change

1. Branch from `main`: `<model-name>/<feature|fix|refactor>/<slug>` (e.g.
   `sonnet-5/fix/scan-symlinks`).
2. Make your change. If it touches `internal/site/templates/` or
   `internal/site/assets/`, check [DESIGN.md](../../DESIGN.md) and
   [Pitfalls](./pitfalls.md) first — several things that look like bugs
   are deliberate.
3. Run the full check before committing:

   ```bash
   go build ./... && go vet ./... && gofmt -l . && go test ./...
   ```

4. Commit with a message that explains _why_, not just _what_ — see this
   repo's own `git log` for examples of the expected level of detail.
5. Open a merge request against `main`. Nothing merges to `main` without
   explicit review approval, regardless of whether checks pass.

## Release process

There is no tagged release process yet — `internal/site.Version` in
`internal/site/render.go` (shown in every generated site's footer) is a
single hand-edited constant, currently `"0.1.0"`. If you're introducing
versioned releases, that's the one place the version string is currently
sourced from.

## Related

- [Standards](./standards.md) — the code style and git conventions referenced above
- [Testing](./testing.md) — what each check in step 3 actually verifies
