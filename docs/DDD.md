# ARCHITECTURE

## How to write code in this project (DDD workflow)

When adding new code, your first question should be:

> “Where should this rule live so the domain stays consistent even if we change the UI/DB/framework?”

We follow DDD + Clean Architecture:

- **Domain**: business rules (no DB, no HTTP, no frameworks)
- **Application**: use-cases / orchestration
- **Infrastructure**: DB/Redis/HTTP clients, implementations
- **Interface (Delivery)**: HTTP handlers/controllers, DTOs

### Priority: where to place logic

#### 1) Value Objects (smallest block — “brick”)

A Value Object represents a meaningful domain concept **without identity**, usually immutable, and responsible for its own invariants.

**Put logic here if it is about that value only:**
- validation (format/range)
- normalization
- comparisons
- small domain behavior

**Examples:**
- `Email`, `PhoneNumber`, `Money`, `Password`, `TenantID`

✅ Example:
- `user.Password.Verify("test123")`

Rule of thumb:
- If you can describe it as “rules of this value”, it belongs in the Value Object.

---

#### 2) Entities (identity + consistency across multiple fields)

Entities have identity and lifecycle. They protect invariants that involve multiple fields / VOs.

**Put logic here if:**
- it changes the entity’s state
- it must keep the entity consistent
- it uses multiple fields/VOs of the same entity

✅ Examples:
- `User.ChangePassword(...)`
- `Order.AddStop(...)`
- `Invoice.Approve()` / `Reject(reason)`

Rule of thumb:
- If the entity must remain valid after the change, the method belongs on the entity.

---

#### 3) Domain Services (domain logic across multiple entities)

Sometimes a rule doesn’t belong to one entity/VO.

**Use a Domain Service when:**
- the logic spans multiple entities
- it’s still pure domain (no DB/HTTP)
- there is no clear single owner

✅ Examples:
- `PricingCalculator`
- `SettlementCalculator`
- `DispatchMatchingService`

Rule of thumb:
- If it’s domain logic but “ownerless”, make it a domain service.

---

#### 4) Application Services / Use-cases (orchestration)

This layer coordinates steps for a user action:
- load aggregate(s) from repository
- call domain methods (VO/Entity/Domain Service)
- persist changes
- publish events (if you use them)
- manage transaction boundaries

✅ Examples:
- `CreateOrder`
- `AssignDriverToTrip`
- `ApproveInvoice`

Rule of thumb:
- Application service = “what happens when user does X”

Keep this layer thin: **no deep business rules**, those belong to the domain.

#### Optional style: isolate use-case-specific helper methods

When a use-case needs internal helper methods (for example, user lookup) and the logic differs from other use-cases, isolate that logic inside a dedicated use-case service struct.

- Keep the public entrypoint as `func UseCase(ctx core.IContext, ...)`.
- Create an internal struct that stores `ctx` for the use-case.
- Use intent-specific method names such as `loadUserForPasswordChange` or `loadUserForSignIn` (avoid generic shared names like `getUser` across different use-cases).

Simple example:

```go
func ChangeUserPassword(ctx core.IContext, id int64, req dto.ChangePasswordUserRequestDTO) (*dto.ChangePasswordUserResponseDTO, error) {
	s := &changeUserPassword{ctx: ctx}
	user, err := s.loadUserForPasswordChange(id)
	if err != nil {
		return nil, err
	}

	if err := s.changePassword(user, req.NewPassword); err != nil {
		return nil, err
	}

	return &dto.ChangePasswordUserResponseDTO{ID: user.ID}, nil
}

type changeUserPassword struct {
	ctx core.IContext
}

func (s *changeUserPassword) loadUserForPasswordChange(id int64) (*entity.User, error) {
	user, err := s.ctx.Storage().User().GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if user.Status != value_object.StatusActive {
		return nil, exception.New("User is inactive. Please contact support.")
	}

	return user, nil
}

func (s *changeUserPassword) changePassword(user *entity.User, newPassword string) error {
	password, err := value_object.NewPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = password
	_, err = s.ctx.Storage().User().ChangePasswordUser(user)
	return err
}
```

Use this pattern only when it improves clarity; for simple flows, a single function is preferred.

---

#### 5) Repositories (persistence abstraction — NOT business logic)

Repositories exist to **load/save aggregates** and hide persistence details.

**Repositories should:**
- fetch and persist aggregates/entities
- map storage models ↔ domain models (often in infra)
- provide query methods (careful: avoid leaking persistence details into domain)

**Repositories should NOT:**
- contain business decisions (`if status == ... then ...`)
- enforce domain invariants (that’s domain’s job)

Rule of thumb:
- If the code needs DB/Redis/SQL/ORM knowledge, it’s infrastructure, not domain logic.

---

## Quick decision checklist

- Is it a rule about a single value? → **Value Object**
- Does it keep one entity consistent while changing state? → **Entity**
- Does it involve multiple domain objects but no infrastructure? → **Domain Service**
- Is it a use-case that coordinates steps + persistence? → **Application Service**
- Is it about storage/query/DB/Redis? → **Repository (infrastructure)**

---

## Error Handling

### Rule: always use `exception.New` for domain errors

Every error that originates from a **business rule violation** must be created with `exception.New` (or `exception.Errorf`) from `domain/exception`:

```go
import "github.com/hzmat24/api/domain/exception"

return exception.New("Your password is too weak. Please add special characters.")
return exception.Errorf("Status %q is not valid. Use \"active\" or \"inactive\".", value)
```

**Why:**
`exception.New` wraps the message in the `DomainError` sentinel. Handlers detect this with `exception.IsDomainException(err)` and automatically convert it to an HTTP **400 Bad Request** with the message surfaced directly to the client:

```json
{ "message": "Your password is too weak. Please add special characters." }
```

If you return a plain `errors.New` or `fmt.Errorf` (without wrapping `DomainError`), the handler will treat it as an infrastructure failure, log it, and return an opaque **HTTP 500** — the user sees nothing useful.

---

### The two-part message format

Every domain error message must answer two questions in one sentence:

| Part | Question answered | Example |
|---|---|---|
| **Why** | What rule was violated? | "Your password is too weak." |
| **How** | What should the user do to fix it? | "Please add special characters." |

Structure: **`"<What went wrong>. <How to fix it>."`**

✅ Good:

```
"Your password is too weak. Please add at least one special character."
"Status \"archived\" is not valid. Use \"active\" or \"inactive\"."
"Invalid date format. Please use YYYY-MM-DD."
"Email is already taken. Please use a different email address."
```

❌ Bad:

```
"invalid status"              // no how-to-fix
"error"                       // meaningless
"something went wrong"        // use this only for unexpected infra errors, not domain errors
"validation failed at field X" // developer-facing, not user-facing
```

---

### Error type reference

| Type / Constructor | Use when | HTTP result |
|---|---|---|
| `exception.New("…")` | Business rule violated — user can fix it | 400 Bad Request |
| `exception.Errorf("…", args)` | Same, with dynamic values in the message | 400 Bad Request |
| `exception.NotFoundError` | Aggregate/entity not found by ID | 404 Not Found |
| Plain `error` (no wrapping) | Infrastructure failure (DB down, network error) | 500 Internal Server Error |

#### `NotFoundError` example

```go
// domain/repository interface — implementation returns this on missing rows
return nil, exception.NotFoundError
```

The handler checks it separately:

```go
if exception.IsNotFoundException(err) {
    c.NotFound()
    return
}
```

Do **not** put a human message in `NotFoundError` — the 404 response body is always `{"message": "not found"}`.

---

### Where to create domain errors

Create errors **as close to the violated rule as possible**:

- **Value Object constructor** — if the rule is about a single value's format or range.
- **Entity method** — if the rule involves the entity's own state.
- **Application use-case** — if the rule requires reading from the repository first (e.g., uniqueness check).

```go
// ✅ Value Object — rule about the value itself
func ParseStatus(s string) (Status, error) {
    if s != "active" && s != "inactive" {
        return "", exception.New(`Invalid status. Use "active" or "inactive".`)
    }
    return Status(s), nil
}

// ✅ Application use-case — rule that needs DB lookup
func CreateUser(ctx core.IContext, req dto.CreateUserRequestDTO) (*dto.CreateUserResponseDTO, error) {
    existing, _ := ctx.Storage().User().GetUserByEmail(email)
    if existing != nil {
        return nil, exception.New("Email is already taken. Please use a different email address.")
    }
    // …
}
```

---

### Handler pattern (for reference)

Handlers must follow this exact error-check order:

```go
result, err := command.DoSomething(c, req)

if exception.IsNotFoundException(err) {
    c.NotFound()       // 404 — entity not found
    return
}
if exception.IsDomainException(err) {
    c.BadRequest(err)  // 400 — message goes straight to the client
    return
}
if err != nil {
    c.Logger().Error(err.Error())
    c.InternalServerError() // 500 — log it, hide details from client
    return
}

c.OK(result)
```

Never expose raw infrastructure error messages to the client.
