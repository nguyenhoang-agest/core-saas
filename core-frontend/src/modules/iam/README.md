# IAM Module (Identity & Access Management)

## Purpose
Provides user lifecycle management, role-based access control (RBAC), invitations, and security audit log capabilities.

## Architecture
- `index.ts`: Registers routes (`/iam`), menu items, and permissions (`iam:view`, `iam:users:manage`, `iam:roles:manage`).
- `pages/iam-page.tsx`: Route entry point rendering user and team management views.
- `store/iam-store.tsx`: Context-safe Zustand store for tracking selection and dialog state.
- `api/iam-api.ts`: API clients for tenant user management.

## TODOs (Domain Implementation)
- [ ] Connect with backend `/api/v1/iam/users` and `/api/v1/iam/roles` endpoints.
- [ ] Implement User GetAllTenant table with pagination, search, and status badges.
- [ ] Implement Invite User modal form with email validation and role selector.
- [ ] Implement Role & Permission matrix editor with `<Can>` guards for write operations.
- [ ] Add session revocation and MFA status badges.
