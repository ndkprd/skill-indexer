---
title: Contributing
description: How to make and submit a change to Skillstore.
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

   If you touched `installer/`, also run `cd installer && bun test`.

4. Commit with a message that explains _why_, not just _what_ — see this
   repo's own `git log` for examples of the expected level of detail.
5. Open a merge request against `main`. Nothing merges to `main` without
   explicit review approval, regardless of whether checks pass.

## Release process

There is no tagged release process yet. Two version strings currently
exist, hand-edited independently, that happen to match by coincidence
(both `"0.1.0"`) rather than by anything keeping them in sync:

- `internal/site.Version` in `internal/site/render.go`, shown in every
  generated site's footer.
- `"version"` in `installer/package.json`, the separate
  `skillstore-install` npm package's own version.

If you're introducing versioned releases, update both, and consider
whether they should be tied together going forward.

## Related

- [Standards](./standards.md) — the code style and git conventions referenced above
- [Testing](./testing.md) — what each check in step 3 actually verifies
