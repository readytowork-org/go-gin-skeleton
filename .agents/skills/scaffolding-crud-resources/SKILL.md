---
name: scaffolding-crud-resources
description: Scaffolds a new CRUD resource (repository, service, controller, routes, fx module) in this Go Gin skeleton and wires it into the app. Use when adding a new resource, entity, or API module under api/.
---

# Scaffolding a CRUD resource

Full architecture and conventions live in `AGENTS.md` at the repo root. Read
it first if this session has not already loaded it.

Copy this checklist and check items off as you go:

```
- [ ] 1. Run `make crud` for the target folder
- [ ] 2. Fill in model, repository, service, controller, routes
- [ ] 3. Register the fx.Module in its parent module
- [ ] 4. Add a migration + `make dao`, if a new table is needed
- [ ] 5. Run `make swagger`
- [ ] 6. Add tests
- [ ] 7. Run the verifying-changes-before-commit skill
```

**Step 1: Run `make crud`**

Answer its prompt with the target folder, grouped by area under `api/`
(for example `api/product/product` for a user-facing resource, or
`api/admin/product` for an admin-only one). This generates
`repository.go`, `service.go`, `controller.go`, `routes.go`, and
`modules.go` from `automate/automate-templates/*.txt`.

**Step 2: Fill in the generated files**

- `model.go` / `dto.go`: GORM model and request/response DTOs.
- `repository.go`: implement `IRepository` using `config.Database`.
- `service.go`: implement `IService`, depending on `IRepository`.
- `controller.go`: gin handlers on `Controller`, with full swag annotations
  (`@Tags`, `@Summary`, `@Description`, `@Security`, `@Produce`, `@Success`,
  `@Failure`, `@Router`, `@Id`), following `api/user/user/controller.go`.
  Use one consistent receiver name across every method on the struct.
- `routes.go`: register real routes via `router.V1...`.

**Step 3: Register the module**

Add the new `fx.Module` to the `fx.Options(...)` list in the parent module,
for example `api/module.go` or `api/admin/module.go`.

**Step 4: Migration, if the resource needs a table**

Write an Atlas migration under `database/migration/`, then run `make dao` to
regenerate `database/dao/*.gen.go`. See the creating-database-migrations
skill for the full flow. Never hand-edit generated DAO files.

**Step 5: Regenerate swagger**

Run `make swagger` to refresh `swagger/docs.go`, `swagger.json`, and
`swagger.yaml`. Never hand-edit those files.

**Step 6: Add tests**

Add tests next to the code (for example `service_test.go`) using a
hand-written mock of `IRepository` with overridable `*Fn` fields, in the
style of `api/admin/user/repository_mock.go`.

**Step 7: Verify**

Use the verifying-changes-before-commit skill (lint, build, test, and the
generated-file checks) before treating the resource as done.
