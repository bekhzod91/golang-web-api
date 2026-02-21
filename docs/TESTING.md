# TESTING

## Test Layout

Integration-style tests are organized in `src/tests/` by domain area:

- `test_auth`
- `test_role`
- `test_user`
- `test_permissions`
- `test_upload_file`

## Run Tests

Use Makefile target:

```bash
make test
```

This currently runs core suites listed in the Makefile.

## Testing Guidelines

- Add tests for every behavior change
- Prefer black-box/API-level tests for endpoint behavior
- Keep fixtures near each suite (`fixtures/` folders)
- Update tests together with OpenAPI or DB changes
