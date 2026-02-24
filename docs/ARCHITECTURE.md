# Architecture

This project is a **multi-tenant REST API** built in Go, following **Clean Architecture** and **Domain-Driven Design (DDD)** principles.

> For DDD-specific code placement rules (Value Objects, Entities, Domain Services, etc.), see [DDD.md](DDD.md).

---

## Tech Stack

| Concern | Technology |
|---------|-----------|
| Language | Go 1.23 |
| HTTP Router | [chi](https://github.com/go-chi/chi) |
| Database | PostgreSQL (`pgx` / `database/sql`) |
| Cache / Token store | Redis |
| DB query generation | [SQLC](https://sqlc.dev) |
| API contract | OpenAPI 2.0 (Swagger) |
| Message broker | Apache Pulsar _(optional)_ |
| File storage | AWS S3 |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) |

---

## Directory Structure

```
.
├── Dockerfile
├── Makefile
├── docker-compose.yml
├── entrypoint.sh
└── src/
    ├── cmd/
    │   ├── server/          # Entrypoint: starts the HTTP server
    │   ├── migrate/         # Entrypoint: runs database migrations
    │   ├── consumer/        # Entrypoint: Pulsar message consumer
    │   └── producer/        # Entrypoint: Pulsar message producer
    │
    ├── config/              # App configuration (env → Config struct)
    │
    ├── domain/              # ── DOMAIN LAYER (pure Go, zero framework deps) ──
    │   ├── entity/          # Aggregates with identity (User, Role, Tenant)
    │   ├── value_object/    # Immutable typed values (Email, Password, Token…)
    │   ├── repository/      # Repository interfaces — no implementations here
    │   └── exception/       # Domain error types (DomainError, NotFoundError)
    │
    ├── application/         # ── APPLICATION LAYER (use-cases) ──
    │   ├── command/         # Write: SignIn, CreateUser, DeleteRole, UploadFile…
    │   └── query/           # Read: GetMe, GetUsers, GetRoles, GetPermissions…
    │
    ├── infrastructure/      # ── INFRASTRUCTURE LAYER (I/O, HTTP, DB) ──
    │   ├── api/
    │   │   ├── dto/         # Request/response structs (Swagger → generated Go)
    │   │   ├── handler/     # HTTP handlers — thin, delegate to application layer
    │   │   ├── helper/      # Pagination & string conversion helpers
    │   │   └── router/      # Route registration (public + tenant-scoped groups)
    │   ├── core/            # App bootstrap, IMux, IContext, middlewares
    │   ├── migrations/
    │   │   ├── shared/      # Migrations for the public/shared schema
    │   │   └── tenant/      # Migrations applied to every tenant schema
    │   ├── sqlc/            # SQLC-generated type-safe query code
    │   └── storage/         # IStorage + concrete repository implementations
    │
    ├── pkg/                 # ── SHARED PACKAGES (any layer may import) ──
    │   ├── aws/             # AWS S3 client wrapper
    │   ├── env/             # .env loader (godotenv)
    │   ├── funcutils/       # Generic slice helpers
    │   ├── logger/          # Structured logger factory (slog / httplog)
    │   ├── multi_tenency/   # Tenant DB routing + X-Tenant middleware
    │   ├── postgres/        # *sql.DB connection pool factory
    │   ├── pulsar/          # Apache Pulsar client wrapper
    │   ├── randutil/        # Cryptographically-safe random string generation
    │   └── redis/           # *redis.Client factory
    │
    ├── openapi/             # OpenAPI 2.0 YAML specs (one file per endpoint)
    ├── static/swagger/      # Compiled swagger.json (served at /swagger/)
    ├── tests/               # Black-box API test suite
    │   ├── conf.go          # Test harness (NewTestApp, BindJSON helpers)
    │   ├── test_auth/
    │   ├── test_user/
    │   ├── test_role/
    │   ├── test_permissions/
    │   └── test_upload_file/
    ├── query.sql            # Raw SQL queries (SQLC input)
    └── sqlc.yaml            # SQLC configuration
```

---

## Architecture Layers

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

### Domain Layer — `src/domain/`

The innermost layer. **No external dependencies** (no DB drivers, no HTTP, no framework).

| Package | Contents |
|---------|---------|
| `entity/` | `User`, `Role`, `Tenant` — objects with identity and lifecycle. Own their own state invariants. |
| `value_object/` | `Email`, `Password`, `Token`, `Status`, `PhoneNumber`, `Date`, `Image`, `Rating`, `Tags`, … Each VO validates itself on construction and carries domain behaviour (e.g., `Password.VerifyPassword`). |
| `repository/` | `IUserRepository`, `IRoleRepository`, `ITokenRepository` — pure Go interfaces that define *what* persistence operations exist. Implementations live in `infrastructure/storage/`. |
| `exception/` | `DomainError`, `NotFoundError`, `exception.New()`, `IsDomainException()` — lets callers distinguish domain failures from infrastructure errors. |

### Application Layer — `src/application/`

Orchestrates use-cases. Depends **only** on the domain layer and the `core.IContext` interface.

- **Commands** (writes): `SignIn`, `SignUp`, `SignOut`, `CreateUser`, `UpdateUser`, `DeleteUser`, `ChangeUserPassword`, `CreateRole`, `UpdateRole`, `DeleteRole`, `UploadFile`
- **Queries** (reads): `GetMe`, `GetUsers`, `GetUserByID`, `GetRoles`, `GetRoleByID`, `GetPermissions`

Every function follows the same signature pattern:
```go
func CreateUser(ctx core.IContext, req dto.CreateUserRequestDTO) (*dto.CreateUserResponseDTO, error)
```

No framework imports. No direct DB calls. Business rules live here and in the domain.

### Infrastructure Layer — `src/infrastructure/`

Implements everything that touches I/O.

| Sub-package | Responsibility |
|-------------|---------------|
| `api/handler/` | HTTP handlers — parse request DTO, call command/query, render response. Deliberately thin. |
| `api/dto/` | Request/response structs auto-generated by `swagger generate model`. Never hand-edit. |
| `api/router/` | Two route groups: `PublicRoutes` (no auth) and `TenantRoutes` (requires `X-Tenant` + JWT). |
| `core/` | `App`, `IMux` (chi wrapper), `IContext`, `Authorization` middleware, `AppMiddleware`. |
| `storage/` | `IStorage` + concrete implementations using `database/sql` and Redis. |
| `sqlc/` | Type-safe query code generated by SQLC from `query.sql`. |
| `migrations/` | `*.up.sql` / `*.down.sql` files, split into `shared/` and `tenant/` directories. |

### Shared Packages — `src/pkg/`

Framework-agnostic packages importable by any layer.

| Package | Purpose |
|---------|---------|
| `multi_tenency` | `MultiTenancy` HTTP middleware + `DB` struct that manages per-tenant `search_path` connections |
| `postgres` | Creates a `*sql.DB` connection pool from config |
| `redis` | Creates a `*redis.Client` from config |
| `aws` | Thin S3 upload/download wrapper |
| `logger` | Structured logger factory (`slog` + `httplog`) |
| `env` | Loads `.env` via `godotenv` |
| `pulsar` | Apache Pulsar producer/consumer client |
| `funcutils` | Generic slice helpers (e.g., `Uniq`) |
| `randutil` | Cryptographically-safe random string generation |

---

## Request Lifecycle

Flow for a typical **authenticated, tenant-scoped** API call:

```
Client
  │  POST /api/v1/users/create/
  │  Authorization: Bearer <token>
  │  X-Tenant: <schema_name>
  ▼
chi Router
  ├─ AppMiddleware       — injects *App (config, DB pool, Redis, AWS) into context
  ├─ RequestLogger       — structured request/response logging
  ├─ RealIP             — resolves real client IP
  ├─ Recoverer          — catches panics, returns 500
  │
  │  [Tenant-scoped group]
  ├─ MultiTenancy        — reads X-Tenant header
  │                        validates schema exists in shared tenant table
  │                        opens (or reuses) *sql.DB with SET search_path = <schema>
  │                        injects *Tenant into context
  │
  │  [Private sub-group]
  └─ Authorization       — reads Bearer token from header
                           validates token via Redis
                           injects *entity.User into context
  ▼
Handler  (infrastructure/api/handler/)
  │  Parses & validates request DTO
  │  Calls application-layer command or query
  ▼
Command / Query  (application/command/ or application/query/)
  │  Constructs/validates domain Value Objects
  │  Calls domain entity methods
  │  Reads / writes via ctx.Storage()
  ▼
IStorage  (infrastructure/storage/)
  ├─ UserRepository   → tenant PostgreSQL (*sql.DB, schema-scoped)
  ├─ RoleRepository   → tenant PostgreSQL
  └─ TokenRepository  → Redis
```

---

## Multi-Tenancy

The API uses **PostgreSQL schema-based tenancy**: each tenant has its own schema inside a single database instance.

| Schema | Contents |
|--------|---------|
| `public` (shared) | `tenant` registry table, global configuration |
| `<tenant_name>` | Per-tenant: `user`, `role`, `permission` tables |

**Runtime flow:**

1. Client includes `X-Tenant: <schema_name>` on every request.
2. `MultiTenancy` middleware looks up the schema name in the shared `tenant` table.
3. Opens a `*sql.DB` connection and runs `SET search_path TO <schema_name>` — this connection is cached in a `map[string]*sql.DB` inside the `DB` struct for the lifetime of the process.
4. The tenant connection is injected into the request context; `ctx.Storage()` uses it transparently.
5. Handlers and commands are **completely unaware** of tenancy — they just call `ctx.Storage().User()`.

---

## Code Generation

### SQLC

`src/query.sql` contains raw SQL. Running `make sqlc` regenerates:

| Generated file | Contents |
|----------------|---------|
| `infrastructure/sqlc/models.go` | Go structs mirroring every DB table |
| `infrastructure/sqlc/query.sql.go` | Type-safe query functions |
| `infrastructure/sqlc/db.go` | `*Queries` struct + `DBTX` interface |

> Always run `make sqlc` after editing `query.sql` or migration files.

### OpenAPI / Swagger

Specs live in `src/openapi/` — one YAML file per endpoint. Running `make swag`:

1. Merges all YAMLs and validates the combined spec.
2. Generates Go DTO structs into `infrastructure/api/dto/` via `swagger generate model`.
3. Flattens the spec to `src/static/swagger/swagger.json` (served at `/swagger/`).

> Never hand-edit files in `infrastructure/api/dto/` — they are overwritten on every `make swag`.

---

## Key Interfaces

### `IContext` — `infrastructure/core/context.go`

Passed to every handler and application-layer function. Acts as the single seam between the HTTP framework and business logic.

```go
type IContext interface {
    Storage() storage.IStorage   // tenant-scoped repositories
    User()    *entity.User       // authenticated user (panics if not set)
    Logger()  logger.ILogger     // structured request logger
    AWS()     *aws.Client        // S3 client

    BindJSON(target any) error                           // parse raw JSON
    ShouldBindJSON(target interface{ Validate(...) }) error // parse + validate via Swagger

    URLParam(key string) string
    QueryParam(key string) string
    TenantSchemaName() string

    JSON(status int, v any)
    OK(v any) | Created(v any) | NoContent()
    BadRequest(err) | NotFound() | Unauthorized() | Forbidden()
}
```

### `IStorage` — `infrastructure/storage/storage.go`

```go
type IStorage interface {
    Token() repository.ITokenRepository
    User()  repository.IUserRepository
    Role()  repository.IRoleRepository
}
```

A new `IStorage` is constructed on every request, wiring the tenant-specific `*sql.DB` into each repository implementation.

### Repository interfaces — `domain/repository/`

Defined in the domain layer; implemented in `infrastructure/storage/`. This inversion keeps the domain free of any persistence details.

```go
type IUserRepository interface {
    GetUsers() ([]*entity.User, error)
    GetUserByID(id int64) (*entity.User, error)
    GetUserByEmail(email value_object.Email) (*entity.User, error)
    CreateUser(user *entity.User) (*entity.User, error)
    UpdateUser(user *entity.User) error
    ChangePasswordUser(user *entity.User) (*entity.User, error)
    DeleteUser(user *entity.User) error
}
```

