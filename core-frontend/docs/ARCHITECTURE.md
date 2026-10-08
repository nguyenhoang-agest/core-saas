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
│   ├── iam/              # STUB: users, roles, invitations
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

- **Edge Resolution (`src/proxy.ts`)**: Next.js 16 edge proxy inspects the incoming `Host` header (e.g. `acme.saas.com`), resolves subdomain tenant IDs, and propagates `x-tenant-id` downstream via headers and cookies.
- **Axios HTTP Client (`src/apis/axios-client.ts`)**: Automatically attaches `X-Tenant-ID` header to all outgoing API requests.
- **Query Cache Isolation (`src/core/tenant/tenant-key.ts`)**: `tenantKey(tenantId, ['users', id])` prefixes all TanStack Query cache keys to prevent cross-tenant cache leakage in multi-tenant sessions.
- **Dynamic Theming**: Tenant colors dynamically set CSS custom properties (`--primary-color`, `--border-radius`) and drive Ant Design's `theme.token`.

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

- **`<Can permission="iam:manage" fallback={...}>`**: Client-side UI guard for conditionally rendering actions based on NextAuth session permissions (wildcard `*` supported).
- **Backend Authoritative Rule**: Client-side `<Can>` is solely a UX presentation convenience. Authoritative access control **must always** be validated and enforced by backend API services.
