CREATE TABLE IF NOT EXISTS "role" (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    permissions jsonb NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO "role" (
    "id",
    "code",
    "name",
    "description",
    "permissions"
) VALUES (
    1,
    'admin',
    'Admin',
    'An Admin manages user access, system settings, security, and overall operational integrity. Responsibilities include user management, permissions, monitoring, troubleshooting, and ensuring compliance.',
    '["view_user","create_user","update_user","delete_user","view_role","create_role","update_role","delete_role","view_client","create_client","update_client","delete_client","view_driver","create_driver","update_driver","delete_driver","view_settings","create_settings","update_settings","delete_settings","view_location","create_location","update_location","delete_location","view_vehicle","create_vehicle","update_vehicle","delete_vehicle","view_vendor","create_vendor","update_vendor","delete_vendor","create_tariff","update_tariff","view_tariff","delete_tariff","create_inquiry","update_inquiry","view_inquiry","delete_inquiry"]'::jsonb
)