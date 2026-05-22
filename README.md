# Go Gin Skeleton

> **Go Gin Skeleton template includes all the common packages and setup used for API development
> using [gin](https://gin-gonic.com).**

## Configured With

- Dependency Injection: [fx](https://github.com/uber-go/fx)
- Routing: [gin web framework](https://gin-gonic.com)
- Logging: [zap](https://github.com/uber-go/zap)
- Database: ([mysql](https://gorm.io/driver/mysql) / [sqlmock](https://github.com/DATA-DOG/go-sqlmock))
- ORM: [gorm](https://gorm.io/docs)
- API documentation: [gin-swagger](https://github.com/swaggo/gin-swagger)
- Middlewares
  - CORS (env-driven allowlist via `CORS_ALLOWED_ORIGINS`)
  - Rate Limit (per-user when authenticated, per-IP otherwise)
  - DB Transaction (scoped to write methods only)
  - Request ID (`X-Request-ID`, surfaces in error envelope `trace_id`)
  - Structured error envelope + central handler (`lib/api_errors`)
  - Idempotency-Key (replays POST responses; MySQL or Redis backend, see below)
- CLI tools
    - [atlas](https://atlasgo.io/): for DB migrations
  - [gentool](https://gorm.io/gen/): to generate dao objects from database
  - [swag](https://github.com/swaggo/swag): to generate swagger docs
  - [air](https://github.com/air-verse/air): hot-reload

**For Debugging 🐞** Debugger runs at `5002`. Vs code configuration is at `.vscode/launch.json` which will attach
debugger to remote application.

## Running

- Copy `.env.example` to `.env` and update according to requirement.
- ### using Docker
  - run `docker-compose up` (with default configuration will run at `5000` and adminer runs at `5001`)
- ### using Gin Watch
  - ### make sure to add go/bin path to your .bashrc/.zshrc file
    - run `make run`

## Commands 🛳

| Command               | Desc                                                                         |
|-----------------------|------------------------------------------------------------------------------|
| `make install`        | installs goalngci-lint and change the hooks config                           |
| `make run`            | runs the project using gin watcher                                           |
| `make migrate`        | runs [atlas](https://atlasgo.io/) migrate command with env configs from .env |
| `make crud`           | Create crud template                                                         |
| `make swagger`        | Run this command to generate swag docs                                       |
| `make dao`            | Generates go structs from database                                           |
| `make context-upload` | Upload Context to Circle CI Directly                                         |

[//]: # "TODO :: Need a proper name ⬇️"

## External Services

> Run `go get <package name>` to install.

#### Firebase

- Package name: github.com/readytowork-org/go_firebase_service
- Github: https://github.com/readytowork-org/go_firebase_service

#### GCP

- Package name: github.com/readytowork-org/go_gcp_service
- Github: https://github.com/readytowork-org/go_gcp_service

## Swagger docs config

> Please refer to [SWAGGER.md](https://github.com/readytowork-org/go-gin-skeleton/blob/develop/SWAGGER.md)

## Run CLI 🖥

The CLI mode is gated by `os.Args[1] == "cli"` (see `lib/utils/isCli.go`). Any
binary built from this project enters interactive CLI mode when launched with
that argument; otherwise it serves HTTP.

### Local

```sh
go run . cli
```

### Inside Docker

```sh
docker-compose exec web sh
./tmp/main cli       # if built via `make run` (air)
# or
go run . cli
```

You'll get an interactive menu (powered by `promptui`) listing available
commands. Currently:

- `CREATE_SEED_DATA` — runs the seed registry defined in `database/seeds/`.
- `EXIT_APPLICATION` — quit.

### Adding a new CLI command

1. Implement the `cli.Command` interface (`Run()` + `Name()`) in a file under
   `cli/`.
2. Provide it via `fx.Provide(NewYourCommand)` in `cli/cli.go`.
3. Add it to the `commands` slice in `NewApplication` (`cli/app.go`).

### Seeds (non-interactive)

Seeds in `database/seeds/` are also run automatically at HTTP server startup
(see `bootstrap/bootstrap.go`). Add a new seed by implementing `seeds.Seed`
and registering it in the `seeds.Module`.

## Hot reload

`make run` uses [air](https://github.com/air-verse/air); config is in
`.air.toml`. The `tmp/` directory holds the dev binary and is gitignored.

## Idempotency

POST endpoints can be made safely retryable by mounting
`middlewares.IdempotencyMiddleware.Handle()` (see
`api/admin/user/route.go` for an example). When a client sends the
`Idempotency-Key` header, the first response is cached for `IDEMPOTENCY_TTL`
and replayed for any subsequent request with the same key — duplicates are
flagged with the response header `Idempotent-Replay: true`.

Backend is chosen by `IDEMPOTENCY_STORE`:

| Value     | Backend                                  | Notes                                  |
|-----------|------------------------------------------|----------------------------------------|
| (empty)   | Noop                                     | Disabled — middleware passes through.  |
| `mysql`   | `idempotency_keys` table (existing DB)   | Durable; survives restarts.            |
| `redis`   | `REDIS_ADDR`                             | Fast; native TTL. Compose includes a redis service. |

The MySQL backend uses the `idempotency_keys` table created by
`database/schema.sql` — run `make migrate` after switching to `mysql`.

## Auth

- `POST /api/v1/login` — issues an access/refresh token pair and persists the
  refresh token (hashed) in `refresh_tokens`.
- `POST /api/v1/login/refresh` — rotates the refresh token; reuse of a
  revoked refresh token revokes all sessions for the user.
- `POST /api/v1/logout` — revokes the supplied refresh token. Always returns
  `200` to avoid leaking token validity.

## Health probes

| Path            | Behaviour                                       |
|-----------------|-------------------------------------------------|
| `/livez`        | Process is up — always `200`.                   |
| `/health-check` | Pings the DB with a 2 s timeout.                |
| `/readyz`       | Alias for `/health-check` for k8s-style probes. |

## Implements Google Cloud Proxy by default

This reduces hassle for developer to update IP in Cloud SQL during IP Change.

Implemented through docker image as follows

- Cloud SQL -> Google Proxy Docker Image -> Web App

#### Points to remember for smooth working

- `ServiceAccountKey.json` requires Cloud SQL Read and Write Permission
- `DB_HOST_NAME` value in `.env` is required
- `DB_HOST=cloud-sql-proxy` instead of `IPV4` or `DB_HOST_NAME` for development environment
- `DB_PORT` will be `3306` by default

## For auto generate of CRUD(Create, ReaD, Update & Delete) api following informations are needed and will be asked in terminal:

- resource-name: name of CRUD in upper camelCase. examples:Food,Puppy,ProductCategory etc.

- resource-table-name: name of CRUD in lower snake case. examples:food,puppy,product_category etc.

- plural-resource-table-name: plural name for the table going to be created. example: foods, puppies,
  product_categories.

- plural-resource-name: plural name of CRUD in Upper camelCase. examples:Foods,Puppies,ProductCategories etc.
