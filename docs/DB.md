# DB

## Database Structure

Migrations are split by scope:

- Shared schema: `src/infrastructure/migrations/shared`
- Tenant schema: `src/infrastructure/migrations/tenant`

Use the matching migration folder based on data ownership and tenancy boundaries.

## Migration Commands

Create migration files:

```bash
make migration-shared title=<name>
make migration-tenant title=<name>
```

Apply migrations:

```bash
make migrate
```

## SQL Code Generation

SQLC config: `src/sqlc.yaml`

Regenerate query code:

```bash
make sqlc
```

If query files or schema references change, regenerate before running tests.
