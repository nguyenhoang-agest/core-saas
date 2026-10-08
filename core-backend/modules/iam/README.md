# IAM Module (STUB)

## Purpose
Manages identities, credentials, role-based access control (RBAC), and authentication tokens (JWT / session cookies).

## Planned Ports & Interfaces
- **Ports (Domain)**:
  - `UserRepository`: Store and query user accounts by ID and email.
  - `PasswordHasher`: Password hashing and verification (bcrypt/argon2).
  - `TokenGenerator`: Issuing and validating signed tokens.
- **Adapters**:
  - `PostgresUserRepository`: Database persistence with tenant isolation.
  - `BcryptHasher`: Password hashing implementation.
  - `JwtTokenService`: EdDSA/RS256/HMAC token signing.
- **Events Emitted**:
  - `iam.user.registered`: On successful user registration.
  - `iam.user.authenticated`: On login.
  - `iam.password.reset`: On password changes.

## TODO Checklist
- [ ] Implement robust password hashing with `golang.org/x/crypto/bcrypt`
- [ ] Implement JWT token generator and validator with expiration claims
- [ ] Implement `PostgresUserRepository` CRUD with pgx
- [ ] Add auth middleware to validate tokens and inject `auth.Claims` into `context.Context`
- [ ] Add password reset flow and email verification hooks
- [ ] Add rate limiting on `/iam/login`
