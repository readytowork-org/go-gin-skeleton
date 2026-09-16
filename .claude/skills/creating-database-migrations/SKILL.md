---
name: creating-database-migrations
description: Creates and applies an Atlas database migration in this Go Gin skeleton, then regenerates the GORM DAO structs. Use when adding, changing, or dropping a table or column, or when asked to update the database schema.
---

# Creating a database migration

This project uses [Atlas](https://atlasgo.io/) for migrations
(`database/migration/`) and `gentool` to generate `database/dao/*.gen.go`
from the live schema. Never hand-edit generated DAO files or `atlas.sum`.

## Steps

1. Change `database/schema.sql` (or the relevant schema source) to reflect
   the desired end state of the table(s).

2. Generate and review a migration:

   ```sh
   make migrate diff
   ```

   Inspect the generated `.sql` file under `database/migration/` before
   applying it. Atlas names migrations by timestamp; do not rename or
   hand-edit an already-applied migration file.

3. Apply the migration:

   ```sh
   make migrate
   ```

   Use `env=local` when running against a host database instead of the
   `web` container, for example `make migrate env=local`.

4. Regenerate DAO structs from the now-updated schema:

   ```sh
   make dao
   ```

   This overwrites `database/dao/*.gen.go`. Do not hand-edit those files;
   if a change is needed, change the schema and regenerate instead.

5. Update any code that depends on the changed schema: repository queries,
   GORM models under `api/**/model.go`, and seed data in `database/seeds/`.

6. If the table backs a new resource, see the scaffolding-crud-resources
   skill for the rest of the wiring.

7. Run the verifying-changes-before-commit skill before treating the change
   as done.
