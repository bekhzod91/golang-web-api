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
