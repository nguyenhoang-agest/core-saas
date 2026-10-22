# Architecture Overview: SaaS Starter Kit

A modular, production-ready frontend starter kit built with Next.js 16 (App Router), Ant Design v6, React Query v5, Zustand v5, NextAuth v4, and Vitest. Designed as a **single application repository** where teams clone the kit, customize branding, and add domain modules.

---

## 1. Directory Structure & Boundaries

```
src/
├── app/                  # Next.js App Router: thin routes rendering module pages
├── core/                 # Framework infra: registry, slots, overrides, tenancy, auth
├── modules/
│   ├── example/          # The only fleshed-out vertical slice
│   ├── user/             # Multi-tenant user directory, RBAC roles, and ACL permissions
│   ├── settings/         # STUB: branding, tenant config
│   ├── billing/          # STUB: plans, subscriptions, payment
│   └── _template/        # Blueprint for scaffolding new modules
├── apis/                 # Axios HTTP client with automatic X-Tenant-ID injection
├── components/           # AppShell layout with dynamic navigation & slots
├── proxy.ts              # Edge request interception (subdomain tenant resolution)
config/app.config.ts      # Central configuration: enabled modules & branding
scripts/                  # Project initialization and module scaffolding CLI
```

### Architectural Boundaries (Enforced via ESLint)
- **`core/`** must never import from `modules/*`. Core only provides infrastructure and contracts.
- **`modules/`** must never import internal paths of other modules (`@/modules/<other>/*`). Cross-module interaction happens strictly through module public entry points (`@/modules/<other>`) or slot contributions.

---

## 2. Core Infrastructure

### Module Registry & Contracts
Each domain module declares a `ModuleDefinition` object:
```ts
export interface ModuleDefinition {
  id: string;                       // Unique identifier (e.g. 'iam')
  name: string;                     // Human-readable title
  description?: string;
  routes?: ModuleRoute[];           // Page components mapped to routes
  menuItems?: ModuleMenuItem[];     // Dynamic navigation items
  permissions?: string[];           // Domain permission keys
  slots?: SlotContribution[];       // Contributions injected into shell slots
}
```
Enabled modules are explicitly enumerated in `config/app.config.ts`. No filesystem scanning. Duplicate module or slot IDs are rejected at registration time.

### Extension Points: Slots & Overrides
- **`<Slot name="header-actions" props={...} />`**: Renders all contributed components for that slot sorted by `order` ascending. Supports fallback nodes.
- **Component Override Map (`<Override name="Logo" defaultComponent={DefaultLogo} />`)**: Allows replacing shell-level UI components without modifying core code.

### Composable `<PlatformProvider>`
Replaces monolithic provider blocks:
1. `SessionProvider` (NextAuth)
2. `QueryClientProvider` (TanStack Query)
3. `TenantProvider` (Runtime tenant config)
4. `<AntdRegistry>` (@ant-design/nextjs-registry): Extracts SSR styling for Ant Design v6
5. `<ConfigProvider>`: Dynamically configured with the tenant's primary colors & tokens
6. `<OverrideProvider>`: Shell component override mappings

---

## 3. Multi-Tenancy Subsystem

- **Edge Resolution (`src/proxy.ts`)**: Next.js 16 edge proxy intercepts incoming requests to `/api/*` and UI pages. It evaluates tenant identity in strict priority:
  1. `x-tenant-id` / `X-Tenant-ID` explicit request headers (from Axios or external callers).
  2. `?tenant=` query parameter override.
  3. Subdomain extraction from `Host` header (e.g. `acme.saas.com` -> `acme`).
  4. `x-tenant-id` cookie persisted across user sessions.
  5. Fallback to `"default"`.
  It injects the resolved `x-tenant-id` into downstream request headers and mirrors it to client cookies.
- **Topbar Tenant Switcher (`src/components/tenant-switcher.tsx`)**: Allows instant workspace switching from any page. Updates both cookies and `localStorage`, re-evaluates theme custom properties, and purges TanStack Query caches to eliminate cross-tenant data bleed.
- **Axios HTTP Client (`src/apis/axios-client.ts`)**: Automatically attaches `X-Tenant-ID` header to all outgoing API requests using client cookie state.
- **Query Cache Isolation (`src/core/tenant/tenant-key.ts`)**: `tenantKey(tenantId, ['users', id])` prefixes all TanStack Query cache keys to prevent cross-tenant cache leakage in multi-tenant sessions.
- **Dynamic Theming**: Tenant colors dynamically set CSS custom properties (`--primary-color`, `--border-radius`) and drive Ant Design's `theme.token`.
- **Isolation Modes (`src/core/tenant/tenant-context.tsx`)**: `TenantConfig` carries `isolateMode: TenantIsolateMode` (`'shared' | 'schema' | 'dedicated'`) indicating the backend data isolation tier for the current workspace.
- **Dedicated Tenant Security Governance (`src/modules/tenant/components/create-tenant-dialog.tsx`)**: When provisioning dedicated database tenants, users specify the target DSN, while connection pool limits (`10 connections`) and TLS security (`sslmode: require`) are strictly locked as read-only platform-managed badges to prevent insecure client overrides.
- **Single Source of Truth Error Translation**: The frontend maps backend machine-readable `message_code` values (e.g. `E031001` - `E031006`, `E032001` - `E032015`) directly to localized copy via `src/locales/{en,vi}/errors.json`, eliminating duplicate translation catalogs on the backend.

---

## 4. State Management (Zustand)

In Next.js App Router & React 19, global singleton stores cause cross-request state leakage during SSR.
Modules define stores via the **`createStoreContext`** factory:
```tsx
const { Provider, useStore } = createStoreContext(createModuleStore, 'MyStore');
```
Stores are instantiated per component tree via React Context, ensuring complete SSR hydration isolation.

---

## 5. Security & Authorization

- **System Role Protection & UI Guards**: Built-in baseline roles (`is_system = true`: `owner`, `admin`, `member`, `viewer`) display with a `System` lock badge. Edit and Delete actions are strictly disabled in the UI; users can only **Clone** system roles into new custom templates.
- **Role-Permission Matrix (`src/modules/user/components/role-permission-matrix.tsx`)**: Visual breakdown of module capabilities across roles with real-time permission counters.
- **`<Can permission="iam:manage" fallback={...}>`**: Client-side UI guard for conditionally rendering actions based on NextAuth session permissions (wildcard `*` supported).
- **Backend Authoritative Rule**: Client-side `<Can>` is solely a UX presentation convenience. Authoritative access control **must always** be validated and enforced by backend API services.

