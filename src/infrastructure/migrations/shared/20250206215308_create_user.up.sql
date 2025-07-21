CREATE TABLE IF NOT EXISTS "user" (
    id BIGSERIAL primary key,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    photo VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS "user_tenant" (
    id BIGSERIAL primary key,
    user_id BIGSERIAL,
    tenant_id BIGSERIAL,
    CONSTRAINT user_tenant_user_fk FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE,
    CONSTRAINT user_tenant_tenant_fk FOREIGN KEY (tenant_id) REFERENCES "tenant" (id) ON DELETE CASCADE
)