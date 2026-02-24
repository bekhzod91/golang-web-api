# API

## API Contract Source

OpenAPI specs live in `src/openapi/`.

- Main spec entry: `src/openapi/_swagger.yaml`
- Endpoint files are split by resource/action (e.g. `users_list.yaml`, `roles_create.yaml`)

## Workflow

1. Update or add OpenAPI spec files under `src/openapi/`
2. Regenerate DTO/models:

```bash
make swag
```

3. Implement handler + use-case changes in Go code
4. Add/update tests in `src/tests/`

## Conventions

- Keep request/response contracts in spec first
- Prefer resource-focused endpoint files
- Keep handler layer thin; business rules should stay in application/domain
