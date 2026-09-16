---
name: verifying-changes-before-commit
description: Runs this Go Gin skeleton's lint, build, test, and generated-file checks before a change is considered done. Use before finishing any code change, opening a PR, or when asked to verify or double-check a change.
---

# Verifying changes before commit

Copy this checklist and check items off as you go:

```
- [ ] 1. golangci-lint run --fast
- [ ] 2. go build ./...
- [ ] 3. go test ./... (for touched packages, at least)
- [ ] 4. make swagger, if controller annotations or DTOs changed
- [ ] 5. make dao, if the DB schema changed
- [ ] 6. Docs updated, if user-facing behavior changed
```

**Step 1-3: Lint, build, test**

```sh
golangci-lint run --fast
go build ./...
go test ./...
```

`golangci-lint run --fast` matches what the repo's pre-commit hook
(`hooks/pre-commit`) runs; use the full `golangci-lint run` (no `--fast`) for
a more thorough pass if time allows.

**Step 4: Swagger**

If any controller's swag comments or request/response DTOs changed, run
`make swagger` and include the resulting changes to `swagger/docs.go`,
`swagger.json`, and `swagger.yaml`. These files are generated; do not hand-
edit them directly.

**Step 5: Generated DAO structs**

If the database schema changed, run `make dao` and include the resulting
changes to `database/dao/*.gen.go`. Do not hand-edit those files.

**Step 6: Documentation**

If the change alters documented, user-facing behavior, update `README.md`
or `SWAGGER.md` to match. See the "Documentation map" in `AGENTS.md` for
which file covers what.
