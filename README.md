# golang-web-api

A **multi-tenant REST API** built in Go, following **Clean Architecture** and **Domain-Driven Design (DDD)** principles.

---

## Tech Stack

| Concern | Technology |
|---|---|
| Language | Go 1.26 |
| HTTP Router | [chi](https://github.com/go-chi/chi) |
| Database | PostgreSQL (`pgx` / `database/sql`) |
| Cache / Token store | Redis |
| DB query generation | [SQLC](https://sqlc.dev) |
| API contract | OpenAPI 2.0 (Swagger) |
| Message broker | Apache Pulsar _(optional)_ |
| File storage | AWS S3 |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) |

---

## Quick Start

### Prerequisites

| Tool | Version |
|---|---|
| [Go](https://go.dev/dl/) | ≥ 1.26 |
| [Docker](https://docs.docker.com/get-docker/) | any recent |
| `make` | any |
| `psql` | any |

### 1. Clone & configure

```bash
git clone <repository-url>
cd golang-web-api
cp .env.example .env   # fill in POSTGRES_* and REDIS_* at minimum
```

### 2. Start infrastructure

```bash
docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres
docker run --name redis -p 6379:6379 -d redis
```

### 3. Install CLI tools

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install gotest.tools/gotestsum@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```

### 4. Run test

```bash
make test
```

### 5. Run migrations & start the server

```bash
make migrate
make run        # generates Swagger DTOs + SQLC code, then starts the server
```



API is available at **http://localhost:8000** · Swagger UI at **/static/swagger/**.

> See **[docs/INSTALL.md](docs/INSTALL.md)** for the full setup guide including Docker production build.

---

## Makefile Reference

| Target | Description |
|---|---|
| `make run` | Generate code (swagger + sqlc) and start the server |
| `make swag` | Regenerate Go DTOs from the OpenAPI spec |
| `make sqlc` | Regenerate DB query code from `src/query.sql` |
| `make migrate` | Apply all pending database migrations |
| `make test` | Run the full black-box API test suite |
| `make migration-shared title=<name>` | Create a new shared-schema migration file |
| `make migration-tenant title=<name>` | Create a new tenant-schema migration file |

---

## Architecture

Dependencies flow **inward only** — outer layers depend on inner layers, never the reverse.

```
┌──────────────────────────────────────────────┐
│              Infrastructure                  │  HTTP, DB, Redis, AWS, SQLC
│   ┌──────────────────────────────────────┐   │
│   │            Application               │   │  Commands + Queries (use-cases)
│   │   ┌──────────────────────────────┐   │   │
│   │   │           Domain             │   │   │  Entities, VOs, repo interfaces
│   │   └──────────────────────────────┘   │   │
│   └──────────────────────────────────────┘   │
└──────────────────────────────────────────────┘
         pkg/ (shared utilities — any layer)
```

| Layer | Location | Responsibility |
|---|---|---|
| Domain | `src/domain/` | Entities, Value Objects, repository interfaces, domain errors. Zero external deps. |
| Application | `src/application/` | Commands (writes) and Queries (reads) — orchestrate use-cases. |
| Infrastructure | `src/infrastructure/` | HTTP handlers, DB/Redis implementations, SQLC, migrations. |
| Shared | `src/pkg/` | Framework-agnostic utilities importable by any layer. |

> See **[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)** for the full layer breakdown, request lifecycle, and key interfaces.

---

## Multi-Tenancy

Each tenant has its own PostgreSQL schema inside a single database instance. The `X-Tenant` HTTP header identifies the tenant on every request. The `MultiTenancy` middleware validates the schema and injects a scoped DB connection — handlers and commands are completely unaware of tenancy.

---

## Project Structure

```
.
├── Dockerfile
├── Makefile
├── docker-compose.yml
├── entrypoint.sh
└── src/
    ├── cmd/            # Entrypoints: server, migrate, consumer, producer
    ├── config/         # Env → Config struct
    ├── domain/         # Entities, Value Objects, repo interfaces, exceptions
    ├── application/    # Commands (writes) and Queries (reads)
    ├── infrastructure/ # HTTP handlers, storage, SQLC, migrations
    ├── pkg/            # Shared utilities (postgres, redis, aws, logger, …)
    ├── openapi/        # OpenAPI 2.0 YAML specs (one file per endpoint)
    ├── static/swagger/ # Compiled swagger.json
    └── tests/          # Black-box API test suite
```

---

## Testing

Tests are black-box HTTP tests that exercise the full stack through the public API.

```bash
make test
```

Test suites: `test_auth`, `test_role`, `test_user`, `test_permissions`.

> See **[docs/TESTING.md](docs/TESTING.md)** for conventions, fixture rules, and examples.

---

## API Contract

OpenAPI specs live in `src/openapi/`. After editing any spec file, regenerate the Go DTOs:

```bash
make swag
```

> See **[docs/OPENAPI.md](docs/OPENAPI.md)** for the full workflow.

---

## Database Migrations

```bash
make migration-shared title=<name>   # shared (public) schema
make migration-tenant title=<name>   # per-tenant schema
make migrate                         # apply all pending migrations
```

> See **[docs/DATABASE_MIGRATION.md](docs/DATABASE_MIGRATION.md)** for migration rules and the recommended new-table checklist.

---

## Contributing

Branch naming: `feat/<topic>`, `fix/<topic>`, `docs/<topic>`, `chore/<topic>`

Commit style (conventional commits): `feat: ...`, `fix: ...`, `refactor: ...`, `test: ...`

> See **[docs/COMMIT.md](docs/COMMIT.md)** for PR conventions.

---

## Documentation Index

| Doc | Contents |
|---|---|
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | System design, layers, request flow, key interfaces |
| [docs/INSTALL.md](docs/INSTALL.md) | Local setup, environment variables, Docker build |
| [docs/TESTING.md](docs/TESTING.md) | Test strategy, fixture rules, examples |
| [docs/OPENAPI.md](docs/OPENAPI.md) | API contract workflow |
| [docs/DATABASE_MIGRATION.md](docs/DATABASE_MIGRATION.md) | Schema rules, migration commands |
| [docs/DDD.md](docs/DDD.md) | Where to place logic (Value Object → Entity → Domain Service → Use-case) |
| [docs/COMMIT.md](docs/COMMIT.md) | Branch naming and commit conventions |
