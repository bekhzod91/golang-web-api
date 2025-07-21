CREATE TABLE IF NOT EXISTS "user" (
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