----------------------Tenant-------------------------
CREATE TABLE "tenant" (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    schema_name VARCHAR(255) UNIQUE NOT NULL,
    is_deleted BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);


-- name: GetTenantBySchemaName :one
SELECT * FROM "tenant"
WHERE "schema_name" = $1 LIMIT 1;

-- name: GetAllTenants :many
SELECT * FROM "tenant";

-- name: CreateTenant :exec
INSERT INTO "tenant" ("name", "status", "schema_name", "is_deleted", updated_at, created_at)
VALUES ($1, "active", $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
---------------------User----------------------------
CREATE TABLE "user" (
    id BIGSERIAL primary key,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    photo VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    birth_date DATE NOT NULL,
    last_login TIMESTAMP NOT NULL,
    role_ids   jsonb NOT NULL,
    is_deleted BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
-- name: GetUsers :many
SELECT * FROM "user" WHERE "is_deleted" = false ORDER BY "id" DESC ;

-- name: GetUserByID :one
SELECT * FROM "user" WHERE "is_deleted" = false AND "id" = $1;

-- name: GetUserByEmail :one
SELECT * FROM "user" WHERE "is_deleted" = false AND "email" = $1;

-- name: CreateUser :one
INSERT INTO "user" (
    "email",
    "password",
    "first_name",
    "last_name",
    "photo",
    "status",
    "last_login",
    "role_ids",
    "is_deleted",
    "updated_at",
    "created_at",
    "phone",
    "birth_date"
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    current_timestamp,
    $7,
    false,
    current_timestamp,
    current_timestamp,
    $8,
    $9
)
RETURNING "id";

-- name: UpdateUser :exec
UPDATE "user"
SET
    "first_name" = $1,
    "last_name" = $2,
    "photo" = $3,
    "status" = $4,
    "role_ids" = $5,
    "phone" = $6,
    "birth_date" = $7,
    "updated_at" = current_timestamp
WHERE "is_deleted" = false AND "id" = $8;

-- name: DeleteUserByID :exec
UPDATE "user" SET "is_deleted" = true WHERE "is_deleted" = false AND "id" = $1;

-- name: ChangePasswordUser :exec
UPDATE "user" SET
    "password" = $1,
    "updated_at" = current_timestamp
WHERE "is_deleted" = false AND "id" = $2;

--------------------Role------------------------------
CREATE TABLE "role" (
    id BIGSERIAL primary key,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    permissions jsonb NOT NULL,
    is_deleted BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- name: GetRoles :many
SELECT * FROM "role"
WHERE "is_deleted" = false AND
      (NOT @has_search::bool OR (role.name ilike '%' || @search::varchar || '%' OR role.code ilike '%' || @search::varchar || '%')) AND
      (NOT @has_created_at__lte::bool OR role.created_at >= @created_at__lte::timestamp) AND
      (NOT @has_created_at__gte::bool OR role.created_at <= @created_at__gte::timestamp)
ORDER BY
    CASE WHEN @order_by::varchar = 'id_asc' THEN role.id END,
    CASE WHEN @order_by::varchar = 'id_desc' THEN role.id END DESC,
    CASE WHEN @order_by::varchar = 'created_at_asc' THEN role.created_at END,
    CASE WHEN @order_by::varchar = 'created_at_desc' THEN role.created_at END DESC,
    CASE WHEN @order_by::varchar = 'code_asc' THEN role.id END,
    CASE WHEN @order_by::varchar = 'code_desc' THEN role.id END DESC,
    CASE WHEN @order_by::varchar = 'name_asc' THEN role.id END,
    CASE WHEN @order_by::varchar = 'name_desc' THEN role.id END DESC
LIMIT $1 OFFSET $2;

-- name: GetRoleCount :one
SELECT count(*) FROM "role"
WHERE "is_deleted" = false AND
    (NOT @has_search::bool OR (role.name ilike '%' || @search::varchar || '%' OR role.code ilike '%' || @search::varchar || '%')) AND
    (NOT @has_created_at__lte::bool OR role.created_at >= @created_at__lte::timestamp) AND
    (NOT @has_created_at__gte::bool OR role.created_at <= @created_at__gte::timestamp);

-- name: GetRoleByID :one
SELECT * FROM "role"
WHERE "is_deleted" = false AND "id" = $1 LIMIT 1;

-- name: GetRolesByIDs :many
SELECT * FROM "role"
WHERE "is_deleted" = false AND "id" = ANY (@ids::bigint[]) ORDER BY "id";

-- name: GetRoleByCode :one
SELECT * FROM "role"
WHERE "is_deleted" = false AND "code" = $1 LIMIT 1;


-- name: CreateRole :one
INSERT INTO "role" ("name", "code", "description", "permissions") VALUES ($1, $2, $3, $4)
RETURNING "id";

-- name: UpdateRole :exec
UPDATE "role"
SET "name" = $1, "code" = $2, "description" = $3, "permissions" = $4
WHERE "is_deleted" = false AND "id" = $5;

-- name: DeleteRoleByID :exec
UPDATE "role" SET "is_deleted" = true WHERE "is_deleted" = false AND "id" = $1;
