# SETUP

## Prerequisites

- Go (see `src/go.mod` for target version)
- Docker + Docker Compose
- Migration/sql tools used by the Makefile workflow (`migrate`, `sqlc`, `swagger`)

## Local Environment

1. Create env file:

```bash
cp .env.example .env
```

2. Start dependencies:

```bash
docker compose up -d
```

3. Generate API models and SQL code, then run server:

```bash
make run
```

> `make run` executes swagger generation + sqlc generation before starting `cmd/server/main.go`.

## Migrations

Run migrations:

```bash
make migrate
```

Create new migration files:

```bash
make migration-shared title=<name>
make migration-tenant title=<name>
```
