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

---

## Database and Migrations

This project follows a **domain-first**, **migration-safe** approach to database design. The database is a persistence layer only: it stores state reliably and consistently, while **all business rules live in Go code**.

---

## Database rules

### 1) Enforce consistency with foreign keys and constraints

Always add **FOREIGN KEY** relations and constraints where they represent real domain relationships.

Use constraints for:

- parent-child relationships (`user → role`, `order → customer`, etc.)
- uniqueness (`email`, `tenant + external_id`, etc.)
- non-negative values (when applicable)

**Goal:** invalid data should be hard to store.

---

### 2) Standard audit + soft delete columns

Every table must include:

- `is_deleted` (soft delete)
- `created_at`
- `updated_at`

Add these when it makes sense:

- `created_by`
- `updated_by`

Use `created_by/updated_by` when:

- changes are made by users/operators
- audit trail matters for the feature

If the table is purely technical (e.g., background job locks), you may omit `created_by/updated_by`.

---

### 3) No business logic in the database

Database must **not** contain business logic such as:

- stored procedures
- triggers
- scheduled jobs inside DB
- “smart” computed logic that changes business meaning

Business rules belong in:

- Value Objects / Entities / Domain Services (domain layer)
- Use-cases (application layer)

DB should only enforce **integrity constraints** (FK/unique/check).

---

### 4) Avoid `NULL` fields (prefer defaults)

Avoid nullable columns whenever possible.

Prefer defaults such as:

- `0` for numeric fields
- `''` for strings
- `FALSE` for booleans
- empty arrays / JSON default (if supported/appropriate)

Why:

- simpler queries
- fewer runtime edge cases
- easier JSON/DTO mapping

> If a field is truly optional by domain meaning (e.g., `middle_name`, `photo_url`), `NULL` can be acceptable — but this should be the exception, not the default.

---

### 5) Use business-domain naming

Table and column names must match the domain language.

✅ Good:

- `driver_settlement`, `trip_stop`, `invoice_item`, `fuel_card_transaction`

❌ Bad:

- `tbl_user2`, `data1`, `misc`, `temp`, `mapping_table_abc`

---

### 6) Use `snake_case` everywhere

All database identifiers must use `snake_case`:

- table names: `user_role`, `driver_settlement`
- column names: `first_name`, `created_at`

---

## Migration rules

### 7) Migrations must be backward compatible

Migrations must not break older versions of the application during rollout.

**Do not:**

- delete columns
- rename columns
- change column types in a breaking way
- change semantics of existing fields

**Prefer additive changes:**

- add new columns with safe defaults
- add new tables
- add indexes
- add new constraints carefully (validate existing data first)

If you must “replace” a column:

1. add a new column
2. backfill data
3. update code to use the new column
4. (optional later) remove old column only after a long deprecation window

---

### 8) SQL migrations style: write keywords in UPPERCASE

SQL migration files must use uppercase SQL keywords:

✅

- `SELECT`, `INSERT`, `UPDATE`, `DELETE`
- `CREATE TABLE`, `ALTER TABLE`, `ADD CONSTRAINT`
- `NOT NULL`, `DEFAULT`, `FOREIGN KEY`

❌

- `select`, `create table`, `alter table`

This keeps migrations consistent and readable.

---

## Recommended checklist for new tables

When creating a new table, verify:

- [ ] Table name is domain-based and `snake_case`
- [ ] Columns use `snake_case`
- [ ] Includes `is_deleted`, `created_at`, `updated_at`
- [ ] Adds `created_by/updated_by` when meaningful
- [ ] Uses `NOT NULL` + `DEFAULT` instead of `NULL` where possible
- [ ] Adds FK + constraints for consistency
- [ ] No procedures/triggers/business logic
- [ ] Migration is backward compatible
- [ ] SQL keywords are uppercase
