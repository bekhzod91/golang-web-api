CREATE TABLE IF NOT EXISTS "tenant" (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    schema_name VARCHAR(255) UNIQUE NOT NULL,
    is_deleted BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

INSERT INTO "tenant" (
        "name",
        "status",
        "schema_name",
        "is_deleted",
        "created_at",
        "updated_at"
) VALUES (
        'Tenant 1',
        'active',
        'test_tenant',
        false,
        current_timestamp,
        current_timestamp
);