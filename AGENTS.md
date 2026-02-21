# AGENTS.md

## Project Brief

`golang-web-api` is a Go backend boilerplate focused on:
- DDD + clean architecture layering
- multi-tenancy support
- JWT-based auth flows (sign up/sign in/sign out)
- RBAC-style roles/permissions APIs
- file upload support (S3-compatible)
- async/event integration via Pulsar

Core code lives in `src/` with clear boundaries:
- `domain/` entities, value objects, repository interfaces
- `application/` use cases (commands/queries)
- `infrastructure/` HTTP handlers, router, storage, migrations
- `cmd/` entrypoints (`server`, `migrate`, `producer`, `consumer`)
- `openapi/` endpoint specs
- `tests/` integration-style endpoint tests

## Quick Start

1. Copy env template:
   - `cp .env.example .env`
2. Start dependencies:
   - `docker compose up -d`
3. Run app/migrations using `Makefile` targets (see `Makefile`).

## Notes for Agents

- Keep architecture boundaries intact (`domain` should not depend on `infrastructure`).
- Prefer adding new API contracts under `src/openapi/` first, then handlers/use-cases.
- When changing DB models, include matching migration updates under:
  - `src/infrastructure/migrations/shared`
  - `src/infrastructure/migrations/tenant` (if tenant-scoped)
- Add/update tests in `src/tests/` for behavior changes.
