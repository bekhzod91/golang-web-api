# Testing Guide

This project uses **black-box API tests**: we test behavior through HTTP requests the same way a real client would.

## Goals

- Test features **through the public API**, not by calling internal services/SDKs directly.
- Keep tests **deterministic** using fixtures.
- Keep tests **declarative** and readable (Arrange → Act → Assert).
- Prevent tests from affecting each other (isolation by fixture sets).

---

## Test layout

All tests live under `src/tests/`.

```text
src/
  tests/
    test_auth/
      fixtures/
        TestMe/
          role.json
          user.json
        TestSignIn/
          ...
      me_test.go
      sign_in_test.go
    test_user/
      fixtures/
        ...
      create_test.go
      list_test.go
    test_role/
      fixtures/
        ...
```

### Folder rules

- One folder per feature/module:
  - `src/tests/test_auth`
  - `src/tests/test_user`
  - `src/tests/test_role`
- One file per endpoint/use-case:
  - `me_test.go`, `sign_in_test.go`, `create_test.go`, `list_test.go`, etc.
- Each test (or test file) must use **its own fixture folder**:
  - `fixtures/TestMe/`, `fixtures/TestSignIn/`, etc.

Why separate fixtures per test?

- Changing fixtures for one scenario should **not** break other tests.
- Each test should communicate **exactly what data it needs**.

---

## Test app (test harness)

Use the provided helper to spin up the test environment:

```go
testApp := tests.NewTestApp(t)
```

The test app is responsible for:

- starting the API server/router
- preparing test storage
- executing HTTP requests
- loading fixtures (including tenant-scoped fixtures)
- helper methods like `Authenticate(...)`

---

## Fixtures

Fixtures are JSON files loaded into storage before the test runs.

Example fixture structure:

```json
[
  {
    "table": "user",
    "pk": 1,
    "fields": {
      "id": 1,
      "email": "admin@example.com",
      "password": "$2a$14$...",
      "status": "active",
      "first_name": "Adam",
      "last_name": "Smith",
      "role_ids": [1, 2, 3]
    }
  }
]
```

### Fixture rules

- Fixtures must be **minimal**: only include fields required for the scenario.
- Prefer stable IDs (`pk`, `id`) so assertions remain predictable.
- Keep fixture folders named by test scenario:
  - `fixtures/TestMe/`
  - `fixtures/TestCreateUser/`
  - `fixtures/TestUserList/`

Load fixtures explicitly inside each test:

```go
testApp.LoadFixtureTenant([]string{
  "fixtures/TestMe/user.json",
  "fixtures/TestMe/role.json",
})
```

---

## How to write a test (Arrange → Act → Assert)

### Example: `GET /api/v1/me/`

```go
func TestMe(t *testing.T) {
  // Arrange
  testApp := tests.NewTestApp(t)
  testApp.LoadFixtureTenant([]string{
    "fixtures/TestMe/user.json",
    "fixtures/TestMe/role.json",
  })
  token := testApp.Authenticate("admin@example.com")

  // Act
  req, _ := http.NewRequest(http.MethodGet, "/api/v1/me/", nil)
  req.Header.Set("Authorization", token)
  req.Header.Set("X-Tenant", "test_tenant")
  rr := testApp.ExecuteRequest(req)

  // Assert (HTTP)
  require.Equal(t, http.StatusOK, rr.Code)

  // Assert (Response DTO)
  response := dto.MeResponseDTO{}
  tests.BindJSON(rr, &response)
  require.EqualValues(t, "Adam", response.FirstName)
  require.EqualValues(t, "Smith", response.LastName)
  require.EqualValues(t, "admin@example.com", response.Email)
}
```

---

## Golden rules

### 1) Test as a black box (API-first)

✅ Good:

- call HTTP endpoints via `ExecuteRequest`
- assert status code + response JSON

❌ Avoid:

- calling internal service methods directly
- calling SDKs directly
- writing “unit tests” against infrastructure details in `src/tests`

Example (good):

```go
req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/create/", bytes.NewBuffer(body))
req.Header.Set("Authorization", token)
req.Header.Set("X-Tenant", "test_tenant")
rr := testApp.ExecuteRequest(req)
require.Equal(t, http.StatusCreated, rr.Code)
```

---

### 2) Tests must be declarative

Tests should read like a spec. Prefer explicit assertions.

✅ Do **not** use loops and `if` statements for assertions.

Instead of:

```go
for _, role := range user.Roles {
  if role.ID == 1 { ... }
}
```

Do:

```go
require.EqualValues(t, 1, user.Roles[0].ID)
require.EqualValues(t, 2, user.Roles[1].ID)
require.EqualValues(t, 3, user.Roles[2].ID)
```

---

### 3) Always set tenant headers when required

Every tenant-scoped request must include:

- `X-Tenant: test_tenant`
- `Authorization: <token>` (if endpoint is protected)

This keeps multi-tenant behavior consistent and prevents false positives.

---

### 4) Verify database state only when it adds value

Prefer verifying outcomes through the API response. If you need to ensure persistence/state changes, you may query storage using repositories:

```go
dbUser := testApp.Storage().User().GetUserByID(response.ID)
```

Use this when:

- you need to validate side effects not visible in response
- you want to confirm DB fields were written correctly

---

## Naming conventions

- Test function names start with `Test` and describe the scenario:
  - `TestMe`
  - `TestSignIn_Success`
  - `TestCreateUser_ValidationError`
- Fixture folder names match test scenario name:
  - `fixtures/TestSignInSuccess/`
  - `fixtures/TestCreateUserValidationError/`

---

## When adding a new endpoint test

1. Create a folder: `src/tests/test_<module>/`
2. Add a test file: `<endpoint>_test.go`
3. Add fixtures:
   - `fixtures/<TestScenarioName>/...json`
4. Write the test:
   - Arrange fixtures + auth
   - Act with HTTP request
   - Assert status + DTO
   - (Optional) assert DB state
