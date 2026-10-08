# Adding a Domain Module

This guide details how to scaffold, implement, and register a new domain module in the SaaS Starter Kit.

---

## 1. Scaffold Automatically

Use the CLI scaffolding script:
```bash
# Using pnpm
pnpm new-module notifications

# Using npm
npm run new-module notifications
```

This single command:
1. Copies `src/modules/_template/` to `src/modules/notifications/`.
2. Renames and transforms tokens across stores, APIs, and pages.
3. Generates a smoke test in `src/modules/notifications/notifications.test.tsx`.
4. Creates a thin route wrapper in `src/app/notifications/page.tsx`.
5. Automatically registers the module in `config/app.config.ts`.

---

## 2. Module Anatomy

Each module is self-contained within `src/modules/<name>/`:

```
src/modules/<name>/
├── index.ts              # Module definition & public exports
├── pages/                # Page components (rendered by thin Next.js route wrappers)
│   └── <name>-page.tsx
├── components/           # Module-internal components & slot widgets
├── api/                  # API fetchers & TanStack Query hooks
│   └── <name>-api.ts
├── store/                # Hydration-safe Zustand store factory
│   └── <name>-store.tsx
├── <name>.test.tsx       # Smoke and unit tests
└── README.md             # Domain purpose, API contracts & TODOs
```

---

## 3. Implementation Steps

### A. Define Data Fetching (`api/<name>-api.ts`)
Always scope TanStack Query keys using the `tenantKey` helper and `apiClient`:

```ts
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/apis/axios-client';
import { tenantKey, useTenant } from '@/core/tenant';

export function useNotifications() {
  const { tenant } = useTenant();

  return useQuery({
    queryKey: tenantKey(tenant.id, ['notifications', 'list']),
    queryFn: async () => {
      const res = await apiClient.get('/notifications');
      return res.data;
    },
  });
}
```

### B. Define State Store (`store/<name>-store.tsx`)
Avoid global module singletons. Use `createStoreContext`:

```tsx
import { createStore } from 'zustand/vanilla';
import { createStoreContext } from '@/core/store';

interface NotificationState {
  unreadCount: number;
  setUnreadCount: (c: number) => void;
}

export function createNotificationStore() {
  return createStore<NotificationState>((set) => ({
    unreadCount: 0,
    setUnreadCount: (unreadCount) => set({ unreadCount }),
  }));
}

export const {
  Provider: NotificationStoreProvider,
  useStore: useNotificationStore,
} = createStoreContext(createNotificationStore, 'NotificationStore');
```

Wrap your page component in `NotificationStoreProvider` to ensure per-tree state isolation.

### C. Contribute to Shell Slots (`components/` & `index.ts`)
If your module contributes widgets to shell slots (such as `header-actions`):

```tsx
// components/notification-bell.tsx
export function NotificationBell() {
  return <Badge count={3}><BellOutlined /></Badge>;
}

// index.ts
export const notificationsModule: ModuleDefinition = {
  id: 'notifications',
  name: 'Notifications',
  slots: [
    {
      id: 'notifications-header-bell',
      slotName: 'header-actions',
      order: 15,
      component: NotificationBell,
    },
  ],
  // ...
};
```

### D. UI Permission Guards (`<Can>`)
Show or hide UI actions based on user permissions:

```tsx
import { Can } from '@/core/auth';

<Can permission="notifications:broadcast" fallback={<span>Read-only</span>}>
  <Button type="primary">Broadcast Notification</Button>
</Can>
```

---

## 4. Enable or Disable Modules

To toggle features, edit `config/app.config.ts`:

```ts
import { notificationsModule } from '@/modules/notifications';

export const appConfig = {
  enabledModules: [
    exampleModule,
    iamModule,
    settingsModule,
    billingModule,
    notificationsModule, // Add or remove here
  ],
};
```

---

## 5. Verify Your Module

Run the verification suite to ensure strict TypeScript, linting boundaries, and tests pass:

```bash
npx tsc --noEmit
npm run lint
npm test
npm run build
```
