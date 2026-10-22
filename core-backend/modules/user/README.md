# User & IAM Module

## Purpose
Manages tenant members, user authentication identities, multi-tenant Role-Based Access Control (RBAC), fine-grained Access Control Lists (ACL), and real-time distributed authorization cache invalidation.

---

## Architecture & Layout

```
modules/user/
├── api/                        # Public inter-module contracts & EventBus topics
│   └── events.go               # RoleCreated, RoleUpdated, RoleDeleted, UserRolesChanged
├── internal/
│   ├── dto/                    # Request & response DTOs (user_dto.go, role_dto.go)
│   ├── errcode/                # Module-scoped error code offsets (E032001+)
│   ├── handler/                # Fiber HTTP controllers (handler.go, role_handler.go)
│   ├── model/                  # GORM models (User, Role, RolePermission, UserRole, AuthzVersion)
│   ├── repository/             # Database access layer with pessimistic locking & transaction management
│   │   ├── user_repo.go        # User CRUD & active status management
│   │   └── role_repo.go        # Role CRUD, clone, assignment, and bumpAuthzVersionTx
│   └── service/                # Business orchestration and policy enforcement
│       ├── service.go          # UserService
│       ├── role_service.go     # RoleService, Subset Rule, Scope Validation
│       └── sync.go             # SyncSystemRoles and BackfillLegacyUserRoles
├── migrations/                 # Embedded Goose SQL migrations (00001, 00002, 00003)
└── module.go                   # Kernel module definition & event subscriptions
```

---

## Role-Based Access Control (RBAC) Architecture

### 1. The `is_system` Field in `roles`

The `roles` table features a boolean column:
```sql
is_system BOOLEAN NOT NULL DEFAULT FALSE
```

| `is_system` | Classification | Created By | Mutability | Auto-Sync on Upgrade |
| :---: | :--- | :--- | :---: | :---: |
| **`true`** | **System Baseline Role** (`owner`, `admin`, `member`, `viewer`) | System (via `SyncSystemRoles`) | **Immutable** (Protected from Edit & Delete) | Yes (via `SyncSystemRoles`) |
| **`false`** | **Custom Tenant Role** (e.g. `billing_auditor`, `editor`) | Tenant Administrator | **Fully Editable** (Edit, Delete, Clone) | No (Preserves tenant customization) |

#### Key Protections Enforced by `is_system = true`:
1. **Immutability & Anti-Lockout**:
   - `UpdateRole` rejects attempts to alter system roles (`CannotModifySystemRole`, `E032007`).
   - `DeleteRole` rejects attempts to remove system roles (`CannotModifySystemRole`, `E032007`).
   - `CreateRole` blocks reserving system role slugs (`IsSystemRoleSlug`).
2. **Auto-Sync & Evolution (`SyncSystemRoles`)**:
   - When new tenant workspaces are provisioned, or when new platform modules introduce new permissions, `SyncSystemRoles` automatically provisions or updates the permissions of system roles across all tenants.
   - Custom roles (`is_system = false`) are never altered during system updates.
3. **UI / UX Differentiation**:
   - System roles render with a `System` badge and lock icon (`<Lock>`), with Edit/Delete buttons hidden.
   - Users can **Clone** system roles to create editable custom variants.
   - In role listings, system roles are sorted first via `ORDER BY is_system DESC, name ASC`.

---

## Security Governance & Safety Invariants

### 1. Anti-Lockout Owner Protection
A tenant must **always** maintain at least one active user with the `owner` role.
- When assigning or revoking roles via `AssignUserRoles`:
  - If a user currently has the `owner` role and the update attempts to revoke it, the repository executes:
    ```sql
    SELECT count(DISTINCT ur.user_id)
    FROM user_roles ur JOIN roles r ON ur.role_id = r.id
    WHERE ur.tenant_id = ? AND r.slug = 'owner' AND r.delete_at IS NULL
    ```
  - If total remaining owners $\le 1$, the operation is aborted with `errcode.LastOwnerViolation` (`E032013`).

### 2. Privilege Escalation Prevention (Subset Rule)
An administrator can only grant permissions that they themselves possess.
- When creating, updating, or cloning roles, or when assigning roles to users:
  - The actor's effective permissions are compared against the target permissions.
  - If target permissions contain actions outside the actor's scope, the request is rejected with `errcode.PrivilegeEscalation` (`E032009`).
  - Wildcard `*` grants full escalation immunity (for `owner`).

### 3. Scope Boundaries (`ScopeTenant` vs `ScopePlatform`)
- `ScopeTenant` permissions can be assigned by any tenant administrator.
- `ScopePlatform` permissions (such as `tenant:create`, `tenant:delete`) can only be granted within host/platform-level tenant contexts (`isHost = true`).

---

## Real-Time Distributed Cache Invalidation

To achieve high-performance authorization checks without stale permission windows across distributed nodes:

1. **`authz_version` Table**:
   - Each tenant maintains a monotonic epoch counter in `authz_version`.
   - Any modification to roles, permissions, or user assignments increments `version++` within the same transaction.
2. **PostgreSQL LISTEN/NOTIFY**:
   - Inside the transaction, Go executes:
     ```sql
     SELECT pg_notify('authz_invalidation', tenant_id)
     ```
   - All backend monolith or microservice nodes listening to the channel instantly flush their in-memory authorization caches for that tenant.

---

## Inter-Module Lifecycle Integration

The user module subscribes to event bus topics published by the `tenant` module:
- **`tenant.created` (`TenantCreatedTopic`)**:
  - Automatically executes `SyncSystemRoles(ctx, db, tenantID, isHost)`.
  - Provisions the initial `owner` user account for the tenant workspace.
- **`tenant.suspended` / `tenant.deleted`**:
  - Invalidates active user sessions and clears cached authorization records.

---

## API Endpoints

### User Management
| Method | Endpoint | Permission Required | Description |
| :--- | :--- | :--- | :--- |
| `GET` | `/users` | `iam:view` | List tenant users with pagination & role filtering |
| `POST` | `/users` | `iam:users:manage` | Invite / create a new workspace user |
| `GET` | `/users/:id` | `iam:view` | Get user details and assigned roles |
| `PUT` | `/users/:id` | `iam:users:manage` | Update user metadata, name, and status |
| `DELETE` | `/users/:id` | `iam:users:manage` | Deactivate / soft-delete user (Anti-lockout enforced) |
| `PUT` | `/users/:id/roles` | `iam:users:manage` | Assign one or more roles to user (Subset rule enforced) |

### Role Management
| Method | Endpoint | Permission Required | Description |
| :--- | :--- | :--- | :--- |
| `GET` | `/roles` | `iam:roles:read` | List tenant roles (sorted by `is_system DESC, name ASC`) |
| `GET` | `/roles/:id` | `iam:roles:read` | Get role details, user count, and permissions |
| `POST` | `/roles` | `iam:roles:manage` | Create a custom role (`is_system = false`) |
| `PUT` | `/roles/:id` | `iam:roles:manage` | Update custom role (System roles protected) |
| `DELETE` | `/roles/:id` | `iam:roles:manage` | Delete custom role (Must have 0 assigned members) |
| `POST` | `/roles/:id/clone` | `iam:roles:manage` | Clone any role into a new custom role |
| `GET` | `/roles/meta` | None (Authenticated) | Get baseline system role definitions |
| `GET` | `/permissions/catalog` | None (Authenticated) | Inspect complete permission registry |
