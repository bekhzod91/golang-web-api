# Installation & Bootstrap Guide

This guide walks you through setting up the project for local development from scratch.

---

## Prerequisites

Ensure the following tools are installed before proceeding:

| Tool | Version | Notes |
|------|---------|-------|
| [Go](https://go.dev/dl/) | ≥ 1.26 | Required to build and run the server |
| [Docker](https://docs.docker.com/get-docker/) | any recent | Runs PostgreSQL and Redis locally |
| `make` | any | Used for all dev workflow commands |
| `psql` | any | PostgreSQL CLI client (for schema setup) |

> **macOS:** Install `psql` via `brew install libpq && brew link libpq`.

---

## 1. Clone the Repository

```bash
git clone <repository-url>
cd golang-web-api
```

---

## 2. Start Infrastructure

Start PostgreSQL and Redis as Docker containers.

**PostgreSQL**
```bash
docker run --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  -d postgres
```

**Redis**
```bash
docker run --name redis \
  -p 6379:6379 \
  -d redis
```

> To stop and remove the containers later:
> ```bash
> docker rm -f postgres redis
> ```

---

## 3. Install CLI Tools

Install the required code-generation and testing tools:

```bash
# SQL query code generator
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# OpenAPI / Swagger model generator
go install github.com/go-swagger/go-swagger/cmd/swagger@latest

# Database migration runner
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Improved test runner (used by `make test`)
go install gotest.tools/gotestsum@latest
```

Make sure `$(go env GOPATH)/bin` is on your `$PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

---

## 4. Configure Environment

Copy the example environment file and fill in your values:

```bash
cp .env.example .env
```

Below is the full list of supported variables with their defaults:

```dotenv
# Application
ENVIRONMENT=development        # development | production
APP_HOST=0.0.0.0
APP_PORT=8000

# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DATABASE=postgres
POSTGRES_SCHEMA=public
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_TLS_ENABLED=false

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_USER=default
REDIS_PASSWORD=
REDIS_NAME=0                   # DB index (0–15)
REDIS_TLS_ENABLED=false

# AWS S3
AWS_BASE_ENDPOINT=
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
AWS_BUCKET=

# Logging
LOG_LEVEL=DEBUG                # DEBUG | INFO | WARN | ERROR
LOG_JSON=false

# Pagination
PAGINATION_LIMIT=10
```

> For a minimal local setup only `POSTGRES_*` and `REDIS_*` variables are required.

---

## 5. Run Database Migrations

Apply all shared and tenant schema migrations:

```bash
make migrate
```

Then create the **test schema** used by the test suite:

```bash
psql -U postgres -d postgres -c "CREATE SCHEMA test_tenant;"
```

> See [DATABASE_MIGRATION.md](DATABASE_MIGRATION.md) for migration authoring rules and commands.

---

## 6. Generate Code & Start the Server

`make run` performs the full generation pipeline and starts the server in one command:

```bash
make run
```

This is equivalent to:
```bash
make swag    # Generates Go DTOs from src/openapi/_swagger.yaml
make sqlc    # Generates type-safe DB query code from src/query.sql
cd src && go run cmd/server/main.go
```

The API server will be available at: **http://localhost:8000**

Swagger UI (if enabled) is served from the `/static/swagger/` path.

---

## 7. Run Tests

```bash
make test
```

This runs the full black-box API test suite across all modules:
- `tests/test_auth/`
- `tests/test_role/`
- `tests/test_user/`
- `tests/test_permissions/`

> See [TESTING.md](TESTING.md) for test conventions and how to write new tests.

---

## Makefile Targets — Quick Reference

| Target | Description |
|--------|-------------|
| `make run` | Generate code (swagger + sqlc) and start the server |
| `make swag` | Regenerate Go DTOs from the OpenAPI spec |
| `make sqlc` | Regenerate DB query code from `src/query.sql` |
| `make migrate` | Apply all pending database migrations |
| `make test` | Run the full API test suite |
| `make migration-shared title=<name>` | Create a new shared-schema migration file |
| `make migration-tenant title=<name>` | Create a new tenant-schema migration file |

---

## Docker Build (Production)

Build the production Docker image:

```bash
docker build -t golang-web-api .
```

Run the container (requires an external `app` network and a config file):

```bash
docker run --rm \
  --env-file .env \
  -p 8000:8000 \
  golang-web-api
```

The `entrypoint.sh` supports the following modes:

| Command | Behaviour |
|---------|-----------|
| _(no argument)_ | Run migrations **then** start the server |
| `server` | Start the server only |
| `migrate` | Run migrations only |

```bash
# Example: run migrations only
docker run --rm --env-file .env golang-web-api migrate
```
