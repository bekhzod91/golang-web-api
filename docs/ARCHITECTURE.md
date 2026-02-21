# ARCHITECTURE

## Overview

This project follows DDD + clean architecture style with clear separation of concerns.

Main layers inside `src/`:

- `domain/` — core business entities, value objects, repository interfaces, domain exceptions
- `application/` — use-cases (commands/queries), orchestration of domain logic
- `infrastructure/` — adapters (HTTP handlers/router, storage implementations, migrations)
- `cmd/` — executable entrypoints (`server`, `migrate`, `producer`, `consumer`)

## Dependency Direction

Keep dependencies flowing inward:

- `domain` should not import `infrastructure`
- `application` depends on `domain` abstractions
- `infrastructure` implements `domain`/`application` contracts
- `cmd` composes and wires dependencies

## Request Flow (HTTP)

1. Router receives request (`infrastructure/api/router`)
2. Handler validates/parses input (`infrastructure/api/handler`)
3. Application use-case executes (`application/command` or `application/query`)
4. Domain rules/entities are applied (`domain/*`)
5. Repository/storage adapter persists/reads (`infrastructure/storage`)
6. Handler maps output to API response

## Multi-Tenancy

The project includes shared and tenant migration directories:

- `src/infrastructure/migrations/shared`
- `src/infrastructure/migrations/tenant`

When schema changes are tenant-specific, update tenant migrations accordingly.
