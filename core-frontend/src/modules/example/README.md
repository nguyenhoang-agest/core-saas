# Example Module (Vertical Slice)

This module serves as the blueprint and only fleshed-out vertical slice in the Starter Kit. It illustrates how domain modules should be organized without introducing excessive business domain logic.

## Architecture

- **`index.ts`**: Module definition adhering to `ModuleDefinition`. Declares route (`/example`), navigation menu item, permissions (`example:view`, `example:write`), and slot contributions (`header-actions`).
- **`store/`**: Context-bound Zustand store created via `createStoreContext` from `@/core/store`. Ensures per-tree hydration safety and prevents state pollution across SSR requests.
- **`api/`**: Data fetching hooks via `@tanstack/react-query`. Every query key is partitioned by tenant ID via `tenantKey(tenant.id, key)`. Outgoing HTTP calls via `apiClient` automatically attach `X-Tenant-ID`.
- **`components/`**: Module-scoped UI components and slot widgets (`ExampleWidget`).
- **`pages/`**: Page views rendered by Next.js thin route wrappers.
- **`example.test.tsx`**: Unit & integration test verifying module rendering, store reactivity, and query hook.

## Integration

Registered in `config/app.config.ts`:
```ts
import { exampleModule } from '@/modules/example';

export const appConfig = {
  enabledModules: [exampleModule],
};
```
Rendered by thin Next.js route:
`src/app/example/page.tsx`
